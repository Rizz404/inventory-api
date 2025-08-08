package web

import (
	"errors"
	"log"
	"math"

	"github.com/Rizz404/inventory-api/domain"
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

// * Semuanya jadi return error biar bisa return di same line
func Success(c *fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessWithPageInfo(c *fiber.Ctx, code int, message string, data any, total int, perPage int, currentPage int) error {
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

func SuccessWithCursor(c *fiber.Ctx, code int, message string, data any, nextCursor string, hasNextPage bool, perPage int, total ...int) error {
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

// * Otomatis handle error berdasarkan domain error (abstraksi terus sampe mampus)
func HandleError(c *fiber.Ctx, err error) error {
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
