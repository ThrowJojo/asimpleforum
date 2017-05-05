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
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *gorm.DB = database.MakeConnection()

func readLatestThreads(context *gin.Context) {

	timestampParam := context.DefaultQuery("timestamp", strconv.FormatInt(database.MakeTimestamp(), 10))
	timestamp, err := strconv.ParseInt(timestampParam, 10, 64)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userValue, exists := context.Get("user")
	var threads []database.Thread

	if exists {
		user := userValue.(*database.User)
		database.GetLatestThreadsForUser(db, user, timestamp, &threads)
	} else {
		database.GetLatestThreads(db, timestamp, &threads)
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": threads,
	})

}

func readLatestPosts(context *gin.Context) {

	timestampParam := context.DefaultQuery("timestamp", strconv.FormatInt(database.MakeTimestamp(), 10))
	timestamp, timestampErr := strconv.ParseInt(timestampParam, 10, 64)

	limitParam := context.DefaultQuery("limit", "10")
	limit, limitErr := strconv.Atoi(limitParam)

	threadId, threadIdErr := strconv.ParseUint(context.Param("id"), 10, 64)

	if timestampErr != nil || threadIdErr != nil || limitErr != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var posts []database.Post
	database.GetPostsForThread(db, timestamp, limit, uint(threadId), &posts)

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
		session.Set("user_id", user.ID)
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

		user, err := database.FindUser(db, userID.(uint))
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
			user, err := database.FindUser(db, userID.(uint))
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

	// TODO: Maybe there should be a message to go with binding errors
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	value := context.MustGet("user")
	user := value.(*database.User)

	_, createErr := database.CreateThread(db, user, data.Title, data.Content)
	if createErr != nil {
		renderError(context, createErr)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
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

	_, replyErr := database.ReplyToThread(db, user, uint(threadId), data.Content)
	if replyErr != nil {
		renderError(context, replyErr)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
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

func renderError(context *gin.Context, err error) {
	context.JSON(http.StatusBadRequest, gin.H{
		"status": http.StatusBadRequest,
		"message": err.Error(),
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
	}

	users := ginRouter.Group("/api/v1/users")
	{
		users.POST("/block/:id", authMiddleware(), blockUser)
		users.POST("/unblock/:id", authMiddleware(), unblockUser)
		users.POST("/new", register)
	}

	return ginRouter

}
