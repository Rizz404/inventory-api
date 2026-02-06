package client

import (
	"github.com/Rizz404/inventory-api/internal/client/gtranslate"
)

// InitGTranslate initializes the Google Translate client
func InitGTranslate() *gtranslate.Client {
	return gtranslate.NewClient()
}
