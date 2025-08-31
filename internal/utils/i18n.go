package utils

import (
	"fmt"
	"strings"
)

// * MessageKey represents a message key for localization
type MessageKey string

// * Error message keys
const (
	// * Common error keys
	ErrBadRequestKey   MessageKey = "error.bad_request"
	ErrUnauthorizedKey MessageKey = "error.unauthorized"
	ErrForbiddenKey    MessageKey = "error.forbidden"
	ErrNotFoundKey     MessageKey = "error.not_found"
	ErrConflictKey     MessageKey = "error.conflict"
	ErrInternalKey     MessageKey = "error.internal"
	ErrValidationKey   MessageKey = "error.validation"

	// * User-specific error keys
	ErrUserNotFoundKey      MessageKey = "error.user.not_found"
	ErrUserNameExistsKey    MessageKey = "error.user.name_exists"
	ErrUserEmailExistsKey   MessageKey = "error.user.email_exists"
	ErrUserIDRequiredKey    MessageKey = "error.user.id_required"
	ErrUserNameRequiredKey  MessageKey = "error.user.name_required"
	ErrUserEmailRequiredKey MessageKey = "error.user.email_required"

	// * Category-specific error keys
	ErrCategoryNotFoundKey     MessageKey = "error.category.not_found"
	ErrCategoryCodeExistsKey   MessageKey = "error.category.code_exists"
	ErrCategoryIDRequiredKey   MessageKey = "error.category.id_required"
	ErrCategoryCodeRequiredKey MessageKey = "error.category.code_required"
	ErrCategoryNameRequiredKey MessageKey = "error.category.name_required"

	// * Auth-specific error keys
	ErrInvalidCredentialsKey MessageKey = "error.auth.invalid_credentials"
	ErrTokenExpiredKey       MessageKey = "error.auth.token_expired"
	ErrTokenInvalidKey       MessageKey = "error.auth.token_invalid"
)

// * Success message keys
const (
	// * Common success keys
	SuccessCreatedKey   MessageKey = "success.created"
	SuccessUpdatedKey   MessageKey = "success.updated"
	SuccessDeletedKey   MessageKey = "success.deleted"
	SuccessRetrievedKey MessageKey = "success.retrieved"
	SuccessCountedKey   MessageKey = "success.counted"
	SuccessCheckedKey   MessageKey = "success.checked"

	// * User-specific success keys
	SuccessUserCreatedKey               MessageKey = "success.user.created"
	SuccessUserUpdatedKey               MessageKey = "success.user.updated"
	SuccessUserDeletedKey               MessageKey = "success.user.deleted"
	SuccessUserRetrievedKey             MessageKey = "success.user.retrieved"
	SuccessUserRetrievedByNameKey       MessageKey = "success.user.retrieved_by_name"
	SuccessUserRetrievedByEmailKey      MessageKey = "success.user.retrieved_by_email"
	SuccessUserCountedKey               MessageKey = "success.user.counted"
	SuccessUserExistenceCheckedKey      MessageKey = "success.user.existence_checked"
	SuccessUserNameExistenceCheckedKey  MessageKey = "success.user.name_existence_checked"
	SuccessUserEmailExistenceCheckedKey MessageKey = "success.user.email_existence_checked"

	// * Category-specific success keys
	SuccessCategoryCreatedKey              MessageKey = "success.category.created"
	SuccessCategoryUpdatedKey              MessageKey = "success.category.updated"
	SuccessCategoryDeletedKey              MessageKey = "success.category.deleted"
	SuccessCategoryRetrievedKey            MessageKey = "success.category.retrieved"
	SuccessCategoryRetrievedByCodeKey      MessageKey = "success.category.retrieved_by_code"
	SuccessCategoryHierarchyRetrievedKey   MessageKey = "success.category.hierarchy_retrieved"
	SuccessCategoryCountedKey              MessageKey = "success.category.counted"
	SuccessCategoryExistenceCheckedKey     MessageKey = "success.category.existence_checked"
	SuccessCategoryCodeExistenceCheckedKey MessageKey = "success.category.code_existence_checked"

	// * Auth-specific success keys
	SuccessLoginKey   MessageKey = "success.auth.login"
	SuccessLogoutKey  MessageKey = "success.auth.logout"
	SuccessRefreshKey MessageKey = "success.auth.refresh"
)

