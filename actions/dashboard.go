package actions

import (
	"assessment/models"
	"assessment/utils"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

func Dashboard(c buffalo.Context) error {
	// Get the DB connection from the context
	user, err := GetUser(c)

	if err != nil {
		response := utils.ResponseUtils[string]{
			Status:  http.StatusUnauthorized,
			Message: "User not logged in",
		}
		return c.Render(http.StatusUnauthorized, r.JSON(response))
	}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	invoicesList := &models.Invoices{}
	invoicesGroup := &models.Invoices{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Invoices from the DB
	if err := q.Where("created_by = ?", &user.ID).GroupBy("created_at", "created_by", "customer_id", "due_date", "id").Order(`created_at desc`).All(invoicesList); err != nil {
		return err
	}

	if err := tx.Where("created_by = ?", &user.ID).All(invoicesGroup); err != nil {
		return err
	}
	// sorting the invoice group by status, and returning the total amount for each status
	finalGroup := invoicesList.SortByStatus()
	type ResponseType struct {
		Invoices *models.Invoices
		Group    map[string]float64
	}

	response := utils.ResponseUtils[ResponseType]{
		Status:  http.StatusOK,
		Message: "Dashboard data retrieved successfully",
		Data: &ResponseType{
			Invoices: invoicesList,
			Group:    finalGroup,
		},
	}
	return c.Render(200, r.JSON(response))

}
