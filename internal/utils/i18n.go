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
	ErrUserNotFoundKey       MessageKey = "error.user.not_found"
	ErrUserNameExistsKey     MessageKey = "error.user.name_exists"
	ErrUserEmailExistsKey    MessageKey = "error.user.email_exists"
	ErrUserIDRequiredKey     MessageKey = "error.user.id_required"
	ErrUserNameRequiredKey   MessageKey = "error.user.name_required"
	ErrUserEmailRequiredKey  MessageKey = "error.user.email_required"
	ErrInvalidOldPasswordKey MessageKey = "error.user.invalid_old_password"

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
	ErrNotificationNotFoundKey         MessageKey = "error.notification.not_found"
	ErrNotificationIDRequiredKey       MessageKey = "error.notification.id_required"
	ErrNotificationUserIDRequiredKey   MessageKey = "error.notification.user_id_required"
	ErrNotificationTypeRequiredKey     MessageKey = "error.notification.type_required"
	ErrNotificationPriorityRequiredKey MessageKey = "error.notification.priority_required"
	ErrNotificationTitleRequiredKey    MessageKey = "error.notification.title_required"
	ErrNotificationMessageRequiredKey  MessageKey = "error.notification.message_required"

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
	ErrInvalidCredentialsKey    MessageKey = "error.auth.invalid_credentials"
	ErrTokenExpiredKey          MessageKey = "error.auth.token_expired"
	ErrTokenInvalidKey          MessageKey = "error.auth.token_invalid"
	ErrAPIKeyMissingKey         MessageKey = "error.auth.api_key_missing"
	ErrAPIKeyInvalidKey         MessageKey = "error.auth.api_key_invalid"
	ErrResetCodeInvalidKey      MessageKey = "error.auth.reset_code_invalid"
	ErrResetCodeExpiredKey      MessageKey = "error.auth.reset_code_expired"
	ErrEmailSendFailedKey       MessageKey = "error.auth.email_send_failed"
	ErrResetCodeNotFoundKey     MessageKey = "error.auth.reset_code_not_found"

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
	SuccessUsersBulkCreatedKey          MessageKey = "success.users.bulk_created"
	SuccessUsersBulkDeletedKey          MessageKey = "success.users.bulk_deleted"
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
	SuccessCategoriesBulkCreatedKey        MessageKey = "success.categories.bulk_created"
	SuccessCategoriesBulkDeletedKey        MessageKey = "success.categories.bulk_deleted"
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
	SuccessLocationsBulkCreatedKey         MessageKey = "success.locations.bulk_created"
	SuccessLocationsBulkDeletedKey         MessageKey = "success.locations.bulk_deleted"
	SuccessLocationRetrievedKey            MessageKey = "success.location.retrieved"
	SuccessLocationRetrievedByCodeKey      MessageKey = "success.location.retrieved_by_code"
	SuccessLocationCountedKey              MessageKey = "success.location.counted"
	SuccessLocationStatisticsRetrievedKey  MessageKey = "success.location.statistics_retrieved"
	SuccessLocationExistenceCheckedKey     MessageKey = "success.location.existence_checked"
	SuccessLocationCodeExistenceCheckedKey MessageKey = "success.location.code_existence_checked"

	// * Asset-specific success keys
	SuccessAssetCreatedKey                      MessageKey = "success.asset.created"
	SuccessAssetsBulkCreatedKey                 MessageKey = "success.assets.bulk_created"
	SuccessAssetUpdatedKey                      MessageKey = "success.asset.updated"
	SuccessAssetDeletedKey                      MessageKey = "success.asset.deleted"
	SuccessAssetsBulkDeletedKey                 MessageKey = "success.assets.bulk_deleted"
	SuccessAssetRetrievedKey                    MessageKey = "success.asset.retrieved"
	SuccessAssetRetrievedByTagKey               MessageKey = "success.asset.retrieved_by_tag"
	SuccessAssetRetrievedByDataMatrixKey        MessageKey = "success.asset.retrieved_by_datamatrix"
	SuccessAssetCountedKey                      MessageKey = "success.asset.counted"
	SuccessAssetStatisticsRetrievedKey          MessageKey = "success.asset.statistics_retrieved"
	SuccessAssetExistenceCheckedKey             MessageKey = "success.asset.existence_checked"
	SuccessAssetTagExistenceCheckedKey          MessageKey = "success.asset.tag_existence_checked"
	SuccessAssetDataMatrixExistenceCheckedKey   MessageKey = "success.asset.datamatrix_existence_checked"
	SuccessAssetSerialNumberExistenceCheckedKey MessageKey = "success.asset.serial_number_existence_checked"
	SuccessAssetTagGeneratedKey                 MessageKey = "success.asset.tag_generated"
	SuccessBulkAssetTagsGeneratedKey            MessageKey = "success.asset.bulk_tags_generated"
	SuccessBulkDataMatrixUploadedKey            MessageKey = "success.asset.bulk_datamatrix_uploaded"

	// * Scan log-specific success keys
	SuccessScanLogCreatedKey             MessageKey = "success.scan_log.created"
	SuccessScanLogDeletedKey             MessageKey = "success.scan_log.deleted"
	SuccessScanLogsBulkCreatedKey        MessageKey = "success.scan_logs.bulk_created"
	SuccessScanLogsBulkDeletedKey        MessageKey = "success.scan_logs.bulk_deleted"
	SuccessScanLogRetrievedKey           MessageKey = "success.scan_log.retrieved"
	SuccessScanLogCountedKey             MessageKey = "success.scan_log.counted"
	SuccessScanLogStatisticsRetrievedKey MessageKey = "success.scan_log.statistics_retrieved"
	SuccessScanLogExistenceCheckedKey    MessageKey = "success.scan_log.existence_checked"

	// * Notification-specific success keys
	SuccessNotificationCreatedKey             MessageKey = "success.notification.created"
	SuccessNotificationUpdatedKey             MessageKey = "success.notification.updated"
	SuccessNotificationDeletedKey             MessageKey = "success.notification.deleted"
	SuccessNotificationsBulkCreatedKey        MessageKey = "success.notifications.bulk_created"
	SuccessNotificationsBulkDeletedKey        MessageKey = "success.notifications.bulk_deleted"
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
	SuccessIssueReportsBulkCreatedKey        MessageKey = "success.issue_reports.bulk_created"
	SuccessIssueReportsBulkDeletedKey        MessageKey = "success.issue_reports.bulk_deleted"
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
	SuccessAssetMovementsBulkCreatedKey        MessageKey = "success.asset_movements.bulk_created"
	SuccessAssetMovementsBulkDeletedKey        MessageKey = "success.asset_movements.bulk_deleted"
	SuccessAssetMovementRetrievedKey           MessageKey = "success.asset_movement.retrieved"
	SuccessAssetMovementCountedKey             MessageKey = "success.asset_movement.counted"
	SuccessAssetMovementStatisticsRetrievedKey MessageKey = "success.asset_movement.statistics_retrieved"
	SuccessAssetMovementExistenceCheckedKey    MessageKey = "success.asset_movement.existence_checked"

	// * Maintenance-specific success keys
	SuccessMaintenanceScheduleCreatedKey             MessageKey = "success.maintenance.schedule_created"
	SuccessMaintenanceScheduleUpdatedKey             MessageKey = "success.maintenance.schedule_updated"
	SuccessMaintenanceScheduleDeletedKey             MessageKey = "success.maintenance.schedule_deleted"
	SuccessMaintenanceSchedulesBulkCreatedKey        MessageKey = "success.maintenance.schedules_bulk_created"
	SuccessMaintenanceSchedulesBulkDeletedKey        MessageKey = "success.maintenance.schedules_bulk_deleted"
	SuccessMaintenanceScheduleRetrievedKey           MessageKey = "success.maintenance.schedule_retrieved"
	SuccessMaintenanceScheduleCountedKey             MessageKey = "success.maintenance.schedule_counted"
	SuccessMaintenanceScheduleStatisticsRetrievedKey MessageKey = "success.maintenance.schedule_statistics_retrieved"
	SuccessMaintenanceRecordCreatedKey               MessageKey = "success.maintenance.record_created"
	SuccessMaintenanceRecordUpdatedKey               MessageKey = "success.maintenance.record_updated"
	SuccessMaintenanceRecordDeletedKey               MessageKey = "success.maintenance.record_deleted"
	SuccessMaintenanceRecordsBulkCreatedKey          MessageKey = "success.maintenance.records_bulk_created"
	SuccessMaintenanceRecordsBulkDeletedKey          MessageKey = "success.maintenance.records_bulk_deleted"
	SuccessMaintenanceRecordRetrievedKey             MessageKey = "success.maintenance.record_retrieved"
	SuccessMaintenanceRecordCountedKey               MessageKey = "success.maintenance.record_counted"
	SuccessMaintenanceRecordStatisticsRetrievedKey   MessageKey = "success.maintenance.record_statistics_retrieved"

	// * Auth-specific success keys
	SuccessLoginKey              MessageKey = "success.auth.login"
	SuccessLogoutKey             MessageKey = "success.auth.logout"
	SuccessRefreshKey            MessageKey = "success.auth.refresh"
	SuccessTokenRefreshedKey     MessageKey = "success.auth.token_refreshed"
	SuccessResetCodeSentKey      MessageKey = "success.auth.reset_code_sent"
	SuccessResetCodeVerifiedKey  MessageKey = "success.auth.reset_code_verified"
	SuccessPasswordResetKey      MessageKey = "success.auth.password_reset"

	// * File upload success keys
	SuccessFileUploadedKey          MessageKey = "success.file.uploaded"
	SuccessAvatarUploadedKey        MessageKey = "success.file.avatar_uploaded"
	SuccessFileDeletedKey           MessageKey = "success.file.deleted"
	SuccessMultipleFilesUploadedKey MessageKey = "success.file.multiple_uploaded"

	// * Asset PDF Export labels
	PDFAssetListReportKey       MessageKey = "pdf.asset_list_report"
	PDFAssetGeneratedOnKey      MessageKey = "pdf.generated_on"
	PDFAssetTotalAssetsKey      MessageKey = "pdf.total_assets"
	PDFAssetAssetTagKey         MessageKey = "pdf.asset_tag"
	PDFAssetAssetNameKey        MessageKey = "pdf.asset_name"
	PDFAssetCategoryKey         MessageKey = "pdf.category"
	PDFAssetBrandKey            MessageKey = "pdf.brand"
	PDFAssetModelKey            MessageKey = "pdf.model"
	PDFAssetStatusKey           MessageKey = "pdf.status"
	PDFAssetConditionKey        MessageKey = "pdf.condition"
	PDFAssetLocationKey         MessageKey = "pdf.location"
	PDFAssetSerialNumberKey     MessageKey = "pdf.serial_number"
	PDFAssetPurchaseDateKey     MessageKey = "pdf.purchase_date"
	PDFAssetPurchasePriceKey    MessageKey = "pdf.purchase_price"
	PDFAssetVendorKey           MessageKey = "pdf.vendor"
	PDFAssetWarrantyEndKey      MessageKey = "pdf.warranty_end"
	PDFAssetAssignedToKey       MessageKey = "pdf.assigned_to"
	PDFAssetStatisticsReportKey MessageKey = "pdf.asset_statistics_report"
	PDFAssetDataMatrixReportKey MessageKey = "pdf.asset_datamatrix_report"

	// * Asset Movement PDF Export labels
	PDFAssetMovementReportKey       MessageKey = "pdf.asset_movement_report"
	PDFAssetMovementTotalKey        MessageKey = "pdf.total_movements"
	PDFAssetMovementFromLocationKey MessageKey = "pdf.from_location"
	PDFAssetMovementToLocationKey   MessageKey = "pdf.to_location"
	PDFAssetMovementFromUserKey     MessageKey = "pdf.from_user"
	PDFAssetMovementToUserKey       MessageKey = "pdf.to_user"
	PDFAssetMovementMovedByKey      MessageKey = "pdf.moved_by"
	PDFAssetMovementDateKey         MessageKey = "pdf.movement_date"

	// * Maintenance Record PDF Export labels
	PDFMaintenanceRecordReportKey     MessageKey = "pdf.maintenance_record_report"
	PDFMaintenanceRecordTotalKey      MessageKey = "pdf.total_records"
	PDFMaintenanceRecordDateKey       MessageKey = "pdf.maintenance_date"
	PDFMaintenanceRecordCompletionKey MessageKey = "pdf.completion_date"
	PDFMaintenanceRecordPerformerKey  MessageKey = "pdf.performer"
	PDFMaintenanceRecordCostKey       MessageKey = "pdf.cost"

	// * Issue Report PDF Export labels
	PDFIssueReportReportKey     MessageKey = "pdf.issue_report_report"
	PDFIssueReportTotalKey      MessageKey = "pdf.total_reports"
	PDFIssueReportTypeKey       MessageKey = "pdf.issue_type"
	PDFIssueReportReportedByKey MessageKey = "pdf.reported_by"

	// * Scan Log PDF Export labels
	PDFScanLogReportKey      MessageKey = "pdf.scan_log_report"
	PDFScanLogTotalKey       MessageKey = "pdf.total_logs"
	PDFScanLogMethodKey      MessageKey = "pdf.scan_method"
	PDFScanLogTimestampKey   MessageKey = "pdf.scan_timestamp"
	PDFScanLogResultKey      MessageKey = "pdf.scan_result"
	PDFScanLogScannedByKey   MessageKey = "pdf.scanned_by"
	PDFScanLogCoordinatesKey MessageKey = "pdf.coordinates"

	// * Maintenance Schedule PDF Export labels
	PDFMaintenanceScheduleReportKey    MessageKey = "pdf.maintenance_schedule_report"
	PDFMaintenanceScheduleTotalKey     MessageKey = "pdf.total_schedules"
	PDFMaintenanceScheduleTypeKey      MessageKey = "pdf.maintenance_type"
	PDFMaintenanceScheduleNextDateKey  MessageKey = "pdf.next_date"
	PDFMaintenanceScheduleRecurringKey MessageKey = "pdf.recurring"
	PDFMaintenanceScheduleCostKey      MessageKey = "pdf.estimated_cost"

	// * User PDF Export labels
	PDFUserReportKey MessageKey = "pdf.user_report"
	PDFUserTotalKey  MessageKey = "pdf.total_users_export"

	// * User PDF Export labels (keep existing)
	PDFUserListReportKey    MessageKey = "pdf.user_list_report"
	PDFUserIDKey            MessageKey = "pdf.user.id"
	PDFUserNameKey          MessageKey = "pdf.user.name"
	PDFUserEmailKey         MessageKey = "pdf.user.email"
	PDFUserFullNameKey      MessageKey = "pdf.user.full_name"
	PDFUserRoleKey          MessageKey = "pdf.user.role"
	PDFUserEmployeeIDKey    MessageKey = "pdf.user.employee_id"
	PDFUserPreferredLangKey MessageKey = "pdf.user.preferred_lang"
	PDFUserIsActiveKey      MessageKey = "pdf.user.is_active"
	PDFUserCreatedAtKey     MessageKey = "pdf.user.created_at"
	PDFUserUpdatedAtKey     MessageKey = "pdf.user.updated_at"

	PDFUserTotalUsersKey MessageKey = "pdf.total_users"

	// * Scan Log PDF Export labels
	PDFScanLogListReportKey      MessageKey = "pdf.scan_log_list_report"
	PDFScanLogIDKey              MessageKey = "pdf.scan_log.id"
	PDFScanLogAssetIDKey         MessageKey = "pdf.scan_log.asset_id"
	PDFScanLogScannedValueKey    MessageKey = "pdf.scan_log.scanned_value"
	PDFScanLogScanMethodKey      MessageKey = "pdf.scan_log.scan_method"
	PDFScanLogScannedByIDKey     MessageKey = "pdf.scan_log.scanned_by_id"
	PDFScanLogScanTimestampKey   MessageKey = "pdf.scan_log.scan_timestamp"
	PDFScanLogScanLocationLatKey MessageKey = "pdf.scan_log.scan_location_lat"
	PDFScanLogScanLocationLngKey MessageKey = "pdf.scan_log.scan_location_lng"
	PDFScanLogScanResultKey      MessageKey = "pdf.scan_log.scan_result"

	PDFScanLogTotalScanLogsKey MessageKey = "pdf.total_scan_logs"
	PDFScanLogPageKey          MessageKey = "pdf.page"
	PDFScanLogOfKey            MessageKey = "pdf.of"

	// * Maintenance Schedule PDF Export labels
	PDFMaintenanceScheduleListReportKey        MessageKey = "pdf.maintenance_schedule_list_report"
	PDFMaintenanceScheduleIDKey                MessageKey = "pdf.maintenance_schedule.id"
	PDFMaintenanceScheduleAssetIDKey           MessageKey = "pdf.maintenance_schedule.asset_id"
	PDFMaintenanceScheduleMaintenanceTypeKey   MessageKey = "pdf.maintenance_schedule.maintenance_type"
	PDFMaintenanceScheduleIsRecurringKey       MessageKey = "pdf.maintenance_schedule.is_recurring"
	PDFMaintenanceScheduleIntervalValueKey     MessageKey = "pdf.maintenance_schedule.interval_value"
	PDFMaintenanceScheduleIntervalUnitKey      MessageKey = "pdf.maintenance_schedule.interval_unit"
	PDFMaintenanceScheduleScheduledTimeKey     MessageKey = "pdf.maintenance_schedule.scheduled_time"
	PDFMaintenanceScheduleNextScheduledDateKey MessageKey = "pdf.maintenance_schedule.next_scheduled_date"
	PDFMaintenanceScheduleLastExecutedDateKey  MessageKey = "pdf.maintenance_schedule.last_executed_date"
	PDFMaintenanceScheduleStateKey             MessageKey = "pdf.maintenance_schedule.state"
	PDFMaintenanceScheduleAutoCompleteKey      MessageKey = "pdf.maintenance_schedule.auto_complete"
	PDFMaintenanceScheduleEstimatedCostKey     MessageKey = "pdf.maintenance_schedule.estimated_cost"
	PDFMaintenanceScheduleCreatedByIDKey       MessageKey = "pdf.maintenance_schedule.created_by_id"
	PDFMaintenanceScheduleCreatedAtKey         MessageKey = "pdf.maintenance_schedule.created_at"
	PDFMaintenanceScheduleUpdatedAtKey         MessageKey = "pdf.maintenance_schedule.updated_at"
	PDFMaintenanceScheduleTitleKey             MessageKey = "pdf.maintenance_schedule.title"
	PDFMaintenanceScheduleDescriptionKey       MessageKey = "pdf.maintenance_schedule.description"

	PDFMaintenanceScheduleTotalSchedulesKey MessageKey = "pdf.total_maintenance_schedules"
	PDFMaintenanceSchedulePageKey           MessageKey = "pdf.page"
	PDFMaintenanceScheduleOfKey             MessageKey = "pdf.of"

	// * Maintenance Record PDF Export labels
	PDFMaintenanceRecordListReportKey        MessageKey = "pdf.maintenance_record_list_report"
	PDFMaintenanceRecordIDKey                MessageKey = "pdf.maintenance_record.id"
	PDFMaintenanceRecordScheduleIDKey        MessageKey = "pdf.maintenance_record.schedule_id"
	PDFMaintenanceRecordAssetIDKey           MessageKey = "pdf.maintenance_record.asset_id"
	PDFMaintenanceRecordMaintenanceDateKey   MessageKey = "pdf.maintenance_record.maintenance_date"
	PDFMaintenanceRecordCompletionDateKey    MessageKey = "pdf.maintenance_record.completion_date"
	PDFMaintenanceRecordDurationMinutesKey   MessageKey = "pdf.maintenance_record.duration_minutes"
	PDFMaintenanceRecordPerformedByUserIDKey MessageKey = "pdf.maintenance_record.performed_by_user_id"
	PDFMaintenanceRecordPerformedByVendorKey MessageKey = "pdf.maintenance_record.performed_by_vendor"
	PDFMaintenanceRecordResultKey            MessageKey = "pdf.maintenance_record.result"
	PDFMaintenanceRecordActualCostKey        MessageKey = "pdf.maintenance_record.actual_cost"
	PDFMaintenanceRecordTitleKey             MessageKey = "pdf.maintenance_record.title"
	PDFMaintenanceRecordNotesKey             MessageKey = "pdf.maintenance_record.notes"
	PDFMaintenanceRecordCreatedAtKey         MessageKey = "pdf.maintenance_record.created_at"
	PDFMaintenanceRecordUpdatedAtKey         MessageKey = "pdf.maintenance_record.updated_at"

	PDFMaintenanceRecordTotalRecordsKey MessageKey = "pdf.total_maintenance_records"
	PDFMaintenanceRecordPageKey         MessageKey = "pdf.page"
	PDFMaintenanceRecordOfKey           MessageKey = "pdf.of"

	// * Issue Report PDF Export labels
	PDFIssueReportListReportKey      MessageKey = "pdf.issue_report_list_report"
	PDFIssueReportIDKey              MessageKey = "pdf.issue_report.id"
	PDFIssueReportAssetIDKey         MessageKey = "pdf.issue_report.asset_id"
	PDFIssueReportReportedByIDKey    MessageKey = "pdf.issue_report.reported_by_id"
	PDFIssueReportReportedDateKey    MessageKey = "pdf.issue_report.reported_date"
	PDFIssueReportIssueTypeKey       MessageKey = "pdf.issue_report.issue_type"
	PDFIssueReportPriorityKey        MessageKey = "pdf.issue_report.priority"
	PDFIssueReportStatusKey          MessageKey = "pdf.issue_report.status"
	PDFIssueReportResolvedDateKey    MessageKey = "pdf.issue_report.resolved_date"
	PDFIssueReportResolvedByIDKey    MessageKey = "pdf.issue_report.resolved_by_id"
	PDFIssueReportTitleKey           MessageKey = "pdf.issue_report.title"
	PDFIssueReportDescriptionKey     MessageKey = "pdf.issue_report.description"
	PDFIssueReportResolutionNotesKey MessageKey = "pdf.issue_report.resolution_notes"
	PDFIssueReportCreatedAtKey       MessageKey = "pdf.issue_report.created_at"
	PDFIssueReportUpdatedAtKey       MessageKey = "pdf.issue_report.updated_at"

	PDFIssueReportTotalReportsKey MessageKey = "pdf.total_issue_reports"
	PDFIssueReportPageKey         MessageKey = "pdf.page"
	PDFIssueReportOfKey           MessageKey = "pdf.of"

	// * Asset Movement PDF Export labels
	PDFAssetMovementListReportKey     MessageKey = "pdf.asset_movement_list_report"
	PDFAssetMovementIDKey             MessageKey = "pdf.asset_movement.id"
	PDFAssetMovementAssetIDKey        MessageKey = "pdf.asset_movement.asset_id"
	PDFAssetMovementFromLocationIDKey MessageKey = "pdf.asset_movement.from_location_id"
	PDFAssetMovementToLocationIDKey   MessageKey = "pdf.asset_movement.to_location_id"
	PDFAssetMovementFromUserIDKey     MessageKey = "pdf.asset_movement.from_user_id"
	PDFAssetMovementToUserIDKey       MessageKey = "pdf.asset_movement.to_user_id"
	PDFAssetMovementMovedByIDKey      MessageKey = "pdf.asset_movement.moved_by_id"
	PDFAssetMovementMovementDateKey   MessageKey = "pdf.asset_movement.movement_date"
	PDFAssetMovementNotesKey          MessageKey = "pdf.asset_movement.notes"
	PDFAssetMovementCreatedAtKey      MessageKey = "pdf.asset_movement.created_at"
	PDFAssetMovementUpdatedAtKey      MessageKey = "pdf.asset_movement.updated_at"

	PDFAssetMovementTotalMovementsKey MessageKey = "pdf.asset_movement.total_movements"
	PDFAssetMovementPageKey           MessageKey = "pdf.page"
	PDFAssetMovementOfKey             MessageKey = "pdf.of"
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
	ErrInvalidOldPasswordKey: {
		"en-US": "Old password is incorrect",
		"id-ID": "Kata sandi lama tidak cocok",
		"ja-JP": "古いパスワードが正しくありません",
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

	// * Asset-specific error messages
	ErrAssetNotFoundKey: {
		"en-US": "Asset not found",
		"id-ID": "Aset tidak ditemukan",
		"ja-JP": "アセットが見つかりません",
	},
	ErrAssetTagExistsKey: {
		"en-US": "Asset tag already exists",
		"id-ID": "Tag aset sudah ada",
		"ja-JP": "アセットタグは既に存在します",
	},
	ErrAssetDataMatrixExistsKey: {
		"en-US": "Data matrix already exists",
		"id-ID": "Data matrix sudah ada",
		"ja-JP": "データマトリックスは既に存在します",
	},
	ErrAssetSerialNumberExistsKey: {
		"en-US": "Serial number already exists",
		"id-ID": "Nomor seri sudah ada",
		"ja-JP": "シリアル番号は既に存在します",
	},
	ErrAssetIDRequiredKey: {
		"en-US": "Asset ID is required",
		"id-ID": "ID aset diperlukan",
		"ja-JP": "アセットIDが必要です",
	},
	ErrAssetTagRequiredKey: {
		"en-US": "Asset tag is required",
		"id-ID": "Tag aset diperlukan",
		"ja-JP": "アセットタグが必要です",
	},
	ErrAssetDataMatrixRequiredKey: {
		"en-US": "Data matrix is required",
		"id-ID": "Data matrix diperlukan",
		"ja-JP": "データマトリックスが必要です",
	},
	ErrAssetSerialNumberRequiredKey: {
		"en-US": "Serial number is required",
		"id-ID": "Nomor seri diperlukan",
		"ja-JP": "シリアル番号が必要です",
	},

	// * Scan log-specific error messages
	ErrScanLogNotFoundKey: {
		"en-US": "Scan log not found",
		"id-ID": "Log scan tidak ditemukan",
		"ja-JP": "スキャンログが見つかりません",
	},
	ErrScanLogIDRequiredKey: {
		"en-US": "Scan log ID is required",
		"id-ID": "ID log scan diperlukan",
		"ja-JP": "スキャンログIDが必要です",
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
	ErrAPIKeyMissingKey: {
		"en-US": "API key is required",
		"id-ID": "API key diperlukan",
		"ja-JP": "APIキーが必要です",
	},
	ErrAPIKeyInvalidKey: {
		"en-US": "Invalid API key",
		"id-ID": "API key tidak valid",
		"ja-JP": "無効なAPIキー",
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
	SuccessUsersBulkCreatedKey: {
		"en-US": "Users created successfully",
		"id-ID": "Pengguna berhasil dibuat secara massal",
		"ja-JP": "複数のユーザーが正常に作成されました",
	},
	SuccessUsersBulkDeletedKey: {
		"en-US": "Users deleted successfully",
		"id-ID": "Pengguna berhasil dihapus secara massal",
		"ja-JP": "複数のユーザーが正常に削除されました",
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
	SuccessCategoriesBulkCreatedKey: {
		"en-US": "Categories created successfully",
		"id-ID": "Kategori berhasil dibuat secara massal",
		"ja-JP": "複数のカテゴリが正常に作成されました",
	},
	SuccessCategoriesBulkDeletedKey: {
		"en-US": "Categories bulk deleted successfully",
		"id-ID": "Kategori berhasil dihapus secara massal",
		"ja-JP": "カテゴリが一括削除されました",
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
	SuccessLocationsBulkCreatedKey: {
		"en-US": "Locations created successfully",
		"id-ID": "Lokasi berhasil dibuat secara massal",
		"ja-JP": "複数のロケーションが正常に作成されました",
	},
	SuccessLocationsBulkDeletedKey: {
		"en-US": "Locations deleted successfully",
		"id-ID": "Lokasi berhasil dihapus secara massal",
		"ja-JP": "複数のロケーションが正常に削除されました",
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

	// * Asset-specific success messages
	SuccessAssetCreatedKey: {
		"en-US": "Asset created successfully",
		"id-ID": "Aset berhasil dibuat",
		"ja-JP": "アセットが正常に作成されました",
	},
	SuccessAssetsBulkCreatedKey: {
		"en-US": "Assets created successfully",
		"id-ID": "Aset berhasil dibuat secara massal",
		"ja-JP": "複数のアセットが正常に作成されました",
	},
	SuccessAssetUpdatedKey: {
		"en-US": "Asset updated successfully",
		"id-ID": "Aset berhasil diperbarui",
		"ja-JP": "アセットが正常に更新されました",
	},
	SuccessAssetDeletedKey: {
		"en-US": "Asset deleted successfully",
		"id-ID": "Aset berhasil dihapus",
		"ja-JP": "アセットが正常に削除されました",
	},
	SuccessAssetsBulkDeletedKey: {
		"en-US": "Assets deleted successfully",
		"id-ID": "Aset berhasil dihapus secara massal",
		"ja-JP": "複数のアセットが正常に削除されました",
	},
	SuccessAssetRetrievedKey: {
		"en-US": "Assets retrieved successfully",
		"id-ID": "Aset berhasil diambil",
		"ja-JP": "アセットが正常に取得されました",
	},
	SuccessAssetRetrievedByTagKey: {
		"en-US": "Asset retrieved successfully by tag",
		"id-ID": "Aset berhasil diambil berdasarkan tag",
		"ja-JP": "タグでアセットが正常に取得されました",
	},
	SuccessAssetRetrievedByDataMatrixKey: {
		"en-US": "Asset retrieved successfully by data matrix",
		"id-ID": "Aset berhasil diambil berdasarkan data matrix",
		"ja-JP": "データマトリックスでアセットが正常に取得されました",
	},
	SuccessAssetCountedKey: {
		"en-US": "Assets counted successfully",
		"id-ID": "Aset berhasil dihitung",
		"ja-JP": "アセットが正常にカウントされました",
	},
	SuccessAssetStatisticsRetrievedKey: {
		"en-US": "Asset statistics retrieved successfully",
		"id-ID": "Statistik aset berhasil diambil",
		"ja-JP": "アセット統計が正常に取得されました",
	},
	SuccessAssetExistenceCheckedKey: {
		"en-US": "Asset existence checked successfully",
		"id-ID": "Keberadaan aset berhasil diperiksa",
		"ja-JP": "アセットの存在が正常に確認されました",
	},
	SuccessAssetTagExistenceCheckedKey: {
		"en-US": "Asset tag existence checked successfully",
		"id-ID": "Keberadaan tag aset berhasil diperiksa",
		"ja-JP": "アセットタグの存在が正常に確認されました",
	},
	SuccessAssetDataMatrixExistenceCheckedKey: {
		"en-US": "Data matrix existence checked successfully",
		"id-ID": "Keberadaan data matrix berhasil diperiksa",
		"ja-JP": "データマトリックスの存在が正常に確認されました",
	},
	SuccessAssetSerialNumberExistenceCheckedKey: {
		"en-US": "Serial number existence checked successfully",
		"id-ID": "Keberadaan nomor seri berhasil diperiksa",
		"ja-JP": "シリアル番号の存在が正常に確認されました",
	},
	SuccessAssetTagGeneratedKey: {
		"en-US": "Asset tag suggestion generated successfully",
		"id-ID": "Saran tag aset berhasil dibuat",
		"ja-JP": "アセットタグの提案が正常に生成されました",
	},
	SuccessBulkAssetTagsGeneratedKey: {
		"en-US": "Bulk asset tags generated successfully",
		"id-ID": "Tag aset massal berhasil dibuat",
		"ja-JP": "一括アセットタグが正常に生成されました",
	},
	SuccessBulkDataMatrixUploadedKey: {
		"en-US": "Bulk data matrix images uploaded successfully",
		"id-ID": "Gambar data matrix massal berhasil diunggah",
		"ja-JP": "一括データマトリックス画像が正常にアップロードされました",
	},

	// * Scan log-specific success messages
	SuccessScanLogCreatedKey: {
		"en-US": "Scan log created successfully",
		"id-ID": "Log scan berhasil dibuat",
		"ja-JP": "スキャンログが正常に作成されました",
	},
	SuccessScanLogDeletedKey: {
		"en-US": "Scan log deleted successfully",
		"id-ID": "Log scan berhasil dihapus",
		"ja-JP": "スキャンログが正常に削除されました",
	},
	SuccessScanLogsBulkCreatedKey: {
		"en-US": "Scan logs created successfully",
		"id-ID": "Log scan berhasil dibuat secara massal",
		"ja-JP": "複数のスキャンログが正常に作成されました",
	},
	SuccessScanLogsBulkDeletedKey: {
		"en-US": "Scan logs deleted successfully",
		"id-ID": "Log scan berhasil dihapus secara massal",
		"ja-JP": "複数のスキャンログが正常に削除されました",
	},
	SuccessScanLogRetrievedKey: {
		"en-US": "Scan logs retrieved successfully",
		"id-ID": "Log scan berhasil diambil",
		"ja-JP": "スキャンログが正常に取得されました",
	},
	SuccessScanLogCountedKey: {
		"en-US": "Scan logs counted successfully",
		"id-ID": "Log scan berhasil dihitung",
		"ja-JP": "スキャンログが正常にカウントされました",
	},
	SuccessScanLogStatisticsRetrievedKey: {
		"en-US": "Scan log statistics retrieved successfully",
		"id-ID": "Statistik log scan berhasil diambil",
		"ja-JP": "スキャンログ統計が正常に取得されました",
	},
	SuccessScanLogExistenceCheckedKey: {
		"en-US": "Scan log existence checked successfully",
		"id-ID": "Keberadaan log scan berhasil diperiksa",
		"ja-JP": "スキャンログの存在が正常に確認されました",
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

	// * PDF Export labels
	PDFAssetListReportKey: {
		"en-US": "Asset List Report",
		"id-ID": "Laporan Daftar Aset",
		"ja-JP": "資産一覧レポート",
	},
	PDFAssetGeneratedOnKey: {
		"en-US": "Generated on",
		"id-ID": "Dibuat pada",
		"ja-JP": "生成日時",
	},
	PDFAssetTotalAssetsKey: {
		"en-US": "Total Assets",
		"id-ID": "Total Aset",
		"ja-JP": "総資産数",
	},
	PDFAssetAssetTagKey: {
		"en-US": "Asset Tag",
		"id-ID": "Tag Aset",
		"ja-JP": "資産タグ",
	},
	PDFAssetAssetNameKey: {
		"en-US": "Asset Name",
		"id-ID": "Nama Aset",
		"ja-JP": "資産名",
	},
	PDFAssetCategoryKey: {
		"en-US": "Category",
		"id-ID": "Kategori",
		"ja-JP": "カテゴリ",
	},
	PDFAssetBrandKey: {
		"en-US": "Brand",
		"id-ID": "Merek",
		"ja-JP": "ブランド",
	},
	PDFAssetModelKey: {
		"en-US": "Model",
		"id-ID": "Model",
		"ja-JP": "モデル",
	},
	PDFAssetStatusKey: {
		"en-US": "Status",
		"id-ID": "Status",
		"ja-JP": "ステータス",
	},
	PDFAssetConditionKey: {
		"en-US": "Condition",
		"id-ID": "Kondisi",
		"ja-JP": "状態",
	},
	PDFAssetLocationKey: {
		"en-US": "Location",
		"id-ID": "Lokasi",
		"ja-JP": "場所",
	},
	PDFAssetSerialNumberKey: {
		"en-US": "Serial Number",
		"id-ID": "Nomor Seri",
		"ja-JP": "シリアル番号",
	},
	PDFAssetPurchaseDateKey: {
		"en-US": "Purchase Date",
		"id-ID": "Tanggal Pembelian",
		"ja-JP": "購入日",
	},
	PDFAssetPurchasePriceKey: {
		"en-US": "Purchase Price",
		"id-ID": "Harga Pembelian",
		"ja-JP": "購入価格",
	},
	PDFAssetVendorKey: {
		"en-US": "Vendor",
		"id-ID": "Vendor",
		"ja-JP": "ベンダー",
	},
	PDFAssetWarrantyEndKey: {
		"en-US": "Warranty End",
		"id-ID": "Akhir Garansi",
		"ja-JP": "保証終了日",
	},
	PDFAssetAssignedToKey: {
		"en-US": "Assigned To",
		"id-ID": "Ditugaskan Ke",
		"ja-JP": "割り当て先",
	},
	PDFAssetStatisticsReportKey: {
		"en-US": "Asset Statistics Report",
		"id-ID": "Laporan Statistik Aset",
		"ja-JP": "資産統計レポート",
	},
	PDFAssetDataMatrixReportKey: {
		"en-US": "Asset Data Matrix Codes",
		"id-ID": "Kode Data Matrix Aset",
		"ja-JP": "資産データマトリックスコード",
	},

	// * Asset Movement PDF Export labels
	PDFAssetMovementReportKey: {
		"en-US": "Asset Movement Report",
		"id-ID": "Laporan Pergerakan Aset",
		"ja-JP": "資産移動レポート",
	},
	PDFAssetMovementTotalKey: {
		"en-US": "Total Movements",
		"id-ID": "Total Pergerakan",
		"ja-JP": "総移動数",
	},
	PDFAssetMovementFromLocationKey: {
		"en-US": "From Location",
		"id-ID": "Dari Lokasi",
		"ja-JP": "元の場所",
	},
	PDFAssetMovementToLocationKey: {
		"en-US": "To Location",
		"id-ID": "Ke Lokasi",
		"ja-JP": "移動先の場所",
	},
	PDFAssetMovementFromUserKey: {
		"en-US": "From User",
		"id-ID": "Dari Pengguna",
		"ja-JP": "元の担当者",
	},
	PDFAssetMovementToUserKey: {
		"en-US": "To User",
		"id-ID": "Ke Pengguna",
		"ja-JP": "移動先の担当者",
	},
	PDFAssetMovementMovedByKey: {
		"en-US": "Moved By",
		"id-ID": "Dipindahkan Oleh",
		"ja-JP": "移動者",
	},
	PDFAssetMovementDateKey: {
		"en-US": "Movement Date",
		"id-ID": "Tanggal Pergerakan",
		"ja-JP": "移動日",
	},

	// * Maintenance Record PDF Export labels
	PDFMaintenanceRecordReportKey: {
		"en-US": "Maintenance Record Report",
		"id-ID": "Laporan Catatan Pemeliharaan",
		"ja-JP": "メンテナンス記録レポート",
	},
	PDFMaintenanceRecordTotalKey: {
		"en-US": "Total Records",
		"id-ID": "Total Catatan",
		"ja-JP": "総記録数",
	},
	PDFMaintenanceRecordDateKey: {
		"en-US": "Maintenance Date",
		"id-ID": "Tanggal Pemeliharaan",
		"ja-JP": "メンテナンス日",
	},
	PDFMaintenanceRecordCompletionKey: {
		"en-US": "Completion Date",
		"id-ID": "Tanggal Penyelesaian",
		"ja-JP": "完了日",
	},
	PDFMaintenanceRecordPerformerKey: {
		"en-US": "Performer",
		"id-ID": "Pelaksana",
		"ja-JP": "実施者",
	},
	PDFMaintenanceRecordCostKey: {
		"en-US": "Cost",
		"id-ID": "Biaya",
		"ja-JP": "費用",
	},

	// * Issue Report PDF Export labels
	PDFIssueReportReportKey: {
		"en-US": "Issue Report List",
		"id-ID": "Daftar Laporan Masalah",
		"ja-JP": "問題報告一覧",
	},
	PDFIssueReportTotalKey: {
		"en-US": "Total Reports",
		"id-ID": "Total Laporan",
		"ja-JP": "総報告数",
	},
	PDFIssueReportTypeKey: {
		"en-US": "Issue Type",
		"id-ID": "Jenis Masalah",
		"ja-JP": "問題タイプ",
	},
	PDFIssueReportReportedByKey: {
		"en-US": "Reported By",
		"id-ID": "Dilaporkan Oleh",
		"ja-JP": "報告者",
	},

	// * Scan Log PDF Export labels
	PDFScanLogReportKey: {
		"en-US": "Scan Log Report",
		"id-ID": "Laporan Log Scan",
		"ja-JP": "スキャンログレポート",
	},
	PDFScanLogTotalKey: {
		"en-US": "Total Scans",
		"id-ID": "Total Scan",
		"ja-JP": "総スキャン数",
	},
	PDFScanLogMethodKey: {
		"en-US": "Scan Method",
		"id-ID": "Metode Scan",
		"ja-JP": "スキャン方法",
	},
	PDFScanLogTimestampKey: {
		"en-US": "Timestamp",
		"id-ID": "Stempel Waktu",
		"ja-JP": "タイムスタンプ",
	},
	PDFScanLogResultKey: {
		"en-US": "Result",
		"id-ID": "Hasil",
		"ja-JP": "結果",
	},
	PDFScanLogScannedByKey: {
		"en-US": "Scanned By",
		"id-ID": "Dipindai Oleh",
		"ja-JP": "スキャン者",
	},
	PDFScanLogCoordinatesKey: {
		"en-US": "Coordinates",
		"id-ID": "Koordinat",
		"ja-JP": "座標",
	},

	// * Maintenance Schedule PDF Export labels
	PDFMaintenanceScheduleReportKey: {
		"en-US": "Maintenance Schedule Report",
		"id-ID": "Laporan Jadwal Pemeliharaan",
		"ja-JP": "メンテナンススケジュールレポート",
	},
	PDFMaintenanceScheduleTotalKey: {
		"en-US": "Total Schedules",
		"id-ID": "Total Jadwal",
		"ja-JP": "総スケジュール数",
	},
	PDFMaintenanceScheduleTypeKey: {
		"en-US": "Type",
		"id-ID": "Tipe",
		"ja-JP": "種類",
	},
	PDFMaintenanceScheduleNextDateKey: {
		"en-US": "Next Date",
		"id-ID": "Tanggal Berikutnya",
		"ja-JP": "次回日",
	},
	PDFMaintenanceScheduleRecurringKey: {
		"en-US": "Recurring",
		"id-ID": "Berulang",
		"ja-JP": "繰り返し",
	},
	PDFMaintenanceScheduleCostKey: {
		"en-US": "Estimated Cost",
		"id-ID": "Biaya Perkiraan",
		"ja-JP": "予定費用",
	},

	// * User PDF Export labels
	PDFUserReportKey: {
		"en-US": "User List Report",
		"id-ID": "Laporan Daftar Pengguna",
		"ja-JP": "ユーザー一覧レポート",
	},
	PDFUserTotalKey: {
		"en-US": "Total Users",
		"id-ID": "Total Pengguna",
		"ja-JP": "総ユーザー数",
	},

	// * User PDF Export labels (keep existing)
	PDFUserListReportKey: {
		"en-US": "User List Report",
		"id-ID": "Laporan Daftar Pengguna",
		"ja-JP": "ユーザー一覧レポート",
	},
	PDFUserIDKey: {
		"en-US": "ID",
		"id-ID": "ID",
		"ja-JP": "ID",
	},
	PDFUserNameKey: {
		"en-US": "Name",
		"id-ID": "Nama",
		"ja-JP": "名前",
	},
	PDFUserEmailKey: {
		"en-US": "Email",
		"id-ID": "Email",
		"ja-JP": "メール",
	},
	PDFUserFullNameKey: {
		"en-US": "Full Name",
		"id-ID": "Nama Lengkap",
		"ja-JP": "フルネーム",
	},
	PDFUserRoleKey: {
		"en-US": "Role",
		"id-ID": "Peran",
		"ja-JP": "役割",
	},
	PDFUserEmployeeIDKey: {
		"en-US": "Employee ID",
		"id-ID": "ID Karyawan",
		"ja-JP": "従業員ID",
	},
	PDFUserPreferredLangKey: {
		"en-US": "Preferred Language",
		"id-ID": "Bahasa Pilihan",
		"ja-JP": "優先言語",
	},
	PDFUserIsActiveKey: {
		"en-US": "Is Active",
		"id-ID": "Aktif",
		"ja-JP": "アクティブ",
	},
	PDFUserCreatedAtKey: {
		"en-US": "Created At",
		"id-ID": "Dibuat Pada",
		"ja-JP": "作成日時",
	},
	PDFUserUpdatedAtKey: {
		"en-US": "Updated At",
		"id-ID": "Diperbarui Pada",
		"ja-JP": "更新日時",
	},

	PDFUserTotalUsersKey: {
		"en-US": "Total Users",
		"id-ID": "Total Pengguna",
		"ja-JP": "総ユーザー数",
	},

	// * Scan Log PDF Export labels
	PDFScanLogListReportKey: {
		"en-US": "Scan Log List Report",
		"id-ID": "Laporan Daftar Log Scan",
		"ja-JP": "スキャンログ一覧レポート",
	},
	PDFScanLogIDKey: {
		"en-US": "ID",
		"id-ID": "ID",
		"ja-JP": "ID",
	},
	PDFScanLogAssetIDKey: {
		"en-US": "Asset ID",
		"id-ID": "ID Aset",
		"ja-JP": "資産ID",
	},
	PDFScanLogScannedValueKey: {
		"en-US": "Scanned Value",
		"id-ID": "Nilai Scan",
		"ja-JP": "スキャン値",
	},
	PDFScanLogScanMethodKey: {
		"en-US": "Scan Method",
		"id-ID": "Metode Scan",
		"ja-JP": "スキャン方法",
	},
	PDFScanLogScannedByIDKey: {
		"en-US": "Scanned By ID",
		"id-ID": "ID Pemindai",
		"ja-JP": "スキャン者ID",
	},
	PDFScanLogScanTimestampKey: {
		"en-US": "Scan Timestamp",
		"id-ID": "Stempel Waktu Scan",
		"ja-JP": "スキャンタイムスタンプ",
	},
	PDFScanLogScanLocationLatKey: {
		"en-US": "Scan Location Latitude",
		"id-ID": "Lintang Lokasi Scan",
		"ja-JP": "スキャン位置緯度",
	},
	PDFScanLogScanLocationLngKey: {
		"en-US": "Scan Location Longitude",
		"id-ID": "Bujur Lokasi Scan",
		"ja-JP": "スキャン位置経度",
	},
	PDFScanLogScanResultKey: {
		"en-US": "Scan Result",
		"id-ID": "Hasil Scan",
		"ja-JP": "スキャン結果",
	},

	PDFScanLogTotalScanLogsKey: {
		"en-US": "Total Scan Logs",
		"id-ID": "Total Log Scan",
		"ja-JP": "総スキャンログ数",
	},

	// * Maintenance Schedule PDF Export labels
	PDFMaintenanceScheduleListReportKey: {
		"en-US": "Maintenance Schedule List Report",
		"id-ID": "Laporan Daftar Jadwal Pemeliharaan",
		"ja-JP": "保守スケジュール一覧レポート",
	},
	PDFMaintenanceScheduleIDKey: {
		"en-US": "ID",
		"id-ID": "ID",
		"ja-JP": "ID",
	},
	PDFMaintenanceScheduleAssetIDKey: {
		"en-US": "Asset ID",
		"id-ID": "ID Aset",
		"ja-JP": "資産ID",
	},
	PDFMaintenanceScheduleMaintenanceTypeKey: {
		"en-US": "Maintenance Type",
		"id-ID": "Tipe Pemeliharaan",
		"ja-JP": "保守タイプ",
	},
	PDFMaintenanceScheduleIsRecurringKey: {
		"en-US": "Is Recurring",
		"id-ID": "Berulang",
		"ja-JP": "繰り返し",
	},
	PDFMaintenanceScheduleIntervalValueKey: {
		"en-US": "Interval Value",
		"id-ID": "Nilai Interval",
		"ja-JP": "間隔値",
	},
	PDFMaintenanceScheduleIntervalUnitKey: {
		"en-US": "Interval Unit",
		"id-ID": "Unit Interval",
		"ja-JP": "間隔単位",
	},
	PDFMaintenanceScheduleScheduledTimeKey: {
		"en-US": "Scheduled Time",
		"id-ID": "Waktu Terjadwal",
		"ja-JP": "予定時間",
	},
	PDFMaintenanceScheduleNextScheduledDateKey: {
		"en-US": "Next Scheduled Date",
		"id-ID": "Tanggal Terjadwal Berikutnya",
		"ja-JP": "次回予定日",
	},
	PDFMaintenanceScheduleLastExecutedDateKey: {
		"en-US": "Last Executed Date",
		"id-ID": "Tanggal Eksekusi Terakhir",
		"ja-JP": "最終実行日",
	},
	PDFMaintenanceScheduleStateKey: {
		"en-US": "State",
		"id-ID": "Negara",
		"ja-JP": "状態",
	},
	PDFMaintenanceScheduleAutoCompleteKey: {
		"en-US": "Auto Complete",
		"id-ID": "Penyelesaian Otomatis",
		"ja-JP": "自動完了",
	},
	PDFMaintenanceScheduleEstimatedCostKey: {
		"en-US": "Estimated Cost",
		"id-ID": "Biaya Perkiraan",
		"ja-JP": "見積もりコスト",
	},
	PDFMaintenanceScheduleCreatedByIDKey: {
		"en-US": "Created By ID",
		"id-ID": "Dibuat Oleh ID",
		"ja-JP": "作成者ID",
	},
	PDFMaintenanceScheduleCreatedAtKey: {
		"en-US": "Created At",
		"id-ID": "Dibuat Pada",
		"ja-JP": "作成日時",
	},
	PDFMaintenanceScheduleUpdatedAtKey: {
		"en-US": "Updated At",
		"id-ID": "Diperbarui Pada",
		"ja-JP": "更新日時",
	},
	PDFMaintenanceScheduleTitleKey: {
		"en-US": "Title",
		"id-ID": "Judul",
		"ja-JP": "タイトル",
	},
	PDFMaintenanceScheduleDescriptionKey: {
		"en-US": "Description",
		"id-ID": "Deskripsi",
		"ja-JP": "説明",
	},

	PDFMaintenanceScheduleTotalSchedulesKey: {
		"en-US": "Total Maintenance Schedules",
		"id-ID": "Total Jadwal Pemeliharaan",
		"ja-JP": "総保守スケジュール数",
	},

	// * Maintenance Record PDF Export labels
	PDFMaintenanceRecordListReportKey: {
		"en-US": "Maintenance Record List Report",
		"id-ID": "Laporan Daftar Catatan Pemeliharaan",
		"ja-JP": "保守レコード一覧レポート",
	},
	PDFMaintenanceRecordIDKey: {
		"en-US": "ID",
		"id-ID": "ID",
		"ja-JP": "ID",
	},
	PDFMaintenanceRecordScheduleIDKey: {
		"en-US": "Schedule ID",
		"id-ID": "ID Jadwal",
		"ja-JP": "スケジュールID",
	},
	PDFMaintenanceRecordAssetIDKey: {
		"en-US": "Asset ID",
		"id-ID": "ID Aset",
		"ja-JP": "資産ID",
	},
	PDFMaintenanceRecordMaintenanceDateKey: {
		"en-US": "Maintenance Date",
		"id-ID": "Tanggal Pemeliharaan",
		"ja-JP": "保守日",
	},
	PDFMaintenanceRecordCompletionDateKey: {
		"en-US": "Completion Date",
		"id-ID": "Tanggal Penyelesaian",
		"ja-JP": "完了日",
	},
	PDFMaintenanceRecordDurationMinutesKey: {
		"en-US": "Duration (Minutes)",
		"id-ID": "Durasi (Menit)",
		"ja-JP": "期間（分）",
	},
	PDFMaintenanceRecordPerformedByUserIDKey: {
		"en-US": "Performed By User ID",
		"id-ID": "Dilakukan Oleh ID Pengguna",
		"ja-JP": "実行者ユーザーID",
	},
	PDFMaintenanceRecordPerformedByVendorKey: {
		"en-US": "Performed By Vendor",
		"id-ID": "Dilakukan Oleh Vendor",
		"ja-JP": "実行者ベンダー",
	},
	PDFMaintenanceRecordResultKey: {
		"en-US": "Result",
		"id-ID": "Hasil",
		"ja-JP": "結果",
	},
	PDFMaintenanceRecordActualCostKey: {
		"en-US": "Actual Cost",
		"id-ID": "Biaya Aktual",
		"ja-JP": "実際のコスト",
	},
	PDFMaintenanceRecordTitleKey: {
		"en-US": "Title",
		"id-ID": "Judul",
		"ja-JP": "タイトル",
	},
	PDFMaintenanceRecordNotesKey: {
		"en-US": "Notes",
		"id-ID": "Catatan",
		"ja-JP": "メモ",
	},
	PDFMaintenanceRecordCreatedAtKey: {
		"en-US": "Created At",
		"id-ID": "Dibuat Pada",
		"ja-JP": "作成日時",
	},
	PDFMaintenanceRecordUpdatedAtKey: {
		"en-US": "Updated At",
		"id-ID": "Diperbarui Pada",
		"ja-JP": "更新日時",
	},

	PDFMaintenanceRecordTotalRecordsKey: {
		"en-US": "Total Maintenance Records",
		"id-ID": "Total Catatan Pemeliharaan",
		"ja-JP": "総保守レコード数",
	},
	PDFMaintenanceRecordPageKey: {
		"en-US": "Page",
		"id-ID": "Halaman",
		"ja-JP": "ページ",
	},
	PDFMaintenanceRecordOfKey: {
		"en-US": "of",
		"id-ID": "dari",
		"ja-JP": "/",
	},

	// * Issue Report PDF Export labels
	PDFIssueReportListReportKey: {
		"en-US": "Issue Report List Report",
		"id-ID": "Laporan Daftar Laporan Masalah",
		"ja-JP": "問題レポート一覧レポート",
	},
	PDFIssueReportIDKey: {
		"en-US": "ID",
		"id-ID": "ID",
		"ja-JP": "ID",
	},
	PDFIssueReportAssetIDKey: {
		"en-US": "Asset ID",
		"id-ID": "ID Aset",
		"ja-JP": "資産ID",
	},
	PDFIssueReportReportedByIDKey: {
		"en-US": "Reported By ID",
		"id-ID": "Dilaporkan Oleh ID",
		"ja-JP": "報告者ID",
	},
	PDFIssueReportReportedDateKey: {
		"en-US": "Reported Date",
		"id-ID": "Tanggal Dilaporkan",
		"ja-JP": "報告日",
	},
	PDFIssueReportIssueTypeKey: {
		"en-US": "Issue Type",
		"id-ID": "Jenis Masalah",
		"ja-JP": "問題タイプ",
	},
	PDFIssueReportPriorityKey: {
		"en-US": "Priority",
		"id-ID": "Prioritas",
		"ja-JP": "優先度",
	},
	PDFIssueReportStatusKey: {
		"en-US": "Status",
		"id-ID": "Status",
		"ja-JP": "ステータス",
	},
	PDFIssueReportResolvedDateKey: {
		"en-US": "Resolved Date",
		"id-ID": "Tanggal Diselesaikan",
		"ja-JP": "解決日",
	},
	PDFIssueReportResolvedByIDKey: {
		"en-US": "Resolved By ID",
		"id-ID": "Diselesaikan Oleh ID",
		"ja-JP": "解決者ID",
	},
	PDFIssueReportTitleKey: {
		"en-US": "Title",
		"id-ID": "Judul",
		"ja-JP": "タイトル",
	},
	PDFIssueReportDescriptionKey: {
		"en-US": "Description",
		"id-ID": "Deskripsi",
		"ja-JP": "説明",
	},
	PDFIssueReportResolutionNotesKey: {
		"en-US": "Resolution Notes",
		"id-ID": "Catatan Penyelesaian",
		"ja-JP": "解決メモ",
	},
	PDFIssueReportCreatedAtKey: {
		"en-US": "Created At",
		"id-ID": "Dibuat Pada",
		"ja-JP": "作成日時",
	},
	PDFIssueReportUpdatedAtKey: {
		"en-US": "Updated At",
		"id-ID": "Diperbarui Pada",
		"ja-JP": "更新日時",
	},

	PDFIssueReportTotalReportsKey: {
		"en-US": "Total Issue Reports",
		"id-ID": "Total Laporan Masalah",
		"ja-JP": "総問題レポート数",
	},

	// * Asset Movement PDF Export labels
	PDFAssetMovementListReportKey: {
		"en-US": "Asset Movement List Report",
		"id-ID": "Laporan Daftar Pergerakan Aset",
		"ja-JP": "資産移動一覧レポート",
	},
	PDFAssetMovementIDKey: {
		"en-US": "ID",
		"id-ID": "ID",
		"ja-JP": "ID",
	},
	PDFAssetMovementAssetIDKey: {
		"en-US": "Asset ID",
		"id-ID": "ID Aset",
		"ja-JP": "資産ID",
	},
	PDFAssetMovementFromLocationIDKey: {
		"en-US": "From Location ID",
		"id-ID": "ID Lokasi Asal",
		"ja-JP": "移動元場所ID",
	},
	PDFAssetMovementToLocationIDKey: {
		"en-US": "To Location ID",
		"id-ID": "ID Lokasi Tujuan",
		"ja-JP": "移動先場所ID",
	},
	PDFAssetMovementFromUserIDKey: {
		"en-US": "From User ID",
		"id-ID": "ID Pengguna Asal",
		"ja-JP": "移動元ユーザーID",
	},
	PDFAssetMovementToUserIDKey: {
		"en-US": "To User ID",
		"id-ID": "ID Pengguna Tujuan",
		"ja-JP": "移動先ユーザーID",
	},
	PDFAssetMovementMovedByIDKey: {
		"en-US": "Moved By ID",
		"id-ID": "ID Pemindah",
		"ja-JP": "移動者ID",
	},
	PDFAssetMovementMovementDateKey: {
		"en-US": "Movement Date",
		"id-ID": "Tanggal Pergerakan",
		"ja-JP": "移動日",
	},
	PDFAssetMovementNotesKey: {
		"en-US": "Notes",
		"id-ID": "Catatan",
		"ja-JP": "メモ",
	},
	PDFAssetMovementCreatedAtKey: {
		"en-US": "Created At",
		"id-ID": "Dibuat Pada",
		"ja-JP": "作成日時",
	},
	PDFAssetMovementUpdatedAtKey: {
		"en-US": "Updated At",
		"id-ID": "Diperbarui Pada",
		"ja-JP": "更新日時",
	},

	PDFAssetMovementTotalMovementsKey: {
		"en-US": "Total Asset Movements",
		"id-ID": "Total Pergerakan Aset",
		"ja-JP": "総資産移動数",
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
	ErrNotificationPriorityRequiredKey: {
		"en-US": "Notification priority is required",
		"id-ID": "Prioritas notifikasi diperlukan",
		"ja-JP": "通知の優先度が必要です",
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
	SuccessNotificationsBulkCreatedKey: {
		"en-US": "Notifications created successfully",
		"id-ID": "Notifikasi berhasil dibuat secara massal",
		"ja-JP": "複数の通知が正常に作成されました",
	},
	SuccessNotificationsBulkDeletedKey: {
		"en-US": "Notifications deleted successfully",
		"id-ID": "Notifikasi berhasil dihapus secara massal",
		"ja-JP": "複数の通知が正常に削除されました",
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
	SuccessIssueReportsBulkCreatedKey: {
		"en-US": "Issue reports created successfully",
		"id-ID": "Laporan masalah berhasil dibuat secara massal",
		"ja-JP": "複数の問題レポートが正常に作成されました",
	},
	SuccessIssueReportsBulkDeletedKey: {
		"en-US": "Issue reports deleted successfully",
		"id-ID": "Laporan masalah berhasil dihapus secara massal",
		"ja-JP": "複数の問題レポートが正常に削除されました",
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
	SuccessAssetMovementsBulkCreatedKey: {
		"en-US": "Asset movements created successfully",
		"id-ID": "Pergerakan aset berhasil dibuat secara massal",
		"ja-JP": "複数のアセット移動が正常に作成されました",
	},
	SuccessAssetMovementsBulkDeletedKey: {
		"en-US": "Asset movements deleted successfully",
		"id-ID": "Pergerakan aset berhasil dihapus secara massal",
		"ja-JP": "複数のアセット移動が正常に削除されました",
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
	SuccessMaintenanceSchedulesBulkCreatedKey: {
		"en-US": "Maintenance schedules created successfully",
		"id-ID": "Jadwal pemeliharaan berhasil dibuat secara massal",
		"ja-JP": "複数の保守スケジュールが正常に作成されました",
	},
	SuccessMaintenanceSchedulesBulkDeletedKey: {
		"en-US": "Maintenance schedules deleted successfully",
		"id-ID": "Jadwal pemeliharaan berhasil dihapus secara massal",
		"ja-JP": "複数の保守スケジュールが正常に削除されました",
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
	SuccessMaintenanceRecordsBulkCreatedKey: {
		"en-US": "Maintenance records created successfully",
		"id-ID": "Catatan pemeliharaan berhasil dibuat secara massal",
		"ja-JP": "複数の保守レコードが正常に作成されました",
	},
	SuccessMaintenanceRecordsBulkDeletedKey: {
		"en-US": "Maintenance records deleted successfully",
		"id-ID": "Catatan pemeliharaan berhasil dihapus secara massal",
		"ja-JP": "複数の保守レコードが正常に削除されました",
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
