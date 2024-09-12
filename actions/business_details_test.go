package actions

import "assessment/models"

const Base_url = "/api/v1"

func (as *ActionSuite) Test_BusinessDetailsResource_List() {
	res := as.JSON(Base_url + "/business_details").Get()

	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Business details retrieved successfully")

}

func (as *ActionSuite) Test_BusinessDetailsResource_Show() {
	res := as.JSON(Base_url + "/business_details/1").Get()

	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_BusinessDetailsResource_Create() {
	Body := &models.BusinessDetail{
		BusinessName: "Test Business",
		Email:        "Test@tets.com",
		PhoneNumber:  "1234567890",
		Address:      "Test Address",
		CreatedByID:  1,
	}
	res := as.JSON(Base_url + "/business_details").Post(Body)

	as.Equal(201, res.Code)
	as.Contains(res.Body.String(), "Business details created successfully")
}

// func (as *ActionSuite) Test_BusinessDetailsResource_Update() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_BusinessDetailsResource_Destroy() {
// 	as.Fail("Not Implemented!")
// }
