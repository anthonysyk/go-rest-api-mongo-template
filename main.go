package main

import (
	"context"
	"github.com/anthonysyk/go-rest-api-mongo-template/api/user"
	"github.com/anthonysyk/go-rest-api-mongo-template/config"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/auth"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/db"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/fs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	conf := config.New()
	dbClient, err := db.NewMongoDB(ctx, conf.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(ctx)

	fsClient, err := fs.NewOSClient(conf.FileStorageRootPath)
	if err != nil {
		log.Fatal(err)
	}

	// Database
	apiDB := dbClient.Database("api")

	// Repositories
	userRepository := user.NewRepository(apiDB)

	// Services
	userService := user.NewService(userRepository, fsClient)

	// Handlers
	userHandler := user.NewHandler(userService)

	r := gin.Default()
	healthAPI := r.Group("/")
	{
		healthAPI.GET("health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		})
	}

	authAPI := r.Group("/")
	{
		authAPI.GET("/login", userHandler.Login(conf.AuthSecret))
	}

	// requires user authentication to access
	userAPI := r.Group("/")
	{
		userAPI.POST("/add/users", userHandler.CreateUsers)
		userAPI.GET("/users/list", userHandler.ListUsers)
		userAPI.GET("/users/:id", auth.Middleware(conf.AuthSecret), userHandler.GetUser)
		userAPI.PUT("/users/:id", auth.Middleware(conf.AuthSecret), userHandler.UpdateUser)
		userAPI.DELETE("/delete/user/:id", auth.Middleware(conf.AuthSecret), userHandler.DeleteUser)
	}

	log.Fatal(http.ListenAndServe(":8080", r.Handler()))
}
