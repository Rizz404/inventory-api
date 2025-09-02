package domain

import "time"

// --- Structs ---

type Location struct {
	ID           string                `json:"id"`
	LocationCode string                `json:"locationCode"`
	Building     *string               `json:"building"`
	Floor        *string               `json:"floor"`
	Latitude     *float64              `json:"latitude"`
	Longitude    *float64              `json:"longitude"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
	Translations []LocationTranslation `json:"translations,omitempty"`
}

type LocationTranslation struct {
	ID           string `json:"id"`
	LocationID   string `json:"locationId"`
	LangCode     string `json:"langCode"`
	LocationName string `json:"locationName"`
}

type LocationResponse struct {
	ID           string   `json:"id"`
	LocationCode string   `json:"locationCode"`
	Building     *string  `json:"building"`
	Floor        *string  `json:"floor"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Name         string   `json:"name"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}

// --- Payloads ---

type CreateLocationPayload struct {
	LocationCode string                             `json:"locationCode" validate:"required,max=20"`
	Building     *string                            `json:"building,omitempty" validate:"omitempty,max=100"`
	Floor        *string                            `json:"floor,omitempty" validate:"omitempty,max=20"`
	Latitude     *float64                           `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude    *float64                           `json:"longitude,omitempty" validate:"omitempty,longitude"`
	Translations []CreateLocationTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateLocationTranslationPayload struct {
	LangCode     string `json:"langCode" validate:"required,max=5"`
	LocationName string `json:"locationName" validate:"required,max=100"`
}

type UpdateLocationPayload struct {
	LocationCode *string                            `json:"locationCode,omitempty" validate:"omitempty,max=20"`
	Building     *string                            `json:"building,omitempty" validate:"omitempty,max=100"`
	Floor        *string                            `json:"floor,omitempty" validate:"omitempty,max=20"`
	Latitude     *float64                           `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude    *float64                           `json:"longitude,omitempty" validate:"omitempty,longitude"`
	Translations []UpdateLocationTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type UpdateLocationTranslationPayload struct {
	LangCode     string  `json:"langCode" validate:"required,max=5"`
	LocationName *string `json:"locationName,omitempty" validate:"omitempty,max=100"`
}