// * messageTranslations contains all message translations
var messageTranslations = map[MessageKey]map[string]string{
	// * Error messages
	ErrBadRequestKey: {
		"en-US": "Bad request",
		"id-ID": "Permintaan tidak valid",
		"ja-JP": "不正なリクエスト",
	},
	ErrUnauthorizedKey: {
		"en-US": "Unauthorized access",
		"id-ID": "Akses tidak diotorisasi",
		"ja-JP": "認証されていないアクセス",
	},
	ErrForbiddenKey: {
		"en-US": "Access forbidden",
		"id-ID": "Akses dilarang",
		"ja-JP": "アクセス禁止",
	},
	ErrNotFoundKey: {
		"en-US": "Resource not found",
		"id-ID": "Sumber daya tidak ditemukan",
		"ja-JP": "リソースが見つかりません",
	},
	ErrConflictKey: {
		"en-US": "Resource conflict",
		"id-ID": "Konflik sumber daya",
		"ja-JP": "リソースの競合",
	},
	ErrInternalKey: {
		"en-US": "An unexpected internal error occurred",
		"id-ID": "Terjadi kesalahan internal yang tidak terduga",
		"ja-JP": "予期しない内部エラーが発生しました",
	},
	ErrValidationKey: {
		"en-US": "Validation failed",
		"id-ID": "Validasi gagal",
		"ja-JP": "検証に失敗しました",
	},

	// * User-specific error messages
	ErrUserNotFoundKey: {
		"en-US": "User not found",
		"id-ID": "Pengguna tidak ditemukan",
		"ja-JP": "ユーザーが見つかりません",
	},
	ErrUserNameExistsKey: {
		"en-US": "Username already exists",
		"id-ID": "Nama pengguna sudah ada",
		"ja-JP": "ユーザー名は既に存在します",
	},
	ErrUserEmailExistsKey: {
		"en-US": "Email already exists",
		"id-ID": "Email sudah ada",
		"ja-JP": "メールアドレスは既に存在します",
	},
	ErrUserIDRequiredKey: {
		"en-US": "User ID is required",
		"id-ID": "ID pengguna diperlukan",
		"ja-JP": "ユーザーIDが必要です",
	},
	ErrUserNameRequiredKey: {
		"en-US": "Username is required",
		"id-ID": "Nama pengguna diperlukan",
		"ja-JP": "ユーザー名が必要です",
	},
	ErrUserEmailRequiredKey: {
		"en-US": "Email is required",
		"id-ID": "Email diperlukan",
		"ja-JP": "メールアドレスが必要です",
	},

	// * Category-specific error messages
	ErrCategoryNotFoundKey: {
		"en-US": "Category not found",
		"id-ID": "Kategori tidak ditemukan",
		"ja-JP": "カテゴリが見つかりません",
	},
	ErrCategoryCodeExistsKey: {
		"en-US": "Category code already exists",
		"id-ID": "Kode kategori sudah ada",
		"ja-JP": "カテゴリコードは既に存在します",
	},
	ErrCategoryIDRequiredKey: {
		"en-US": "Category ID is required",
		"id-ID": "ID kategori diperlukan",
		"ja-JP": "カテゴリIDが必要です",
	},
	ErrCategoryCodeRequiredKey: {
		"en-US": "Category code is required",
		"id-ID": "Kode kategori diperlukan",
		"ja-JP": "カテゴリコードが必要です",
	},
	ErrCategoryNameRequiredKey: {
		"en-US": "Category name is required",
		"id-ID": "Nama kategori diperlukan",
		"ja-JP": "カテゴリ名が必要です",
	},

	// * Auth-specific error messages
	ErrInvalidCredentialsKey: {
		"en-US": "Invalid credentials",
		"id-ID": "Kredensial tidak valid",
		"ja-JP": "無効な資格情報",
	},
	ErrTokenExpiredKey: {
		"en-US": "Token has expired",
		"id-ID": "Token telah kedaluwarsa",
		"ja-JP": "トークンの有効期限が切れています",
	},
	ErrTokenInvalidKey: {
		"en-US": "Invalid token",
		"id-ID": "Token tidak valid",
		"ja-JP": "無効なトークン",
	},

	// * Success messages
	SuccessCreatedKey: {
		"en-US": "Created successfully",
		"id-ID": "Berhasil dibuat",
		"ja-JP": "正常に作成されました",
	},
	SuccessUpdatedKey: {
		"en-US": "Updated successfully",
		"id-ID": "Berhasil diperbarui",
		"ja-JP": "正常に更新されました",
	},
	SuccessDeletedKey: {
		"en-US": "Deleted successfully",
		"id-ID": "Berhasil dihapus",
		"ja-JP": "正常に削除されました",
	},
	SuccessRetrievedKey: {
		"en-US": "Retrieved successfully",
		"id-ID": "Berhasil diambil",
		"ja-JP": "正常に取得されました",
	},
	SuccessCountedKey: {
		"en-US": "Counted successfully",
		"id-ID": "Berhasil dihitung",
		"ja-JP": "正常にカウントされました",
	},
	SuccessCheckedKey: {
		"en-US": "Checked successfully",
		"id-ID": "Berhasil diperiksa",
		"ja-JP": "正常にチェックされました",
	},

	// * User-specific success messages
	SuccessUserCreatedKey: {
		"en-US": "User created successfully",
		"id-ID": "Pengguna berhasil dibuat",
		"ja-JP": "ユーザーが正常に作成されました",
	},
	SuccessUserUpdatedKey: {
		"en-US": "User updated successfully",
		"id-ID": "Pengguna berhasil diperbarui",
		"ja-JP": "ユーザーが正常に更新されました",
	},
	SuccessUserDeletedKey: {
		"en-US": "User deleted successfully",
		"id-ID": "Pengguna berhasil dihapus",
		"ja-JP": "ユーザーが正常に削除されました",
	},
	SuccessUserRetrievedKey: {
		"en-US": "User retrieved successfully",
		"id-ID": "Pengguna berhasil diambil",
		"ja-JP": "ユーザーが正常に取得されました",
	},
	SuccessUserRetrievedByNameKey: {
		"en-US": "User retrieved successfully by name",
		"id-ID": "Pengguna berhasil diambil berdasarkan nama",
		"ja-JP": "名前でユーザーが正常に取得されました",
	},
	SuccessUserRetrievedByEmailKey: {
		"en-US": "User retrieved successfully by email",
		"id-ID": "Pengguna berhasil diambil berdasarkan email",
		"ja-JP": "メールでユーザーが正常に取得されました",
	},
	SuccessUserCountedKey: {
		"en-US": "Users counted successfully",
		"id-ID": "Pengguna berhasil dihitung",
		"ja-JP": "ユーザーが正常にカウントされました",
	},
	SuccessUserExistenceCheckedKey: {
		"en-US": "User existence checked successfully",
		"id-ID": "Keberadaan pengguna berhasil diperiksa",
		"ja-JP": "ユーザーの存在が正常に確認されました",
	},
	SuccessUserNameExistenceCheckedKey: {
		"en-US": "Username existence checked successfully",
		"id-ID": "Keberadaan nama pengguna berhasil diperiksa",
		"ja-JP": "ユーザー名の存在が正常に確認されました",
	},
	SuccessUserEmailExistenceCheckedKey: {
		"en-US": "Email existence checked successfully",
		"id-ID": "Keberadaan email berhasil diperiksa",
		"ja-JP": "メールの存在が正常に確認されました",
	},

	// * Category-specific success messages
	SuccessCategoryCreatedKey: {
		"en-US": "Category created successfully",
		"id-ID": "Kategori berhasil dibuat",
		"ja-JP": "カテゴリが正常に作成されました",
	},
	SuccessCategoryUpdatedKey: {
		"en-US": "Category updated successfully",
		"id-ID": "Kategori berhasil diperbarui",
		"ja-JP": "カテゴリが正常に更新されました",
	},
	SuccessCategoryDeletedKey: {
		"en-US": "Category deleted successfully",
		"id-ID": "Kategori berhasil dihapus",
		"ja-JP": "カテゴリが正常に削除されました",
	},
	SuccessCategoryRetrievedKey: {
		"en-US": "Categories retrieved successfully",
		"id-ID": "Kategori berhasil diambil",
		"ja-JP": "カテゴリが正常に取得されました",
	},
	SuccessCategoryRetrievedByCodeKey: {
		"en-US": "Category retrieved successfully by code",
		"id-ID": "Kategori berhasil diambil berdasarkan kode",
		"ja-JP": "コードでカテゴリが正常に取得されました",
	},
	SuccessCategoryHierarchyRetrievedKey: {
		"en-US": "Category hierarchy retrieved successfully",
		"id-ID": "Hierarki kategori berhasil diambil",
		"ja-JP": "カテゴリ階層が正常に取得されました",
	},
	SuccessCategoryCountedKey: {
		"en-US": "Categories counted successfully",
		"id-ID": "Kategori berhasil dihitung",
		"ja-JP": "カテゴリが正常にカウントされました",
	},
	SuccessCategoryExistenceCheckedKey: {
		"en-US": "Category existence checked successfully",
		"id-ID": "Keberadaan kategori berhasil diperiksa",
		"ja-JP": "カテゴリの存在が正常に確認されました",
	},
	SuccessCategoryCodeExistenceCheckedKey: {
		"en-US": "Category code existence checked successfully",
		"id-ID": "Keberadaan kode kategori berhasil diperiksa",
		"ja-JP": "カテゴリコードの存在が正常に確認されました",
	},

	// * Auth-specific success messages
	SuccessLoginKey: {
		"en-US": "Login successful",
		"id-ID": "Login berhasil",
		"ja-JP": "ログイン成功",
	},
	SuccessLogoutKey: {
		"en-US": "Logout successful",
		"id-ID": "Logout berhasil",
		"ja-JP": "ログアウト成功",
	},
	SuccessRefreshKey: {
		"en-US": "Token refreshed successfully",
		"id-ID": "Token berhasil diperbarui",
		"ja-JP": "トークンが正常に更新されました",
	},
}

