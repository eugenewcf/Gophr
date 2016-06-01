package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
)

func main() {
	router := NewRouter()

	router.Handle("GET", "/", HandleHome)
	router.Handle("GET", "/register", HandleUserNew)
	router.Handle("POST", "/register", HandleUserCreate)
	router.Handle("GET", "/login", HandleSessionNew)
	router.Handle("POST", "/login", HandleSessionCreate)
	router.Handle("GET", "/image/:imageID", HandleImageShow)
	router.Handle("GET", "/user/:userID", HandleUserShow)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	router.ServeFiles(
		"/im/*filepath",
		http.Dir("data/images"),
	)

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", HandleSessionDestory)
	secureRouter.Handle("GET", "/account", HandleUserEdit)
	secureRouter.Handle("POST", "/account", HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", HandleImageNew)
	secureRouter.Handle("POST", "/images/new", HandleImageCreate)

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

func init() {
	// Assign a user store
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store
	globalUserEmailMapping = map[string]*User{}
	globalUserUsernameMapping = map[string]*User{}
	for _, user := range store.Users {
		globalUserEmailMapping[strings.ToLower(user.Email)] = &user
		globalUserUsernameMapping[strings.ToLower(user.Username)] = &user
	}

	// Assign a session store
	sessionStore, err := NewFileSessionStore("./data/session.json")
	if err != nil {
		panic(fmt.Errorf("Error creating sesion storeL %s", err))
	}
	globalSessionStore = sessionStore

	// Assign a sql database
	db, err := NewMySQLDB("root:@tcp(127.0.0.1:3306)/gophr")
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db

	// Assgin an image store
	globalImageStore = NewDBImageStore()
}
