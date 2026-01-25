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
)

type Client struct {
	cld *cloudinary.Cloudinary
}

type UploadConfig struct {
	AllowedTypes         []string `json:"allowedTypes"`         // e.g., ["image/jpeg", "image/png", "image/gif"]
	FolderName           string   `json:"folderName"`           // e.g., "avatars", "documents"
	InputName            string   `json:"inputName"`            // e.g., "avatar", "file"
	MaxFiles             int      `json:"maxFiles"`             // Maximum number of files for multiple upload
	MaxFileSize          int64    `json:"maxFileSize"`          // Maximum file size in bytes (e.g., 5MB = 5*1024*1024)
	PublicID             *string  `json:"publicId"`             // Optional custom public ID
	Overwrite            bool     `json:"overwrite"`            // Whether to overwrite existing files
	EagerTransformations string   `json:"eagerTransformations"` // Eager transformations to apply on upload (e.g., "f_webp,q_auto")
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

	// ! Apply eager transformations for optimization (WebP conversion, compression)
	// * Eager transformations are processed during upload, no extra response time for clients
	if config.EagerTransformations != "" {
		uploadParams.Eager = config.EagerTransformations
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

// UploadMultipleFiles uploads multiple files to Cloudinary (framework-agnostic)
// For each file, you can customize the PublicID by modifying baseConfig before calling
func (c *Client) UploadMultipleFiles(ctx context.Context, files []*multipart.FileHeader, baseConfig UploadConfig) (*MultiUploadResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	// Validate number of files
	if baseConfig.MaxFiles > 0 && len(files) > baseConfig.MaxFiles {
		return nil, fmt.Errorf("too many files: got %d, max allowed %d", len(files), baseConfig.MaxFiles)
	}

	result := &MultiUploadResult{
		Results: make([]UploadResult, 0, len(files)),
		Failed:  make([]UploadError, 0),
	}

	// Upload each file with the base config
	for _, file := range files {
		uploadResult, err := c.UploadSingleFile(ctx, file, baseConfig)
		if err != nil {
			result.Failed = append(result.Failed, UploadError{
				FileName: file.Filename,
				Error:    err.Error(),
			})
			continue
		}

		result.Results = append(result.Results, *uploadResult)
	}

	// Return error if all uploads failed
	if len(result.Results) == 0 && len(result.Failed) > 0 {
		return nil, fmt.Errorf("all %d file uploads failed", len(files))
	}

	return result, nil
}

// UploadMultipleFilesWithPublicIDs uploads multiple files with custom public IDs for each file
// This is more efficient than calling UploadSingleFile in a loop when you need different public IDs
func (c *Client) UploadMultipleFilesWithPublicIDs(ctx context.Context, files []*multipart.FileHeader, publicIDs []string, baseConfig UploadConfig) (*MultiUploadResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	if len(files) != len(publicIDs) {
		return nil, fmt.Errorf("number of files (%d) must match number of public IDs (%d)", len(files), len(publicIDs))
	}

	// Validate number of files
	if baseConfig.MaxFiles > 0 && len(files) > baseConfig.MaxFiles {
		return nil, fmt.Errorf("too many files: got %d, max allowed %d", len(files), baseConfig.MaxFiles)
	}

	result := &MultiUploadResult{
		Results: make([]UploadResult, 0, len(files)),
		Failed:  make([]UploadError, 0),
	}

	// Upload each file with its corresponding public ID
	for i, file := range files {
		// Create a copy of config with the specific public ID
		fileConfig := baseConfig
		fileConfig.PublicID = &publicIDs[i]

		uploadResult, err := c.UploadSingleFile(ctx, file, fileConfig)
		if err != nil {
			result.Failed = append(result.Failed, UploadError{
				FileName: file.Filename,
				Error:    err.Error(),
			})
			continue
		}

		result.Results = append(result.Results, *uploadResult)
	}

	// Return error if all uploads failed
	if len(result.Results) == 0 && len(result.Failed) > 0 {
		return nil, fmt.Errorf("all %d file uploads failed", len(files))
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

// DeleteMultipleFiles deletes multiple files from Cloudinary by public IDs
func (c *Client) DeleteMultipleFiles(ctx context.Context, publicIDs []string) (int, []string, error) {
	if len(publicIDs) == 0 {
		return 0, nil, fmt.Errorf("no public IDs provided")
	}

	deletedCount := 0
	failedIDs := []string{}

	// Delete each file individually
	for _, publicID := range publicIDs {
		err := c.DeleteFile(ctx, publicID)
		if err != nil {
			failedIDs = append(failedIDs, publicID)
			continue
		}
		deletedCount++
	}

	return deletedCount, failedIDs, nil
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

	// Check against allowed extensions
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return nil
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
// * Eager transformations convert to WebP + auto quality for optimal file size
// * Free tier compatible, transformations processed during upload (no client delay)
func GetAvatarUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{
			".jpg",
			".jpeg",
			".png",
			".gif",
			".webp",
			".bmp",
			".tiff",
			".tif",
			".svg",
			".ico",
			".heic",
			".heif",
			".avif",
		},
		FolderName:           "sigma-asset/avatars",
		InputName:            "avatar",
		MaxFiles:             1,
		MaxFileSize:          5 * 1024 * 1024, // 5MB
		Overwrite:            true,
		EagerTransformations: "f_webp,q_auto", // Convert to WebP + auto quality
	}
}

// GetDataMatrixImageUploadConfig returns a pre-configured upload config for asset data matrix images
// * Eager transformations: optimize PNG for QR/barcode readability
// * PNG format better for barcodes/QR codes than WebP
func GetDataMatrixImageUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{
			".jpg",
			".jpeg",
			".png",
			".gif",
			".webp",
			".bmp",
			".tiff",
			".tif",
			".svg",
			".avif",
		},
		FolderName:           "sigma-asset/datamatrix",
		InputName:            "dataMatrixImage",
		MaxFiles:             1,
		MaxFileSize:          2 * 1024 * 1024, // 2MB
		Overwrite:            true,
		EagerTransformations: "f_png,q_auto:best", // Keep PNG for barcode clarity + best quality compression
	}
}

// GetBulkDataMatrixImageUploadConfig returns a pre-configured upload config for bulk data matrix images
// * Eager transformations: optimize PNG for QR/barcode readability in bulk
func GetBulkDataMatrixImageUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{
			".jpg",
			".jpeg",
			".png",
			".gif",
			".webp",
			".bmp",
			".tiff",
			".tif",
			".svg",
			".avif",
		},
		FolderName:           "sigma-asset/datamatrix",
		InputName:            "dataMatrixImages",
		MaxFiles:             0,                // No limit for bulk upload
		MaxFileSize:          10 * 1024 * 1024, // 10MB per file
		Overwrite:            true,
		EagerTransformations: "f_png,q_auto:best", // Keep PNG for barcode clarity
	}
}

