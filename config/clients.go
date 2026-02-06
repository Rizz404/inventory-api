package config

import (
	"github.com/Rizz404/inventory-api/config/client"
	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
	"github.com/Rizz404/inventory-api/internal/client/fcm"
	"github.com/Rizz404/inventory-api/internal/client/gtranslate"
	"github.com/Rizz404/inventory-api/internal/client/smtp"
)

// Clients holds all external service clients
type Clients struct {
	Cloudinary *cloudinary.Client
	FCM        *fcm.Client
	SMTP       *smtp.Client
	Translator *gtranslate.Client
}

// InitializeClients initializes all external service clients
func InitializeClients() *Clients {
	return &Clients{
		Cloudinary: client.InitCloudinary(),
		FCM:        client.InitFCM(),
		SMTP:       client.InitSMTP(),
		Translator: client.InitGTranslate(),
	}
}
