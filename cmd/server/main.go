package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	var wait time.Duration
	flag.DurationVar(
		&wait,
		"gracefulShutDown",
		time.Second*15,
		"duration during which server will try to gracefully shutdown",
	)
	flag.Parse()

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

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0%v", defaultPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Println("Running on: ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
