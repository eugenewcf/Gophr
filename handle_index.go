package main

import (
  "github.com/julienschmidt/httprouter"
  "net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Display Home Page
	RenderTemplate(w, r, "index/home", nil)
}


func HandleRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  // Display Register Page
  RenderTemplate(w, r, "index/register", nil)
}
