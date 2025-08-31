package main

import (
	"fmt"

	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Example showing how the i18n implementation works
func main() {
	// * Test error messages in different languages
	fmt.Println("=== Error Messages ===")

	// * English (default)
	fmt.Printf("English: %s\n", utils.GetLocalizedMessage(utils.ErrUserNotFoundKey, "en-US"))

	// * Indonesian
	fmt.Printf("Indonesian: %s\n", utils.GetLocalizedMessage(utils.ErrUserNotFoundKey, "id-ID"))

	// * Japanese
	fmt.Printf("Japanese: %s\n", utils.GetLocalizedMessage(utils.ErrUserNotFoundKey, "ja-JP"))

	fmt.Println("\n=== Auth Error Messages ===")

	// * Invalid credentials in different languages
	fmt.Printf("English: %s\n", utils.GetLocalizedMessage(utils.ErrInvalidCredentialsKey, "en-US"))
	fmt.Printf("Indonesian: %s\n", utils.GetLocalizedMessage(utils.ErrInvalidCredentialsKey, "id-ID"))
	fmt.Printf("Japanese: %s\n", utils.GetLocalizedMessage(utils.ErrInvalidCredentialsKey, "ja-JP"))

	fmt.Println("\n=== Success Messages ===")

	// * English (default)
	fmt.Printf("English: %s\n", utils.GetLocalizedMessage(utils.SuccessUserCreatedKey, "en-US"))

	// * Indonesian
	fmt.Printf("Indonesian: %s\n", utils.GetLocalizedMessage(utils.SuccessUserCreatedKey, "id-ID"))

	// * Japanese
	fmt.Printf("Japanese: %s\n", utils.GetLocalizedMessage(utils.SuccessUserCreatedKey, "ja-JP"))

	fmt.Println("\n=== Auth Success Messages ===")

	// * Login success in different languages
	fmt.Printf("English: %s\n", utils.GetLocalizedMessage(utils.SuccessLoginKey, "en-US"))
	fmt.Printf("Indonesian: %s\n", utils.GetLocalizedMessage(utils.SuccessLoginKey, "id-ID"))
	fmt.Printf("Japanese: %s\n", utils.GetLocalizedMessage(utils.SuccessLoginKey, "ja-JP"))

	fmt.Println("\n=== Category Messages ===")

	// * English (default)
	fmt.Printf("English: %s\n", utils.GetLocalizedMessage(utils.SuccessCategoryRetrievedKey, "en-US"))

	// * Indonesian
	fmt.Printf("Indonesian: %s\n", utils.GetLocalizedMessage(utils.SuccessCategoryRetrievedKey, "id-ID"))

	// * Japanese
	fmt.Printf("Japanese: %s\n", utils.GetLocalizedMessage(utils.SuccessCategoryRetrievedKey, "ja-JP"))

	fmt.Println("\n=== Language Code Normalization ===")

	// * Test different language code formats
	testLangCodes := []string{"en", "id", "ja", "en-us", "id_ID", "ja_jp", "unknown"}

	for _, langCode := range testLangCodes {
		normalized := utils.GetLocalizedMessage(utils.SuccessUserCreatedKey, langCode)
		fmt.Printf("Input: %s -> Message: %s\n", langCode, normalized)
	}
}
