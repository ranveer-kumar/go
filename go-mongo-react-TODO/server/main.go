package main

import (
	"log"
	"net/http"
	"ranveer/go/go-mongo-react-TODO/server/consul"
	"ranveer/go/go-mongo-react-TODO/server/router"
)

func main() {
	log.Println("in main..")
	r := router.Router()
	log.Println("Starting server on the port 8080...")
	consul.RegisterServiceWithConsul()
	log.Fatal(http.ListenAndServe(":8080", r))
}
