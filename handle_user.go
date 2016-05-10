package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandlerUserNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Display Home Page
	RenderTemplate(w, r, "users/new", nil)
}
