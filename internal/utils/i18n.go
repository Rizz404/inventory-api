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

	// * Location-specific error keys
	ErrLocationNotFoundKey     MessageKey = "error.location.not_found"
	ErrLocationCodeExistsKey   MessageKey = "error.location.code_exists"
	ErrLocationIDRequiredKey   MessageKey = "error.location.id_required"
	ErrLocationCodeRequiredKey MessageKey = "error.location.code_required"
	ErrLocationNameRequiredKey MessageKey = "error.location.name_required"

	// * Asset-specific error keys
	ErrAssetNotFoundKey             MessageKey = "error.asset.not_found"
	ErrAssetTagExistsKey            MessageKey = "error.asset.tag_exists"
	ErrAssetDataMatrixExistsKey     MessageKey = "error.asset.datamatrix_exists"
	ErrAssetSerialNumberExistsKey   MessageKey = "error.asset.serial_number_exists"
	ErrAssetIDRequiredKey           MessageKey = "error.asset.id_required"
	ErrAssetTagRequiredKey          MessageKey = "error.asset.tag_required"
	ErrAssetDataMatrixRequiredKey   MessageKey = "error.asset.datamatrix_required"
	ErrAssetSerialNumberRequiredKey MessageKey = "error.asset.serial_number_required"

	// * Scan log-specific error keys
	ErrScanLogNotFoundKey   MessageKey = "error.scan_log.not_found"
	ErrScanLogIDRequiredKey MessageKey = "error.scan_log.id_required"

	// * Notification-specific error keys
	ErrNotificationNotFoundKey        MessageKey = "error.notification.not_found"
	ErrNotificationIDRequiredKey      MessageKey = "error.notification.id_required"
	ErrNotificationUserIDRequiredKey  MessageKey = "error.notification.user_id_required"
	ErrNotificationTypeRequiredKey    MessageKey = "error.notification.type_required"
	ErrNotificationTitleRequiredKey   MessageKey = "error.notification.title_required"
	ErrNotificationMessageRequiredKey MessageKey = "error.notification.message_required"

	// * Issue report-specific error keys
	ErrIssueReportNotFoundKey         MessageKey = "error.issue_report.not_found"
	ErrIssueReportIDRequiredKey       MessageKey = "error.issue_report.id_required"
	ErrIssueReportAssetIDRequiredKey  MessageKey = "error.issue_report.asset_id_required"
	ErrIssueReportTypeRequiredKey     MessageKey = "error.issue_report.type_required"
	ErrIssueReportPriorityRequiredKey MessageKey = "error.issue_report.priority_required"
	ErrIssueReportTitleRequiredKey    MessageKey = "error.issue_report.title_required"
	ErrIssueReportAlreadyResolvedKey  MessageKey = "error.issue_report.already_resolved"
	ErrIssueReportCannotReopenKey     MessageKey = "error.issue_report.cannot_reopen"

	// * Asset movement-specific error keys
	ErrAssetMovementNotFoundKey        MessageKey = "error.asset_movement.not_found"
	ErrAssetMovementIDRequiredKey      MessageKey = "error.asset_movement.id_required"
	ErrAssetMovementAssetIDRequiredKey MessageKey = "error.asset_movement.asset_id_required"
	ErrAssetMovementInvalidLocationKey MessageKey = "error.asset_movement.invalid_location"
	ErrAssetMovementInvalidUserKey     MessageKey = "error.asset_movement.invalid_user"
	ErrAssetMovementNoChangeKey        MessageKey = "error.asset_movement.no_change"
	ErrAssetMovementSameLocationKey    MessageKey = "error.asset_movement.same_location"

	// * Maintenance-specific error keys
	ErrMaintenanceScheduleNotFoundKey      MessageKey = "error.maintenance.schedule_not_found"
	ErrMaintenanceRecordNotFoundKey        MessageKey = "error.maintenance.record_not_found"
	ErrMaintenanceScheduleIDRequiredKey    MessageKey = "error.maintenance.schedule_id_required"
	ErrMaintenanceRecordIDRequiredKey      MessageKey = "error.maintenance.record_id_required"
	ErrMaintenanceAssetIDRequiredKey       MessageKey = "error.maintenance.asset_id_required"
	ErrMaintenanceScheduleDateRequiredKey  MessageKey = "error.maintenance.schedule_date_required"
	ErrMaintenanceRecordDateRequiredKey    MessageKey = "error.maintenance.record_date_required"
	ErrMaintenanceScheduleTitleRequiredKey MessageKey = "error.maintenance.schedule_title_required"
	ErrMaintenanceRecordTitleRequiredKey   MessageKey = "error.maintenance.record_title_required"

	// * Auth-specific error keys
	ErrInvalidCredentialsKey MessageKey = "error.auth.invalid_credentials"
	ErrTokenExpiredKey       MessageKey = "error.auth.token_expired"
	ErrTokenInvalidKey       MessageKey = "error.auth.token_invalid"

	// * File upload error keys
	ErrFileRequiredKey       MessageKey = "error.file.required"
	ErrFileTypeNotAllowedKey MessageKey = "error.file.type_not_allowed"
	ErrFileSizeTooLargeKey   MessageKey = "error.file.size_too_large"
	ErrTooManyFilesKey       MessageKey = "error.file.too_many_files"
	ErrFileUploadFailedKey   MessageKey = "error.file.upload_failed"
	ErrFileDeleteFailedKey   MessageKey = "error.file.delete_failed"
	ErrCloudinaryConfigKey   MessageKey = "error.file.cloudinary_config"
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
	SuccessUserStatisticsRetrievedKey   MessageKey = "success.user.statistics_retrieved"
	SuccessUserExistenceCheckedKey      MessageKey = "success.user.existence_checked"
	SuccessUserNameExistenceCheckedKey  MessageKey = "success.user.name_existence_checked"
	SuccessUserEmailExistenceCheckedKey MessageKey = "success.user.email_existence_checked"

	// * Category-specific success keys
	SuccessCategoryCreatedKey              MessageKey = "success.category.created"
	SuccessCategoryUpdatedKey              MessageKey = "success.category.updated"
	SuccessCategoryDeletedKey              MessageKey = "success.category.deleted"
	SuccessCategoryRetrievedKey            MessageKey = "success.category.retrieved"
	SuccessCategoryRetrievedByCodeKey      MessageKey = "success.category.retrieved_by_code"
	SuccessCategoryCountedKey              MessageKey = "success.category.counted"
	SuccessCategoryStatisticsRetrievedKey  MessageKey = "success.category.statistics_retrieved"
	SuccessCategoryExistenceCheckedKey     MessageKey = "success.category.existence_checked"
	SuccessCategoryCodeExistenceCheckedKey MessageKey = "success.category.code_existence_checked"

	// * Location-specific success keys
	SuccessLocationCreatedKey              MessageKey = "success.location.created"
	SuccessLocationUpdatedKey              MessageKey = "success.location.updated"
	SuccessLocationDeletedKey              MessageKey = "success.location.deleted"
	SuccessLocationRetrievedKey            MessageKey = "success.location.retrieved"
	SuccessLocationRetrievedByCodeKey      MessageKey = "success.location.retrieved_by_code"
	SuccessLocationCountedKey              MessageKey = "success.location.counted"
	SuccessLocationStatisticsRetrievedKey  MessageKey = "success.location.statistics_retrieved"
	SuccessLocationExistenceCheckedKey     MessageKey = "success.location.existence_checked"
	SuccessLocationCodeExistenceCheckedKey MessageKey = "success.location.code_existence_checked"

	// * Asset-specific success keys
	SuccessAssetCreatedKey                      MessageKey = "success.asset.created"
	SuccessAssetUpdatedKey                      MessageKey = "success.asset.updated"
	SuccessAssetDeletedKey                      MessageKey = "success.asset.deleted"
	SuccessAssetRetrievedKey                    MessageKey = "success.asset.retrieved"
	SuccessAssetRetrievedByTagKey               MessageKey = "success.asset.retrieved_by_tag"
	SuccessAssetRetrievedByDataMatrixKey        MessageKey = "success.asset.retrieved_by_datamatrix"
	SuccessAssetCountedKey                      MessageKey = "success.asset.counted"
	SuccessAssetStatisticsRetrievedKey          MessageKey = "success.asset.statistics_retrieved"
	SuccessAssetExistenceCheckedKey             MessageKey = "success.asset.existence_checked"
	SuccessAssetTagExistenceCheckedKey          MessageKey = "success.asset.tag_existence_checked"
	SuccessAssetDataMatrixExistenceCheckedKey   MessageKey = "success.asset.datamatrix_existence_checked"
	SuccessAssetSerialNumberExistenceCheckedKey MessageKey = "success.asset.serial_number_existence_checked"

	// * Scan log-specific success keys
	SuccessScanLogCreatedKey             MessageKey = "success.scan_log.created"
	SuccessScanLogDeletedKey             MessageKey = "success.scan_log.deleted"
	SuccessScanLogRetrievedKey           MessageKey = "success.scan_log.retrieved"
	SuccessScanLogCountedKey             MessageKey = "success.scan_log.counted"
	SuccessScanLogStatisticsRetrievedKey MessageKey = "success.scan_log.statistics_retrieved"
	SuccessScanLogExistenceCheckedKey    MessageKey = "success.scan_log.existence_checked"

	// * Notification-specific success keys
	SuccessNotificationCreatedKey             MessageKey = "success.notification.created"
	SuccessNotificationUpdatedKey             MessageKey = "success.notification.updated"
	SuccessNotificationDeletedKey             MessageKey = "success.notification.deleted"
	SuccessNotificationRetrievedKey           MessageKey = "success.notification.retrieved"
	SuccessNotificationCountedKey             MessageKey = "success.notification.counted"
	SuccessNotificationStatisticsRetrievedKey MessageKey = "success.notification.statistics_retrieved"
	SuccessNotificationExistenceCheckedKey    MessageKey = "success.notification.existence_checked"
	SuccessNotificationMarkedAsReadKey        MessageKey = "success.notification.marked_as_read"
	SuccessNotificationMarkedAsUnreadKey      MessageKey = "success.notification.marked_as_unread"

	// * Issue report-specific success keys
	SuccessIssueReportCreatedKey             MessageKey = "success.issue_report.created"
	SuccessIssueReportUpdatedKey             MessageKey = "success.issue_report.updated"
	SuccessIssueReportDeletedKey             MessageKey = "success.issue_report.deleted"
	SuccessIssueReportRetrievedKey           MessageKey = "success.issue_report.retrieved"
	SuccessIssueReportCountedKey             MessageKey = "success.issue_report.counted"
	SuccessIssueReportStatisticsRetrievedKey MessageKey = "success.issue_report.statistics_retrieved"
	SuccessIssueReportExistenceCheckedKey    MessageKey = "success.issue_report.existence_checked"
	SuccessIssueReportResolvedKey            MessageKey = "success.issue_report.resolved"
	SuccessIssueReportReopenedKey            MessageKey = "success.issue_report.reopened"

	// * Asset movement-specific success keys
	SuccessAssetMovementCreatedKey             MessageKey = "success.asset_movement.created"
	SuccessAssetMovementUpdatedKey             MessageKey = "success.asset_movement.updated"
	SuccessAssetMovementDeletedKey             MessageKey = "success.asset_movement.deleted"
	SuccessAssetMovementRetrievedKey           MessageKey = "success.asset_movement.retrieved"
	SuccessAssetMovementCountedKey             MessageKey = "success.asset_movement.counted"
	SuccessAssetMovementStatisticsRetrievedKey MessageKey = "success.asset_movement.statistics_retrieved"
	SuccessAssetMovementExistenceCheckedKey    MessageKey = "success.asset_movement.existence_checked"

	// * Maintenance-specific success keys
	SuccessMaintenanceScheduleCreatedKey             MessageKey = "success.maintenance.schedule_created"
	SuccessMaintenanceScheduleUpdatedKey             MessageKey = "success.maintenance.schedule_updated"
	SuccessMaintenanceScheduleDeletedKey             MessageKey = "success.maintenance.schedule_deleted"
	SuccessMaintenanceScheduleRetrievedKey           MessageKey = "success.maintenance.schedule_retrieved"
	SuccessMaintenanceScheduleCountedKey             MessageKey = "success.maintenance.schedule_counted"
	SuccessMaintenanceScheduleStatisticsRetrievedKey MessageKey = "success.maintenance.schedule_statistics_retrieved"
	SuccessMaintenanceRecordCreatedKey               MessageKey = "success.maintenance.record_created"
	SuccessMaintenanceRecordUpdatedKey               MessageKey = "success.maintenance.record_updated"
	SuccessMaintenanceRecordDeletedKey               MessageKey = "success.maintenance.record_deleted"
	SuccessMaintenanceRecordRetrievedKey             MessageKey = "success.maintenance.record_retrieved"
	SuccessMaintenanceRecordCountedKey               MessageKey = "success.maintenance.record_counted"
	SuccessMaintenanceRecordStatisticsRetrievedKey   MessageKey = "success.maintenance.record_statistics_retrieved"

	// * Auth-specific success keys
	SuccessLoginKey          MessageKey = "success.auth.login"
	SuccessLogoutKey         MessageKey = "success.auth.logout"
	SuccessRefreshKey        MessageKey = "success.auth.refresh"
	SuccessTokenRefreshedKey MessageKey = "success.auth.token_refreshed"

	// * File upload success keys
	SuccessFileUploadedKey          MessageKey = "success.file.uploaded"
	SuccessAvatarUploadedKey        MessageKey = "success.file.avatar_uploaded"
	SuccessFileDeletedKey           MessageKey = "success.file.deleted"
	SuccessMultipleFilesUploadedKey MessageKey = "success.file.multiple_uploaded"
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
		"en-US": "Name already exists",
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
		"en-US": "Name is required",
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

	// * Location-specific error messages
	ErrLocationNotFoundKey: {
		"en-US": "Location not found",
		"id-ID": "Lokasi tidak ditemukan",
		"ja-JP": "ロケーションが見つかりません",
	},
	ErrLocationCodeExistsKey: {
		"en-US": "Location code already exists",
		"id-ID": "Kode lokasi sudah ada",
		"ja-JP": "ロケーションコードは既に存在します",
	},
	ErrLocationIDRequiredKey: {
		"en-US": "Location ID is required",
		"id-ID": "ID lokasi diperlukan",
		"ja-JP": "ロケーションIDが必要です",
	},
	ErrLocationCodeRequiredKey: {
		"en-US": "Location code is required",
		"id-ID": "Kode lokasi diperlukan",
		"ja-JP": "ロケーションコードが必要です",
	},
	ErrLocationNameRequiredKey: {
		"en-US": "Location name is required",
		"id-ID": "Nama lokasi diperlukan",
		"ja-JP": "ロケーション名が必要です",
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

	// * File upload error messages
	ErrFileRequiredKey: {
		"en-US": "File is required",
		"id-ID": "File diperlukan",
		"ja-JP": "ファイルが必要です",
	},
	ErrFileTypeNotAllowedKey: {
		"en-US": "File type not allowed",
		"id-ID": "Tipe file tidak diizinkan",
		"ja-JP": "ファイルタイプは許可されていません",
	},
	ErrFileSizeTooLargeKey: {
		"en-US": "File size too large",
		"id-ID": "Ukuran file terlalu besar",
		"ja-JP": "ファイルサイズが大きすぎます",
	},
	ErrTooManyFilesKey: {
		"en-US": "Too many files",
		"id-ID": "Terlalu banyak file",
		"ja-JP": "ファイルが多すぎます",
	},
	ErrFileUploadFailedKey: {
		"en-US": "File upload failed",
		"id-ID": "Unggah file gagal",
		"ja-JP": "ファイルアップロードに失敗しました",
	},
	ErrFileDeleteFailedKey: {
		"en-US": "File delete failed",
		"id-ID": "Hapus file gagal",
		"ja-JP": "ファイル削除に失敗しました",
	},
	ErrCloudinaryConfigKey: {
		"en-US": "Cloudinary configuration error",
		"id-ID": "Kesalahan konfigurasi Cloudinary",
		"ja-JP": "Cloudinary設定エラー",
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
		"en-US": "Name existence checked successfully",
		"id-ID": "Keberadaan nama pengguna berhasil diperiksa",
		"ja-JP": "ユーザー名の存在が正常に確認されました",
	},
	SuccessUserEmailExistenceCheckedKey: {
		"en-US": "Email existence checked successfully",
		"id-ID": "Keberadaan email berhasil diperiksa",
		"ja-JP": "メールの存在が正常に確認されました",
	},
	SuccessUserStatisticsRetrievedKey: {
		"en-US": "User statistics retrieved successfully",
		"id-ID": "Statistik pengguna berhasil diambil",
		"ja-JP": "ユーザー統計が正常に取得されました",
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
	SuccessCategoryStatisticsRetrievedKey: {
		"en-US": "Category statistics retrieved successfully",
		"id-ID": "Statistik kategori berhasil diambil",
		"ja-JP": "カテゴリ統計が正常に取得されました",
	},

	// * Location-specific success messages
	SuccessLocationCreatedKey: {
		"en-US": "Location created successfully",
		"id-ID": "Lokasi berhasil dibuat",
		"ja-JP": "ロケーションが正常に作成されました",
	},
	SuccessLocationUpdatedKey: {
		"en-US": "Location updated successfully",
		"id-ID": "Lokasi berhasil diperbarui",
		"ja-JP": "ロケーションが正常に更新されました",
	},
	SuccessLocationDeletedKey: {
		"en-US": "Location deleted successfully",
		"id-ID": "Lokasi berhasil dihapus",
		"ja-JP": "ロケーションが正常に削除されました",
	},
	SuccessLocationRetrievedKey: {
		"en-US": "Locations retrieved successfully",
		"id-ID": "Lokasi berhasil diambil",
		"ja-JP": "ロケーションが正常に取得されました",
	},
	SuccessLocationRetrievedByCodeKey: {
		"en-US": "Location retrieved successfully by code",
		"id-ID": "Lokasi berhasil diambil berdasarkan kode",
		"ja-JP": "コードでロケーションが正常に取得されました",
	},
	SuccessLocationCountedKey: {
		"en-US": "Locations counted successfully",
		"id-ID": "Lokasi berhasil dihitung",
		"ja-JP": "ロケーションが正常にカウントされました",
	},
	SuccessLocationExistenceCheckedKey: {
		"en-US": "Location existence checked successfully",
		"id-ID": "Keberadaan lokasi berhasil diperiksa",
		"ja-JP": "ロケーションの存在が正常に確認されました",
	},
	SuccessLocationCodeExistenceCheckedKey: {
		"en-US": "Location code existence checked successfully",
		"id-ID": "Keberadaan kode lokasi berhasil diperiksa",
		"ja-JP": "ロケーションコードの存在が正常に確認されました",
	},
	SuccessLocationStatisticsRetrievedKey: {
		"en-US": "Location statistics retrieved successfully",
		"id-ID": "Statistik lokasi berhasil diambil",
		"ja-JP": "ロケーション統計が正常に取得されました",
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
	SuccessTokenRefreshedKey: {
		"en-US": "Token refreshed successfully",
		"id-ID": "Token berhasil diperbarui",
		"ja-JP": "トークンが正常に更新されました",
	},

	// * File upload success messages
	SuccessFileUploadedKey: {
		"en-US": "File uploaded successfully",
		"id-ID": "File berhasil diunggah",
		"ja-JP": "ファイルが正常にアップロードされました",
	},
	SuccessAvatarUploadedKey: {
		"en-US": "Avatar uploaded successfully",
		"id-ID": "Avatar berhasil diunggah",
		"ja-JP": "アバターが正常にアップロードされました",
	},
	SuccessFileDeletedKey: {
		"en-US": "File deleted successfully",
		"id-ID": "File berhasil dihapus",
		"ja-JP": "ファイルが正常に削除されました",
	},
	SuccessMultipleFilesUploadedKey: {
		"en-US": "Multiple files uploaded successfully",
		"id-ID": "Beberapa file berhasil diunggah",
		"ja-JP": "複数のファイルが正常にアップロードされました",
	},

	// * Notification error messages
	ErrNotificationNotFoundKey: {
		"en-US": "Notification not found",
		"id-ID": "Notifikasi tidak ditemukan",
		"ja-JP": "通知が見つかりません",
	},
	ErrNotificationIDRequiredKey: {
		"en-US": "Notification ID is required",
		"id-ID": "ID notifikasi diperlukan",
		"ja-JP": "通知IDが必要です",
	},
	ErrNotificationUserIDRequiredKey: {
		"en-US": "User ID is required",
		"id-ID": "ID pengguna diperlukan",
		"ja-JP": "ユーザーIDが必要です",
	},
	ErrNotificationTypeRequiredKey: {
		"en-US": "Notification type is required",
		"id-ID": "Jenis notifikasi diperlukan",
		"ja-JP": "通知タイプが必要です",
	},
	ErrNotificationTitleRequiredKey: {
		"en-US": "Notification title is required",
		"id-ID": "Judul notifikasi diperlukan",
		"ja-JP": "通知タイトルが必要です",
	},
	ErrNotificationMessageRequiredKey: {
		"en-US": "Notification message is required",
		"id-ID": "Pesan notifikasi diperlukan",
		"ja-JP": "通知メッセージが必要です",
	},

	// * Notification success messages
	SuccessNotificationCreatedKey: {
		"en-US": "Notification created successfully",
		"id-ID": "Notifikasi berhasil dibuat",
		"ja-JP": "通知が正常に作成されました",
	},
	SuccessNotificationUpdatedKey: {
		"en-US": "Notification updated successfully",
		"id-ID": "Notifikasi berhasil diperbarui",
		"ja-JP": "通知が正常に更新されました",
	},
	SuccessNotificationDeletedKey: {
		"en-US": "Notification deleted successfully",
		"id-ID": "Notifikasi berhasil dihapus",
		"ja-JP": "通知が正常に削除されました",
	},
	SuccessNotificationRetrievedKey: {
		"en-US": "Notification retrieved successfully",
		"id-ID": "Notifikasi berhasil diambil",
		"ja-JP": "通知が正常に取得されました",
	},
	SuccessNotificationCountedKey: {
		"en-US": "Notification counted successfully",
		"id-ID": "Notifikasi berhasil dihitung",
		"ja-JP": "通知が正常にカウントされました",
	},
	SuccessNotificationStatisticsRetrievedKey: {
		"en-US": "Notification statistics retrieved successfully",
		"id-ID": "Statistik notifikasi berhasil diambil",
		"ja-JP": "通知統計が正常に取得されました",
	},
	SuccessNotificationExistenceCheckedKey: {
		"en-US": "Notification existence checked successfully",
		"id-ID": "Keberadaan notifikasi berhasil diperiksa",
		"ja-JP": "通知の存在が正常に確認されました",
	},
	SuccessNotificationMarkedAsReadKey: {
		"en-US": "Notification marked as read successfully",
		"id-ID": "Notifikasi berhasil ditandai sebagai sudah dibaca",
		"ja-JP": "通知が既読として正常にマークされました",
	},
	SuccessNotificationMarkedAsUnreadKey: {
		"en-US": "Notification marked as unread successfully",
		"id-ID": "Notifikasi berhasil ditandai sebagai belum dibaca",
		"ja-JP": "通知が未読として正常にマークされました",
	},

	// * Issue report error messages
	ErrIssueReportNotFoundKey: {
		"en-US": "Issue report not found",
		"id-ID": "Laporan masalah tidak ditemukan",
		"ja-JP": "問題レポートが見つかりません",
	},
	ErrIssueReportIDRequiredKey: {
		"en-US": "Issue report ID is required",
		"id-ID": "ID laporan masalah diperlukan",
		"ja-JP": "問題レポートIDが必要です",
	},
	ErrIssueReportAssetIDRequiredKey: {
		"en-US": "Asset ID is required",
		"id-ID": "ID aset diperlukan",
		"ja-JP": "アセットIDが必要です",
	},
	ErrIssueReportTypeRequiredKey: {
		"en-US": "Issue type is required",
		"id-ID": "Jenis masalah diperlukan",
		"ja-JP": "問題タイプが必要です",
	},
	ErrIssueReportPriorityRequiredKey: {
		"en-US": "Priority is required",
		"id-ID": "Prioritas diperlukan",
		"ja-JP": "優先度が必要です",
	},
	ErrIssueReportTitleRequiredKey: {
		"en-US": "Title is required",
		"id-ID": "Judul diperlukan",
		"ja-JP": "タイトルが必要です",
	},
	ErrIssueReportAlreadyResolvedKey: {
		"en-US": "Issue report is already resolved",
		"id-ID": "Laporan masalah sudah diselesaikan",
		"ja-JP": "問題レポートは既に解決されています",
	},
	ErrIssueReportCannotReopenKey: {
		"en-US": "Cannot reopen closed issue report",
		"id-ID": "Tidak dapat membuka kembali laporan masalah yang sudah ditutup",
		"ja-JP": "閉じられた問題レポートを再開できません",
	},

	// * Issue report success messages
	SuccessIssueReportCreatedKey: {
		"en-US": "Issue report created successfully",
		"id-ID": "Laporan masalah berhasil dibuat",
		"ja-JP": "問題レポートが正常に作成されました",
	},
	SuccessIssueReportUpdatedKey: {
		"en-US": "Issue report updated successfully",
		"id-ID": "Laporan masalah berhasil diperbarui",
		"ja-JP": "問題レポートが正常に更新されました",
	},
	SuccessIssueReportDeletedKey: {
		"en-US": "Issue report deleted successfully",
		"id-ID": "Laporan masalah berhasil dihapus",
		"ja-JP": "問題レポートが正常に削除されました",
	},
	SuccessIssueReportRetrievedKey: {
		"en-US": "Issue report retrieved successfully",
		"id-ID": "Laporan masalah berhasil diambil",
		"ja-JP": "問題レポートが正常に取得されました",
	},
	SuccessIssueReportCountedKey: {
		"en-US": "Issue report counted successfully",
		"id-ID": "Laporan masalah berhasil dihitung",
		"ja-JP": "問題レポートが正常にカウントされました",
	},
	SuccessIssueReportStatisticsRetrievedKey: {
		"en-US": "Issue report statistics retrieved successfully",
		"id-ID": "Statistik laporan masalah berhasil diambil",
		"ja-JP": "問題レポート統計が正常に取得されました",
	},
	SuccessIssueReportExistenceCheckedKey: {
		"en-US": "Issue report existence checked successfully",
		"id-ID": "Keberadaan laporan masalah berhasil diperiksa",
		"ja-JP": "問題レポートの存在が正常に確認されました",
	},
	SuccessIssueReportResolvedKey: {
		"en-US": "Issue report resolved successfully",
		"id-ID": "Laporan masalah berhasil diselesaikan",
		"ja-JP": "問題レポートが正常に解決されました",
	},
	SuccessIssueReportReopenedKey: {
		"en-US": "Issue report reopened successfully",
		"id-ID": "Laporan masalah berhasil dibuka kembali",
		"ja-JP": "問題レポートが正常に再開されました",
	},

	// * Asset movement error messages
	ErrAssetMovementNotFoundKey: {
		"en-US": "Asset movement not found",
		"id-ID": "Pergerakan aset tidak ditemukan",
		"ja-JP": "アセット移動が見つかりません",
	},
	ErrAssetMovementIDRequiredKey: {
		"en-US": "Asset movement ID is required",
		"id-ID": "ID pergerakan aset diperlukan",
		"ja-JP": "アセット移動IDが必要です",
	},
	ErrAssetMovementAssetIDRequiredKey: {
		"en-US": "Asset ID is required",
		"id-ID": "ID aset diperlukan",
		"ja-JP": "アセットIDが必要です",
	},
	ErrAssetMovementInvalidLocationKey: {
		"en-US": "Invalid location specified",
		"id-ID": "Lokasi yang ditentukan tidak valid",
		"ja-JP": "指定された場所が無効です",
	},
	ErrAssetMovementInvalidUserKey: {
		"en-US": "Invalid user specified",
		"id-ID": "Pengguna yang ditentukan tidak valid",
		"ja-JP": "指定されたユーザーが無効です",
	},
	ErrAssetMovementNoChangeKey: {
		"en-US": "No change detected in asset location or assignment",
		"id-ID": "Tidak ada perubahan yang terdeteksi dalam lokasi atau penugasan aset",
		"ja-JP": "アセットの場所または割り当てに変更が検出されませんでした",
	},
	ErrAssetMovementSameLocationKey: {
		"en-US": "Asset is already at the specified location",
		"id-ID": "Aset sudah berada di lokasi yang ditentukan",
		"ja-JP": "アセットは既に指定された場所にあります",
	},

	// * Asset movement success messages
	SuccessAssetMovementCreatedKey: {
		"en-US": "Asset movement created successfully",
		"id-ID": "Pergerakan aset berhasil dibuat",
		"ja-JP": "アセット移動が正常に作成されました",
	},
	SuccessAssetMovementUpdatedKey: {
		"en-US": "Asset movement updated successfully",
		"id-ID": "Pergerakan aset berhasil diperbarui",
		"ja-JP": "アセット移動が正常に更新されました",
	},
	SuccessAssetMovementDeletedKey: {
		"en-US": "Asset movement deleted successfully",
		"id-ID": "Pergerakan aset berhasil dihapus",
		"ja-JP": "アセット移動が正常に削除されました",
	},
	SuccessAssetMovementRetrievedKey: {
		"en-US": "Asset movement retrieved successfully",
		"id-ID": "Pergerakan aset berhasil diambil",
		"ja-JP": "アセット移動が正常に取得されました",
	},
	SuccessAssetMovementCountedKey: {
		"en-US": "Asset movement counted successfully",
		"id-ID": "Pergerakan aset berhasil dihitung",
		"ja-JP": "アセット移動が正常にカウントされました",
	},
	SuccessAssetMovementStatisticsRetrievedKey: {
		"en-US": "Asset movement statistics retrieved successfully",
		"id-ID": "Statistik pergerakan aset berhasil diambil",
		"ja-JP": "アセット移動統計が正常に取得されました",
	},
	SuccessAssetMovementExistenceCheckedKey: {
		"en-US": "Asset movement existence checked successfully",
		"id-ID": "Keberadaan pergerakan aset berhasil diperiksa",
		"ja-JP": "アセット移動の存在が正常に確認されました",
	},

	// * Maintenance error messages
	ErrMaintenanceScheduleNotFoundKey: {
		"en-US": "Maintenance schedule not found",
		"id-ID": "Jadwal pemeliharaan tidak ditemukan",
		"ja-JP": "保守スケジュールが見つかりません",
	},
	ErrMaintenanceRecordNotFoundKey: {
		"en-US": "Maintenance record not found",
		"id-ID": "Catatan pemeliharaan tidak ditemukan",
		"ja-JP": "保守レコードが見つかりません",
	},
	ErrMaintenanceScheduleIDRequiredKey: {
		"en-US": "Maintenance schedule ID is required",
		"id-ID": "ID jadwal pemeliharaan diperlukan",
		"ja-JP": "保守スケジュールIDが必要です",
	},
	ErrMaintenanceRecordIDRequiredKey: {
		"en-US": "Maintenance record ID is required",
		"id-ID": "ID catatan pemeliharaan diperlukan",
		"ja-JP": "保守レコードIDが必要です",
	},
	ErrMaintenanceAssetIDRequiredKey: {
		"en-US": "Asset ID is required",
		"id-ID": "ID aset diperlukan",
		"ja-JP": "アセットIDが必要です",
	},
	ErrMaintenanceScheduleDateRequiredKey: {
		"en-US": "Scheduled date is required",
		"id-ID": "Tanggal terjadwal diperlukan",
		"ja-JP": "予定日が必要です",
	},
	ErrMaintenanceRecordDateRequiredKey: {
		"en-US": "Maintenance date is required",
		"id-ID": "Tanggal pemeliharaan diperlukan",
		"ja-JP": "保守日が必要です",
	},
	ErrMaintenanceScheduleTitleRequiredKey: {
		"en-US": "Schedule title is required",
		"id-ID": "Judul jadwal diperlukan",
		"ja-JP": "スケジュールタイトルが必要です",
	},
	ErrMaintenanceRecordTitleRequiredKey: {
		"en-US": "Record title is required",
		"id-ID": "Judul catatan diperlukan",
		"ja-JP": "レコードタイトルが必要です",
	},

	// * Maintenance success messages
	SuccessMaintenanceScheduleCreatedKey: {
		"en-US": "Maintenance schedule created successfully",
		"id-ID": "Jadwal pemeliharaan berhasil dibuat",
		"ja-JP": "保守スケジュールが正常に作成されました",
	},
	SuccessMaintenanceScheduleUpdatedKey: {
		"en-US": "Maintenance schedule updated successfully",
		"id-ID": "Jadwal pemeliharaan berhasil diperbarui",
		"ja-JP": "保守スケジュールが正常に更新されました",
	},
	SuccessMaintenanceScheduleDeletedKey: {
		"en-US": "Maintenance schedule deleted successfully",
		"id-ID": "Jadwal pemeliharaan berhasil dihapus",
		"ja-JP": "保守スケジュールが正常に削除されました",
	},
	SuccessMaintenanceScheduleRetrievedKey: {
		"en-US": "Maintenance schedules retrieved successfully",
		"id-ID": "Jadwal pemeliharaan berhasil diambil",
		"ja-JP": "保守スケジュールが正常に取得されました",
	},
	SuccessMaintenanceScheduleCountedKey: {
		"en-US": "Maintenance schedules counted successfully",
		"id-ID": "Jadwal pemeliharaan berhasil dihitung",
		"ja-JP": "保守スケジュールが正常にカウントされました",
	},
	SuccessMaintenanceScheduleStatisticsRetrievedKey: {
		"en-US": "Maintenance schedule statistics retrieved successfully",
		"id-ID": "Statistik jadwal pemeliharaan berhasil diambil",
		"ja-JP": "保守スケジュール統計が正常に取得されました",
	},
	SuccessMaintenanceRecordCreatedKey: {
		"en-US": "Maintenance record created successfully",
		"id-ID": "Catatan pemeliharaan berhasil dibuat",
		"ja-JP": "保守レコードが正常に作成されました",
	},
	SuccessMaintenanceRecordUpdatedKey: {
		"en-US": "Maintenance record updated successfully",
		"id-ID": "Catatan pemeliharaan berhasil diperbarui",
		"ja-JP": "保守レコードが正常に更新されました",
	},
	SuccessMaintenanceRecordDeletedKey: {
		"en-US": "Maintenance record deleted successfully",
		"id-ID": "Catatan pemeliharaan berhasil dihapus",
		"ja-JP": "保守レコードが正常に削除されました",
	},
	SuccessMaintenanceRecordRetrievedKey: {
		"en-US": "Maintenance records retrieved successfully",
		"id-ID": "Catatan pemeliharaan berhasil diambil",
		"ja-JP": "保守レコードが正常に取得されました",
	},
	SuccessMaintenanceRecordCountedKey: {
		"en-US": "Maintenance records counted successfully",
		"id-ID": "Catatan pemeliharaan berhasil dihitung",
		"ja-JP": "保守レコードが正常にカウントされました",
	},
	SuccessMaintenanceRecordStatisticsRetrievedKey: {
		"en-US": "Maintenance record statistics retrieved successfully",
		"id-ID": "Statistik catatan pemeliharaan berhasil diambil",
		"ja-JP": "保守レコード統計が正常に取得されました",
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
