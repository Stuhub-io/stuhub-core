// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOrganization = "organizations"

// Organization mapped from table <organizations>
type Organization struct {
	Pkid        int64     `gorm:"column:pkid;type:bigint;primaryKey;autoIncrement:true" json:"pkid"`
	ID          string    `gorm:"column:id;type:uuid;not null;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"column:name;type:character varying(255);not null" json:"name"`
	Slug        string    `gorm:"column:slug;type:character varying(255);not null" json:"slug"`
	Description string    `gorm:"column:description;type:text;not null" json:"description"`
	Avatar      string    `gorm:"column:avatar;type:character varying;not null" json:"avatar"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()" json:"updated_at"`
}

// TableName Organization's table name
func (*Organization) TableName() string {
	return TableNameOrganization
}
