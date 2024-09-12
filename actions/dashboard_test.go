package actions

import (
	"net/http"
)

func (as *ActionSuite) Test_Dashboard() {
	res := as.JSON("/api/v1/dashboard").Get()

	as.Equal(http.StatusOK, res.Code)
	as.Containsf(res.Body.String(), "Dashboard data fetched successfully", "Response body does not match")
	as.Contains(res.Body.String(), "Invoice List")
	as.Contains(res.Body.String(), "Invoice Group By Status")
}
