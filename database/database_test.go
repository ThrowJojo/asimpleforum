package database

import (
	"testing"
	"github.com/jinzhu/gorm"
	"fmt"
)

var db *gorm.DB = MakeConnection()

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
	user, _ := FindUser(db, 1)
	CreateThread(db, user, "POST12342345234", "GDDFSG234234234234")
}

func TestFindUser2(t *testing.T) {
	_, err := FindUser(db, 200000)
	if err == nil {
		t.Error("Expected a not exist error")
	}
}

func TestFindUserByCredentials(t *testing.T) {
	_, err := FindUserByCredentials(db, TEST_USER1.Username, TEST_USER1.Password)
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestFindUserByCredentials2(t *testing.T) {
	_, err := FindUserByCredentials(db, TEST_USER1.Username, TEST_USER1.Password)
	if err != nil {
		t.Error("Expected user object to be " + TEST_USER1.Username)
	}
}

func TestReplyToThread(t *testing.T) {
	user, _ := FindUser(db, 1)
	ReplyToThread(db, user, 1, "The six bears went to town adfsdfsdfsdfsdfdfs")
}

func TestGetLatestThreads(t *testing.T) {
	var threads []Thread
	GetLatestThreads(db, MakeTimestamp(), 5, &threads)
	fmt.Printf("%+v\n", threads)
}

func TestGetPostsForThread(t *testing.T) {
	var posts []Post
	GetPostsForThread(db, MakeTimestamp(), 10, 2, &posts)
	fmt.Printf("%+v\n", posts)
}

func TestNotIn(t *testing.T) {
	var threads []Thread
	db.Limit(10).Where("id not in (?)", []uint{1, 2, 3}).Find(&threads)
	fmt.Printf("%+v\n", threads)
}

func TestBlockUser(t *testing.T) {
	user, _ := FindUser(db, 1)
	BlockUser(db, user, 2)
}

func TestUnblockUser(t *testing.T) {
	user, _ := FindUser(db, 1)
	UnblockUser(db, user, 2)
}

func TestGetBlockedIds(t *testing.T) {
	user, _ := FindUser(db, 1)
	var blockedIDs []int
	GetBlockedIds(db, user, &blockedIDs)
	fmt.Println(blockedIDs)
}

func TestFindUserThread(t *testing.T) {
	user, _ := FindUser(db, 1)
	thread, err := FindUserThread(db, user, 1)
	if err != nil {
		t.Error("Expected thread to be found: ", err)
	} else {
		fmt.Println(thread)
	}
}

func TestFindUserPost(t *testing.T) {
	user, _ := FindUser(db, 1)
	_, err := FindUserPost(db, user, 1000)
	if err == nil {
		t.Error("Expected post to not be found")
	}
}

func TestFindUserPost2(t *testing.T) {
	user, _ := FindUser(db, 1)
	post, err := FindUserPost(db, user, 1)
	if err != nil {
		t.Error("Expected post to be found: ", err)
	} else {
		fmt.Println(post)
	}
}

func TestFindUserThread2(t *testing.T) {
	user, _ := FindUser(db, 2)
	_, err := FindUserThread(db, user, 1)
	if err == nil {
		t.Error("Expected error finding thread")
	} else {
		fmt.Println("Error: ", err)
	}
}

func TestDeleteThread(t *testing.T) {
	user, _ := FindUser(db, 2)
	deleteErr := DeleteThread(db, user, 1)
	if deleteErr == nil {
		t.Error("Expected error deleting thread")
	} else {
		fmt.Println("ERROR: ", deleteErr)
	}
}

func TestDeleteThread2(t *testing.T) {
	user, _ := FindUser(db, 1)
	deleteErr := DeleteThread(db, user, 1)
	if deleteErr != nil {
		t.Error("Unexpected error deleting thread", deleteErr)
	}
}

func TestDeletePost(t *testing.T) {
	user, _ := FindUser(db, 1)
	post, postErr := ReplyToThread(db, user, 2, "I'm not sure what else could possibly go here adsfsfdsdf")
	if postErr != nil {
		t.Error("Unexpected error adding reply: ", postErr)
	} else {
		if deleteErr := DeletePost(db, user, post.ID); deleteErr != nil {
			t.Error("Unexpected error deleting post", deleteErr)
		}
	}
}

func TestJoinQuery(t *testing.T) {
	var threads []Thread
	//err := db.Joins("INNER JOIN user_threads ON user_threads.thread_id = threads.id").Where("user_id NOT IN (?)", []int{1, 3}).Find(&threads)
	db.Joins("INNER JOIN user_threads ON user_threads.thread_id = threads.id").Order("last_update desc").Preload("Authors").Preload("Posts").Preload("Posts.Authors").
			Limit(10).Where("timestamp < ? AND user_id NOT IN (?)", MakeTimestamp(), []int{1, 3}).Find(&threads)
	fmt.Println(len(threads))
}

func TestGetLatestThreadsForUser(t *testing.T) {
	user, _ := FindUser(db, 1)
	var threads []Thread
	GetLatestThreadsForUser(db, user, MakeTimestamp(), 10, &threads)
	fmt.Println(len(threads))
}

func TestGetUsers(t *testing.T) {
	users := GetUsers(db)
	fmt.Printf("%+v\n", *users)
}