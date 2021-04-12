package main

import (
	"net/http"
	"time"
)

func main() {
	hs, logger := setup()

	go func() {
		logger.Println("Up and running")

		if err := hs.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	graceful(hs, logger, 5*time.Second)
}
