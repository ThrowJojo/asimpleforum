package database

import (
	"testing"
	"github.com/jinzhu/gorm"
	"fmt"
	"ForumDatabase/helpers"
)

var db *gorm.DB = MakeConnection(true)

func TestClear(t *testing.T) {
	db.Exec("DROP TABLE block_records, posts, thread_posts, threads, user_posts, user_threads, users")
}

func TestSetup(t *testing.T) {
	Setup(db)
}

func TestCreateUser(t *testing.T) {

	err := CreateUser(db, TEST_USER1.Username, TEST_USER1.Password)
	err2 := CreateUser(db, TEST_USER2.Username, TEST_USER2.Password)

	if err != nil {
		t.Error("Expected first attempt to be successful")
	}

	if err2 != nil {
		t.Error("Expected second attempt to be successful")
	}

}

func TestFindUser(t *testing.T) {
	_, err := FindUser(db, 1)
	if err != nil {
		t.Error("Expected to find user")
	}
}

func TestCountTotalThreads(t *testing.T) {
	count := CountTotalThreads(db)
	if count > 0 {
		t.Error("Expected 0 threads")
	}
}

func TestCreateThread(t *testing.T) {
	user, _ := FindUser(db, 1)
	_, err := CreateThread(db, user, "NEW THREAD 1111111", "dfgdfgffjngdkjdjgfdjkkfgjd")

	if err != nil {
		t.Error("Error creating thread", err)
	}

	count := CountTotalThreads(db)
	if count < 1 {
		t.Error("Expected more than 0 threads")
	}

}

func TestFindUser2(t *testing.T) {
	_, err := FindUser(db, 200000)
	if err == nil {
		t.Error("Expected a not exist error")
	}
}

func TestFindUserByCredentials(t *testing.T) {
	_, err := FindUserByCredentials(db, TEST_USER1.Username, TEST_USER1.Password)
	if err != nil {
		t.Error("Expected user object to be " + TEST_USER1.Username)
	}
}

func TestCountPostsForThread(t *testing.T) {
	count := CountPostsForThread(db, 1)
	if count > 0 {
		t.Error("Expected 0 posts for thread")
	}
}

func TestReplyToThread(t *testing.T) {
	user, _ := FindUser(db, 1)
	_, err := ReplyToThread(db, user, 1, "The six bears went to town adfsdfsdfsdfsdfdfs")
	if err != nil {
		t.Error("Unexpected error creating post", err)
	}
	count := CountPostsForThread(db, 1)
	if count < 1 {
		t.Error("Expected posts to be more than 0")
	}
}

func TestGetLatestThreads(t *testing.T) {
	var threads []Thread
	GetLatestThreads(db, MakeTimestamp(), 5, &threads)
	if len(threads) < 1 {
		t.Error("Expected more than 1 thread")
	}
	if len(threads) > 5 {
		t.Error("Expected less than 5 threads")
	}
}

func TestGetPostsForThread(t *testing.T) {
	var posts []Post
	GetPostsForThread(db, MakeTimestamp(), 10, 1, &posts)
	if len(posts) < 1 {
		t.Error("Expected more than 1 thread")
	}
	if len(posts) > 5 {
		t.Error("Expected less than 5 posts")
	}
}

func TestBlockUser(t *testing.T) {
	user, _ := FindUser(db, 1)
	BlockUser(db, user, 2)
	var blockedIDs []int
	GetBlockedIds(db, user, &blockedIDs)
	if !helpers.IntInSlice(blockedIDs, 2) {
		t.Error("Expected 2 to be in blocked user IDs")
	}
}

func TestUnblockUser(t *testing.T) {
	user, _ := FindUser(db, 1)
	UnblockUser(db, user, 2)
	var blockedIDs []int
	GetBlockedIds(db, user, &blockedIDs)
	if helpers.IntInSlice(blockedIDs, 2) {
		t.Error("Expected 2 to NOT be in blocked user IDs")
	}
}

func TestFindUserThread(t *testing.T) {
	user, _ := FindUser(db, 1)
	_, err := FindUserThread(db, user, 1)
	if err != nil {
		t.Error("Expected thread to be found: ", err)
	}
}

func TestFindUserPost(t *testing.T) {
	user, _ := FindUser(db, 1)
	_, err := FindUserPost(db, user, 1000000)
	if err == nil {
		t.Error("Expected post to not be found")
	}
}

func TestFindUserPost2(t *testing.T) {
	user, _ := FindUser(db, 1)
	_, err := FindUserPost(db, user, 1)
	if err != nil {
		t.Error("Expected post to be found: ", err)
	}
}

func TestFindUserThread2(t *testing.T) {
	user, _ := FindUser(db, 2)
	_, err := FindUserThread(db, user, 1)
	if err == nil {
		t.Error("Expected error finding thread")
	} else {
		fmt.Println("Public Error: ", err)
	}
}

func TestDeleteThread(t *testing.T) {
	user, _ := FindUser(db, 2)
	deleteErr := DeleteThread(db, user, 1)
	if deleteErr == nil {
		t.Error("Expected error deleting thread")
	} else {
		fmt.Println("Public Error: ", deleteErr)
	}
}

func TestDeleteThread2(t *testing.T) {
	user, _ := FindUser(db, 1)
	deleteErr := DeleteThread(db, user, 1)
	if deleteErr != nil {
		t.Error("Unexpected error deleting thread", deleteErr)
	}
	_, postErr := ReplyToThread(db, user, 1, "Some content to reply with to the original post")
	if postErr == nil {
		t.Error("Expected error from attempting to post to deleted thread")
	}
}

func TestDeletePost(t *testing.T) {
	user, _ := FindUser(db, 1)

	thread, threadErr := CreateThread(db, user, "New thread title and what not", "When it feels too hard to hold, When it feels I'm in the dark")

	if threadErr != nil {
		t.Error("Unexpected error creating thread", threadErr)
	}

	if post, postErr := ReplyToThread(db, user, thread.ID, "I'm not sure what else could possibly go here adsfsfdsdf"); postErr != nil {
		t.Error("Unexpected error adding reply: ", postErr)
	} else {
		if deleteErr := DeletePost(db, user, post.ID); deleteErr != nil {
			t.Error("Unexpected error deleting post", deleteErr)
		}
	}

}

func TestGetLatestThreadsForUser(t *testing.T) {
	user, _ := FindUser(db, 1)
	var threads []Thread
	GetLatestThreadsForUser(db, user, MakeTimestamp(), 10, &threads)
	if len(threads) < 1 {
		t.Error("Expected user to have more than 1 thread")
	}
}

func TestGetUsers(t *testing.T) {
	var users []User
	GetUsers(db, &users)
	if len(users) < 1 {
		t.Error("Expected more than 0 users")
	}
}