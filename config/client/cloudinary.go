package client

import (
	"log"
	"os"

	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
)

// InitCloudinary initializes Cloudinary client
func InitCloudinary() *cloudinary.Client {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL != "" {
		client, err := cloudinary.NewClientFromURL(cloudinaryURL)
		if err != nil {
			log.Printf("Warning: Failed to initialize Cloudinary client: %v. File upload will be disabled.", err)
			return nil
		}
		log.Printf("Cloudinary client initialized successfully")
		return client
	}

	// Try individual environment variables
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName != "" && apiKey != "" && apiSecret != "" {
		client, err := cloudinary.NewClient(cloudName, apiKey, apiSecret)
		if err != nil {
			log.Printf("Warning: Failed to initialize Cloudinary client: %v. File upload will be disabled.", err)
			return nil
		}
		log.Printf("Cloudinary client initialized successfully")
		return client
	}

	log.Printf("Cloudinary credentials not provided. File upload will be disabled.")
	return nil
}
