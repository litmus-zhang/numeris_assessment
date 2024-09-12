package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

type Status string

// Define the enum values
const (
	Pending   Status = "Pending"
	Paid      Status = "Paid"
	Cancelled Status = "Cancelled"
	Draft     Status = "Draft"
)

// Invoice is used by pop to map your invoices database table to your go code.
type Invoice struct {
	ID              int           `json:"id" db:"id"`
	Items           Items         `has_many:"items" json:"items,omitempty" db:"-"`
	Total           float64       `json:"total" db:"total_amount"`
	InvoiceNumber   string        `json:"invoice_number" db:"invoice_number"`
	CreatedBy       User          `belongs_to:"user" json:"-" db:"-"`
	CreatedByID     int           `json:"created_by" db:"created_by"`
	DueDate         time.Time     `json:"due_date,omitempty" db:"due_date"`
	PaymentDetail   PaymentDetail `has_one:"payment_detail" json:"payment_detail,omitempty" db:"-"`
	CustomerDetails Customer      `has_one:"customer_details" json:"customer_details,omitempty" db:"-"`
	Customer        Customer      `belongs_to:"customer" json:"-" db:"-"`
	CustomerID      int           `json:"customer_id" db:"customer_id"`
	Status          Status        `json:"status" db:"status"`
	Note            string        `json:"invoice_note" db:"invoice_note"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

func (i Invoice) GenerateInvoiceNumber() string {
	num := fmt.Sprintf("INV-%d", time.Now().Unix())
	return num
}

func (i *Invoice) GetTotal(tx *pop.Connection) {
	// total := 0.0
	for _, item := range i.Items {
		item.Invoice = *i
		item.Price = item.GetPrice()
		tx.Save(&item)
		i.Total += item.UnitPrice * float64(item.Quantity)
	}
	tx.Update(i)
}

func (i *Invoice) GetCustomerDetails(tx *pop.Connection) Customer {
	// get the customer details
	err := tx.Find(&i.Customer, i.CustomerID)
	if err != nil {
		fmt.Println(err)
	}
	i.CustomerDetails = i.Customer
	return i.CustomerDetails

}

func (i *Invoice) GetPaymentDetails(tx *pop.Connection) PaymentDetail {
	// get the user details
	err := tx.Find(&i.PaymentDetail, i.CreatedByID)
	if err != nil {
		fmt.Println(err)
	}
	return i.PaymentDetail

}

func (i *Invoice) SetDueDate(year, months, days int) time.Time {
	i.DueDate = time.Now().AddDate(year, months, days)
	return i.DueDate
}

func (i *Invoice) ChangeInvoiceState(st Status) Status {
	i.Status = st
	return i.Status
}

// String is not required by pop and may be deleted
func (i Invoice) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Invoices is not required by pop and may be deleted
type Invoices []Invoice

// String is not required by pop and may be deleted
func (i Invoices) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// create a nethod to sort the invoices by status and return the total amount for each status
func (i *Invoices) SortByStatus() map[string]float64 {
	group := make(map[string]float64)
	for _, invoice := range *i {
		group[string(invoice.Status)] += invoice.Total
	}
	return group
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Invoice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *Invoice) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *Invoice) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
