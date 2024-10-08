package actions

import (
	"assessment/models"
	"assessment/utils"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Invoice)
// DB Table: Plural (invoices)
// Resource: Plural (Invoices)
// Path: Plural (/invoices)
// View Template Folder: Plural (/templates/invoices/)

// InvoicesResource is the resource for the Invoice model
type InvoicesResource struct {
	buffalo.Resource
}

// List gets all Invoices. This function is mapped to the path
// GET /invoices
func (v InvoicesResource) List(c buffalo.Context) error {
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

	invoices := &models.Invoices{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Invoices from the DB
	if err := q.Where("created_by = ?", &user.ID).GroupBy("created_at", "created_by", "customer_id", "due_date", "id").Order(`created_at desc`).All(invoices); err != nil {
		return err
	}

	response := utils.ResponseUtils[models.Invoices]{
		Status:  http.StatusOK,
		Message: "Invoices retrieved successfully",
		Data:    invoices,
	}
	return c.Render(200, r.JSON(response))

}

// Show gets the data for one Invoice. This function is mapped to
// the path GET /invoices/{invoice_id}
func (v InvoicesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Invoice
	invoice := &models.Invoice{}

	// To find the Invoice the parameter invoice_id is used.
	if err := tx.Find(invoice, c.Param("invoice_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	invoice.CustomerDetails = invoice.GetCustomerDetails(tx)
	invoice.PaymentDetail = invoice.GetPaymentDetails(tx)
	return c.Render(200, r.JSON(invoice))

}

// Create adds a Invoice to the DB. This function is mapped to the
// path POST /invoices
func (v InvoicesResource) Create(c buffalo.Context) error {
	// Allocate an empty Invoice
	user, err := GetUser(c)

	if err != nil {
		response := utils.ResponseUtils[string]{
			Status:  http.StatusUnauthorized,
			Message: "User not logged in",
		}
		return c.Render(http.StatusUnauthorized, r.JSON(response))
	}

	invoice := &models.Invoice{}

	// // Set the user that created the invoice
	invoice.CreatedBy = *user

	// Bind invoice to the html form elements
	if err := c.Bind(invoice); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	invoice.InvoiceNumber = invoice.GenerateInvoiceNumber()
	invoice.DueDate = invoice.SetDueDate(0, 0, 7)
	_, err = tx.ValidateAndCreate(invoice)
	if err != nil {
		response := utils.ResponseUtils[string]{
			Status:  http.StatusBadRequest,
			Message: "Error creating invoice",
		}
		return c.Render(http.StatusBadRequest, r.JSON(response))
	}
	invoice.GetTotal(tx)

	// Get the DB connection from the context
	response := utils.ResponseUtils[models.Invoice]{
		Status:  http.StatusCreated,
		Message: "Invoice created successfully",
		Data:    invoice,
	}

	return c.Render(201, r.JSON(response))

}

// Update changes a Invoice in the DB. This function is mapped to
// the path PUT /invoices/{invoice_id}
func (v InvoicesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Invoice
	invoice := &models.Invoice{}

	if err := tx.Find(invoice, c.Param("invoice_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Invoice to the html form elements
	if err := c.Bind(invoice); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(invoice)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("invoice", invoice)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("invoices/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "invoice.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/invoices/%v", invoice.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(invoice))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(invoice))
	}).Respond(c)
}

// Destroy deletes a Invoice from the DB. This function is mapped
// to the path DELETE /invoices/{invoice_id}
func (v InvoicesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Invoice
	invoice := &models.Invoice{}

	// To find the Invoice the parameter invoice_id is used.
	if err := tx.Find(invoice, c.Param("invoice_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(invoice); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "invoice.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/invoices")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(invoice))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(invoice))
	}).Respond(c)
}
