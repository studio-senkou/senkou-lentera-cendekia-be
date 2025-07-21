package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(data any) map[string]string {
	errors := make(map[string]string)

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := getFieldName(data, err.Field())
			message := getErrorMessage(field, err.Tag(), err.Param())
			errors[field] = message
		}
	}

	return errors
}

func ValidateOnly(c *fiber.Ctx, dto any) (map[string]string, error) {
	if err := c.BodyParser(dto); err != nil {
		return nil, err
	}

	validationErrors := ValidateStruct(dto)
	return validationErrors, nil
}

func ValidateRequest(c *fiber.Ctx, dto any) (map[string]string, error) {
	validationErrors, err := ValidateOnly(c, dto)

	if err != nil {
		return nil, err
	}

	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	return nil, nil
}

func getFieldName(data any, fieldName string) string {
	t := reflect.TypeOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return strings.ToLower(fieldName)
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		return strings.Split(jsonTag, ",")[0]
	}

	return strings.ToLower(fieldName)
}

func getErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("The %s field is required", field)
	case "email":
		return fmt.Sprintf("The %s field must be a valid email address", field)
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("The %s field must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("The %s field must be exactly %s characters long", field, param)
	case "oneof":
		return fmt.Sprintf("The %s field must be one of the following values: %s", field, param)
	case "numeric":
		return fmt.Sprintf("The %s field must be a numeric value", field)
	case "uuid":
		return fmt.Sprintf("The %s field must be a valid UUID", field)
	case "alpha":
		return fmt.Sprintf("The %s field must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("The %s field must contain only alphanumeric characters", field)
	case "url":
		return fmt.Sprintf("The %s field must be a valid URL", field)
	case "gte":
		return fmt.Sprintf("The %s field must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("The %s field must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("The %s field must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("The %s field must be less than %s", field, param)
	default:
		return fmt.Sprintf("The %s field is invalid", field)
	}
}
