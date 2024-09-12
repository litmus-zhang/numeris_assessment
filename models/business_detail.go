package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// BusinessDetail is used by pop to map your business_details database table to your go code.
type BusinessDetail struct {
	ID           int       `json:"id" db:"id"`
	BusinessName string    `json:"business_name" db:"business_name"`
	Email        string    `json:"email" db:"email"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	Address      string    `json:"address" db:"address"`
	CreatedBy    User      `belongs_to:"user" json:"-" db:"-"`
	CreatedByID  int       `json:"created_by" db:"created_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (b BusinessDetail) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BusinessDetails is not required by pop and may be deleted
type BusinessDetails []BusinessDetail

// String is not required by pop and may be deleted
func (b BusinessDetails) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BusinessDetail) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *BusinessDetail) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *BusinessDetail) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
