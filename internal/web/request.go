package web

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"slices"
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
	case reflect.Pointer:
		elemKind := field.Type().Elem().Kind()

		// Handle pointer to string or custom string types
		if elemKind == reflect.String {
			// Create new value of the correct type
			newVal := reflect.New(field.Type().Elem())
			newVal.Elem().SetString(value)
			field.Set(newVal)
		} else {
			// Handle other pointer types recursively
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

// *===========================FILE VALIDATION===========================*

// FileValidationError represents a file validation error with detailed information
type FileValidationError struct {
	Field   string
	Message string
}

func (e *FileValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateImageFile validates an image file with detailed error messages
func ValidateImageFile(file *multipart.FileHeader, fieldName string, maxSizeMB int) error {
	if file == nil {
		return nil // Optional file
	}

	// Check file size
	maxSizeBytes := int64(maxSizeMB * 1024 * 1024)
	if file.Size > maxSizeBytes {
		sizeMB := float64(file.Size) / (1024 * 1024)
		return &FileValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("File size too large (%.2f MB). Maximum allowed size is %d MB", sizeMB, maxSizeMB),
		}
	}

	if file.Size == 0 {
		return &FileValidationError{
			Field:   fieldName,
			Message: "File is empty (0 bytes)",
		}
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".tif", ".svg", ".ico", ".heic", ".heif", ".avif"}

	isValidExtension := slices.Contains(allowedExtensions, ext)

	if !isValidExtension {
		return &FileValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("Invalid file type '%s'. Allowed types: JPG, JPEG, PNG, GIF, WEBP, BMP, TIFF, SVG, ICO, HEIC, HEIF, AVIF", ext),
		}
	}

	// Check if filename is too long
	if len(file.Filename) > 255 {
		return &FileValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("Filename too long (%d characters). Maximum allowed is 255 characters", len(file.Filename)),
		}
	}

	// Try to open file to ensure it's readable
	src, err := file.Open()
	if err != nil {
		return &FileValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("Cannot read file: %s", err.Error()),
		}
	}
	defer src.Close()

	// Read first few bytes to verify it's actually an image
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil {
		return &FileValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("Cannot read file content: %s", err.Error()),
		}
	}

	// Check magic numbers for common image formats
	if n > 0 {
		// JPEG: FF D8 FF
		isJPEG := n >= 3 && buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF
		// PNG: 89 50 4E 47
		isPNG := n >= 4 && buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47
		// GIF: 47 49 46 38
		isGIF := n >= 4 && buffer[0] == 0x47 && buffer[1] == 0x49 && buffer[2] == 0x46 && buffer[3] == 0x38
		// WEBP: 52 49 46 46 ... 57 45 42 50
		isWEBP := n >= 12 && buffer[0] == 0x52 && buffer[1] == 0x49 && buffer[2] == 0x46 && buffer[3] == 0x46 &&
			buffer[8] == 0x57 && buffer[9] == 0x45 && buffer[10] == 0x42 && buffer[11] == 0x50
		// BMP: 42 4D
		isBMP := n >= 2 && buffer[0] == 0x42 && buffer[1] == 0x4D
		// TIFF: 49 49 2A 00 (little-endian) or 4D 4D 00 2A (big-endian)
		isTIFF := n >= 4 && ((buffer[0] == 0x49 && buffer[1] == 0x49 && buffer[2] == 0x2A && buffer[3] == 0x00) ||
			(buffer[0] == 0x4D && buffer[1] == 0x4D && buffer[2] == 0x00 && buffer[3] == 0x2A))
		// SVG: starts with < or <?xml
		isSVG := n >= 5 && ((buffer[0] == 0x3C && buffer[1] == 0x3F && buffer[2] == 0x78 && buffer[3] == 0x6D && buffer[4] == 0x6C) || // <?xml
			(buffer[0] == 0x3C && buffer[1] == 0x73 && buffer[2] == 0x76 && buffer[3] == 0x67)) // <svg
		// ICO: 00 00 01 00
		isICO := n >= 4 && buffer[0] == 0x00 && buffer[1] == 0x00 && buffer[2] == 0x01 && buffer[3] == 0x00
		// HEIC/HEIF: starts with various patterns, checking for 'ftyp' at offset 4
		isHEIC := n >= 12 && buffer[4] == 0x66 && buffer[5] == 0x74 && buffer[6] == 0x79 && buffer[7] == 0x70 &&
			(buffer[8] == 0x68 && buffer[9] == 0x65 && buffer[10] == 0x69 && buffer[11] == 0x63) // heic
		isHEIF := n >= 12 && buffer[4] == 0x66 && buffer[5] == 0x74 && buffer[6] == 0x79 && buffer[7] == 0x70 &&
			(buffer[8] == 0x6D && buffer[9] == 0x69 && buffer[10] == 0x66 && buffer[11] == 0x31) // mif1
		// AVIF: similar to HEIC but with 'avif' identifier
		isAVIF := n >= 12 && buffer[4] == 0x66 && buffer[5] == 0x74 && buffer[6] == 0x79 && buffer[7] == 0x70 &&
			(buffer[8] == 0x61 && buffer[9] == 0x76 && buffer[10] == 0x69 && buffer[11] == 0x66) // avif

		if !isJPEG && !isPNG && !isGIF && !isWEBP && !isBMP && !isTIFF && !isSVG && !isICO && !isHEIC && !isHEIF && !isAVIF {
			return &FileValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("File '%s' is not a valid image file. The file content does not match any supported image format", file.Filename),
			}
		}
	}

	return nil
}

// FormatFileValidationError formats a file validation error into a user-friendly message
func FormatFileValidationError(err error) string {
	if validationErr, ok := err.(*FileValidationError); ok {
		return validationErr.Error()
	}
	return err.Error()
}
