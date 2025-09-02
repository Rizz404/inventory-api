package model

import (
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
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

func (u *IssueReport) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 IssueReport.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for IssueReport: %s", u.ID.String())
	}

	return nil
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
	return "issue_report_translations"
}

func (u *IssueReportTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 IssueReportTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for IssueReportTranslation: %s", u.ID.String())
	}

	return nil
}
