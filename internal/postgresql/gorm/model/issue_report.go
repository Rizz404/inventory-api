package model

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
)

type IssueReport struct {
	ID             SQLULID              `gorm:"primaryKey;type:varchar(26)"`
	AssetID        SQLULID              `gorm:"type:varchar(26);not null"`
	ReportedBy     SQLULID              `gorm:"type:varchar(26);not null"`
	ReportedDate   time.Time            `gorm:"default:CURRENT_TIMESTAMP"`
	IssueType      string               `gorm:"type:varchar(50);not null"`
	Priority       domain.IssuePriority `gorm:"type:issue_priority;default:'Medium'"`
	Status         domain.IssueStatus   `gorm:"type:issue_status;default:'Open'"`
	ResolvedDate   *time.Time
	ResolvedBy     *SQLULID                 `gorm:"type:varchar(26)"`
	Asset          Asset                    `gorm:"foreignKey:AssetID"`
	ReportedByUser User                     `gorm:"foreignKey:ReportedBy"`
	ResolvedByUser *User                    `gorm:"foreignKey:ResolvedBy"`
	Translations   []IssueReportTranslation `gorm:"foreignKey:ReportID"`
}

func (IssueReport) TableName() string {
	return "issue_reports"
}

type IssueReportTranslation struct {
	ID              SQLULID `gorm:"primaryKey;type:varchar(26)"`
	ReportID        SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_rep_lang"`
	LangCode        string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_rep_lang"`
	Title           string  `gorm:"type:varchar(200);not null"`
	Description     *string `gorm:"type:text"`
	ResolutionNotes *string `gorm:"type:text"`
}

func (IssueReportTranslation) TableName() string {
	return "issue_reports_translation"
}
