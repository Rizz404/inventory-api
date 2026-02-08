package utils

import (
	"context"
	"log"
	"strings"

	"github.com/Rizz404/inventory-api/internal/client/gtranslate"
)

// *===========================GENERIC INTERFACES===========================*

// Translatable interface for types that can be translated
type Translatable interface {
	GetLangCode() string
	GetPrimaryText() string      // Main text to translate (e.g., CategoryName, LocationName, Title)
	GetSecondaryText() *string   // Optional text (e.g., Description)
	GetTertiaryText() *string    // Optional text (e.g., ResolutionNotes)
}

// *===========================SUPPORTED LANGUAGES===========================*

// Supported languages mapping: database lang code -> gtranslate lang code
var supportedLangs = map[string]string{
	"en-US": "en",
	"id-ID": "id",
	"ja-JP": "ja",
}

// GetAllSupportedLangCodes returns all supported database lang codes
func GetAllSupportedLangCodes() []string {
	return []string{"en-US", "id-ID", "ja-JP"}
}

// NormalizeToGTranslateLang converts database lang code (e.g., "en-US") to gtranslate format (e.g., "en")
func NormalizeToGTranslateLang(langCode string) string {
	if normalized, ok := supportedLangs[langCode]; ok {
		return normalized
	}
	// Fallback: take first part before hyphen
	parts := strings.Split(langCode, "-")
	return strings.ToLower(parts[0])
}

// GetMissingTranslationLangCodes returns language codes that are missing from translations
func GetMissingTranslationLangCodes(existingLangCodes []string) []string {
	existingMap := make(map[string]bool)
	for _, code := range existingLangCodes {
		existingMap[code] = true
	}

	allLangs := GetAllSupportedLangCodes()
	var missing []string
	for _, lang := range allLangs {
		if !existingMap[lang] {
			missing = append(missing, lang)
		}
	}
	return missing
}

// *===========================CATEGORY TRANSLATION TYPES===========================*

// Translation payload types for Category (to avoid circular dependency with domain)
type CategoryCreateTranslation struct {
	LangCode     string
	CategoryName string
	Description  *string
}

func (t CategoryCreateTranslation) GetLangCode() string     { return t.LangCode }
func (t CategoryCreateTranslation) GetPrimaryText() string  { return t.CategoryName }
func (t CategoryCreateTranslation) GetSecondaryText() *string { return t.Description }
func (t CategoryCreateTranslation) GetTertiaryText() *string  { return nil }

type CategoryUpdateTranslation struct {
	LangCode     string
	CategoryName *string
	Description  *string
}

type CategoryExistingTranslation struct {
	LangCode     string
	CategoryName string
	Description  *string
}

// *===========================LOCATION TRANSLATION TYPES===========================*

type LocationCreateTranslation struct {
	LangCode     string
	LocationName string
}

func (t LocationCreateTranslation) GetLangCode() string     { return t.LangCode }
func (t LocationCreateTranslation) GetPrimaryText() string  { return t.LocationName }
func (t LocationCreateTranslation) GetSecondaryText() *string { return nil }
func (t LocationCreateTranslation) GetTertiaryText() *string  { return nil }

type LocationUpdateTranslation struct {
	LangCode     string
	LocationName *string
}

type LocationExistingTranslation struct {
	LangCode     string
	LocationName string
}

// *===========================ISSUE REPORT TRANSLATION TYPES===========================*

type IssueReportCreateTranslation struct {
	LangCode        string
	Title           string
	Description     *string
	ResolutionNotes *string
}

func (t IssueReportCreateTranslation) GetLangCode() string       { return t.LangCode }
func (t IssueReportCreateTranslation) GetPrimaryText() string    { return t.Title }
func (t IssueReportCreateTranslation) GetSecondaryText() *string { return t.Description }
func (t IssueReportCreateTranslation) GetTertiaryText() *string  { return t.ResolutionNotes }

type IssueReportUpdateTranslation struct {
	LangCode        string
	Title           *string
	Description     *string
	ResolutionNotes *string
}

type IssueReportExistingTranslation struct {
	LangCode        string
	Title           string
	Description     *string
	ResolutionNotes *string
}

// *===========================MAINTENANCE SCHEDULE TRANSLATION TYPES===========================*

type MaintenanceScheduleCreateTranslation struct {
	LangCode    string
	Title       string
	Description *string
}

