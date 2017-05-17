package router

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io"
	"io/ioutil"
	"fmt"
	"bytes"
	"encoding/json"
	"strconv"
	"ForumDatabase/database"
	"net/http/cookiejar"
)

type Response struct {
	Status int `json:"status"`
}

var (
	server *httptest.Server = httptest.NewServer(Create(true))
	TYPE_JSON = "application/json"
)

func TestClear(t *testing.T) {
	db := database.MakeConnection(true)
	db.Exec("DROP TABLE block_records, posts, thread_posts, threads, user_posts, user_threads, users")
	db.Close()
}

func registerUser(user *database.User) Response {
	client := createClient()
	data := createJson(map[string]string{"username": user.Username, "password": user.Password})
	httpRes, _ := client.Post(server.URL + "/api/v1/users/new", TYPE_JSON, data)
	var response Response
	bindResponse(httpRes.Body, &response)
	return response
}

func loginWithCredentials(t *testing.T, client *http.Client, user *database.User) {
	credentials := createJson(map[string]string{"username": user.Username, "password": user.Password})
	_, err := client.Post(server.URL + "/auth/login", "application/json", credentials)
	if err != nil {
		t.Error("Error logging in: ", err)
	}
}

func TestRegister(t *testing.T) {
	response := registerUser(&database.TEST_USER1)
	if response.Status != http.StatusOK {
		t.Error("Expected new user registration to be successful")
	}
}

func TestRegister2(t *testing.T) {
	response := registerUser(&database.TEST_USER2)
	if response.Status != http.StatusOK {
		t.Error("Expected second user registration to be successful")
	}
}

func TestRegister3(t *testing.T) {
	response := registerUser(&database.TEST_USER1)
	if response.Status == http.StatusOK {
		t.Error("Expected duplicate user registration to fail")
	}
}

func TestLogin(t *testing.T) {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
}

func TestLoginAndCreateThread(t *testing.T) {

	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	postData := createJson(map[string]string{"title": "Eyyyyyyyyyyyyy", "content": "I'm not too sure what should go here but here's some more input"})
	postResponse, postError := client.Post(server.URL + "/api/v1/threads/new", "application/json", postData)

	if postError != nil {
		t.Error("Error creating new thread", postError)
	} else {
		bodyString := getBodyString(postResponse.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestReadLastThreads(t *testing.T) {
	timestampString := strconv.FormatInt(database.MakeTimestamp(), 10)
	response, err := http.Get(server.URL + "/api/v1/threads/latest?timestamp=" + timestampString + "&limit=1")
	defer response.Body.Close()
	if err != nil {
		t.Error("Error getting latest threads: ", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}
}

func TestReadLatestThreads2(t *testing.T) {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	timestampString := strconv.FormatInt(database.MakeTimestamp(), 10)
	response, err := client.Get(server.URL + "/api/v1/threads/latest?timestamp=" + timestampString + "&limit=10")

	if err != nil {
		t.Error("Error getting latest threads ", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}
}

func TestReadPostsForThread(t *testing.T) {

	timestampString := strconv.FormatInt(database.MakeTimestamp(), 10)
	response, err := http.Get(server.URL + "/api/v1/threads/responses/2?timestamp=" + timestampString + "&limit=4")
	defer response.Body.Close()

	fmt.Println(response.StatusCode)

	if err != nil {
		t.Error("Error getting posts for thread: ", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestLoginAndReply(t *testing.T)  {

	client := createClient()

	loginWithCredentials(t, client, &database.TEST_USER1)
	postData := createJson(map[string]string{"content": "Shut it down! Now!"})
	postResponse, err := client.Post(server.URL + "/api/v1/threads/reply/4", "application/json", postData)

	if err != nil {
		t.Error("Error responding to thread", err)
	} else {
		bodyString := getBodyString(postResponse.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestBlockUser(t *testing.T) {

	client := createClient()

	loginWithCredentials(t, client, &database.TEST_USER1)
	response, err := client.Post(server.URL + "/api/v1/users/block/7", "application/json", nil)

	if err != nil {
		t.Error("Error blocking user", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestUnblockUser(t *testing.T) {

	client := createClient()

	loginWithCredentials(t, client, &database.TEST_USER1)
	response, err := client.Post(server.URL + "/api/v1/users/unblock/2", "application/json", nil)

	if err != nil {
		t.Error("Error unblocking user", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestDeleteThread(t *testing.T) {

	client := createClient()

	loginWithCredentials(t, client, &database.TEST_USER1)
	if response, err := client.Post(server.URL + "/api/v1/threads/delete/3", "application/json", nil); err != nil {
		t.Error("Error deleting thread: ", err)
	} else {
		jsonPrintBody(response.Body)
	}

}

func TestDeletePost(t *testing.T) {

	client := createClient()

	loginWithCredentials(t, client, &database.TEST_USER1)
	if response, err := client.Post(server.URL + "/api/v1/posts/delete/3", "application/json", nil); err != nil {
		t.Error("Request error: ", err)
	} else {
		jsonPrintBody(response.Body)
	}

}

func createClient() *http.Client {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	return client
}

func jsonPrintBody(byteBody io.ReadCloser) {
	bodyString := getBodyString(byteBody)
	fmt.Print(jsonPrettyPrint(bodyString))
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func getBodyString(byteBody io.ReadCloser) string {
	body, err := ioutil.ReadAll(byteBody)
	if err == nil {
		return string(body)
	} else {
		return ""
	}
}

func createJson(data interface{}) *bytes.Buffer {
	jsonBytes, _ := json.Marshal(data)
	return bytes.NewBuffer(jsonBytes)
}

func bindResponse(data io.ReadCloser, response *Response) error {
	if byteArray, readErr := ioutil.ReadAll(data); readErr != nil {
		return readErr
	} else {
		return json.Unmarshal(byteArray, &response)
	}
}