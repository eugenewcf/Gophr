package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func main() {
	unathenticatedRouter := NewRouter()
	unathenticatedRouter.GET("/", HandleHome)
	unathenticatedRouter.GET("/register", HandleRegister)

	authenticatedRouter := NewRouter()
	authenticatedRouter.GET("/images/new", HandleImageNew)

	staticFilesRouter := NewRouter()
	staticFilesRouter.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	middleware := Middleware{}
	middleware.Add(unathenticatedRouter)
	middleware.Add(staticFilesRouter)
	middleware.Add(http.HandlerFunc(AuthenticateRequest))
	middleware.Add(authenticatedRouter)

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", middleware)
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
