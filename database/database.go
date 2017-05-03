package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/pkg/errors"
	"ForumDatabase/helpers"
	"ForumDatabase/config"
)

var (
	ErrNotExist = errors.New("Record does not exist")
	ErrBadRecord = errors.New("Bad record")
	ErrSystem = errors.New("System error")
	ErrExists = errors.New("Record already exists")
)

var (
	TEST_USER1 = User{
		Username: "Sholomobo2",
		Password: "backthemall97",
	}
	TEST_USER2 = User{
		Username: "goldtime34",
		Password: "breakingdown434",
	}
)

type BaseModel struct {
	ID uint `gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
}

type User struct {
	BaseModel
	Username string `json:"username"`
	Password string `json:"-"`
	Threads []Thread `json:"-" gorm:"many2many:user_threads;"`
	Posts []Post `json:"-" gorm:"many2many:user_posts;"`
	BlockRecords []BlockRecord `json:"-"`
}

type BlockRecord struct {
	BaseModel
	Target User
	TargetID int
	UserID uint
}

type Thread struct {
	BaseModel
	Title string `json:"title"`
	Content string `json:"content"`
	Timestamp int64 `json:"timestamp"`
	LastUpdate int64 `json:"lastUpdate"`
	Deleted bool `json:"-"`
	Authors []User `json:"authors" gorm:"many2many:user_threads;"`
	Posts []Post `json:"posts" gorm:"many2many:thread_posts"`
}

type Post struct {
	BaseModel
	Threads []Thread `json:"threads" gorm:"many2many:thread_posts"`
	Authors []User `json:"authors" gorm:"many2many:user_posts;"`
	Content string `json:"content"`
	Deleted bool `json:"-"`
	Timestamp int64 `json:"timestamp"`
}

// Gets a connection to the database
func MakeConnection() *gorm.DB {

	connectionString, stringErr := config.GetConnectionString()
	if stringErr != nil {
		panic(stringErr)
	}

	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	// TODO: Need to check if these settings are appropriate
	db.DB().SetConnMaxLifetime(time.Hour * 10)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(20)

	return db
}

// Does the auto-migrations, sets up the unique constraint indexes
func Setup(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Thread{}, &Post{}, &BlockRecord{})
	db.Model(&BlockRecord{}).AddUniqueIndex("BlockRecordIndex", "target_id", "user_id")
	db.Table("thread_posts").AddUniqueIndex("ThreadPostsIndex", "thread_id", "post_id")
	db.Table("user_threads").AddUniqueIndex("UserThreadsIndex", "user_id", "thread_id")
	db.Table("user_posts").AddUniqueIndex("UserPostsIndex", "user_id", "post_id")
}

// Creates a new user from the username and password(which gets encrypted)
func CreateUser(db *gorm.DB, username string, password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return ErrSystem }

	var existingUser User
	db.First(&existingUser, "username = ?", username)

	if existingUser.ID > 0 {
		return ErrExists
	}

	if usernameError := helpers.ValidateUsername(username); usernameError != nil {
		return usernameError
	}

	if passwordError := helpers.ValidatePassword(password); passwordError != nil {
		return passwordError
	}

	newUser := User{Username: username, Password: string(hash)}
	db.Save(&newUser)
	return nil

}

// Finds a user by id
func FindUser(db *gorm.DB, id uint) (*User, error) {
	var user User
	db.First(&user, id)
	if user.ID > 0 {
		return &user, nil
	} else {
		return nil, ErrNotExist
	}
}

// Finds a user with the given credentials, returns error if user can't be found/credentials are incorrect
func FindUserByCredentials(db *gorm.DB, username string, password string) (*User, error) {
	var user User
	db.Where("username = ?", username).First(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrNotExist
	}
	return &user, nil
}

// Finds a thread by id if it hasn't been deleted
func FindThread(db *gorm.DB, id uint) (*Thread, error) {
	var thread Thread
	db.Where("deleted = ?", false).First(&thread, id)
	if thread.ID > 0 {
		return &thread, nil
	} else {
		return nil, ErrNotExist
	}
}

// Returns a thread by id but only if the user is the author of it
func FindUserThread(db *gorm.DB, user *User, id uint) (*Thread, error) {
	var thread Thread
	db.Joins("INNER JOIN user_threads ON user_threads.thread_id = threads.id").Where("id = ? AND user_id = ?", id, user.ID).First(&thread)
	if thread.ID > 0 {
		return &thread, nil
	} else {
		return nil, ErrNotExist
	}
}

// Returns a post by id but only if the user is the author of it
func FindUserPost(db *gorm.DB, user *User, id uint) (*Post, error) {
	var post Post
	db.Joins("INNER JOIN user_posts ON user_posts.post_id = posts.id").Where("id = ? AND user_id = ?", id, user.ID).First(&post)
	if post.ID > 0 {
		return &post, nil
	} else {
		return nil, ErrNotExist
	}
}

// Creates a thread for the specified user
func CreateThread(db *gorm.DB, user *User, title string, content string) (*Thread, error) {

	if titleError := helpers.ValidateTitle(title); titleError != nil {
		return nil, titleError
	}

	if contentError := helpers.ValidateContent(content); contentError != nil {
		return nil, contentError
	}

	timestamp := MakeTimestamp()
	thread := Thread{Title: title, Content: content, Timestamp: timestamp, LastUpdate: timestamp}
	db.Model(&user).Association("Threads").Append(&thread)
	return &thread, nil

}

// Finds the block records for a user and writes the blocked user IDs to userIDs pointer
func GetBlockedIds(db *gorm.DB, user *User, userIDs *[]int) {
	rows, err := db.Table("block_records").Select("target_id").Where("user_id = ?", user.ID).Rows()
	if err == nil {
		for rows.Next() {
			var id int
			rows.Scan(&id)
			*userIDs = append(*userIDs, id)
		}
	}
}

// Blocks a user with the targetID
func BlockUser(db *gorm.DB, user *User, targetID uint) error {
	var existingRecord BlockRecord
	target, err := FindUser(db, targetID)
	db.Where("target_id = ? AND user_id = ?", targetID, user.ID).First(&existingRecord)

	if targetID == user.ID {
		return ErrBadRecord
	}

	if err == nil && existingRecord.ID < 1 {
		record := BlockRecord{Target: *target}
		db.Model(&user).Association("BlockRecords").Append(&record)
		return nil
	} else {
		return err
	}

}

// Unblocks a user with the targetID
func UnblockUser(db *gorm.DB, user *User, targetID uint) error {
	var record BlockRecord
	db.Where("target_id = ? AND user_id = ?", targetID, user.ID).First(&record)
	if record.ID > 0 {
		db.Unscoped().Delete(&record)
		return nil
	} else {
		return ErrNotExist
	}
}

// Replies with content using the threadId
func ReplyToThread(db *gorm.DB, user *User, threadId uint, content string) (*Post, error) {

	if contentErr := helpers.ValidateContent(content); contentErr != nil {
		return nil, contentErr
	}

	if thread, err := FindThread(db, threadId); err != nil {
		return nil, err
	} else {
		timestamp := MakeTimestamp()
		post := Post{Content: content, Timestamp: timestamp}
		db.Model(&thread).Association("Posts").Append(&post)
		db.Model(&user).Association("Posts").Append(&post)
		thread.LastUpdate = timestamp
		db.Save(&thread)
		return &post, nil
	}

}

// Marks the thread with supplied id as deleted if it can be found/the user has permission
func DeleteThread(db *gorm.DB, user *User, threadId uint) error {
	if thread, err := FindUserThread(db, user, threadId); err != nil {
		return err
	} else {
		thread.Deleted = true
		db.Save(&thread)
		return nil
	}
}

// Marks the post with the supplied id as a deleted if it can be found/the user has permission
func DeletePost(db *gorm.DB, user *User, postId uint) error {
	if post, err := FindUserPost(db, user, postId); err != nil {
		return err
	} else {
		post.Deleted = true
		db.Save(&post)
		return nil
	}
}

// Gets latest threads if the user isn't authenticated/user has no block records
func GetLatestThreads(db *gorm.DB, timestamp int64, threads *[]Thread) {
	db.Preload("Authors").Preload("Posts").Preload("Posts.Authors").Order("last_update desc").Limit(10).Where("timestamp < ? AND deleted = ?", timestamp, false).Find(&threads)
}

// Gets latest threads
// TODO: Needs a limit parameter
func GetLatestThreadsForUser(db *gorm.DB, user *User, timestamp int64, threads *[]Thread) {
	var blockedIDs []int
	GetBlockedIds(db, user, &blockedIDs)
	if len(blockedIDs) > 0 {
		db.Joins("INNER JOIN user_threads ON user_threads.thread_id = threads.id").Order("last_update desc").Preload("Authors").Preload("Posts").Preload("Posts.Authors").
				Limit(10).Where("timestamp < ? AND user_id NOT IN (?) AND deleted = ?", timestamp, blockedIDs, false).Find(&threads)
	} else {
		GetLatestThreads(db, timestamp, threads)
	}
}

// TODO: Remove or change to the same format as above and below
func GetLatestPosts(db *gorm.DB, timestamp int64) *[]Post {
	var posts []Post
	db.Preload("Authors").Order("timestamp desc").Limit(10).Where("timestamp < ? AND deleted = ?", timestamp, false).Find(&posts)
	return &posts
}

// Gets posts for the thread with the supplied id
// TODO: Needs a join query to find the posts in a thread, this isn't actually using threadId in it's current state
// TODO: Needs limit parameter
func GetPostsForThread(db *gorm.DB, timestamp int64, threadId uint, posts *[]Post) {
	db.Preload("Authors").Order("timestamp").Limit(10).Where("timestamp < ? AND deleted = ?", timestamp, false).Find(&posts)
}

// Creates a epoch millisecond timestamp
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Gets all the current users from the database
func GetUsers(db *gorm.DB) *[]User {
	var users []User
	db.Preload("Threads").Preload("Posts").Find(&users)
	return &users
}