package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// DecodeAndValidate reads JSON from the request body and validates the struct.
func DecodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := validate.Struct(dst); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			messages := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				messages = append(messages, formatValidationError(e))
			}
			return fmt.Errorf("%s", strings.Join(messages, "; "))
		}
		return err
	}

	return nil
}

func formatValidationError(e validator.FieldError) string {
	field := toSnakeCase(e.Field())
	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, e.Param())
	case "uuid":
		return field + " must be a valid UUID"
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "gtfield":
		return fmt.Sprintf("%s must be after %s", field, toSnakeCase(e.Param()))
	default:
		return fmt.Sprintf("%s failed validation: %s", field, e.Tag())
	}
}

func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(c+'a'-'A'))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
