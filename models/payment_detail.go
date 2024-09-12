package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// PaymentDetail is used by pop to map your payment_details database table to your go code.
type PaymentDetail struct {
	ID            int       `json:"id" db:"id"`
	BankName      string    `json:"bank_name" db:"bank_name"`
	AccountName   string    `json:"account_name" db:"account_name"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	CreatedBy     User      `belongs_to:"user" json:"-" db:"-"`
	CreatedByID   int       `json:"created_by" db:"created_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (p PaymentDetail) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PaymentDetails is not required by pop and may be deleted
type PaymentDetails []PaymentDetail

// String is not required by pop and may be deleted
func (p PaymentDetails) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PaymentDetail) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PaymentDetail) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PaymentDetail) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
