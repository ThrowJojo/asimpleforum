package router

import (
	"github.com/gin-gonic/gin"
	"ForumDatabase/database"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"github.com/gin-contrib/sessions"
	"fmt"
	"ForumDatabase/config"
	"ForumDatabase/errors"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type QueryRequest struct {
	Timestamp int64 `form:"timestamp" binding:"required"`
	Limit int `form:"limit" binding:"required"`
}

var db *gorm.DB = database.MakeConnection()

func readLatestThreads(context *gin.Context) {

	data := new (QueryRequest)
	if bindErr := context.Bind(data); bindErr != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userValue, exists := context.Get("user")
	var threads []database.Thread

	if exists {
		user := userValue.(*database.User)
		database.GetLatestThreadsForUser(db, user, data.Timestamp, data.Limit, &threads)
	} else {
		database.GetLatestThreads(db, data.Timestamp, data.Limit, &threads)
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": threads,
	})

}

func readLatestPosts(context *gin.Context) {

	threadId, threadIdErr := strconv.ParseUint(context.Param("id"), 10, 64)
	data := new (QueryRequest)

	if bindErr := context.Bind(data); threadIdErr != nil || bindErr != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var posts []database.Post
	database.GetPostsForThread(db, data.Timestamp, data.Limit, uint(threadId), &posts)

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": posts,
	})
}

func login(context *gin.Context) {

	session := sessions.Default(context)
	data := new(AuthRequest)
	err := context.BindJSON(data)

	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, userErr := database.FindUserByCredentials(db, data.Username, data.Password)
	if userErr != nil {
		context.JSON(http.StatusOK, gin.H {
			"status": http.StatusOK,
			"error": userErr.Error(),
		})
	} else {
		session.Set("user_id", user.UniqueID)
		session.Save()
		context.JSON(http.StatusOK, gin.H {
			"status": http.StatusOK,
		})
	}

}

func register(context *gin.Context) {

	data := new(AuthRequest)
	err := context.BindJSON(data)

	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	createErr := database.CreateUser(db, data.Username, data.Password)
	if createErr != nil {
		renderError(context, createErr)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}

func authMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		userID := session.Get("user_id")
		if userID == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := database.FindUserByUnique(db, userID.(string))
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Set("user", user)
		context.Next()
	}
}

func softAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		userID := session.Get("user_id")

		if userID != nil {
			user, err := database.FindUserByUnique(db, userID.(string))
			if err != nil {
				context.AbortWithStatus(http.StatusUnauthorized)
				return
			} else {
				context.Set("user", user)
			}
		}
		context.Next()
	}
}

// TODO: Not sure if this is even needed anymore
func errorMiddleware(context *gin.Context) {
	context.Next()
	fmt.Println("PRINT ERRORS")
	errorToPrint := context.Errors.ByType(gin.ErrorTypePublic).Last()
	fmt.Println("Errors ", errorToPrint)
	if errorToPrint != nil {
		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": errorToPrint.Error(),
		})
	}
}

func createThread(context *gin.Context) {

	data := new (database.Thread)
	err := context.BindJSON(data)

	// TODO: Maybe binding errors should just be blank
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	value := context.MustGet("user")
	user := value.(*database.User)

	thread, createErr := database.CreateThread(db, user, data.Title, data.Content)
	if createErr != nil {
		renderError(context, createErr)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": thread,
	})

}

func addPost(context *gin.Context) {

	data := new (database.Post)
	err := context.BindJSON(data)
	threadId, convertErr := strconv.ParseUint(context.Param("id"), 10, 64)

	if convertErr != nil || err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	value := context.MustGet("user")
	user := value.(*database.User)

	post, replyErr := database.ReplyToThread(db, user, uint(threadId), data.Content)
	if replyErr != nil {
		renderError(context, replyErr)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": post,
	})

}

func blockUser(context *gin.Context) {

	userId, err := strconv.ParseUint(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	value := context.MustGet("user")
	user := value.(*database.User)

	blockError := database.BlockUser(db, user, uint(userId))
	if blockError != nil {
		renderError(context, blockError)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}

func unblockUser(context *gin.Context) {

	userId, err := strconv.ParseUint(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	value := context.MustGet("user")
	user := value.(*database.User)

	unblockError := database.UnblockUser(db, user, uint(userId))
	if unblockError != nil {
		renderError(context, unblockError)
		return
	}

	context.JSON(http.StatusOK, gin.H {
		"status": http.StatusOK,
	})

}

func deleteThread(context *gin.Context) {

	threadId, err := strconv.ParseUint(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	value := context.MustGet("user")
	user := value.(*database.User)

	if deleteErr := database.DeleteThread(db, user, uint(threadId)); deleteErr != nil {
		renderError(context, deleteErr)
		return
	}

	context.JSON(http.StatusOK, gin.H {
		"status": http.StatusOK,
	})

}

func deletePost(context *gin.Context) {

	postId, err := strconv.ParseUint(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	value := context.MustGet("user")
	user := value.(*database.User)

	if deleteErr := database.DeletePost(db, user, uint(postId)); deleteErr != nil {
		renderError(context, deleteErr)
		return
	}

	context.JSON(http.StatusOK, gin.H {
		"status": http.StatusOK,
	})

}


func renderError(context *gin.Context, err *errors.UserError) {
	context.JSON(http.StatusBadRequest, gin.H{
		"status": http.StatusBadRequest,
		"message": err.Error(),
		"error": err.Code,
	})
}

func Create() *gin.Engine {

	configData, err := config.LoadConfigData()
	if err != nil {
		panic("Issue loading config file")
	}

	// TODO: Maybe change to a memcache or redis store
	store := sessions.NewCookieStore([]byte(configData.Secret))
	database.Setup(db)
	ginRouter := gin.Default()
	ginRouter.Use(sessions.Sessions("mysession", store))

	auth := ginRouter.Group("/auth/login")
	{
		auth.POST("/", login)
	}

	threads := ginRouter.Group("/api/v1/threads")
	{
		threads.GET("/latest", softAuthMiddleware(), readLatestThreads)
		threads.GET("/responses/:id", readLatestPosts)
		threads.POST("/new", authMiddleware(), createThread)
		threads.POST("/reply/:id", authMiddleware(), addPost)
		threads.POST("/delete/:id", authMiddleware(), deleteThread)
	}

	posts := ginRouter.Group("/api/v1/posts")
	{
		posts.POST("/delete/:id", authMiddleware(), deletePost)
	}

	users := ginRouter.Group("/api/v1/users")
	{
		users.POST("/block/:id", authMiddleware(), blockUser)
		users.POST("/unblock/:id", authMiddleware(), unblockUser)
		users.POST("/new", register)
	}

	return ginRouter

}
