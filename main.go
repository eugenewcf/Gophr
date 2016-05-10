package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
)

func main() {
	router := NewRouter()

	router.Handle("GET", "/", HandleHome)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	middleware := Middleware{}
	middleware.Add(router)
	fmt.Println("Listening on port :3000...")
	log.Fatal(http.ListenAndServe(":3000", middleware))
}

type NotFound struct {}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request){
}

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	notFound := new(NotFound)
	router.NotFound = notFound
	return router
}
