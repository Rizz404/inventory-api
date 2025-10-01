package client

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/Rizz404/inventory-api/internal/client/fcm"
	"google.golang.org/api/option"
)

// InitFCM initializes Firebase Cloud Messaging client
func InitFCM() *fcm.Client {
	enableFCM := os.Getenv("ENABLE_FCM") == "true"
	if !enableFCM {
		log.Printf("Firebase services disabled via ENABLE_FCM environment variable")
		return nil
	}

	credentialsMap := map[string]string{
		"type":                        os.Getenv("FIREBASE_TYPE"),
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    os.Getenv("FIREBASE_AUTH_URI"),
		"token_uri":                   os.Getenv("FIREBASE_TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
		"universe_domain":             os.Getenv("FIREBASE_UNIVERSE_DOMAIN"),
	}

	credentialsJSON, err := json.Marshal(credentialsMap)
	if err != nil {
		log.Printf("Warning: Failed to marshal Firebase credentials: %v. Firebase services will be disabled.", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opt := option.WithCredentialsJSON(credentialsJSON)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("Warning: Failed to initialize Firebase app: %v. Firebase services will be disabled.", err)
		return nil
	}

	// Initialize FCM Client
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("Warning: Failed to get FCM messaging client: %v. FCM will be disabled.", err)
		return nil
	}

	fcmClient := fcm.NewClientFromMessaging(client)
	log.Printf("FCM client initialized successfully")

	return fcmClient
}
