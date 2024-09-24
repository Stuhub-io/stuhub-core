package domain

import "time"

type OrganizationInvite struct {
	PkId             int64  `json:"pk_id"`
	ID               string `json:"id"`
	UserPkID         int64  `json:"user_pkid"`
	OrganizationPkID int64  `json:"organization_pkid"`
	CreatedAt        string `json:"created_at"`
	ExpiredAt        string `json:"expired_at"`
}

const OrgInvitationExpiredTime time.Duration = time.Minute * 15 //15m
