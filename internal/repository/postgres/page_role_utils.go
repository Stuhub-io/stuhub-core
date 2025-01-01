package postgres

import (
	"slices"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	sliceutils "github.com/Stuhub-io/utils/slice"
	"gorm.io/gorm"
)

type InheritPageRolesParams struct {
	ParentFolder model.Page
	NewPagePkID  int64
}

func inheritPageRoles(tx *gorm.DB, input InheritPageRolesParams) error {
	var parentPageRoles []model.PageRole
	if err := tx.Where("page_pkid = ?", input.ParentFolder.Pkid).Find(&parentPageRoles).Error; err != nil {
		return err
	}
	if len(parentPageRoles) == 0 {
		return nil
	}

	// Inherit direct parentPageRoles from parent
	pageRoles := make([]model.PageRole, 0, len(parentPageRoles))
	for _, permission := range parentPageRoles {
		pageRoles = append(pageRoles, model.PageRole{
			PagePkid: input.NewPagePkID,
			Email:    permission.Email,
			UserPkid: permission.UserPkid,
			Role:     domain.PageInherit.String(),
		})
	}

	if err := tx.Create(&pageRoles).Error; err != nil {
		return err
	}

	return nil
}

type PageRoleResult struct {
	model.PageRole
	User *model.User `gorm:"foreignKey:user_pkid" json:"user"` // Define foreign key relationship
	Page *model.Page `gorm:"foreignKey:page_pkid" json:"page"` // Define foreign key relationship
}

type queryPageRolesPreloadOption struct {
	User bool
	Page bool
}
type queryPageRolesParams struct {
	PagePkIDs    []int64
	Emails       []string
	Roles        []domain.PageRole
	Preload      queryPageRolesPreloadOption
	ExcludeRoles []domain.PageRole
}

func queryPageRoles(tx *gorm.DB, params queryPageRolesParams) ([]PageRoleResult, *domain.Error) {
	pageRoles := make([]PageRoleResult, 0, len(params.PagePkIDs))

	if len(params.PagePkIDs) == 0 && len(params.Emails) == 0 && len(params.Roles) == 0 {
		return pageRoles, nil
	}
	query := buildQueryPageRoles(tx, params)

	if err := query.Find(&pageRoles).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	slices.SortFunc(pageRoles, func(i, j PageRoleResult) int {
		return slices.Index(params.PagePkIDs, i.Pkid) - slices.Index(params.PagePkIDs, j.Pkid)
	})

	return pageRoles, nil
}

func buildQueryPageRoles(tx *gorm.DB, params queryPageRolesParams) *gorm.DB {
	PagePkIDs := params.PagePkIDs
	Emails := params.Emails
	Roles := params.Roles
	Preload := params.Preload
	ExcludeRoles := params.ExcludeRoles

	if len(PagePkIDs) == 0 && len(Emails) == 0 && len(Roles) == 0 {
		return tx
	}

	// Preload
	if Preload.User {
		tx = tx.Preload("User")
	}
	if Preload.Page {
		tx = tx.Preload("Page")
	}

	// tx
	if len(PagePkIDs) != 0 {
		if len(PagePkIDs) == 1 {
			tx = tx.Where("page_pkid = ?", PagePkIDs[0])
		} else {
			tx = tx.Where("page_pkid IN (?)", PagePkIDs)
		}
	}

	if len(Roles) != 0 {
		if len(Roles) == 1 {
			tx = tx.Where("role = ?", Roles[0].String())
		} else {
			tx = tx.Where("role IN (?)", sliceutils.Map(Roles, func(role domain.PageRole) string {
				return role.String()
			}))
		}
	}

	if len(ExcludeRoles) != 0 {
		if len(ExcludeRoles) == 1 {
			tx = tx.Where("role != ?", ExcludeRoles[0].String())
		} else {
			tx = tx.Where("role NOT IN (?)", sliceutils.Map(ExcludeRoles, func(role domain.PageRole) string {
				return role.String()
			}))
		}
	}

	if len(Emails) != 0 {
		if len(Emails) == 1 {
			tx = tx.Where("email = ?", Emails[0])
		} else {
			tx = tx.Where("email IN (?)", Emails)
		}
	}
	return tx
}
