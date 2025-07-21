package validator_test

import (
	"testing"

	. "github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

func TestValidateStruct(t *testing.T) {
	type UserDTO struct {
		Name     string `json:"name" validate:"required,min=3"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Age      int    `json:"age" validate:"gte=18,lte=99"`
	}

	tests := []struct {
		name     string
		input    UserDTO
		expected map[string]string
	}{
		{
			name: "all fields valid",
			input: UserDTO{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "secret123",
				Age:      25,
			},
			expected: map[string]string{},
		},
		{
			name: "missing required fields",
			input: UserDTO{
				Name:     "",
				Email:    "",
				Password: "",
				Age:      0,
			},
			expected: map[string]string{
				"name":     "The name field is required",
				"email":    "The email field is required",
				"password": "The password field is required",
				"age":      "The age field must be greater than or equal to 18",
			},
		},
		{
			name: "invalid email and short password",
			input: UserDTO{
				Name:     "Jo",
				Email:    "not-an-email",
				Password: "123",
				Age:      17,
			},
			expected: map[string]string{
				"name":     "The name field must be at least 3 characters long",
				"email":    "The email field must be a valid email address",
				"password": "The password field must be at least 6 characters long",
				"age":      "The age field must be greater than or equal to 18",
			},
		},
		{
			name: "age too high",
			input: UserDTO{
				Name:     "Valid Name",
				Email:    "valid@email.com",
				Password: "validpass",
				Age:      120,
			},
			expected: map[string]string{
				"age": "The age field must be less than or equal to 99",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateStruct(tt.input)
			if len(errs) != len(tt.expected) {
				t.Errorf("expected %d errors, got %d: %v", len(tt.expected), len(errs), errs)
			}
			for k, v := range tt.expected {
				if errs[k] != v {
					t.Errorf("expected error for %q: %q, got: %q", k, v, errs[k])
				}
			}
		})
	}
}
