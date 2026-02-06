package utils

import (
	"context"
	"log"
	"strings"

	"github.com/Rizz404/inventory-api/internal/client/gtranslate"
)

// Translation payload types (to avoid circular dependency with domain)
type CreateTranslationPayload struct {
	LangCode     string
	CategoryName string
	Description  *string
}

type UpdateTranslationPayload struct {
	LangCode     string
	CategoryName *string
	Description  *string
}

type ExistingTranslation struct {
	LangCode     string
	CategoryName string
	Description  *string
}

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

// AutoTranslateCategoryCreate automatically translates missing category translations
// Returns: translated payloads for all languages
func AutoTranslateCategoryCreate(ctx context.Context, translator *gtranslate.Client, translations []CreateTranslationPayload) ([]CreateTranslationPayload, error) {
	if len(translations) >= 3 {
		return translations, nil // All languages provided
	}

	// Map existing lang codes
	existingLangs := make(map[string]CreateTranslationPayload)
	for _, t := range translations {
		existingLangs[t.LangCode] = t
	}

	// Find missing languages
	allLangs := GetAllSupportedLangCodes()
	var missingLangs []string
	for _, lang := range allLangs {
		if _, exists := existingLangs[lang]; !exists {
			missingLangs = append(missingLangs, lang)
		}
	}

	if len(missingLangs) == 0 {
		return translations, nil
	}

	// Use first existing translation as source
	var sourceTranslation CreateTranslationPayload
	var sourceLang string
	for lang, trans := range existingLangs {
		sourceTranslation = trans
		sourceLang = NormalizeToGTranslateLang(lang)
		break
	}

	// Translate to missing languages
	result := make([]CreateTranslationPayload, 0, 3)
	result = append(result, translations...)

	for _, targetLangCode := range missingLangs {
		targetLang := NormalizeToGTranslateLang(targetLangCode)

		// Translate category name
		translatedName, err := translator.Translate(ctx, sourceTranslation.CategoryName, sourceLang, targetLang)
		if err != nil {
			log.Printf("Failed to translate category name to %s: %v", targetLangCode, err)
			continue
		}

		// Translate description if exists
		var translatedDesc *string
		if sourceTranslation.Description != nil && *sourceTranslation.Description != "" {
			desc, err := translator.Translate(ctx, *sourceTranslation.Description, sourceLang, targetLang)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		result = append(result, CreateTranslationPayload{
			LangCode:     targetLangCode,
			CategoryName: translatedName,
			Description:  translatedDesc,
		})
	}

	return result, nil
}

// AutoTranslateCategoryUpdate automatically translates missing category update translations
func AutoTranslateCategoryUpdate(ctx context.Context, translator *gtranslate.Client, translations []UpdateTranslationPayload, existingTranslations []ExistingTranslation) ([]UpdateTranslationPayload, error) {
	if len(translations) == 0 {
		return translations, nil // No updates provided
	}

	// Map existing translations from database
	existingLangs := make(map[string]ExistingTranslation)
	for _, t := range existingTranslations {
		existingLangs[t.LangCode] = t
	}

	// Map updated translations
	updatedLangs := make(map[string]UpdateTranslationPayload)
	for _, t := range translations {
		updatedLangs[t.LangCode] = t
	}

	// Find which translations actually changed
	var changedTranslation *UpdateTranslationPayload
	var changedLangCode string
	for lang, updated := range updatedLangs {
		existing, exists := existingLangs[lang]
		if !exists {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}

		// Check if name or description changed
		if (updated.CategoryName != nil && *updated.CategoryName != existing.CategoryName) ||
			(updated.Description != nil && ((existing.Description == nil && *updated.Description != "") ||
				(existing.Description != nil && *updated.Description != *existing.Description))) {
			changedTranslation = &updated
			changedLangCode = lang
			break
		}
	}

	// No changes detected, return as is
	if changedTranslation == nil {
		return translations, nil
	}

	// Find missing languages to translate
	allLangs := GetAllSupportedLangCodes()
	var missingLangs []string
	for _, lang := range allLangs {
		if _, exists := updatedLangs[lang]; !exists {
			missingLangs = append(missingLangs, lang)
		}
	}

	if len(missingLangs) == 0 {
		return translations, nil
	}

	// Use changed translation as source
	sourceLang := NormalizeToGTranslateLang(changedLangCode)
	result := make([]UpdateTranslationPayload, 0, 3)
	result = append(result, translations...)

	// Determine source text for translation
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

	for _, targetLangCode := range missingLangs {
		targetLang := NormalizeToGTranslateLang(targetLangCode)

		// Translate category name if source exists
		var translatedName *string
		if sourceName != "" {
			name, err := translator.Translate(ctx, sourceName, sourceLang, targetLang)
			if err != nil {
				log.Printf("Failed to translate category name to %s: %v", targetLangCode, err)
			} else {
				translatedName = &name
			}
		}

		// Translate description if exists
		var translatedDesc *string
		if sourceDesc != nil && *sourceDesc != "" {
			desc, err := translator.Translate(ctx, *sourceDesc, sourceLang, targetLang)
			if err != nil {
				log.Printf("Failed to translate description to %s: %v", targetLangCode, err)
			} else {
				translatedDesc = &desc
			}
		}

		result = append(result, UpdateTranslationPayload{
			LangCode:     targetLangCode,
			CategoryName: translatedName,
			Description:  translatedDesc,
		})
	}

	return result, nil
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
