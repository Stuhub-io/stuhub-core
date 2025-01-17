package postgres

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	"gorm.io/gorm"
)

type PageResult struct {
	model.Page
	Asset        *model.Asset        `gorm:"foreignKey:page_pkid"`
	Doc          *model.Document     `gorm:"foreignKey:page_pkid"`
	Author       *model.User         `gorm:"foreignKey:author_pkid"`
	Organization *model.Organization `gorm:"foreignKey:org_pkid"`
	PageRoles    []model.PageRole    `gorm:"foreignKey:page_pkid"`
}

type initPageModelResults struct {
	Page         *model.Page
	ParentFolder *PageResult
}

func (r *PageRepository) initPageModel(
	tx *gorm.DB,
	pageInput domain.PageInput,
) (initPageModelResults, *domain.Error) {
	path := ""
	// get path
	var parentFolder *PageResult

	if pageInput.ParentPagePkID != nil {
		if err := tx.Where("pkid = ?", pageInput.ParentPagePkID).First(&parentFolder).Error; err != nil {
			return initPageModelResults{}, domain.NewErr("Parent Page not found", domain.BadRequestCode)
		}
		path = pageutils.AppendPath(parentFolder.Path, strconv.FormatInt(parentFolder.Pkid, 10))
	}

	GeneralRole := domain.PageInherit.String()

	// Default General Access for Root Page
	if parentFolder == nil {
		GeneralRole = domain.PageRestrict.String()
	}

	newPage := model.Page{
		Name:           pageInput.Name,
		CoverImage:     pageInput.CoverImage,
		OrgPkid:        &pageInput.OrganizationPkID,
		AuthorPkid:     &pageInput.AuthorPkID,
		ParentPagePkid: pageInput.ParentPagePkID,
		ViewType:       pageInput.ViewType.String(),
		Path:           path,
		GeneralRole:    GeneralRole,
	}
	return initPageModelResults{
		Page:         &newPage,
		ParentFolder: parentFolder,
	}, nil
}

type PreloadPageResultParams struct {
	Asset        bool
	Doc          bool
	Author       bool
	Organization bool
	PageRoles    bool
}

func preloadPageResult(q *gorm.DB, option PreloadPageResultParams) *gorm.DB {
	// Preload Asset and Doc
	if option.Asset {
		q = q.Preload("Asset")
	}
	if option.Doc {
		q = q.Preload("Doc")
	}
	if option.Author {
		q = q.Preload("Author")
	}
	if option.Organization {
		q = q.Preload("Organization")
	}
	if option.PageRoles {
		q = q.Preload("PageRoles")
	}
	return q
}

// Alway include isAll -- to search for multiple (different parents).
func buildPageQuery(
	tx *gorm.DB,
	q domain.PageListQuery,
) *gorm.DB {

	query := tx

	if q.OrgPkID != nil {
		query = query.Where("org_pkid = ?", q.OrgPkID)
	}

	if q.IsArchived != nil {
		if *q.IsArchived {
			query = query.Where("pages.archived_at IS NOT NULL")
		} else {
			query = query.Where("pages.archived_at IS NULL")
		}
	}
	if !q.IsAll {
		if q.ParentPagePkID != nil {
			query = query.Where("pages.parent_page_pkid = ?", *q.ParentPagePkID)
		} else {
			query = query.Where("pages.parent_page_pkid IS NULL")
		}
	}
	if (len(q.ViewTypes)) > 0 {
		query = query.Where("pages.view_type IN ?", q.ViewTypes)
	}

	if q.PathBeginWith != "" {
		query = query.Where("pages.path LIKE ?", q.PathBeginWith+"%")
	}

	query = query.Order("pages.updated_at desc").Offset(q.Offset)
	if q.Limit > 0 {
		query = query.Limit(q.Limit)
	}
	return query
}

func buildOrderByValuesClause(columnName string, ids []int64) string {
	caseStatements := make([]string, len(ids))
	for i, value := range ids {
		caseStatements[i] = fmt.Sprintf("WHEN %d THEN %d", value, i)
	}
	return fmt.Sprintf("CASE %s %s ELSE %d END", columnName, strings.Join(caseStatements, " "), len(ids))
}
