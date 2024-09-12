package actions

import (
	"net/http"

	"assessment/models"
)

func GetSampleUser() *models.User {
	return &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
}

func (as *ActionSuite) Test_Users_Create() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	u := GetSampleUser()

	res := as.JSON(Base_url + "/signin").Post(u)
	as.Equal(http.StatusFound, res.Code)

	count, err = as.DB.Count(&models.User{})
	as.NoError(err)
	as.Equal(1, count)
}
