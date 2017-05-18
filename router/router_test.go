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
	database.Setup(db)
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

func createNewThread(client *http.Client, thread *database.Thread) Response {
	data := createJson(map[string]string{"title": thread.Title, "content": thread.Content})
	httpRes, _ := client.Post(server.URL + "/api/v1/threads/new", TYPE_JSON, data)
	var response Response
	bindResponse(httpRes.Body, &response)
	return response
}

func respondToThread(client *http.Client, threadId int, post *database.Post) Response {
	data := createJson(map[string]string{"content": post.Content})
	httpRes, _ := client.Post(server.URL + "/api/v1/threads/reply/" + strconv.Itoa(threadId), TYPE_JSON, data)
	var response Response
	bindResponse(httpRes.Body, &response)
	return response
}

func blockUserWithId(client *http.Client, userId int) Response {
	httpRes, _ := client.Post(server.URL + "/api/v1/users/block/" + strconv.Itoa(userId), TYPE_JSON, nil)
	var response Response
	bindResponse(httpRes.Body, &response)
	return response
}

func unblockUserWithId(client *http.Client, userId int) Response {
	httpRes, _ := client.Post(server.URL + "/api/v1/users/unblock/" + strconv.Itoa(userId), TYPE_JSON, nil)
	var response Response
	bindResponse(httpRes.Body, &response)
	return response
}

func deleteThreadWithId(client *http.Client, threadId int) Response {
	httpRes, _ := client.Post(server.URL + "/api/v1/threads/delete/" + strconv.Itoa(threadId), TYPE_JSON, nil)
	var response Response
	bindResponse(httpRes.Body, &response)
	return response
}

func deletePostWithId(client *http.Client, postId int) Response {
	httpRes, _ := client.Post(server.URL + "/api/v1/posts/delete/" + strconv.Itoa(postId), TYPE_JSON, nil)
	var response Response
	bindResponse(httpRes.Body, &response)
	return response
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

func TestCreateThread(t *testing.T) {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	thread := database.Thread{Title: "Eyyyyyyyyyyyyyyyy!", Content: "Not too sure what should be going here but here's some kind of text anyway"}
	response := createNewThread(client, &thread)
	if response.Status != http.StatusOK {
		t.Error("Unexpected issue creating a new thread")
	}
}

func TestReplyToThread(t *testing.T)  {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	post := database.Post{Content: "Shut it down! Now!"}
	response := respondToThread(client, 1, &post)
	if response.Status != http.StatusOK {
		t.Error("Unexpected issue responding to thread")
	}
}

func TestBlockUser(t *testing.T) {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	response := blockUserWithId(client, 2)
	if response.Status != http.StatusOK {
		t.Error("Unexpected issue blocking user")
	}
}

func TestUnblockUser(t *testing.T) {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	response := unblockUserWithId(client, 2)
	if response.Status != http.StatusOK {
		t.Error("Unexpected issue unblocking user")
	}
}

func TestDeletePost(t *testing.T) {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	response := deletePostWithId(client, 1)
	if response.Status != http.StatusOK {
		t.Error("Unexpected issue deleting post")
	}
}

/*
func TestDeleteThread(t *testing.T) {
	client := createClient()
	loginWithCredentials(t, client, &database.TEST_USER1)
	response := deleteThreadWithId(client, 1)
	if response.Status != http.StatusOK {
		t.Error("Unexpected issue deleting thread")
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
*/

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