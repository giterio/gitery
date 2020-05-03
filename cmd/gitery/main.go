package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"gitery/config"
	"gitery/internal/controllers"
	"gitery/internal/middlewares"
	"gitery/internal/models"
)

func wrapMiddlewares(h http.Handler) http.Handler {
	h = middlewares.Authentication(h)
	h = middlewares.Constraint(h)
	return middlewares.LoadContext(h)
}

func main() {
	// init project configuration
	appConfig, err := config.Init(config.Development)
	if err != nil {
		panic(err)
	}
	// connect to the Db
	dbConfig := appConfig.Database
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// init root handler with sub-handlers
	rootHandler := &controllers.RootHandler{
		AuthHandler:    &controllers.AuthHandler{Model: &models.AuthService{DB: db, JwtSecret: appConfig.JwtSecret}},
		UserHandler:    &controllers.UserHandler{Model: &models.UserService{DB: db}},
		PostHandler:    &controllers.PostHandler{Model: &models.PostService{DB: db}},
		CommentHandler: &controllers.CommentHandler{Model: &models.CommentService{DB: db}},
	}
	// config the server
	server := http.Server{
		Addr:    appConfig.HTTP.Host + ":" + appConfig.HTTP.Port,
		Handler: wrapMiddlewares(rootHandler),
	}
	// start the server
	server.ListenAndServe()
}
