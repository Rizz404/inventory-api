package gtranslate

import (
	"context"

	"github.com/bregydoc/gtranslate"
	"golang.org/x/text/language"
)

// Client wraps the gtranslate functionality
type Client struct{}

// NewClient creates a new gtranslate client
func NewClient() *Client {
	return &Client{}
}

// langCodeToTag converts language code string to language.Tag
func langCodeToTag(code string) language.Tag {
	switch code {
	case "en":
		return language.English
	case "id":
		return language.Indonesian
	case "ja":
		return language.Japanese
	case "zh":
		return language.Chinese
	case "es":
		return language.Spanish
	case "fr":
		return language.French
	case "de":
		return language.German
	default:
		tag, err := language.Parse(code)
		if err != nil {
			return language.English // default fallback
		}
		return tag
	}
}

// Translate translates text from source language to target language
// sourceLang: language code (e.g., "en", "id", "ja")
// targetLang: language code (e.g., "en", "id", "ja")
func (c *Client) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	sourceTag := langCodeToTag(sourceLang)
	targetTag := langCodeToTag(targetLang)

	translated, err := gtranslate.Translate(text, sourceTag, targetTag)
	if err != nil {
		return "", err
	}
	return translated, nil
}

// TranslateAuto translates text with auto-detection of source language
func (c *Client) TranslateAuto(ctx context.Context, text, targetLang string) (string, error) {
	targetTag := langCodeToTag(targetLang)

	translated, err := gtranslate.Translate(text, language.Und, targetTag)
	if err != nil {
		return "", err
	}
	return translated, nil
}
