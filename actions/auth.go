package actions

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"assessment/models"
	"assessment/utils"
)

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {

	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(u.Email))).First(u)

	// helper function to handle bad attempts
	bad := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password")
		response := utils.ResponseUtils[string]{
			Status:  http.StatusUnauthorized,
			Message: "Invalid email/password",
		}

		c.Set("errors", verrs)
		c.Set("user", u)

		return c.Render(http.StatusUnauthorized, r.JSON(response))
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// couldn't find an user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}
	token := u.CreateJwtToken()
	response := utils.ResponseUtils[models.Token]{
		Status:  http.StatusFound,
		Message: "User logged in successfully",
		Data:    &token,
	}

	return c.Render(http.StatusFound, r.JSON(response))
}

func AuthGetMe(c buffalo.Context) error {

	user, err := GetUser(c)
	if err != nil {
		return err
	}
	response := utils.ResponseUtils[models.User]{
		Status:  http.StatusOK,
		Message: "User found",
		Data:    user,
	}
	return c.Render(http.StatusOK, r.JSON(response))
}

func GetUser(c buffalo.Context) (*models.User, error) {
	u := &models.User{}
	jwtToken := c.Request().Header.Get("Authorization")
	slice := strings.Split(jwtToken, " ")

	user, _ := u.GetUserFromToken(slice[1])
	if user == nil {
		response := utils.ResponseUtils[string]{
			Status:  http.StatusUnauthorized,
			Message: "User not logged in",
		}
		return nil, c.Render(http.StatusUnauthorized, r.JSON(response))
	}
	return user, nil
}
