package ports

import (
	"time"

	"github.com/Stuhub-io/core/domain"
)

type TokenMaker interface {
	CreateToken(pkid int64, email string, duration time.Duration) (string, error)
	DecodeToken(token string) (*domain.TokenAuthPayload, error)
	CreateOrgInviteToken(userPkId, orgPkId int64, duration time.Duration) (string, error)
	DecodeOrgInviteToken(token string) (*domain.TokenOrgInvitePayload, error)
}
