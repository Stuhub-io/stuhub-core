package postgres

import (
	"context"
	"strconv"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	"gorm.io/gorm"
)

type PageResult struct {
	model.Page
	Asset  *model.Asset    `gorm:"foreignKey:page_pkid"`
	Doc    *model.Document `gorm:"foreignKey:page_pkid"`
	Author *model.User     `gorm:"foreignKey:author_pkid"`
}

type initPageModelResults struct {
	Page         *model.Page
	ParentFolder *model.Page
}

func (r *PageRepository) initPageModel(
	ctx context.Context,
	pageInput domain.PageInput,
) (initPageModelResults, *domain.Error) {
	path := ""
	// get path
	var parentFolder *model.Page

	if pageInput.ParentPagePkID != nil {
		if err := r.store.DB().Where("pkid = ?", pageInput.ParentPagePkID).First(&parentFolder).Error; err != nil {
			return initPageModelResults{}, domain.NewErr("Parent Page not found", domain.BadRequestCode)
		}
		path = pageutils.AppendPath(parentFolder.Path, strconv.FormatInt(parentFolder.Pkid, 10))
	}

	GeneralAccess := false
	GeneralRole := domain.PageInherit.String()

	// Default General Access for Root Page
	if parentFolder == nil {
		GeneralRole = domain.PageViewer.String()
	}

	newPage := model.Page{
		Name:            pageInput.Name,
		CoverImage:      pageInput.CoverImage,
		OrgPkid:         &pageInput.OrganizationPkID,
		AuthorPkid:      &pageInput.AuthorPkID,
		ParentPagePkid:  pageInput.ParentPagePkID,
		ViewType:        pageInput.ViewType.String(),
		Path:            path,
		GeneralRole:     GeneralRole,
		IsGeneralAccess: GeneralAccess,
	}
	return initPageModelResults{
		Page:         &newPage,
		ParentFolder: parentFolder,
	}, nil
}

type PreloadPageResultParams struct {
	Asset  bool
	Doc    bool
	Author bool
}

func (*PageRepository) preloadPageResult(q *gorm.DB, option PreloadPageResultParams) *gorm.DB {
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
	return q
}
