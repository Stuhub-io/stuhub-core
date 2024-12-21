package postgres

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	"gorm.io/gorm"
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

func (r *PageRepository) GetPageRoleByUserPkId(
	ctx context.Context,
	pagePkID, userPkID int64,
) (*domain.PageRoleUser, *domain.Error) {
	var pageRole model.PageRole
	if err := r.buildFilterPageRoleQuery(pagePkID, userPkID).First(&pageRole).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return pageutils.TransformPageRoleModelToDomain(
		pageutils.PageRoleWithUser{
			PageRole: pageRole,
		},
	), nil
}

func (r *PageRepository) GetPageRoles(
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

func (r *PageRepository) UpdatePageRole(
	ctx context.Context,
	updateInput domain.PageRoleUpdateInput,
) *domain.Error {
	if err := r.buildFilterPageRoleQuery(updateInput.PagePkID, updateInput.UserPkID).Model(&model.PageRole{}).Update("role", updateInput.Role.String()).Error; err != nil {
		return domain.ErrDatabaseMutation
	}

	return nil
}

func (r *PageRepository) DeletePageRole(
	ctx context.Context,
	updateInput domain.PageRoleDeleteInput,
) *domain.Error {
	if err := r.buildFilterPageRoleQuery(updateInput.PagePkID, updateInput.UserPkID).Delete(&model.PageRole{}).Error; err != nil {
		return domain.ErrDatabaseMutation
	}

	return nil
}

func (r *PageRepository) buildFilterPageRoleQuery(pagePkID, userPkID int64) *gorm.DB {
	return r.store.DB().Where("page_pkid = ?", pagePkID).Where("user_pkid = ?", userPkID)
}
