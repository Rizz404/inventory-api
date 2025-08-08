package domain

// --- Structs ---

type Location struct {
	ID           string                `json:"id"`
	LocationCode string                `json:"locationCode"`
	Building     *string               `json:"building"`
	Floor        *string               `json:"floor"`
	Translations []LocationTranslation `json:"translations,omitempty"`
}

type LocationTranslation struct {
	ID           string `json:"id"`
	LocationID   string `json:"locationId"`
	LangCode     string `json:"langCode"`
	LocationName string `json:"locationName"`
}

type LocationResponse struct {
	ID           string  `json:"id"`
	LocationCode string  `json:"locationCode"`
	Building     *string `json:"building,omitempty"`
	Floor        *string `json:"floor,omitempty"`
	Name         string  `json:"name"`
}

// --- Payloads ---

type CreateLocationPayload struct {
	LocationCode string                             `json:"locationCode" validate:"required,max=20"`
	Building     *string                            `json:"building,omitempty" validate:"omitempty,max=100"`
	Floor        *string                            `json:"floor,omitempty" validate:"omitempty,max=20"`
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
	Translations []UpdateLocationTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type UpdateLocationTranslationPayload struct {
	LangCode     string  `json:"langCode" validate:"required,max=5"`
	LocationName *string `json:"locationName,omitempty" validate:"omitempty,max=100"`
}
