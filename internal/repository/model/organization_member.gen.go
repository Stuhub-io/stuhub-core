// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOrganizationMember = "organization_member"

// OrganizationMember mapped from table <organization_member>
type OrganizationMember struct {
	Pkid             int64      `gorm:"column:pkid;type:bigint;primaryKey;autoIncrement:true" json:"pkid"`
	OrganizationPkid int64      `gorm:"column:organization_pkid;type:bigint;not null;index:idx_organization_member_organization_pkid,priority:1" json:"organization_pkid"`
	UserPkid         *int64     `gorm:"column:user_pkid;type:bigint;index:idx_organization_member_user_pkid,priority:1" json:"user_pkid"`
	Role             string     `gorm:"column:role;type:character varying(50);not null" json:"role"`
	ActivatedAt      *time.Time `gorm:"column:activated_at;type:timestamp with time zone" json:"activated_at"`
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()" json:"updated_at"`
}

// TableName OrganizationMember's table name
func (*OrganizationMember) TableName() string {
	return TableNameOrganizationMember
}
