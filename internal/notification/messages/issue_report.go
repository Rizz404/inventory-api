package messages

// Issue Report notification message keys
const (
	// Issue Report Created
	NotifIssueReportedTitleKey   NotificationMessageKey = "notification.issue_report.reported.title"
	NotifIssueReportedMessageKey NotificationMessageKey = "notification.issue_report.reported.message"

	// Issue Report Updated
	NotifIssueUpdatedTitleKey   NotificationMessageKey = "notification.issue_report.updated.title"
	NotifIssueUpdatedMessageKey NotificationMessageKey = "notification.issue_report.updated.message"

	// Issue Report Resolved
	NotifIssueResolvedTitleKey   NotificationMessageKey = "notification.issue_report.resolved.title"
	NotifIssueResolvedMessageKey NotificationMessageKey = "notification.issue_report.resolved.message"

	// Issue Report Reopened
	NotifIssueReopenedTitleKey   NotificationMessageKey = "notification.issue_report.reopened.title"
	NotifIssueReopenedMessageKey NotificationMessageKey = "notification.issue_report.reopened.message"
)

// issueReportNotificationTranslations contains all issue report notification message translations
var issueReportNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== ISSUE REPORT CREATED ====================
	NotifIssueReportedTitleKey: {
		"en-US": "New Issue Reported",
		"id-ID": "Masalah Baru Dilaporkan",
		"ja-JP": "新しい問題が報告されました",
	},
	NotifIssueReportedMessageKey: {
		"en-US": "A new issue has been reported for asset \"{assetName}\".",
		"id-ID": "Masalah baru telah dilaporkan untuk aset \"{assetName}\".",
		"ja-JP": "資産 \"{assetName}\" に対して新しい問題が報告されました。",
	},

	// ==================== ISSUE REPORT UPDATED ====================
	NotifIssueUpdatedTitleKey: {
		"en-US": "Issue Updated",
		"id-ID": "Masalah Diperbarui",
		"ja-JP": "問題が更新されました",
	},
	NotifIssueUpdatedMessageKey: {
		"en-US": "Issue report for asset \"{assetName}\" has been updated.",
		"id-ID": "Laporan masalah untuk aset \"{assetName}\" telah diperbarui.",
		"ja-JP": "資産 \"{assetName}\" の問題レポートが更新されました。",
	},

	// ==================== ISSUE REPORT RESOLVED ====================
	NotifIssueResolvedTitleKey: {
		"en-US": "Issue Resolved",
		"id-ID": "Masalah Diselesaikan",
		"ja-JP": "問題が解決されました",
	},
	NotifIssueResolvedMessageKey: {
		"en-US": "Issue report for asset \"{assetName}\" has been resolved. Resolution: \"{resolutionNotes}\".",
		"id-ID": "Laporan masalah untuk aset \"{assetName}\" telah diselesaikan. Resolusi: \"{resolutionNotes}\".",
		"ja-JP": "資産 \"{assetName}\" の問題レポートが解決されました。解決策: \"{resolutionNotes}\"。",
	},

	// ==================== ISSUE REPORT REOPENED ====================
	NotifIssueReopenedTitleKey: {
		"en-US": "Issue Reopened",
		"id-ID": "Masalah Dibuka Kembali",
		"ja-JP": "問題が再開されました",
	},
	NotifIssueReopenedMessageKey: {
		"en-US": "Issue report for asset \"{assetName}\" has been reopened.",
		"id-ID": "Laporan masalah untuk aset \"{assetName}\" telah dibuka kembali.",
		"ja-JP": "資産 \"{assetName}\" の問題レポートが再開されました。",
	},
}

// GetIssueReportNotificationMessage returns the localized issue report notification message
func GetIssueReportNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	return GetNotificationMessage(key, langCode, params, issueReportNotificationTranslations)
}

// GetIssueReportNotificationTranslations returns all translations for an issue report notification
func GetIssueReportNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	return GetNotificationTranslations(titleKey, messageKey, params, issueReportNotificationTranslations)
}

// ==================== ISSUE REPORT NOTIFICATION HELPER FUNCTIONS ====================

// IssueReportedNotification creates notification for new issue report
func IssueReportedNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifIssueReportedTitleKey, NotifIssueReportedMessageKey, params
}

// IssueUpdatedNotification creates notification for issue report update
func IssueUpdatedNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifIssueUpdatedTitleKey, NotifIssueUpdatedMessageKey, params
}

// IssueResolvedNotification creates notification for resolved issue report
func IssueResolvedNotification(assetName, assetTag, resolutionNotes string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":       assetName,
		"assetTag":        assetTag,
		"resolutionNotes": resolutionNotes,
	}
	return NotifIssueResolvedTitleKey, NotifIssueResolvedMessageKey, params
}

// IssueReopenedNotification creates notification for reopened issue report
func IssueReopenedNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifIssueReopenedTitleKey, NotifIssueReopenedMessageKey, params
}
