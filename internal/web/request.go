package web

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", ve[0].Message)
}

// * Validasi pakai playground validator
func Validate(s any) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors ValidationErrors

		for _, err := range err.(validator.ValidationErrors) {
			validationError := ValidationError{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: fmt.Sprintf("%v", err.Value()),
			}

			validationError.Message = generateValidationMessage(err)
			validationErrors = append(validationErrors, validationError)
		}

		return validationErrors
	}
	return nil
}

func generateValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	default:
		return fmt.Sprintf("%s failed validation for tag '%s'", field, tag)
	}
}

// * Fiber helpers - lebih simple karena Fiber sudah handle banyak hal
func ParseAndValidate(c *fiber.Ctx, dst any) error {
	// Parse body menggunakan built-in Fiber BodyParser
	if err := c.BodyParser(dst); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to parse request body: %v", err))
	}

	// Validasi menggunakan validator
	if err := Validate(dst); err != nil {
		if validationErrors, ok := err.(ValidationErrors); ok {
			return &FiberValidationError{Errors: validationErrors}
		}
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err))
	}

	return nil
}

// * Custom error type untuk validation errors di Fiber
type FiberValidationError struct {
	Errors ValidationErrors
}

func (e *FiberValidationError) Error() string {
	return e.Errors.Error()
}

// * Parse form data untuk multipart/form-data
func ParseFormAndValidate(c *fiber.Ctx, dst any) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to parse multipart form: %v", err))
	}

	// Map form values ke struct
	if err := mapFormValuesToStruct(form.Value, dst); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to map form values: %v", err))
	}

	// Validasi
	if err := Validate(dst); err != nil {
		if validationErrors, ok := err.(ValidationErrors); ok {
			return &FiberValidationError{Errors: validationErrors}
		}
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err))
	}

	return nil
}

// * Helper untuk mapping form values ke struct
func mapFormValuesToStruct(values map[string][]string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		fieldName := getFieldName(fieldType)
		formValues, exists := values[fieldName]
		if !exists || len(formValues) == 0 {
			continue
		}

		// Ambil value pertama
		formValue := formValues[0]
		if err := setFieldValue(field, formValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return nil
}

func getFieldName(fieldType reflect.StructField) string {
	fieldName := fieldType.Name

	// Prioritas: form tag > json tag > field name
	if formTag := fieldType.Tag.Get("form"); formTag != "" {
		if parts := strings.Split(formTag, ","); parts[0] != "" && parts[0] != "-" {
			fieldName = parts[0]
		}
	} else if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
		if parts := strings.Split(jsonTag, ","); parts[0] != "" && parts[0] != "-" {
			fieldName = parts[0]
		}
	}

	return fieldName
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err != nil {
			return err
		} else {
			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintVal, err := strconv.ParseUint(value, 10, 64); err != nil {
			return err
		} else {
			field.SetUint(uintVal)
		}
	case reflect.Float32, reflect.Float64:
		if floatVal, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			field.SetFloat(floatVal)
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(value); err != nil {
			return err
		} else {
			field.SetBool(boolVal)
		}
	case reflect.Ptr:
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(&value))
		} else {
			newVal := reflect.New(field.Type().Elem())
			if err := setFieldValue(newVal.Elem(), value); err != nil {
				return err
			}
			field.Set(newVal)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
