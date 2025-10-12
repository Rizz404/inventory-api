package rest

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/asset"
	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	Service asset.AssetService
}

func NewAssetHandler(app fiber.Router, s asset.AssetService) {
	handler := &AssetHandler{
		Service: s,
	}

	// * Asset routes group
	assets := app.Group("/assets")

	// * Create
	assets.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! uncomment in production
		handler.CreateAsset,
	)

	assets.Get("/", handler.GetAssetsPaginated)
	assets.Get("/statistics", handler.GetAssetStatistics)
	assets.Get("/cursor", handler.GetAssetsCursor)
	assets.Get("/count", handler.CountAssets)
	assets.Get("/tag/:tag", handler.GetAssetByAssetTag)
	assets.Get("/check/tag/:tag", handler.CheckAssetTagExists)
	assets.Get("/check/serial/:serial", handler.CheckSerialNumberExists)
	assets.Get("/check/:id", handler.CheckAssetExists)
	assets.Post("/generate-tag", handler.GenerateAssetTagSuggestion)
	assets.Get("/:id", handler.GetAssetById)
	assets.Patch("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleUser), // ! uncomment in production
		handler.UpdateAsset,
	)
	assets.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! uncomment in production
		handler.DeleteAsset,
	)

	// * Export endpoints
	assets.Post("/export/list",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleUser), // ! uncomment in production
		handler.ExportAssetList,
	)
	assets.Get("/export/statistics",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleUser), // ! uncomment in production
		handler.ExportAssetStatistics,
	)
}

func (h *AssetHandler) parseAssetFiltersAndSort(c *fiber.Ctx) (domain.AssetParams, error) {
	params := domain.AssetParams{}

	// * Parse search query
	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	// * Parse sorting options
	sortBy := c.Query("sortBy")
	if sortBy != "" {
		sortOrder := c.Query("sortOrder", "desc")
		params.Sort = &domain.AssetSortOptions{
			Field: domain.AssetSortField(sortBy),
			Order: domain.SortOrder(sortOrder),
		}
	}

	// * Parse filtering options
	statusStr := c.Query("status")
	var status *domain.AssetStatus
	if statusStr != "" {
		s := domain.AssetStatus(statusStr)
		status = &s
	}

	conditionStr := c.Query("condition")
	var condition *domain.AssetCondition
	if conditionStr != "" {
		cond := domain.AssetCondition(conditionStr)
		condition = &cond
	}

	filters := &domain.AssetFilterOptions{
		Status:    status,
		Condition: condition,
	}

	if categoryID := c.Query("categoryId"); categoryID != "" {
		filters.CategoryID = &categoryID
	}

	if locationID := c.Query("locationId"); locationID != "" {
		filters.LocationID = &locationID
	}

	if assignedTo := c.Query("assignedTo"); assignedTo != "" {
		filters.AssignedTo = &assignedTo
	}

	if brand := c.Query("brand"); brand != "" {
		filters.Brand = &brand
	}

	if model := c.Query("model"); model != "" {
		filters.Model = &model
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	var payload domain.CreateAssetPayload

	// Check content type to determine parsing method
	contentType := c.Get("Content-Type")
	var dataMatrixImageFile *multipart.FileHeader

	if strings.Contains(contentType, "multipart/form-data") {
		// Parse multipart form data
		if err := web.ParseFormAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}

		// Try to get data matrix image file (optional)
		file, err := c.FormFile("dataMatrixImage")
		if err == nil {
			// Validate data matrix image file before processing (max 10MB for QR/barcode images)
			if validationErr := web.ValidateImageFile(file, "dataMatrixImage", 10); validationErr != nil {
				return web.HandleError(c, domain.ErrBadRequest(web.FormatFileValidationError(validationErr)))
			}
			dataMatrixImageFile = file
		}
		// Note: We don't return error if data matrix image file is missing since it's optional
	} else {
		// Parse JSON/form-urlencoded data
		if err := web.ParseAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}
	}

	asset, err := h.Service.CreateAsset(c.Context(), &payload, dataMatrixImageFile, web.GetLanguageFromContext(c))
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessAssetCreatedKey, asset)
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetIDRequiredKey))
	}

	var payload domain.UpdateAssetPayload

	// Check content type to determine parsing method
	contentType := c.Get("Content-Type")
	var dataMatrixImageFile *multipart.FileHeader

	if strings.Contains(contentType, "multipart/form-data") {
		// Parse multipart form data
		if err := web.ParseFormAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}

		// Try to get data matrix image file (optional)
		file, err := c.FormFile("dataMatrixImage")
		if err == nil {
			// Validate data matrix image file before processing (max 10MB for QR/barcode images)
			if validationErr := web.ValidateImageFile(file, "dataMatrixImage", 10); validationErr != nil {
				return web.HandleError(c, domain.ErrBadRequest(web.FormatFileValidationError(validationErr)))
			}
			dataMatrixImageFile = file
		}
		// Note: We don't return error if data matrix image file is missing since it's optional
	} else {
		// Parse JSON/form-urlencoded data
		if err := web.ParseAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}
	}

	asset, err := h.Service.UpdateAsset(c.Context(), id, &payload, dataMatrixImageFile, web.GetLanguageFromContext(c))
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetUpdatedKey, asset)
}

