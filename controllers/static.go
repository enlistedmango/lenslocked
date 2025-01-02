package controllers

import (
	"net/http"

	"github.com/enlistedmango/lenslocked/models"
	"github.com/enlistedmango/lenslocked/views"
)

func StaticHandler(app *App, tpl views.Template, currentRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.UserFromContext(r.Context())
		data := views.TemplateData{
			Nav: views.NavigationData{
				CurrentRoute: currentRoute,
				User:         views.ConvertToViewUser(user),
			},
			Title: getTitle(currentRoute),
		}

		// If user is signed in and we're on the home page, fetch their galleries
		if user != nil && currentRoute == "/" {
			galleries, err := app.GalleryService.GetByUserID(user.ID)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			data.Galleries = galleries
		}

		tpl.Execute(w, data)
	}
}

// Helper function to get page titles
func getTitle(route string) string {
	switch route {
	case "/":
		return "Home"
	case "/contact":
		return "Contact Us"
	case "/faq":
		return "FAQ"
	default:
		return ""
	}
}
