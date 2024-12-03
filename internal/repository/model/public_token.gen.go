// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePublicToken = "public_token"

// PublicToken mapped from table <public_token>
type PublicToken struct {
	Pkid       int64      `gorm:"column:pkid;type:bigint;primaryKey;autoIncrement:true" json:"pkid"`
	ID         string     `gorm:"column:id;type:uuid;not null;default:uuid_generate_v4()" json:"id"`
	PagePkid   int64      `gorm:"column:page_pkid;type:bigint;not null" json:"page_pkid"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()" json:"created_at"`
	ArchivedAt *time.Time `gorm:"column:archived_at;type:timestamp with time zone" json:"archived_at"`
}

// TableName PublicToken's table name
func (*PublicToken) TableName() string {
	return TableNamePublicToken
}