// * GetLocalizedMessage returns the localized message for the given key and language
func GetLocalizedMessage(key MessageKey, langCode string, params ...string) string {
	translations, exists := messageTranslations[key]
	if !exists {
		return string(key) // * Return the key itself if no translation found
	}

	// * Normalize language code to match our translation keys
	normalizedLang := normalizeLanguageCode(langCode)

	message, exists := translations[normalizedLang]
	if !exists {
		// * Fallback to English if the requested language is not available
		message, exists = translations["en-US"]
		if !exists {
			return string(key) // * Return the key itself if no translation found
		}
	}

	// * Replace parameters if provided
	if len(params) > 0 {
		for i, param := range params {
			placeholder := fmt.Sprintf("{%d}", i)
			message = strings.ReplaceAll(message, placeholder, param)
		}
	}

	return message
}

// * normalizeLanguageCode normalizes language codes to match our translation keys
func normalizeLanguageCode(langCode string) string {
	langCode = strings.ToLower(langCode)

	if strings.HasPrefix(langCode, "en") {
		return "en-US"
	} else if strings.HasPrefix(langCode, "id") {
		return "id-ID"
	} else if strings.HasPrefix(langCode, "ja") {
		return "ja-JP"
	}

	// * Return the original if no match found, or default to English
	switch langCode {
	case "en-us", "en_us":
		return "en-US"
	case "id-id", "id_id":
		return "id-ID"
	case "ja-jp", "ja_jp":
		return "ja-JP"
	default:
		return "en-US" // * Default fallback
	}
}

// * GetAvailableLanguages returns the list of available language codes
func GetAvailableLanguages() []string {
	return []string{"en-US", "id-ID", "ja-JP"}
}
