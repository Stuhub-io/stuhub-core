package postgres

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
)

func (r *PageRepository) CreatePageRole(
	ctx context.Context,
	createInput domain.PageRoleCreateInput,
) (*domain.PageRoleUser, *domain.Error) {
	pageRole := model.PageRole{
		PagePkid: createInput.PagePkID,
		UserPkid: createInput.UserPkID,
		Role:     createInput.Role.String(),
	}
	if err := r.store.DB().Create(&pageRole).Error; err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	user := model.User{
		Pkid: pageRole.UserPkid,
	}
	if err := r.store.DB().First(&user).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return pageutils.TransformPageRoleModelToDomain(
		pageutils.PageRoleWithUser{
			PageRole: pageRole,
			User:     user,
		},
	), nil
}

func (r *PageRepository) GetOneRoleUserByUserPkId(
	ctx context.Context,
	pagePkID, userPkID int64,
) (*domain.PageRoleUser, *domain.Error) {
	pageRole := model.PageRole{
		Pkid:     pagePkID,
		UserPkid: userPkID,
	}
	if err := r.store.DB().First(&pageRole).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return pageutils.TransformPageRoleModelToDomain(
		pageutils.PageRoleWithUser{
			PageRole: pageRole,
		},
	), nil
}

func (r *PageRepository) GetAllRoleUsersByPagePkId(
	ctx context.Context,
	pagePkID int64,
) ([]domain.PageRoleUser, *domain.Error) {
	var pageRole []pageutils.PageRoleWithUser

	if err := r.store.DB().Where("page_pkid = ?", pagePkID).Preload("User").Find(&pageRole).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var pageRoleUsers []domain.PageRoleUser
	for _, roleUser := range pageRole {
		pageRoleUsers = append(pageRoleUsers, *pageutils.TransformPageRoleModelToDomain(roleUser))
	}

	return pageRoleUsers, nil
}
