package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortValidation(t *testing.T) {
	validPort := Port{
		Name:    "Jebel Ali",
		City:    "Dubai",
		Country: "United Arab Emirates",
		UNLOC:   "AEJEA",
	}

	tests := []struct {
		name             string
		port             Port
		expectedErrorMsg string
	}{
		{
			name:             "Valid port",
			port:             validPort,
			expectedErrorMsg: "",
		},
		{
			name:             "Empty name",
			port:             Port{Name: "", City: "Dubai", Country: "United Arab Emirates", UNLOC: "AEJEA"},
			expectedErrorMsg: "Field validation for 'Name'",
		},
		{
			name:             "Empty city",
			port:             Port{Name: "Jebel Ali", City: "", Country: "United Arab Emirates", UNLOC: "AEJEA"},
			expectedErrorMsg: "Field validation for 'City'",
		},
		{
			name:             "Empty country",
			port:             Port{Name: "Jebel Ali", City: "Dubai", Country: "", UNLOC: "AEJEA"},
			expectedErrorMsg: "Field validation for 'Country'",
		},
		{
			name:             "Invalid UNLOC length",
			port:             Port{Name: "Jebel Ali", City: "Dubai", Country: "United Arab Emirates", UNLOC: "AEJE"},
			expectedErrorMsg: "Field validation for 'UNLOC'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.port.Validate()
			if test.expectedErrorMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.expectedErrorMsg)
			}
		})
	}
}
