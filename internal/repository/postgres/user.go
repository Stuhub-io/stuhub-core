package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/userutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	return userutils.TransformUserModelToDomain(&user), nil
}

func (r *UserRepository) GetUserByPkID(ctx context.Context, pkId int64) (*domain.User, *domain.Error) {
	cachedUser := r.store.Cache().GetUser(pkId)
	if cachedUser != nil {
		return cachedUser, nil
	}

	var userModel model.User
	err := r.store.DB().Where("pkid = ?", pkId).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, domain.ErrDatabaseQuery
	}

	user := userutils.TransformUserModelToDomain(&userModel)

	// go func() {
	// 	r.store.Cache().SetUser(user, time.Hour)
	// }()

	return user, nil
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

	return userutils.TransformUserModelToDomain(&user), nil
}

func (r *UserRepository) GetOrCreateUserByEmail(ctx context.Context, email string, salt string) (*domain.User, *domain.Error, bool) {
	var user model.User
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	created := false
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDatabaseQuery, false
		}

		user = model.User{
			Email: email,
			Salt:  salt,
		}

		err = r.store.DB().Create(&user).Error
		created = true
		if err != nil {
			return nil, domain.ErrDatabaseQuery, false
		}
	}

	return userutils.TransformUserModelToDomain(&user), nil, created
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

	return userutils.TransformUserModelToDomain(&user), nil
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

func (r *UserRepository) UpdateUserInfo(ctx context.Context, PkID int64, firstName, lastName, avatar string) (*domain.User, *domain.Error) {
	var user = model.User{
		FirstName: firstName,
		LastName:  lastName,
		Avatar:    avatar,
	}
	err := r.store.DB().Model(&model.User{}).Where("pkid = ?", PkID).Updates(&user).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return userutils.TransformUserModelToDomain(&user), nil
}

func (r *UserRepository) SetUserActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.User, *domain.Error) {
	var user model.User

	err := r.store.DB().Model(&user).Clauses(clause.Returning{}).Where("pkid = ?", pkID).Update("activated_at", activatedAt).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return userutils.TransformUserModelToDomain(&user), nil
}

func (r *UserRepository) Search(ctx context.Context, input domain.UserSearchQuery, currentUser *domain.User) ([]domain.User, *domain.Error) {
	var users []model.User

	query := r.store.DB().Model(&model.User{})
	if input.Search != "" {
		query = query.Where("unaccent(CONCAT(first_name, ' ', last_name)) ILIKE unaccent(?) OR unaccent(email) ILIKE unaccent(?)", "%"+input.Search+"%", "%"+input.Search+"%")
	}

	if currentUser != nil {
		query = query.Where("pkid != ?", currentUser.PkID)
	}

	if input.OrganizationPkID != nil {
		query = query.Joins("JOIN organization_members ON organization_members.user_pkid = users.pkid").Where("organization_members.organization_pkid = ?", *input.OrganizationPkID)
	}

	if len(input.Emails) > 0 {
		query = query.Where("email IN (?)", input.Emails)
	}

	if err := query.Limit(input.Limit).Offset(input.Offset).Find(&users).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	resultUsers := make([]domain.User, 0, len(users))
	for _, user := range users {
		resultUsers = append(resultUsers, *userutils.TransformUserModelToDomain(&user))
	}

	// return domainUsers, nil
	return resultUsers, nil
}

func (r *UserRepository) UnsafeListUsers(ctx context.Context, q domain.UserListQuery) ([]domain.User, *domain.Error) {

	users := []model.User{}
	query := r.store.DB().Model(&model.User{})
	if err := query.Where("pkid in ?", q.UserPkIDs).Find(&users).Error; err != nil {
		return nil, domain.NewErr(err.Error(), domain.ErrDatabaseQuery.Code)
	}

	resultUsers := make([]domain.User, 0, len(users))
	for _, user := range users {
		resultUsers = append(resultUsers, *userutils.TransformUserModelToDomain(&user))
	}

	return resultUsers, nil
}
