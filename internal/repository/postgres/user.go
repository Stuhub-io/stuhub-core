package postgres

import (
	"context"
	"errors"

	"github.com/Stuhub-io/core/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, *domain.Error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFoundById(id)
		}

		return nil, domain.ErrDatabaseQuery
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFoundByEmail(email)
		}

		return nil, domain.ErrDatabaseQuery
	}

	return &user, nil
}

func (r *UserRepository) GetOrCreateUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error) {

	var user domain.User
	// Try to find the user by email
	err := r.db.Where("email = ?", email).First(&user).Error
	if err == nil {
		// User found, return the existing user
		return &user, nil
	}

	// If the error is not "record not found", return the error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDatabaseQuery
	}

	// User not found, create a new user
	user = domain.User{
		Email: email,
	}
	err = r.db.Create(&user).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return &user, nil
}
