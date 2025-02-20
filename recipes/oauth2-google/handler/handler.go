package handler

import (
	"velocity-oauth-google/auth"

	"go.khulnasoft.com/velocity"
)

// Auth velocity handler
func Auth(c *velocity.Ctx) error {
	path := auth.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return c.Redirect(url)
}

// Callback to receive google's response
func Callback(c *velocity.Ctx) error {
	token, error := auth.ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
	if error != nil {
		panic(error)
	}
	email := auth.GetEmail(token.AccessToken)
	return c.Status(200).JSON(velocity.Map{"email": email, "login": true})
}
