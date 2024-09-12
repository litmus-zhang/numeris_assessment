package actions

import "assessment/models"

func (as *ActionSuite) Test_PaymentDetailsResource_List() {
	res := as.JSON(Base_url + "/payment_details").Get()

	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Payment details retrieved successfully")
}

func (as *ActionSuite) Test_PaymentDetailsResource_Show() {
	res := as.JSON(Base_url + "/payment_details/1").Get()

	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_PaymentDetailsResource_Create() {
	u := GetSampleUser()
	Body := &models.PaymentDetail{
		BankName:      "Test Bank",
		AccountNumber: "1234567890",
		AccountName:   "Test Account",
		CreatedBy:     *u,
	}
	res := as.JSON(Base_url + "/payment_details").Post(Body)

	as.Equal(201, res.Code)
	as.Contains(res.Body.String(), "Payment details created successfully")
}

// func (as *ActionSuite) Test_PaymentDetailsResource_Update() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_PaymentDetailsResource_Destroy() {
// 	as.Fail("Not Implemented!")
// }
