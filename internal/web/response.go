package web

import (
	"errors"
	"log"
	"math"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type PageInfo struct {
	Total       int  `json:"total"`
	PerPage     int  `json:"per_page"`
	CurrentPage int  `json:"current_page"`
	TotalPages  int  `json:"total_pages"`
	HasPrevPage bool `json:"has_prev_page"`
	HasNextPage bool `json:"has_next_page"`
}

// * Buat yang cursor base pagination
type CursorInfo struct {
	NextCursor  string `json:"next_cursor"`
	HasNextPage bool   `json:"has_next_page"`
	PerPage     int    `json:"per_page"`
	Total       int    `json:"total,omitempty"`
}

type JSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	// * Tergantung datanya jadi bisa gak ada
	PageInfo   *PageInfo   `json:"pagination,omitempty"`
	CursorInfo *CursorInfo `json:"cursor,omitempty"`
}

// * Success creates a localized success response
func Success(c *fiber.Ctx, code int, messageKey utils.MessageKey, data any) error {
	// * Get language from context
	langCode := GetLanguageFromContext(c)
	message := utils.GetLocalizedMessage(messageKey, langCode)

	return c.Status(code).JSON(JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// * SuccessWithMessage creates a success response with custom message (for backward compatibility)
func SuccessWithMessage(c *fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// * SuccessWithPageInfo creates a localized success response with pagination info
func SuccessWithPageInfo(c *fiber.Ctx, code int, messageKey utils.MessageKey, data any, total int, perPage int, currentPage int) error {
	// * Get language from context
	langCode := GetLanguageFromContext(c)
	message := utils.GetLocalizedMessage(messageKey, langCode)

	totalPages := 0
	if perPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	}

	pageInfo := &PageInfo{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		HasPrevPage: currentPage > 1,
		HasNextPage: currentPage < totalPages,
	}

	return c.Status(code).JSON(JSONResponse{
		Status:   "success",
		Message:  message,
		Data:     data,
		PageInfo: pageInfo,
	})
}

// * SuccessWithPageInfoAndMessage creates a success response with pagination info and custom message (for backward compatibility)
func SuccessWithPageInfoAndMessage(c *fiber.Ctx, code int, message string, data any, total int, perPage int, currentPage int) error {
	totalPages := 0
	if perPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	}

	pageInfo := &PageInfo{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		HasPrevPage: currentPage > 1,
		HasNextPage: currentPage < totalPages,
	}

	return c.Status(code).JSON(JSONResponse{
		Status:   "success",
		Message:  message,
		Data:     data,
		PageInfo: pageInfo,
	})
}

// * SuccessWithCursor creates a localized success response with cursor pagination info
func SuccessWithCursor(c *fiber.Ctx, code int, messageKey utils.MessageKey, data any, nextCursor string, hasNextPage bool, perPage int, total ...int) error {
	// * Get language from context
	langCode := GetLanguageFromContext(c)
	message := utils.GetLocalizedMessage(messageKey, langCode)

	cursorInfo := &CursorInfo{
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
		PerPage:     perPage,
	}

	if len(total) > 0 {
		cursorInfo.Total = total[0]
	}

	return c.Status(code).JSON(JSONResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		CursorInfo: cursorInfo,
	})
}

// * SuccessWithCursorAndMessage creates a success response with cursor pagination info and custom message (for backward compatibility)
func SuccessWithCursorAndMessage(c *fiber.Ctx, code int, message string, data any, nextCursor string, hasNextPage bool, perPage int, total ...int) error {
	cursorInfo := &CursorInfo{
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
		PerPage:     perPage,
	}

	if len(total) > 0 {
		cursorInfo.Total = total[0]
	}

	return c.Status(code).JSON(JSONResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		CursorInfo: cursorInfo,
	})
}

func Error(c *fiber.Ctx, code int, message string, errorDetails any) error {
	return c.Status(code).JSON(JSONResponse{
		Status:  "error",
		Message: message,
		Error:   errorDetails,
	})
}

// * HandleError creates a localized error response
func HandleError(c *fiber.Ctx, err error) error {
	// * Get language from context
	langCode := GetLanguageFromContext(c)

	// * Handle validation errors first
	var validationErr *FiberValidationError
	if errors.As(err, &validationErr) {
		message := utils.GetLocalizedMessage(utils.ErrValidationKey, langCode)
		return Error(c, fiber.StatusBadRequest, message, validationErr.Errors)
	}

	// * Handle domain app errors
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		if appErr.Code >= 500 {
			// * Log error aslinya dari AppError
			log.Printf("error: internal server error handled: %v", appErr.Unwrap())
		}

		// * Use localized message if available
		message := appErr.GetLocalizedMessage(langCode)
		return Error(c, appErr.Code, message, nil)
	}

	// * Handle fiber errors
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(c, fiberErr.Code, fiberErr.Message, nil)
	}

	// * Untuk error yang tidak terduga, ini adalah fallback
	log.Printf("error: unexpected internal error: %v", err)
	message := utils.GetLocalizedMessage(utils.ErrInternalKey, langCode)
	return Error(c, fiber.StatusInternalServerError, message, nil)
}

// * HandleErrorWithMessage creates an error response with custom message (for backward compatibility)
func HandleErrorWithMessage(c *fiber.Ctx, err error) error {
	// * Handle validation errors first
	var validationErr *FiberValidationError
	if errors.As(err, &validationErr) {
		return Error(c, fiber.StatusBadRequest, "Validation failed", validationErr.Errors)
	}

	// * Handle domain app errors
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		if appErr.Code >= 500 {
			// * Log error aslinya dari AppError
			log.Printf("error: internal server error handled: %v", appErr.Unwrap())
		}
		return Error(c, appErr.Code, appErr.Error(), nil)
	}

	// * Handle fiber errors
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(c, fiberErr.Code, fiberErr.Message, nil)
	}

	// * Untuk error yang tidak terduga, ini adalah fallback
	log.Printf("error: unexpected internal error: %v", err)
	return Error(c, fiber.StatusInternalServerError, "An unexpected error occurred", nil)
}
