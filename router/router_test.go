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


func TestReadLastThreads(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()

	timestampString := strconv.FormatInt(database.MakeTimestamp(), 10)
	response, err := http.Get(server.URL + "/api/v1/threads/latest?timestamp=" + timestampString)
	defer response.Body.Close()

	if err != nil {
		t.Error("Error getting latest threads: ", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestReadLatestThreads2(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()
	client := createClient()

	loginWithCredentials(t, server, client, database.TEST_USER1.Username, database.TEST_USER1.Password)

	timestampString := strconv.FormatInt(database.MakeTimestamp(), 10)
	response, err := client.Get(server.URL + "/api/v1/threads/latest?timestamp=" + timestampString)

	if err != nil {
		t.Error("Error getting latest threads ", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}


func TestReadPostsForThread(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()

	timestampString := strconv.FormatInt(database.MakeTimestamp(), 10)
	response, err := http.Get(server.URL + "/api/v1/threads/responses/2?timestamp=" + timestampString + "&limit=5")
	defer response.Body.Close()

	fmt.Println(response.StatusCode)

	if err != nil {
		t.Error("Error getting posts for thread: ", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestSessionManager(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	timestampString := strconv.FormatInt(database.MakeTimestamp(), 10)
	client.Get(server.URL + "/api/v1/threads/latest?timestamp=" + timestampString)
	response, err := client.Get(server.URL + "/api/v1/posts/latest")

	if err != nil {
		t.Error("Error testing session manager", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}


func loginWithCredentials(t *testing.T, server *httptest.Server, client *http.Client, username string, password string) {

	credentials := createJson(map[string]string{"username": username, "password": password})
	_, err := client.Post(server.URL + "/auth/login", "application/json", credentials)

	if err != nil {
		t.Error("Error logging in: ", err)
	}

}

func TestAbortBehaviour(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()
	client := createClient()

	response, _ := client.Post(server.URL + "/api/v1/users/test", "application/json", nil)

	fmt.Println("Status code: ", response.StatusCode)
	bodyString := getBodyString(response.Body)
	fmt.Print(jsonPrettyPrint(bodyString))

}

func TestRegister(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()
	client := createClient()

	data := createJson(map[string]string{"username": "Moreboy", "password": "more543543"})
	response, _ := client.Post(server.URL + "/api/v1/users/new", "application/json", data)

	bodyString := getBodyString(response.Body)
	fmt.Print(jsonPrettyPrint(bodyString))

}

func TestLoginAndCreateThread(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()
	client := createClient()

	loginWithCredentials(t, server, client, database.TEST_USER1.Username, database.TEST_USER1.Password)

	postData := createJson(map[string]string{"title": "Eyyyyyyyyyyyyy", "content": "I'm not too sure what should go here but here's some more input"})
	postResponse, postError := client.Post(server.URL + "/api/v1/threads/new", "application/json", postData)

	if postError != nil {
		t.Error("Error creating new thread", postError)
	} else {
		bodyString := getBodyString(postResponse.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestLoginAndReply(t *testing.T)  {

	server := httptest.NewServer(Create())
	defer server.Close()
	client := createClient()

	loginWithCredentials(t, server, client, database.TEST_USER1.Username, database.TEST_USER2.Password)

	postData := createJson(map[string]string{"content": "Shut it down! Now!"})
	postResponse, err := client.Post(server.URL + "/api/v1/threads/reply/7", "application/json", postData)

	if err != nil {
		t.Error("Error responding to thread", err)
	} else {
		bodyString := getBodyString(postResponse.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestBlockUser(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()
	client := createClient()

	loginWithCredentials(t, server, client, database.TEST_USER1.Username, database.TEST_USER2.Password)
	response, err := client.Post(server.URL + "/api/v1/users/block/7", "application/json", nil)

	if err != nil {
		t.Error("Error blocking user", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func TestUnblockUser(t *testing.T) {

	server := httptest.NewServer(Create())
	defer server.Close()
	client := createClient()

	loginWithCredentials(t, server, client, database.TEST_USER1.Username, database.TEST_USER1.Password)
	response, err := client.Post(server.URL + "/api/v1/users/unblock/2", "application/json", nil)

	if err != nil {
		t.Error("Error unblocking user", err)
	} else {
		bodyString := getBodyString(response.Body)
		fmt.Print(jsonPrettyPrint(bodyString))
	}

}

func createClient() *http.Client {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	return client
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