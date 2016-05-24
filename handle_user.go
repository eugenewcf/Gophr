package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleUserNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Display Home Page
	RenderTemplate(w, r, "users/new", nil)
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Process creating a user
	user, errs := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)

	// if err != nil {
	if len(errs) > 0 {
		var errMsgs []string
		for _, err := range errs {
			if IsValidationError(err){
				errMsgs = append(errMsgs, err.Error())
			}else{
				panic(err)
				break
			}
		}
		// if IsValidationError(err) {
		if len(errMsgs) > 0 {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				// "Error": err.Error(),
				"Errors": errMsgs,
				"User":  user,
			})
			return
		}
	}

	err := globalUserStore.Save(user)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
}
