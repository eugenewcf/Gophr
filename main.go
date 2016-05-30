package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

	router.Handle("GET", "/", HandleHome)
	router.Handle("GET", "/register", HandleUserNew)
	router.Handle("POST", "/register", HandleUserCreate)
	router.Handle("GET", "/login", HandleSessionNew)
	router.Handle("POST", "/login", HandleSessionCreate)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", HandleSessionDestory)
	secureRouter.Handle("GET", "/account", HandleUserEdit)
	secureRouter.Handle("POST", "/account", HandleUserUpdate)

	middleware := Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(secureRouter)

	fmt.Println("Listening on port :3000...")
	log.Fatal(http.ListenAndServe(":3000", middleware))
}

type NotFound struct{}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	notFound := new(NotFound)
	router.NotFound = notFound
	return router
}
