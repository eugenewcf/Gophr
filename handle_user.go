package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleUserNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	next := r.URL.Query().Get("next")
	// Display Home Page
	RenderTemplate(w, r, "users/new", map[string]interface{}{
		"Next": next,
	})
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Process creating a user
	user, errs := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)
	next := r.FormValue("next")

	// if err != nil {
	if len(errs) > 0 {
		var errMsgs []string
		for _, err := range errs {
			if IsValidationError(err) {
				errMsgs = append(errMsgs, err.Error())
			} else {
				panic(err)
				break
			}
		}
		// if IsValidationError(err) {
		if len(errMsgs) > 0 {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				// "Error": err.Error(),
				"Errors": errMsgs,
				"User":   user,
				"Next":   next,
			})
			return
		}
	}

	err := globalUserStore.Save(user)
	if err != nil {
		panic(err)
		return
	}

	// Create a new session
	session := NewSession(w)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, next+"?flash=User+created", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := RequestUser(r)

	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currentUser := RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")

	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+updated", http.StatusFound)
}

func HandleUserShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := globalUserStore.Find(params.ByName("userID"))
	if err != nil {
		panic(err)
	}

	// 404
	if user == nil {
		http.NotFound(w, r)
		return
	}

	images, err := globalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "users/show", map[string]interface{}{
		"Images": images,
		"User": user,
	})
}
