package main

import "os"

func main() {
	a := App{}
	a.Initialize(os.Getenv("POSTGRES_URL"))
	a.Run(defaultPort)
}