func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetIDRequiredKey))
	}

	err := h.Service.DeleteAsset(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetDeletedKey, nil)
}

// *===========================QUERY===========================*
func (h *AssetHandler) GetAssetsPaginated(c *fiber.Ctx) error {
	params, err := h.parseAssetFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.PaginationOptions{Limit: limit, Offset: offset}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	assets, total, err := h.Service.GetAssetsPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessAssetRetrievedKey, assets, int(total), limit, (offset/limit)+1)
}

func (h *AssetHandler) GetAssetsCursor(c *fiber.Ctx) error {
	params, err := h.parseAssetFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &domain.PaginationOptions{Limit: limit, Cursor: cursor}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	assets, err := h.Service.GetAssetsCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(assets) == limit
	if hasNextPage {
		nextCursor = assets[len(assets)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessAssetRetrievedKey, assets, nextCursor, hasNextPage, limit)
}

func (h *AssetHandler) GetAssetById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetIDRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	asset, err := h.Service.GetAssetById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetRetrievedKey, asset)
}

func (h *AssetHandler) GetAssetByAssetTag(c *fiber.Ctx) error {
	tag := c.Params("tag")
	if tag == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetTagRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	asset, err := h.Service.GetAssetByAssetTag(c.Context(), tag, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetRetrievedByTagKey, asset)
}

func (h *AssetHandler) CheckAssetExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetIDRequiredKey))
	}

	exists, err := h.Service.CheckAssetExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *AssetHandler) CheckAssetTagExists(c *fiber.Ctx) error {
	tag := c.Params("tag")
	if tag == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetTagRequiredKey))
	}

	exists, err := h.Service.CheckAssetTagExists(c.Context(), tag)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetTagExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *AssetHandler) CheckSerialNumberExists(c *fiber.Ctx) error {
	serial := c.Params("serial")
	if serial == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetSerialNumberRequiredKey))
	}

	exists, err := h.Service.CheckSerialNumberExists(c.Context(), serial)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetSerialNumberExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *AssetHandler) CountAssets(c *fiber.Ctx) error {
	params, err := h.parseAssetFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	count, err := h.Service.CountAssets(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetCountedKey, count)
}

func (h *AssetHandler) GetAssetStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetAssetStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetStatisticsRetrievedKey, stats)
}

func (h *AssetHandler) GenerateAssetTagSuggestion(c *fiber.Ctx) error {
	var payload domain.GenerateAssetTagPayload

	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	response, err := h.Service.GenerateAssetTagSuggestion(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetTagGeneratedKey, response)
}

// *===========================EXPORT===========================*
func (h *AssetHandler) ExportAssetList(c *fiber.Ctx) error {
	var payload domain.ExportAssetListPayload

	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// Get language from headers
	langCode := web.GetLanguageFromContext(c)

	// Export asset list
	data, filename, err := h.Service.ExportAssetList(c.Context(), &payload, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	// Set appropriate content type and headers
	var contentType string
	switch payload.Format {
	case domain.ExportFormatPDF:
		contentType = "application/pdf"
	case domain.ExportFormatExcel:
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(data)
}

func (h *AssetHandler) ExportAssetStatistics(c *fiber.Ctx) error {
	// Get language from headers
	langCode := web.GetLanguageFromContext(c)

	// Export asset statistics
	data, filename, err := h.Service.ExportAssetStatistics(c.Context(), langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	// Set appropriate content type and headers for PDF
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(data)
}
