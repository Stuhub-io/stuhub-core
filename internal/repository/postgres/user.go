package postgres

import (
	"context"
	"errors"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	store store.DBStore
	cfg   config.Config
}

type NewUserRepositoryParams struct {
	Store store.DBStore
	Cfg   config.Config
}

func NewUserRepository(params NewUserRepositoryParams) ports.UserRepository {
	return &UserRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFoundById(id)
		}
		return nil, domain.ErrDatabaseQuery
	}
	var activatedAt string = ""
	if user.ActivatedAt != nil {
		activatedAt = user.ActivatedAt.String()
	}
	return &domain.User{
		PkID:      user.Pkid,
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,

		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil

}

func (r *UserRepository) GetUserByPkID(ctx context.Context, pkId string) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("pkid = ?", pkId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFoundById(pkId)
		}
		return nil, domain.ErrDatabaseQuery
	}

	var activatedAt string = ""
	if user.ActivatedAt != nil {
		activatedAt = user.ActivatedAt.String()
	}
	return &domain.User{
		PkID:      user.Pkid,
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,

		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFoundByEmail(email)
		}

		return nil, domain.ErrDatabaseQuery
	}
	var activatedAt string = ""
	if user.ActivatedAt != nil {
		activatedAt = user.ActivatedAt.String()
	}
	return &domain.User{
		PkID:      user.Pkid,
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,

		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil

}

func (r *UserRepository) GetOrCreateUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error) {

	var user model.User
	// Try to find the user by email
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDatabaseQuery
		}
		// If the user is not found, create a new user
		user = model.User{
			Email: email,
		}

		err = r.store.DB().Create(&user).Error
		if err != nil {
			return nil, domain.ErrDatabaseQuery
		}
	}

	var activatedAt string = ""
	if user.ActivatedAt != nil {
		activatedAt = user.ActivatedAt.String()
	}
	// User found, return the existing user
	return &domain.User{
		PkID:  user.Pkid,
		ID:    user.ID,
		Email: user.Email,

		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,

		// Socials
		OauthGmail:   user.OathGmail,
		HavePassword: user.Password != nil && *user.Password != "",

		ActivatedAt: activatedAt,
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
	}, nil
}

func (r *UserRepository) SetUserPassword(ctx context.Context, pkID int64, password string) *domain.Error {
	err := r.store.DB().Model(&model.User{}).Where("pkid = ?", pkID).Update("password", password).Error
	if err != nil {
		return domain.ErrDatabaseMutation
	}
	return nil
}
