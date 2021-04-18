package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/raisultan/abac/pkg/delete"
	"github.com/raisultan/abac/pkg/http/rest"
	"github.com/raisultan/abac/pkg/jwt_refresh"
	"github.com/raisultan/abac/pkg/list"
	"github.com/raisultan/abac/pkg/login"
	"github.com/raisultan/abac/pkg/register"
	"github.com/raisultan/abac/pkg/retrieve"
	"github.com/raisultan/abac/pkg/storage/postgres"
	"github.com/raisultan/abac/pkg/update"
)

const defaultPort = ":8080"

func main() {
	var registerer register.Service
	var loginer login.Service
	var jwtRefresher jwt_refresh.Service
	var lister list.Service
	var retriever retrieve.Service
	var updater update.Service
	var deleter delete.Service

	s, _ := postgres.NewStorage()

	registerer = register.NewService(s)
	loginer = login.NewService(s)
	jwtRefresher = jwt_refresh.NewService(s)
	lister = list.NewService(s)
	retriever = retrieve.NewService(s)
	updater = update.NewService(s)
	deleter = delete.Service(s)

	router := rest.Handler(
		registerer,
		loginer,
		jwtRefresher,
		lister,
		retriever,
		updater,
		deleter,
	)

	fmt.Println("Running on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(defaultPort, router))
}
