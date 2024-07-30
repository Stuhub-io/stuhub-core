package domain

import "time"

const (
	AccessTokenDuration  = 24 * time.Hour
	RefreshTokenDuration = 24 * 7 * time.Hour
)

const (
	EmailVerificationTokenDuration         = 10 * time.Minute
	NextStepTokenDuration                  = 5 * time.Minute
	OrgInvitationVerificationTokenDuration = 24 * 7 * time.Hour
)

type TokenAuthPayload struct {
	UserPkID  int64     `json:"user_pkid"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type TokenOrgInvitePayload struct {
	UserPkID  int64     `json:"user_pkid"`
	OrgPkID   int64     `json:"org_pkid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
