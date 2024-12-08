package postgres

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	"gorm.io/gorm"
)

type PageResult struct {
	model.Page
	Asset *model.Asset    `gorm:"foreignKey:page_pkid"`
	Doc   *model.Document `gorm:"foreignKey:page_pkid"`
}

func (r *PageRepository) initPageModel(ctx context.Context, pageInput domain.PageInput) (*model.Page, *domain.Error) {
	path := ""
	// get path
	if pageInput.ParentPagePkID != nil {
		var parentPage model.Page
		if err := r.store.DB().Where("pkid = ?", pageInput.ParentPagePkID).First(&parentPage).Error; err != nil {
			return nil, domain.NewErr("Parent Page not found", domain.BadRequestCode)
		}
		path = pageutils.AppendPath(parentPage.Path, parentPage.ID)
	}

	newPage := model.Page{
		Name:           pageInput.Name,
		CoverImage:     pageInput.CoverImage,
		OrgPkid:        &pageInput.OrganizationPkID,
		ParentPagePkid: pageInput.ParentPagePkID,
		ViewType:       pageInput.ViewType.String(),
		Path:           path,
	}
	return &newPage, nil
}

func (*PageRepository) preloadPageResult(q *gorm.DB) *gorm.DB {
	// Preload Asset and Doc
	return q.Preload("Asset").Preload("Doc")
}
