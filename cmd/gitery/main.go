package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"gitery/configs"
	"gitery/internal/controllers"
	"gitery/internal/middlewares"
	"gitery/internal/models"
)

func wrapMiddlewares(h http.Handler) http.Handler {
	h = middlewares.Authentication(h) // need access to RootHandler, should be the first one
	h = middlewares.Constraint(h)
	return middlewares.LoadContext(h)
}

func main() {
	// get environment variable "APP_ENV"
	env := os.Getenv("APP_ENV")
	if env != "production" && env != "development" {
		env = "debug"
	}
	// init project configuration
	appConfig, err := configs.Init(configs.EnvType(env))
	if err != nil {
		log.Fatal(err)
	}
	// connect to the Db
	dbConfig := appConfig.Database
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	// init root handler with sub-handlers
	rootHandler := &controllers.RootHandler{
		AuthHandler: &controllers.AuthHandler{
			Model: &models.AuthService{DB: db, JwtSecret: appConfig.JwtSecret},
		},
		UserHandler: &controllers.UserHandler{
			Model: &models.UserService{DB: db},
			UserPostHandler: &controllers.UserPostHandler{
				Model: &models.UserPostService{DB: db},
			},
		},
		PostHandler: &controllers.PostHandler{
			Model: &models.PostService{DB: db},
			PostLikeHandler: &controllers.PostLikeHandler{
				Model: &models.PostLikeService{DB: db},
			},
		},
		CommentHandler: &controllers.CommentHandler{
			Model: &models.CommentService{DB: db},
			CommentVoteHandler: &controllers.CommentVoteHandler{
				Model: &models.CommentVoteService{DB: db},
			},
		},
		TagHandler: &controllers.TagHandler{
			Model: &models.TagService{DB: db},
		},
	}
	// config the server
	server := http.Server{
		Addr:    appConfig.HTTP.Host + ":" + appConfig.HTTP.Port,
		Handler: wrapMiddlewares(rootHandler),
	}
	// start the server
	server.ListenAndServe()
}
