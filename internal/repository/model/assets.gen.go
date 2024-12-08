// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAsset = "assets"

// Asset mapped from table <assets>
type Asset struct {
	Pkid       int64     `gorm:"column:pkid;type:bigint;primaryKey;autoIncrement:true" json:"pkid"`
	PagePkid   int64     `gorm:"column:page_pkid;type:bigint;not null" json:"page_pkid"`
	URL        string    `gorm:"column:url;type:text;not null" json:"url"`
	Size       *int64    `gorm:"column:size;type:bigint" json:"size"`
	Extension  *string   `gorm:"column:extension;type:character(100)" json:"extension"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()" json:"updated_at"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()" json:"created_at"`
	Thumbnails string    `gorm:"column:thumbnails;type:jsonb;not null;default:{}" json:"thumbnails"`
}

// TableName Asset's table name
func (*Asset) TableName() string {
	return TableNameAsset
}
