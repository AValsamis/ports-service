package domain

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Port represents a port entity.
type Port struct {
	UNLOC       string    `json:"unloc" validate:"required,len=5"`
	Name        string    `json:"name" validate:"required"`
	City        string    `json:"city" validate:"required"`
	Country     string    `json:"country" validate:"required"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	UNLOCs      []string  `json:"unlocs" validate:"dive,len=5"`
	Code        string    `json:"code"`
}

// Validate performs validation on the Port struct.
// If validation fails, it returns an error with the details of the validation errors.
func (p *Port) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		return fmt.Errorf("validation error in struct: %+v, with details: '%v'", *p, strings.Join(validationErrors, ", "))
	}
	return nil
}
