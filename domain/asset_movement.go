package domain

import "time"

// --- Structs ---

type AssetMovement struct {
	ID             string                     `json:"id"`
	AssetID        string                     `json:"assetId"`
	FromLocationID *string                    `json:"fromLocationId"`
	ToLocationID   *string                    `json:"toLocationId"`
	FromUserID     *string                    `json:"fromUserId"`
	ToUserID       *string                    `json:"toUserId"`
	MovementDate   time.Time                  `json:"movementDate"`
	MovedBy        string                     `json:"movedBy"`
	CreatedAt      time.Time                  `json:"createdAt"`
	UpdatedAt      time.Time                  `json:"updatedAt"`
	Translations   []AssetMovementTranslation `json:"translations,omitempty"`
}

type AssetMovementTranslation struct {
	ID         string  `json:"id"`
	MovementID string  `json:"movementId"`
	LangCode   string  `json:"langCode"`
	Notes      *string `json:"notes"`
}

type AssetMovementResponse struct {
	ID             string  `json:"id"`
	AssetID        string  `json:"assetId"`
	FromLocationID *string `json:"fromLocationId,omitempty"`
	ToLocationID   *string `json:"toLocationId,omitempty"`
	FromUserID     *string `json:"fromUserId,omitempty"`
	ToUserID       *string `json:"toUserId,omitempty"`
	MovedByID      string  `json:"movedById"`
	MovementDate   string  `json:"movementDate"`
	Notes          *string `json:"notes,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
	// * Populated
	Asset        AssetResponse     `json:"asset"`
	FromLocation *LocationResponse `json:"fromLocation,omitempty"`
	ToLocation   *LocationResponse `json:"toLocation,omitempty"`
	FromUser     *UserResponse     `json:"fromUser,omitempty"`
	ToUser       *UserResponse     `json:"toUser,omitempty"`
	MovedBy      UserResponse      `json:"movedBy"`
}

// --- Payloads ---

type CreateAssetMovementPayload struct {
	AssetID      string                                  `json:"assetId" validate:"required"`
	ToLocationID *string                                 `json:"toLocationId,omitempty" validate:"omitempty"`
	ToUserID     *string                                 `json:"toUserId,omitempty" validate:"omitempty"`
	Translations []CreateAssetMovementTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type CreateAssetMovementTranslationPayload struct {
	LangCode string `json:"langCode" validate:"required,max=5"`
	Notes    string `json:"notes" validate:"required"`
}
