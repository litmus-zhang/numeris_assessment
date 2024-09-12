package actions

import (
	"assessment/models"
	"log"
)

func (as *ActionSuite) Test_CustomersResource_List() {
	res := as.JSON(Base_url + "/customers").Get()
	log.Println(res.Body.String())
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Customers retrieved successfully")
}

func (as *ActionSuite) Test_CustomersResource_Show() {
	res := as.JSON(Base_url + "/customers/1").Get()
	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_CustomersResource_Create() {
	body := &models.Customer{
		FirstName:   "Test",
		LastName:    "User",
		Email:       "testuser@test.com",
		PhoneNumber: "1234567890",
		CreatedByID: 1,
	}
	res := as.JSON(Base_url + "/customers").Post(body)
	as.Equal(201, res.Code)
	as.Contains(res.Body.String(), "Customer created successfully")
}

// func (as *ActionSuite) Test_CustomersResource_Update() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_CustomersResource_Destroy() {
// 	as.Fail("Not Implemented!")
// }
