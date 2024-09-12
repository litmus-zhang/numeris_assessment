package actions

import (
	"net/http"

	"assessment/models"
)

func (as *ActionSuite) createUser() (*models.User, error) {
	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny(), "validation error: %v", verrs)

	return u, err
}

func (as *ActionSuite) Test_Auth_Signin() {
	res := as.HTML("/auth/").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), `<a href="/auth/new/">Sign In</a>`)
}

func (as *ActionSuite) Test_Auth_New() {
	res := as.HTML("/auth/new").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "Sign In")
}

func (as *ActionSuite) Test_Auth_Create() {
	u, err := as.createUser()
	as.NoError(err)

	tcases := []struct {
		Email    string
		Password string
		Status   int

		Identifier string
	}{
		{u.Email, u.Password, http.StatusFound, "Valid"},
		{"noexist@example.com", "password", http.StatusUnauthorized, "Email Invalid"},
		{u.Email, "invalidPassword", http.StatusUnauthorized, "Password Invalid"},
	}

	for _, tcase := range tcases {
		as.Run(tcase.Identifier, func() {
			res := as.HTML("/auth").Post(&models.User{
				Email:    tcase.Email,
				Password: tcase.Password,
			})

			as.Equal(tcase.Status, res.Code)
		})
	}
}
