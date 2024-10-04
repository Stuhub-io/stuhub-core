// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePage = "pages"

// Page mapped from table <pages>
type Page struct {
	Pkid           int64      `gorm:"column:pkid;type:bigint;primaryKey;autoIncrement:true" json:"pkid"`
	ID             string     `gorm:"column:id;type:uuid;not null;default:uuid_generate_v4()" json:"id"`
	Name           string     `gorm:"column:name;type:character varying(255);not null" json:"name"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()" json:"updated_at"`
	SpacePkid      int64      `gorm:"column:space_pkid;type:bigint;not null" json:"space_pkid"`
	ParentPagePkid *int64     `gorm:"column:parent_page_pkid;type:bigint;index:idx_parent_page_pkid,priority:1" json:"parent_page_pkid"`
	ViewType       string     `gorm:"column:view_type;type:character varying(50);not null" json:"view_type"`
	ArchivedAt     *time.Time `gorm:"column:archived_at;type:timestamp with time zone" json:"archived_at"`
	CoverImage     string     `gorm:"column:cover_image;type:character varying;not null" json:"cover_image"`
}

// TableName Page's table name
func (*Page) TableName() string {
	return TableNamePage
}
