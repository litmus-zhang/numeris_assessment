package actions

import (
	"assessment/utils"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HealthHandler(c buffalo.Context) error {
	response := utils.ResponseUtils[error]{
		Status:  http.StatusOK,
		Message: "All Systems operational",
	}
	return c.Render(http.StatusOK, r.JSON(response))
}
