package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleSessionDestory(w http.ResponseWriter, r * http.Request, _ httprouter.Params) {
  session := RequestSession(r)
  if session != nil {
    err := globalSessionStore.Delete(session)
    if err != nil {
      panic(err)
    }
  }
  RenderTemplate(w, r, "session/destory", nil)
}

func HandleSessionNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  next := r.URL.Query().Get("next")
  RenderTemplate(w, r, "session/new", map[string]interface{}{
    "Next": next,
  })
}

func HandleSessionCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  username := r.FormValue("username")
  password := r.FormValue("password")
  next := r.FormValue("next")

  user, err := FindUser(username, password)
  if err != nil {
    if IsValidationError(err) {
      RenderTemplate(w, r, "session/new", map[string]interface{}{
        "Error": err,
        "User": user,
        "Next": next,
      })
      return
    }
    panic(err)
  }

  session := FindOrCreateSession(w, r)
  session.UserID = user.ID
  err = globalSessionStore.Save(session)
  if err != nil {
    panic(err)
  }

  if next == "" {
    next = "/"
  }

  http.Redirect(w, r, next+"?flash=Signed+in", http.StatusFound)
}
