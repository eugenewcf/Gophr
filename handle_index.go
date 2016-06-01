package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	images, err := globalImageStore.FindAll(0)
	if err != nil {
		panic(err)
	}
	// Display Home Page
	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}

func HandleRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Display Register Page
	RenderTemplate(w, r, "index/register", nil)
}
