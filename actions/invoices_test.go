package actions

import (
	"assessment/models"
	"time"
)

func (as *ActionSuite) Test_InvoicesResource_List() {
	res := as.JSON(Base_url + "/invoices").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Invoices retrieved successfully")
}

func (as *ActionSuite) Test_InvoicesResource_Show() {
	res := as.JSON(Base_url + "/invoices/1").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_InvoicesResource_Create() {
	body := &models.Invoice{
		Items: []models.Item{
			{
				Name:      "Test Item",
				UnitPrice: 100,
				Quantity:  1,
			},
			{
				Name:      "Test Item 2",
				UnitPrice: 200,
				Quantity:  2,
			},
		},
		DueDate: time.Now().Add(time.Hour * 24 * 7),
		Customer: models.Customer{
			FirstName:   "Test",
			LastName:    "User",
			Email:       "test@tr.com",
			PhoneNumber: "1234567890",
		},
		CreatedBy: models.User{
			Email:                "user@test.com",
			Password:             "password",
			PasswordConfirmation: "password",
		},
	}
	res := as.JSON(Base_url + "/invoices").Post(body)
	as.Equal(201, res.Code)
	as.Contains(res.Body.String(), "Invoice created successfully")
}

// func (as *ActionSuite) Test_InvoicesResource_Update() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_InvoicesResource_Destroy() {
// 	as.Fail("Not Implemented!")
// }
