package actions

import (
	"assessment/models"
	"assessment/utils"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		response := utils.ResponseUtils[error]{
			Status:  http.StatusBadRequest,
			Message: verrs.Error(),
		}
		return c.Render(http.StatusBadRequest, r.JSON(response))
	}
	response := utils.ResponseUtils[error]{
		Status:  http.StatusCreated,
		Message: "User created successfully",
	}

	return c.Render(http.StatusCreated, r.JSON(response))
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {

	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			user, err := GetUser(c)
			if err != nil {
				return errors.WithStack(err)
			}

			c.Set("current_user", user)
			c.Data()["current_user"] = user
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		_, err := GetUser(c)
		if err != nil {
			response := utils.ResponseUtils[string]{
				Status:  http.StatusUnauthorized,
				Message: "User not logged in",
			}
			return c.Render(http.StatusUnauthorized, r.JSON(response))
		}
		return next(c)
	}
}
