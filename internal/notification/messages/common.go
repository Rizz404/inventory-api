package messages

import "github.com/Rizz404/inventory-api/internal/utils"

// NotificationMessageKey represents a notification message key
type NotificationMessageKey string

// NotificationTranslation represents a notification translation
type NotificationTranslation struct {
	LangCode string
	Title    string
	Message  string
}

// SupportedLanguages contains all supported language codes
var SupportedLanguages = []string{"en-US", "id-ID", "ja-JP"}

// GetNotificationMessage returns the localized notification message
func GetNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string, translations map[NotificationMessageKey]map[string]string) string {
	translationMap, exists := translations[key]
	if !exists {
		return string(key)
	}

	normalizedLang := normalizeLanguageCode(langCode)
	message, exists := translationMap[normalizedLang]
	if !exists {
		// Fallback to English
		message, exists = translationMap["en-US"]
		if !exists {
			return string(key)
		}
	}

	// Replace placeholders with actual values
	for placeholder, value := range params {
		message = utils.ReplaceAllCaseInsensitive(message, "{"+placeholder+"}", value)
	}

	return message
}

// GetNotificationTranslations returns all translations for a notification message key
func GetNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string, translations map[NotificationMessageKey]map[string]string) []NotificationTranslation {
	result := []NotificationTranslation{}

	for _, lang := range SupportedLanguages {
		title := GetNotificationMessage(titleKey, lang, params, translations)
		message := GetNotificationMessage(messageKey, lang, params, translations)

		result = append(result, NotificationTranslation{
			LangCode: lang,
			Title:    title,
			Message:  message,
		})
	}

	return result
}

// normalizeLanguageCode normalizes language codes to supported format
func normalizeLanguageCode(langCode string) string {
	// Convert common variations to our standard format
	switch langCode {
	case "en", "en-us", "EN", "EN-US":
		return "en-US"
	case "id", "id-id", "ID", "ID-ID":
		return "id-ID"
	case "ja", "ja-jp", "JA", "JA-JP":
		return "ja-JP"
	default:
		return langCode
	}
}
