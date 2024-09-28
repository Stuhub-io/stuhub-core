package domain

import "time"

type OrganizationInvite struct {
	PkID             int64         `json:"pk_id"`
	ID               string        `json:"id"`
	UserPkID         int64         `json:"user_pkid"`
	OrganizationPkID int64         `json:"organization_pkid"`
	IsUsed           bool          `json:"is_used"`
	Organization     *Organization `json:"organization"`
	CreatedAt        time.Time     `json:"created_at"`
	ExpiredAt        time.Time     `json:"expired_at"`
}

const OrgInvitationExpiredTime time.Duration = time.Minute * 15 //15m
