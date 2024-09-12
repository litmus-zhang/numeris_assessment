package actions

import (
	"net/http"
)

func (as *ActionSuite) Test_HealthHandler() {
	res := as.JSON("/api/v1/health").Get()

	as.Equal(http.StatusOK, res.Code)
	as.Containsf(res.Body.String(), "All Systems operational", "Response body does not match")
}
