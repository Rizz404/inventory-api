package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
)

type Client struct {
	cld *cloudinary.Cloudinary
}

type UploadConfig struct {
	AllowedTypes []string `json:"allowedTypes"` // e.g., ["image/jpeg", "image/png", "image/gif"]
	FolderName   string   `json:"folderName"`   // e.g., "avatars", "documents"
	InputName    string   `json:"inputName"`    // e.g., "avatar", "file"
	MaxFiles     int      `json:"maxFiles"`     // Maximum number of files for multiple upload
	MaxFileSize  int64    `json:"maxFileSize"`  // Maximum file size in bytes (e.g., 5MB = 5*1024*1024)
	PublicID     *string  `json:"publicId"`     // Optional custom public ID
	Overwrite    bool     `json:"overwrite"`    // Whether to overwrite existing files
}

type UploadResult struct {
	PublicID     string `json:"publicId"`
	URL          string `json:"url"`
	SecureURL    string `json:"secureUrl"`
	Format       string `json:"format"`
	ResourceType string `json:"resourceType"`
	Bytes        int    `json:"bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	OriginalName string `json:"originalName"`
}

type MultiUploadResult struct {
	Results []UploadResult `json:"results"`
	Failed  []UploadError  `json:"failed"`
}

type UploadError struct {
	FileName string `json:"fileName"`
	Error    string `json:"error"`
}

// NewClient creates a new Cloudinary client
func NewClient(cloudName, apiKey, apiSecret string) (*Client, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudinary client: %w", err)
	}

	return &Client{
		cld: cld,
	}, nil
}

// NewClientFromURL creates a new Cloudinary client from cloudinary URL
func NewClientFromURL(cloudinaryURL string) (*Client, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudinary client from URL: %w", err)
	}

	return &Client{
		cld: cld,
	}, nil
}

// UploadSingleFile uploads a single file to Cloudinary
func (c *Client) UploadSingleFile(ctx context.Context, file *multipart.FileHeader, config UploadConfig) (*UploadResult, error) {
	// Validate file type
	if err := c.validateFileType(file, config.AllowedTypes); err != nil {
		return nil, err
	}

	// Validate file size
	if err := c.validateFileSize(file, config.MaxFileSize); err != nil {
		return nil, err
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Prepare upload parameters
	uploadParams := uploader.UploadParams{
		Folder:    config.FolderName,
		Overwrite: api.Bool(config.Overwrite),
	}

	// Set public ID if provided
	if config.PublicID != nil && *config.PublicID != "" {
		uploadParams.PublicID = *config.PublicID
	}

	// Upload to Cloudinary
	result, err := c.cld.Upload.Upload(ctx, src, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to cloudinary: %w", err)
	}

	return &UploadResult{
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Bytes:        result.Bytes,
		Width:        result.Width,
		Height:       result.Height,
		OriginalName: file.Filename,
	}, nil
}

// UploadMultipleFiles uploads multiple files to Cloudinary
func (c *Client) UploadMultipleFiles(ctx context.Context, fiberCtx *fiber.Ctx, config UploadConfig) (*MultiUploadResult, error) {
	// Parse multipart form
	form, err := fiberCtx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	files := form.File[config.InputName]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found with input name: %s", config.InputName)
	}

	// Validate number of files
	if len(files) > config.MaxFiles {
		return nil, fmt.Errorf("too many files: got %d, max allowed %d", len(files), config.MaxFiles)
	}

	result := &MultiUploadResult{
		Results: make([]UploadResult, 0),
		Failed:  make([]UploadError, 0),
	}

	// Upload each file
	for _, file := range files {
		uploadResult, err := c.UploadSingleFile(ctx, file, config)
		if err != nil {
			result.Failed = append(result.Failed, UploadError{
				FileName: file.Filename,
				Error:    err.Error(),
			})
			continue
		}

		result.Results = append(result.Results, *uploadResult)
	}

	return result, nil
}

// DeleteFile deletes a file from Cloudinary by public ID
func (c *Client) DeleteFile(ctx context.Context, publicID string) error {
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from cloudinary: %w", err)
	}

	return nil
}

// GetFileInfo gets file information from Cloudinary
func (c *Client) GetFileInfo(ctx context.Context, publicID string) (*UploadResult, error) {
	result, err := c.cld.Admin.Asset(ctx, admin.AssetParams{
		PublicID: publicID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info from cloudinary: %w", err)
	}

	return &UploadResult{
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Bytes:        result.Bytes,
		Width:        result.Width,
		Height:       result.Height,
	}, nil
}

// validateFileType validates if the file type is allowed
func (c *Client) validateFileType(file *multipart.FileHeader, allowedTypes []string) error {
	if len(allowedTypes) == 0 {
		return nil // No restrictions
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Check against allowed types
	for _, allowedType := range allowedTypes {
		// Support both MIME types and extensions
		if strings.HasPrefix(allowedType, "image/") || strings.HasPrefix(allowedType, "application/") {
			// Handle MIME type validation by opening and checking the file
			src, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file for type validation: %w", err)
			}
			defer src.Close()

			// Read first 512 bytes to detect content type
			buffer := make([]byte, 512)
			_, err = src.Read(buffer)
			if err != nil {
				return fmt.Errorf("failed to read file for type validation: %w", err)
			}

			// Use simple extension-based validation for now
			if strings.HasSuffix(allowedType, ext[1:]) {
				return nil
			}
		} else {
			// Handle extension validation
			if ext == allowedType || ext == "."+allowedType {
				return nil
			}
		}
	}

	return fmt.Errorf("file type not allowed: %s, allowed types: %v", ext, allowedTypes)
}

// validateFileSize validates if the file size is within limits
func (c *Client) validateFileSize(file *multipart.FileHeader, maxSize int64) error {
	if maxSize <= 0 {
		return nil // No size restriction
	}

	if file.Size > maxSize {
		return fmt.Errorf("file size too large: %d bytes, max allowed: %d bytes", file.Size, maxSize)
	}

	return nil
}

// GenerateTransformationURL generates a URL with Cloudinary transformations
func (c *Client) GenerateTransformationURL(publicID string, transformations string) (string, error) {
	url, err := c.cld.Image(publicID)
	if err != nil {
		return "", fmt.Errorf("failed to generate transformation URL: %w", err)
	}

	urlString, err := url.String()
	if err != nil {
		return "", fmt.Errorf("failed to convert URL to string: %w", err)
	}

	// Apply transformations if provided
	if transformations != "" {
		urlString = urlString + "/" + transformations
	}

	return urlString, nil
}

// GetAvatarUploadConfig returns a pre-configured upload config for user avatars
func GetAvatarUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		FolderName:   "avatars",
		InputName:    "avatar",
		MaxFiles:     1,
		MaxFileSize:  5 * 1024 * 1024, // 5MB
		Overwrite:    true,
	}
}

// GetDataMatrixImageUploadConfig returns a pre-configured upload config for asset data matrix images
func GetDataMatrixImageUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		FolderName:   "datamatrix",
		InputName:    "dataMatrixImage",
		MaxFiles:     1,
		MaxFileSize:  2 * 1024 * 1024, // 2MB
		Overwrite:    true,
	}
}

// GetDocumentUploadConfig returns a pre-configured upload config for documents
func GetDocumentUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{"application/pdf", "image/jpeg", "image/png"},
		FolderName:   "documents",
		InputName:    "documents",
		MaxFiles:     10,
		MaxFileSize:  10 * 1024 * 1024, // 10MB
		Overwrite:    false,
	}
}
