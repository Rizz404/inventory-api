package config

import (
	"github.com/Rizz404/inventory-api/config/client"
	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
	"github.com/Rizz404/inventory-api/internal/client/fcm"
)

// Clients holds all external service clients
type Clients struct {
	Cloudinary *cloudinary.Client
	FCM        *fcm.Client
}

// InitializeClients initializes all external service clients
func InitializeClients() *Clients {
	return &Clients{
		Cloudinary: client.InitCloudinary(),
		FCM:        client.InitFCM(),
	}
}