// GetDocumentUploadConfig returns a pre-configured upload config for documents
// * Eager transformations: auto format + quality for non-PDF images
// * PDFs remain unchanged, images converted to WebP
func GetDocumentUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{
			".pdf",
			".jpg",
			".jpeg",
			".png",
			".gif",
			".webp",
			".tiff",
			".tif",
			".bmp",
		},
		FolderName:           "sigma-asset/documents",
		InputName:            "documents",
		MaxFiles:             10,
		MaxFileSize:          10 * 1024 * 1024, // 10MB
		Overwrite:            false,
		EagerTransformations: "f_auto,q_auto", // Auto format + quality for images
	}
}

// GetCategoryImageUploadConfig returns a pre-configured upload config for category images
// * Eager transformations: resize to max 800px + WebP + auto quality
// * Reduces storage & bandwidth usage on free tier
func GetCategoryImageUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{
			".jpg",
			".jpeg",
			".png",
			".gif",
			".webp",
		},
		FolderName:           "sigma-asset/categories",
		InputName:            "image",
		MaxFiles:             1,
		MaxFileSize:          5 * 1024 * 1024, // 5MB
		Overwrite:            false,
		EagerTransformations: "w_800,c_limit/f_webp,q_auto", // Resize max 800px + WebP + auto quality
	}
}

// GetAssetImageUploadConfig returns a pre-configured upload config for asset images
// * Eager transformations: resize to max 1920px + WebP + auto quality
// * Optimizes for display while reducing storage costs
func GetAssetImageUploadConfig() UploadConfig {
	return UploadConfig{
		AllowedTypes: []string{
			".jpg",
			".jpeg",
			".png",
			".gif",
			".webp",
		},
		FolderName:           "sigma-asset/assets",
		InputName:            "images",
		MaxFiles:             10,               // Allow up to 10 images per asset
		MaxFileSize:          10 * 1024 * 1024, // 10MB per image
		Overwrite:            false,
		EagerTransformations: "w_1920,c_limit/f_webp,q_auto", // Resize max 1920px + WebP + auto quality
	}
}

// ExtractPublicIDFromURL extracts the public ID from a Cloudinary URL
// Example: https://res.cloudinary.com/demo/image/upload/v1234567890/sigma-asset/datamatrix/ASSET-001_01HQXXX.jpg
// Returns: sigma-asset/datamatrix/ASSET-001_01HQXXX
func ExtractPublicIDFromURL(url string) string {
	if url == "" {
		return ""
	}

	// Find the upload segment in the URL
	parts := strings.Split(url, "/upload/")
	if len(parts) < 2 {
		return ""
	}

	// Get everything after /upload/
	afterUpload := parts[1]

	// Split by / to get path segments
	segments := strings.Split(afterUpload, "/")
	if len(segments) < 2 {
		return ""
	}

	// Skip version (v1234567890) if present
	startIdx := 0
	if len(segments) > 0 && strings.HasPrefix(segments[0], "v") {
		startIdx = 1
	}

	// Reconstruct path without version and extension
	pathParts := segments[startIdx:]

	// Remove file extension from last segment
	if len(pathParts) > 0 {
		lastPart := pathParts[len(pathParts)-1]
		if dotIdx := strings.LastIndex(lastPart, "."); dotIdx > 0 {
			pathParts[len(pathParts)-1] = lastPart[:dotIdx]
		}
	}

	return strings.Join(pathParts, "/")
}
