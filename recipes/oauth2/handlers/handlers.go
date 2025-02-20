package handlers

import (
	"oauth2/models"

	"go.khulnasoft.com/velocity"
)

// HTMLPages will render and return "public" pages
func HTMLPages(c *velocity.Ctx) error {
	models.SYSLOG.Tracef("entering HtmlPages; original URL: %v", c.OriginalURL())
	defer models.SYSLOG.Trace("exiting HtmlPages")

	// models.SYSLOG.Tracef("inspecting header: %v", &c.Request().Header)

	p := c.Path()
	switch p {
	case "/index.html":
		return c.Render("index", velocity.Map{
			"ClientID": models.ClientID,
		})

	case "/welcome.html":
		sessData, err := models.MySessionStore.Get(c)
		if err != nil {
			return c.Redirect("/errpage.html", velocity.StatusInternalServerError)
		}

		return c.Render("welcome", velocity.Map{
			"tokenValue": sessData.Get("oauth-token"),
		})

	case "/errpage.html":
		return c.Render("errpage", velocity.Map{})
	}

	return c.Next()
}