func (t MaintenanceScheduleCreateTranslation) GetLangCode() string       { return t.LangCode }
func (t MaintenanceScheduleCreateTranslation) GetPrimaryText() string    { return t.Title }
func (t MaintenanceScheduleCreateTranslation) GetSecondaryText() *string { return t.Description }
func (t MaintenanceScheduleCreateTranslation) GetTertiaryText() *string  { return nil }

type MaintenanceScheduleUpdateTranslation struct {
	LangCode    string
	Title       *string
	Description *string
}

type MaintenanceScheduleExistingTranslation struct {
	LangCode    string
	Title       string
	Description *string
}

// *===========================MAINTENANCE RECORD TRANSLATION TYPES===========================*

type MaintenanceRecordCreateTranslation struct {
	LangCode string
	Title    string
	Notes    *string
}

func (t MaintenanceRecordCreateTranslation) GetLangCode() string       { return t.LangCode }
func (t MaintenanceRecordCreateTranslation) GetPrimaryText() string    { return t.Title }
func (t MaintenanceRecordCreateTranslation) GetSecondaryText() *string { return t.Notes }
func (t MaintenanceRecordCreateTranslation) GetTertiaryText() *string  { return nil }

type MaintenanceRecordUpdateTranslation struct {
	LangCode string
	Title    *string
	Notes    *string
}

type MaintenanceRecordExistingTranslation struct {
	LangCode string
	Title    string
	Notes    *string
}

// *===========================GENERIC TRANSLATION FUNCTIONS===========================*

// TranslateText translates a single text from source to target language
func TranslateText(ctx context.Context, translator *gtranslate.Client, text, sourceLang, targetLang string) (string, error) {
	sourceNormalized := NormalizeToGTranslateLang(sourceLang)
	targetNormalized := NormalizeToGTranslateLang(targetLang)
	return translator.Translate(ctx, text, sourceNormalized, targetNormalized)
}

// *===========================CATEGORY AUTO TRANSLATION===========================*

