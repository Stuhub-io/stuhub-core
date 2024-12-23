// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePageRole = "page_roles"

// PageRole mapped from table <page_roles>
type PageRole struct {
	Pkid      int64     `gorm:"column:pkid;type:bigint;primaryKey;autoIncrement:true" json:"pkid"`
	PagePkid  int64     `gorm:"column:page_pkid;type:bigint;not null" json:"page_pkid"`
	UserPkid  *int64    `gorm:"column:user_pkid;type:bigint" json:"user_pkid"`
	Role      string    `gorm:"column:role;type:character varying(20);not null;default:viewer" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()" json:"updated_at"`
	Email     string    `gorm:"column:email;type:character varying(255);not null" json:"email"`
}

// TableName PageRole's table name
func (*PageRole) TableName() string {
	return TableNamePageRole
}
