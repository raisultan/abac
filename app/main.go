package main

import "os"

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	a.Run(defaultPort)
}
