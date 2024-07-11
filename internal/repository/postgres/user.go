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
	store *store.DBStore
	cfg   config.Config
}

type NewUserRepositoryParams struct {
	Store *store.DBStore
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

		Salt:         user.Salt,
		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

func (r *UserRepository) GetUserByPkID(ctx context.Context, pkId int64) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("pkid = ?", pkId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
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

		Salt:         user.Salt,
		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		OauthGmail:   user.OauthGmail,
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

		Salt:         user.Salt,
		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil

}

func (r *UserRepository) GetOrCreateUserByEmail(ctx context.Context, email string, salt string) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDatabaseQuery
		}

		user = model.User{
			Email: email,
			Salt:  salt,
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

	return &domain.User{
		PkID:         user.Pkid,
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Avatar:       user.Avatar,
		Salt:         user.Salt,
		OauthGmail:   user.OauthGmail,
		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

func (r *UserRepository) CreateUserWithGoogleInfo(ctx context.Context, email, salt, firstName, lastName, avatar string) (*domain.User, *domain.Error) {
	user := model.User{
		Email:      email,
		Salt:       salt,
		FirstName:  firstName,
		LastName:   lastName,
		Avatar:     avatar,
		OauthGmail: email,
	}

	err := r.store.DB().Create(&user).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	activatedAt := ""
	if user.ActivatedAt != nil {
		activatedAt = user.ActivatedAt.String()
	}

	return &domain.User{
		PkID:         user.Pkid,
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Avatar:       user.Avatar,
		Salt:         user.Salt,
		OauthGmail:   user.OauthGmail,
		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  activatedAt,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

func (r *UserRepository) SetUserPassword(ctx context.Context, pkID int64, hashedPassword string) *domain.Error {
	// FIXME: Add password hashing
	err := r.store.DB().Model(&model.User{}).Where("pkid = ?", pkID).Update("password", hashedPassword).Error
	if err != nil {
		return domain.ErrDatabaseMutation
	}

	return nil
}

func (r *UserRepository) CheckPassword(ctx context.Context, email, rawPassword string, hasher ports.Hasher) (bool, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, domain.ErrUserNotFoundByEmail(email)
		}

		return false, domain.ErrDatabaseQuery
	}

	valid := hasher.Compare(rawPassword, *user.Password, user.Salt)

	return valid, nil
}

func (r *UserRepository) UpdateUserInfo(ctx context.Context, PkID int64, firstName, lastName string) (*domain.User, *domain.Error) {
	var user model.User = model.User{
		Pkid: PkID,
	}
	err := r.store.DB().Model(&user).Where("pkid = ?", PkID).Updates(map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
	}).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	

	return &domain.User{
		PkID:      user.Pkid,
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,

		Salt:         user.Salt,
		HavePassword: user.Password != nil && *user.Password != "",
		ActivatedAt:  user.ActivatedAt.String(),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil

}