package actions

import (
	"assessment/models"
	"assessment/utils"
	"fmt"
	"log"
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
// Model: Singular (Customer)
// DB Table: Plural (customers)
// Resource: Plural (Customers)
// Path: Plural (/customers)
// View Template Folder: Plural (/templates/customers/)

// CustomersResource is the resource for the Customer model
type CustomersResource struct {
	buffalo.Resource
}

// List gets all Customers. This function is mapped to the path
// GET /customers
func (v CustomersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	user, err := GetUser(c)
	log.Println(user)
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

	customers := &models.Customers{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Customers from the DB
	if err := q.Where("created_by = ?", &user.ID).All(customers); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("customers", customers)
		return c.Render(http.StatusOK, r.HTML("customers/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		response := utils.ResponseUtils[models.Customers]{
			Status:  http.StatusOK,
			Message: "Customers retrieved successfully",
			Data:    customers,
		}
		return c.Render(200, r.JSON(response))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(customers))
	}).Respond(c)
}

// Show gets the data for one Customer. This function is mapped to
// the path GET /customers/{customer_id}
func (v CustomersResource) Show(c buffalo.Context) error {

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Customer
	customer := &models.Customer{}

	// To find the Customer the parameter customer_id is used.
	if err := tx.Find(customer, c.Param("customer_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("customer", customer)

		return c.Render(http.StatusOK, r.HTML("customers/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(customer))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(customer))
	}).Respond(c)
}

// Create adds a Customer to the DB. This function is mapped to the
// path POST /customers
func (v CustomersResource) Create(c buffalo.Context) error {

	user, err := GetUser(c)
	log.Println(user)
	if err != nil {
		response := utils.ResponseUtils[string]{
			Status:  http.StatusUnauthorized,
			Message: "User not logged in",
		}
		return c.Render(http.StatusUnauthorized, r.JSON(response))
	}
	// Allocate an empty Customer
	customer := &models.Customer{}

	customer.CreatedBy = *user

	// Bind customer to the html form elements
	if err := c.Bind(customer); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(customer)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("customer", customer)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("customers/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "customer.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/customers/%v", customer.ID)
	}).Wants("json", func(c buffalo.Context) error {
		response := utils.ResponseUtils[string]{
			Status:  http.StatusCreated,
			Message: "Customer created successfully",
		}
		return c.Render(http.StatusCreated, r.JSON(response))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(customer))
	}).Respond(c)
}

// Update changes a Customer in the DB. This function is mapped to
// the path PUT /customers/{customer_id}
func (v CustomersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Customer
	customer := &models.Customer{}

	if err := tx.Find(customer, c.Param("customer_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Customer to the html form elements
	if err := c.Bind(customer); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(customer)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("customer", customer)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("customers/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "customer.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/customers/%v", customer.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(customer))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(customer))
	}).Respond(c)
}

// Destroy deletes a Customer from the DB. This function is mapped
// to the path DELETE /customers/{customer_id}
func (v CustomersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Customer
	customer := &models.Customer{}

	// To find the Customer the parameter customer_id is used.
	if err := tx.Find(customer, c.Param("customer_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(customer); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "customer.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/customers")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(customer))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(customer))
	}).Respond(c)
}
