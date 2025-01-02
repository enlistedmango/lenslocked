package controllers

import (
	"fmt"
	"net/http"

	"github.com/enlistedmango/lenslocked/middleware"
	"github.com/enlistedmango/lenslocked/models"
	"github.com/enlistedmango/lenslocked/views"
	"github.com/gorilla/sessions"
)

type Users struct {
	Templates struct {
		New    views.Template
		SignIn views.Template
	}
	UserService  *models.UserService
	SessionStore *sessions.CookieStore
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	data := views.TemplateData{
		Nav: views.NavigationData{
			CurrentRoute: "/signup",
		},
	}
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := views.TemplateData{
		Nav: views.NavigationData{
			CurrentRoute: "/signup",
		},
		Form: &views.FormData{
			Email: email,
		},
	}

	user, err := u.UserService.Create(email, password)
	if err != nil {
		data.Alert = &views.Alert{
			Color:   "error",
			Message: "Something went wrong...",
		}
		switch err.(type) {
		case models.ErrEmailTaken:
			data.Alert.Message = "Email address is already taken"
		default:
			if err.Error() == "password must be at least 8 characters long" {
				data.Alert.Message = err.Error()
			} else {
				fmt.Println("Error creating user:", err)
				data.Alert.Message = "Something went wrong. Please try again later."
			}
		}
		u.Templates.New.Execute(w, data)
		return
	}

	session, _ := u.SessionStore.Get(r, middleware.SessionName)
	session.Values[middleware.UserIDKey] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	session, _ := u.SessionStore.Get(r, middleware.SessionName)
	session.Values = make(map[interface{}]interface{})
	session.Save(r, w)
	http.Redirect(w, r, "/?signedout=true", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	data := views.TemplateData{
		Nav: views.NavigationData{
			CurrentRoute: "/signin",
		},
	}
	u.Templates.SignIn.Execute(w, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		data := views.TemplateData{
			Nav: views.NavigationData{
				CurrentRoute: "/signin",
			},
			Form: &views.FormData{
				Email: email,
			},
			Alert: &views.Alert{
				Color:   "error",
				Message: "Invalid email or password",
			},
		}
		u.Templates.SignIn.Execute(w, data)
		return
	}

	session, _ := u.SessionStore.Get(r, middleware.SessionName)
	session.Values[middleware.UserIDKey] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/?success=true", http.StatusFound)
}
