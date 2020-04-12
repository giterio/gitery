package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	ctl "gitery/internal/controllers"
	"gitery/internal/middlewares"
	"gitery/internal/models"
)

const (
	host     = "157.245.188.163"
	port     = 5432
	user     = "gitery"
	password = "5S8bjCl@30Nq"
	dbname   = "gitery"
)

func wrapMiddlewares(h http.Handler) http.Handler {
	return middlewares.WrapContext(h)
}

func main() {
	// connect to the Db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	router := ctl.Router{
		PostHandler: &ctl.PostHandler{Model: &models.Post{DB: db}},
	}

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: wrapMiddlewares(&router),
	}
	server.ListenAndServe()
}
