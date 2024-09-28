package userutils

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
)

func GetUserFullName(f string, l string) string {
	if f == "" {
		return l
	}

	if l == "" {
		return f
	}

	return f + " " + l
}

func TransformUserModelToDomain(model model.User) *domain.User {
	activatedAt := ""
	if model.ActivatedAt != nil {
		activatedAt = model.ActivatedAt.String()
	}

	return &domain.User{
		PkID:         model.Pkid,
		ID:           model.ID,
		Email:        model.Email,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		Avatar:       model.Avatar,
		Salt:         model.Salt,
		OauthGmail:   model.OauthGmail,
		HavePassword: model.Password != nil && *model.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    model.CreatedAt.String(),
		UpdatedAt:    model.UpdatedAt.String(),
	}
}
