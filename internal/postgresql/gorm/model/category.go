package model

type Category struct {
	ID           SQLULID               `gorm:"primaryKey;type:varchar(26)"`
	ParentID     *SQLULID              `gorm:"type:varchar(26)"`
	CategoryCode string                `gorm:"type:varchar(20);unique;not null"`
	Parent       *Category             `gorm:"foreignKey:ParentID"`
	Children     []Category            `gorm:"foreignKey:ParentID"`
	Translations []CategoryTranslation `gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
	return "categories"
}

type CategoryTranslation struct {
	ID           SQLULID `gorm:"primaryKey;type:varchar(26)"`
	CategoryID   SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_cat_lang"`
	LangCode     string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_cat_lang"`
	CategoryName string  `gorm:"type:varchar(100);not null"`
	Description  *string `gorm:"type:text"`
}

func (CategoryTranslation) TableName() string {
	return "categories_translation"
}
