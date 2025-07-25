package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateFormData(c *fiber.Ctx, dto any) (map[string]string, error) {
	if err := parseFormData(c, dto); err != nil {
		return nil, err
	}

	validationErrors := ValidateFormDataStruct(dto)
	return validationErrors, nil
}

func parseFormData(c *fiber.Ctx, dto any) error {
	v := reflect.ValueOf(dto)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dto must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		formFieldName := getFormFieldName(fieldType)
		formValue := c.FormValue(formFieldName)

		if formValue == "" {
			continue
		}

		if err := setFieldValue(field, formValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", formFieldName, err)
		}
	}

	return nil
}

func getFormFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		return strings.Split(jsonTag, ",")[0]
	}

	formTag := field.Tag.Get("form")
	if formTag != "" && formTag != "-" {
		return strings.Split(formTag, ",")[0]
	}

	return strings.ToLower(field.Name)
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Ptr:
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(&value))
		}
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}

	return nil
}

func ValidateFormDataStruct(dto any) map[string]string {
	validate := validator.New()
	err := validate.Struct(dto)

	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()

		errors[field] = fmt.Sprintf("field %s is not valid", field)
	}

	return errors
}