// AutoTranslateCategoryCreate automatically translates missing category translations
func AutoTranslateCategoryCreate(ctx context.Context, translator *gtranslate.Client, translations []CategoryCreateTranslation) ([]CategoryCreateTranslation, error) {
	if len(translations) >= 3 {
		return translations, nil
	}

	existingLangs := make(map[string]CategoryCreateTranslation)
	for _, t := range translations {
		existingLangs[t.LangCode] = t
	}

	missingLangs := GetMissingTranslationLangCodes(getLanguageCodes(translations))
	if len(missingLangs) == 0 {
		return translations, nil
	}

	var sourceTranslation CategoryCreateTranslation
	var sourceLang string
	for lang, trans := range existingLangs {
		sourceTranslation = trans
		sourceLang = lang
		break
	}

	result := make([]CategoryCreateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		translatedName, err := TranslateText(ctx, translator, sourceTranslation.CategoryName, sourceLang, targetLangCode)
		if err != nil {
			log.Printf("Failed to translate category name to %s: %v", targetLangCode, err)
			continue
		}

		var translatedDesc *string
		if sourceTranslation.Description != nil && *sourceTranslation.Description != "" {
			desc, err := TranslateText(ctx, translator, *sourceTranslation.Description, sourceLang, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		result = append(result, CategoryCreateTranslation{
			LangCode:     targetLangCode,
			CategoryName: translatedName,
			Description:  translatedDesc,
		})
	}

	return result, nil
}

// AutoTranslateCategoryUpdate automatically translates missing category update translations
func AutoTranslateCategoryUpdate(ctx context.Context, translator *gtranslate.Client, translations []CategoryUpdateTranslation, existingTranslations []CategoryExistingTranslation) ([]CategoryUpdateTranslation, error) {
	if len(translations) == 0 {
		return translations, nil
	}

	existingLangs := make(map[string]CategoryExistingTranslation)
	for _, t := range existingTranslations {
		existingLangs[t.LangCode] = t
	}

	updatedLangs := make(map[string]CategoryUpdateTranslation)
	for _, t := range translations {
		updatedLangs[t.LangCode] = t
	}

	var changedTranslation *CategoryUpdateTranslation
	var changedLangCode string
	for lang, updated := range updatedLangs {
		existing, exists := existingLangs[lang]
		if !exists {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}

		if (updated.CategoryName != nil && *updated.CategoryName != existing.CategoryName) ||
			(updated.Description != nil && ((existing.Description == nil && *updated.Description != "") ||
				(existing.Description != nil && *updated.Description != *existing.Description))) {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}
	}

	if changedTranslation == nil {
		return translations, nil
	}

	var existingCodes []string
	for lang := range updatedLangs {
		existingCodes = append(existingCodes, lang)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	sourceName := ""
	if changedTranslation.CategoryName != nil {
		sourceName = *changedTranslation.CategoryName
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceName = existing.CategoryName
	}

	var sourceDesc *string
	if changedTranslation.Description != nil {
		sourceDesc = changedTranslation.Description
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceDesc = existing.Description
	}

	result := make([]CategoryUpdateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		var translatedName *string
		if sourceName != "" {
			name, err := TranslateText(ctx, translator, sourceName, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate category name to %s: %v", targetLangCode, err)
			} else {
				translatedName = &name
			}
		}

		var translatedDesc *string
		if sourceDesc != nil && *sourceDesc != "" {
			desc, err := TranslateText(ctx, translator, *sourceDesc, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		result = append(result, CategoryUpdateTranslation{
			LangCode:     targetLangCode,
			CategoryName: translatedName,
			Description:  translatedDesc,
		})
	}

	return result, nil
}

// getLanguageCodes helper to extract language codes from CategoryCreateTranslation slice
func getLanguageCodes(translations []CategoryCreateTranslation) []string {
	codes := make([]string, len(translations))
	for i, t := range translations {
		codes[i] = t.LangCode
	}
	return codes
}

// *===========================LOCATION AUTO TRANSLATION===========================*

// AutoTranslateLocationCreate automatically translates missing location translations
func AutoTranslateLocationCreate(ctx context.Context, translator *gtranslate.Client, translations []LocationCreateTranslation) ([]LocationCreateTranslation, error) {
	if len(translations) >= 3 {
		return translations, nil
	}

	existingLangs := make(map[string]LocationCreateTranslation)
	for _, t := range translations {
		existingLangs[t.LangCode] = t
	}

	var existingCodes []string
	for _, t := range translations {
		existingCodes = append(existingCodes, t.LangCode)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	var sourceTranslation LocationCreateTranslation
	var sourceLang string
	for lang, trans := range existingLangs {
		sourceTranslation = trans
		sourceLang = lang
		break
	}

	result := make([]LocationCreateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		translatedName, err := TranslateText(ctx, translator, sourceTranslation.LocationName, sourceLang, targetLangCode)
		if err != nil {
			log.Printf("Failed to translate location name to %s: %v", targetLangCode, err)
			continue
		}

		result = append(result, LocationCreateTranslation{
			LangCode:     targetLangCode,
			LocationName: translatedName,
		})
	}

	return result, nil
}

// AutoTranslateLocationUpdate automatically translates missing location update translations
func AutoTranslateLocationUpdate(ctx context.Context, translator *gtranslate.Client, translations []LocationUpdateTranslation, existingTranslations []LocationExistingTranslation) ([]LocationUpdateTranslation, error) {
	if len(translations) == 0 {
		return translations, nil
	}

	existingLangs := make(map[string]LocationExistingTranslation)
	for _, t := range existingTranslations {
		existingLangs[t.LangCode] = t
	}

	updatedLangs := make(map[string]LocationUpdateTranslation)
	for _, t := range translations {
		updatedLangs[t.LangCode] = t
	}

	var changedTranslation *LocationUpdateTranslation
	var changedLangCode string
	for lang, updated := range updatedLangs {
		existing, exists := existingLangs[lang]
		if !exists {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}

		if (updated.LocationName != nil && *updated.LocationName != existing.LocationName) {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}
	}

	if changedTranslation == nil {
		return translations, nil
	}

	var existingCodes []string
	for lang := range updatedLangs {
		existingCodes = append(existingCodes, lang)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	sourceName := ""
	if changedTranslation.LocationName != nil {
		sourceName = *changedTranslation.LocationName
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceName = existing.LocationName
	}

	result := make([]LocationUpdateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		var translatedName *string
		if sourceName != "" {
			name, err := TranslateText(ctx, translator, sourceName, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate location name to %s: %v", targetLangCode, err)
			} else {
				translatedName = &name
			}
		}

		result = append(result, LocationUpdateTranslation{
			LangCode:     targetLangCode,
			LocationName: translatedName,
		})
	}

	return result, nil
}

// *===========================ISSUE REPORT AUTO TRANSLATION===========================*

// AutoTranslateIssueReportCreate automatically translates missing issue report translations
func AutoTranslateIssueReportCreate(ctx context.Context, translator *gtranslate.Client, translations []IssueReportCreateTranslation) ([]IssueReportCreateTranslation, error) {
	if len(translations) >= 3 {
		return translations, nil
	}

	existingLangs := make(map[string]IssueReportCreateTranslation)
	for _, t := range translations {
		existingLangs[t.LangCode] = t
	}

	var existingCodes []string
	for _, t := range translations {
		existingCodes = append(existingCodes, t.LangCode)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	var sourceTranslation IssueReportCreateTranslation
	var sourceLang string
	for lang, trans := range existingLangs {
		sourceTranslation = trans
		sourceLang = lang
		break
	}

	result := make([]IssueReportCreateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		translatedTitle, err := TranslateText(ctx, translator, sourceTranslation.Title, sourceLang, targetLangCode)
		if err != nil {
			log.Printf("Failed to translate issue report title to %s: %v", targetLangCode, err)
			continue
		}

		var translatedDesc *string
		if sourceTranslation.Description != nil && *sourceTranslation.Description != "" {
			desc, err := TranslateText(ctx, translator, *sourceTranslation.Description, sourceLang, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		var translatedNotes *string
		if sourceTranslation.ResolutionNotes != nil && *sourceTranslation.ResolutionNotes != "" {
			notes, err := TranslateText(ctx, translator, *sourceTranslation.ResolutionNotes, sourceLang, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate resolution notes to %s: %v", targetLangCode, err)
			} else {
				translatedNotes = &notes
			}
		}

		result = append(result, IssueReportCreateTranslation{
			LangCode:        targetLangCode,
			Title:           translatedTitle,
			Description:     translatedDesc,
			ResolutionNotes: translatedNotes,
		})
	}

	return result, nil
}

// AutoTranslateIssueReportUpdate automatically translates missing issue report update translations
func AutoTranslateIssueReportUpdate(ctx context.Context, translator *gtranslate.Client, translations []IssueReportUpdateTranslation, existingTranslations []IssueReportExistingTranslation) ([]IssueReportUpdateTranslation, error) {
	if len(translations) == 0 {
		return translations, nil
	}

	existingLangs := make(map[string]IssueReportExistingTranslation)
	for _, t := range existingTranslations {
		existingLangs[t.LangCode] = t
	}

	updatedLangs := make(map[string]IssueReportUpdateTranslation)
	for _, t := range translations {
		updatedLangs[t.LangCode] = t
	}

	var changedTranslation *IssueReportUpdateTranslation
	var changedLangCode string
	for lang, updated := range updatedLangs {
		existing, exists := existingLangs[lang]
		if !exists {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}

		if (updated.Title != nil && *updated.Title != existing.Title) ||
			(updated.Description != nil && ((existing.Description == nil && *updated.Description != "") ||
				(existing.Description != nil && *updated.Description != *existing.Description))) ||
			(updated.ResolutionNotes != nil && ((existing.ResolutionNotes == nil && *updated.ResolutionNotes != "") ||
				(existing.ResolutionNotes != nil && *updated.ResolutionNotes != *existing.ResolutionNotes))) {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}
	}

	if changedTranslation == nil {
		return translations, nil
	}

	var existingCodes []string
	for lang := range updatedLangs {
		existingCodes = append(existingCodes, lang)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	sourceTitle := ""
	if changedTranslation.Title != nil {
		sourceTitle = *changedTranslation.Title
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceTitle = existing.Title
	}

	var sourceDesc *string
	if changedTranslation.Description != nil {
		sourceDesc = changedTranslation.Description
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceDesc = existing.Description
	}

	var sourceNotes *string
	if changedTranslation.ResolutionNotes != nil {
		sourceNotes = changedTranslation.ResolutionNotes
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceNotes = existing.ResolutionNotes
	}

	result := make([]IssueReportUpdateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		var translatedTitle *string
		if sourceTitle != "" {
			title, err := TranslateText(ctx, translator, sourceTitle, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate issue report title to %s: %v", targetLangCode, err)
			} else {
				translatedTitle = &title
			}
		}

		var translatedDesc *string
		if sourceDesc != nil && *sourceDesc != "" {
			desc, err := TranslateText(ctx, translator, *sourceDesc, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		var translatedNotes *string
		if sourceNotes != nil && *sourceNotes != "" {
			notes, err := TranslateText(ctx, translator, *sourceNotes, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate resolution notes to %s: %v", targetLangCode, err)
			} else {
				translatedNotes = &notes
			}
		}

		result = append(result, IssueReportUpdateTranslation{
			LangCode:        targetLangCode,
			Title:           translatedTitle,
			Description:     translatedDesc,
			ResolutionNotes: translatedNotes,
		})
	}

	return result, nil
}

// *===========================MAINTENANCE SCHEDULE AUTO TRANSLATION===========================*

// AutoTranslateMaintenanceScheduleCreate automatically translates missing maintenance schedule translations
func AutoTranslateMaintenanceScheduleCreate(ctx context.Context, translator *gtranslate.Client, translations []MaintenanceScheduleCreateTranslation) ([]MaintenanceScheduleCreateTranslation, error) {
	if len(translations) >= 3 {
		return translations, nil
	}

	existingLangs := make(map[string]MaintenanceScheduleCreateTranslation)
	for _, t := range translations {
		existingLangs[t.LangCode] = t
	}

	var existingCodes []string
	for _, t := range translations {
		existingCodes = append(existingCodes, t.LangCode)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	var sourceTranslation MaintenanceScheduleCreateTranslation
	var sourceLang string
	for lang, trans := range existingLangs {
		sourceTranslation = trans
		sourceLang = lang
		break
	}

	result := make([]MaintenanceScheduleCreateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		translatedTitle, err := TranslateText(ctx, translator, sourceTranslation.Title, sourceLang, targetLangCode)
		if err != nil {
			log.Printf("Failed to translate maintenance schedule title to %s: %v", targetLangCode, err)
			continue
		}

		var translatedDesc *string
		if sourceTranslation.Description != nil && *sourceTranslation.Description != "" {
			desc, err := TranslateText(ctx, translator, *sourceTranslation.Description, sourceLang, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		result = append(result, MaintenanceScheduleCreateTranslation{
			LangCode:    targetLangCode,
			Title:       translatedTitle,
			Description: translatedDesc,
		})
	}

	return result, nil
}

// AutoTranslateMaintenanceScheduleUpdate automatically translates missing maintenance schedule update translations
func AutoTranslateMaintenanceScheduleUpdate(ctx context.Context, translator *gtranslate.Client, translations []MaintenanceScheduleUpdateTranslation, existingTranslations []MaintenanceScheduleExistingTranslation) ([]MaintenanceScheduleUpdateTranslation, error) {
	if len(translations) == 0 {
		return translations, nil
	}

	existingLangs := make(map[string]MaintenanceScheduleExistingTranslation)
	for _, t := range existingTranslations {
		existingLangs[t.LangCode] = t
	}

	updatedLangs := make(map[string]MaintenanceScheduleUpdateTranslation)
	for _, t := range translations {
		updatedLangs[t.LangCode] = t
	}

	var changedTranslation *MaintenanceScheduleUpdateTranslation
	var changedLangCode string
	for lang, updated := range updatedLangs {
		existing, exists := existingLangs[lang]
		if !exists {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}

		if (updated.Title != nil && *updated.Title != existing.Title) ||
			(updated.Description != nil && ((existing.Description == nil && *updated.Description != "") ||
				(existing.Description != nil && *updated.Description != *existing.Description))) {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}
	}

	if changedTranslation == nil {
		return translations, nil
	}

	var existingCodes []string
	for lang := range updatedLangs {
		existingCodes = append(existingCodes, lang)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	sourceTitle := ""
	if changedTranslation.Title != nil {
		sourceTitle = *changedTranslation.Title
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceTitle = existing.Title
	}

	var sourceDesc *string
	if changedTranslation.Description != nil {
		sourceDesc = changedTranslation.Description
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceDesc = existing.Description
	}

	result := make([]MaintenanceScheduleUpdateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		var translatedTitle *string
		if sourceTitle != "" {
			title, err := TranslateText(ctx, translator, sourceTitle, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate maintenance schedule title to %s: %v", targetLangCode, err)
			} else {
				translatedTitle = &title
			}
		}

		var translatedDesc *string
		if sourceDesc != nil && *sourceDesc != "" {
			desc, err := TranslateText(ctx, translator, *sourceDesc, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		result = append(result, MaintenanceScheduleUpdateTranslation{
			LangCode:    targetLangCode,
			Title:       translatedTitle,
			Description: translatedDesc,
		})
	}

	return result, nil
}

// *===========================MAINTENANCE RECORD AUTO TRANSLATION===========================*

// AutoTranslateMaintenanceRecordCreate automatically translates missing maintenance record translations
func AutoTranslateMaintenanceRecordCreate(ctx context.Context, translator *gtranslate.Client, translations []MaintenanceRecordCreateTranslation) ([]MaintenanceRecordCreateTranslation, error) {
	if len(translations) >= 3 {
		return translations, nil
	}

	existingLangs := make(map[string]MaintenanceRecordCreateTranslation)
	for _, t := range translations {
		existingLangs[t.LangCode] = t
	}

	var existingCodes []string
	for _, t := range translations {
		existingCodes = append(existingCodes, t.LangCode)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	var sourceTranslation MaintenanceRecordCreateTranslation
	var sourceLang string
	for lang, trans := range existingLangs {
		sourceTranslation = trans
		sourceLang = lang
		break
	}

	result := make([]MaintenanceRecordCreateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		translatedTitle, err := TranslateText(ctx, translator, sourceTranslation.Title, sourceLang, targetLangCode)
		if err != nil {
			log.Printf("Failed to translate maintenance record title to %s: %v", targetLangCode, err)
			continue
		}

		var translatedNotes *string
		if sourceTranslation.Notes != nil && *sourceTranslation.Notes != "" {
			notes, err := TranslateText(ctx, translator, *sourceTranslation.Notes, sourceLang, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate notes to %s: %v", targetLangCode, err)
			} else {
				translatedNotes = &notes
			}
		}

		result = append(result, MaintenanceRecordCreateTranslation{
			LangCode: targetLangCode,
			Title:    translatedTitle,
			Notes:    translatedNotes,
		})
	}

	return result, nil
}

// AutoTranslateMaintenanceRecordUpdate automatically translates missing maintenance record update translations
func AutoTranslateMaintenanceRecordUpdate(ctx context.Context, translator *gtranslate.Client, translations []MaintenanceRecordUpdateTranslation, existingTranslations []MaintenanceRecordExistingTranslation) ([]MaintenanceRecordUpdateTranslation, error) {
	if len(translations) == 0 {
		return translations, nil
	}

	existingLangs := make(map[string]MaintenanceRecordExistingTranslation)
	for _, t := range existingTranslations {
		existingLangs[t.LangCode] = t
	}

	updatedLangs := make(map[string]MaintenanceRecordUpdateTranslation)
	for _, t := range translations {
		updatedLangs[t.LangCode] = t
	}

	var changedTranslation *MaintenanceRecordUpdateTranslation
	var changedLangCode string
	for lang, updated := range updatedLangs {
		existing, exists := existingLangs[lang]
		if !exists {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}

		if (updated.Title != nil && *updated.Title != existing.Title) ||
			(updated.Notes != nil && ((existing.Notes == nil && *updated.Notes != "") ||
				(existing.Notes != nil && *updated.Notes != *existing.Notes))) {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}
	}

	if changedTranslation == nil {
		return translations, nil
	}

	var existingCodes []string
	for lang := range updatedLangs {
		existingCodes = append(existingCodes, lang)
	}
	missingLangs := GetMissingTranslationLangCodes(existingCodes)
	if len(missingLangs) == 0 {
		return translations, nil
	}

	sourceTitle := ""
	if changedTranslation.Title != nil {
		sourceTitle = *changedTranslation.Title
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceTitle = existing.Title
	}

	var sourceNotes *string
	if changedTranslation.Notes != nil {
		sourceNotes = changedTranslation.Notes
	} else if existing, exists := existingLangs[changedLangCode]; exists {
		sourceNotes = existing.Notes
	}

	result := make([]MaintenanceRecordUpdateTranslation, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		var translatedTitle *string
		if sourceTitle != "" {
			title, err := TranslateText(ctx, translator, sourceTitle, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate maintenance record title to %s: %v", targetLangCode, err)
			} else {
				translatedTitle = &title
			}
		}

		var translatedNotes *string
		if sourceNotes != nil && *sourceNotes != "" {
			notes, err := TranslateText(ctx, translator, *sourceNotes, changedLangCode, targetLangCode)
			if err != nil {
				log.Printf("Failed to translate notes to %s: %v", targetLangCode, err)
			} else {
				translatedNotes = &notes
			}
		}

		result = append(result, MaintenanceRecordUpdateTranslation{
			LangCode: targetLangCode,
			Title:    translatedTitle,
			Notes:    translatedNotes,
		})
	}

	return result, nil
}
