package views

import "github.com/enlistedmango/lenslocked/models"

type User struct {
	Email string
	// Add other fields as necessary
}

// Convert models.User to views.User
func ConvertToViewUser(mUser *models.User) *User {
	if mUser == nil {
		return nil
	}
	return &User{
		Email: mUser.Email, // Map fields as necessary
		// Map other fields as necessary
	}
}

type NavigationData struct {
	CurrentRoute string
	User         *User
}

type Alert struct {
	Color   string
	Message string
}

type FormData struct {
	CSRFToken string
	Email     string
	Alerts    []Alert
}
