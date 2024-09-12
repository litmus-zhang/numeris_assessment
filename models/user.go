package models

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User is a generated model from buffalo-auth, it serves as the base for username/password authentication.
type User struct {
	ID                   int       `json:"id" db:"id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
	Email                string    `json:"email" db:"email"`
	PasswordHash         string    `json:"-" db:"password_hash"`
	Password             string    `json:"password,omitempty" db:"-"`
	PasswordConfirmation string    `json:"confirm_password,omitempty" db:"-"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	return tx.ValidateAndCreate(u)
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *User) CreateJwtToken() Token {
	refresh_exp := time.Now().Add(7 * 24 * time.Hour).Unix() // expire in 1 week
	access_exp := time.Now().Add(1 * 24 * time.Hour).Unix()  // expire in 1 day
	acces_claims := jwt.StandardClaims{
		ExpiresAt: access_exp,
		Id:        strconv.Itoa(u.ID),
	}
	ref_claims := jwt.StandardClaims{
		ExpiresAt: refresh_exp,
		Id:        strconv.Itoa(u.ID),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, acces_claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ref_claims)
	access_signingString := os.Getenv("ACCESS_JWT_SECRET")
	refresh_signingString := os.Getenv("REFRESH_JWT_SECRET")
	access_token, _ := accessToken.SignedString([]byte(access_signingString))
	refresh_token, _ := refreshToken.SignedString([]byte(refresh_signingString))
	return Token{access_token, refresh_token}

}

func (u *User) VerifyJwtToken(tokenString string) (bool, error) {
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_JWT_SECRET")), nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}

// decode the token
func (u *User) DecodeJwtToken(tokenString string) (int, error) {
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, nil
	}
	id, err := strconv.Atoi(claims.Id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserFromToken gets the user from the token
func (u *User) GetUserFromToken(tokenString string) (*User, error) {
	id, err := u.DecodeJwtToken(tokenString)
	if err != nil {
		return nil, err
	}
	user := &User{}
	tx, err := pop.Connect("")
	if err != nil {
		return nil, err
	}
	err = tx.Find(user, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
