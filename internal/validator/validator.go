package validator

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator type - map
type Validator struct {
	Errors map[string]string
}

// Create a new Validator using a factory function
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// checks the data
func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

// error is added to the map, error is then prompted
func (v *Validator) AddError(field string, message string) {
	_, exists := v.Errors[field]
	if !exists {
		v.Errors[field] = message
	}
}
func (v *Validator) Check(ok bool, field string, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// returns true if data is present in the input box
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxLength returns true if the value contains no more than n characters
func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinLength returns true if the value contains at least n characters
func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func IsValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}

// validate option
// AllowedStatus checks if the given status is valid
func AllowedStatus(status string) bool {
	allowed := map[string]bool{
		"pending":     true,
		"in progress": true,
		"completed":   true,
	}
	return allowed[status]
}

func IsValidDate(date time.Time) bool {
	if date.IsZero() {
		return false
	}
	loc, _ := time.LoadLocation("Local") // Use local timezone
	today := time.Now().In(loc).Truncate(24 * time.Hour)
	date = date.In(loc).Truncate(24 * time.Hour)

	return !date.Before(today)
}

func IsPositiveQuantity(quantity int64) bool {
	return quantity > 0
}
func IsPositivePrice(price float64) bool {
	return price > 0
}
