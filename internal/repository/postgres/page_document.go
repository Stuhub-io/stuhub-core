package postgres

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	"gorm.io/gorm/clause"
)

func (r *PageRepository) CreateDocumentPage(
	ctx context.Context,
	pageInput domain.DocumentPageInput,
) (*domain.Page, *domain.Error) {
	result, iErr := r.initPageModel(ctx, pageInput.PageInput)
	newPage := result.Page

	if iErr != nil {
		return nil, iErr
	}

	if pageInput.Document.JsonContent == "" {
		pageInput.Document.JsonContent = "{}"
	}

	// Begin Tx
	tx, doneTx := r.store.NewTransaction()

	err := tx.DB().Create(&newPage).Error
	if err != nil {
		return nil, doneTx(err)
	}

	document := model.Document{
		JSONContent: &pageInput.Document.JsonContent,
		PagePkid:    newPage.Pkid,
	}

	rerr := tx.DB().Create(&document).Error
	if rerr != nil {
		return nil, doneTx(err)
	}

	// Inherit Parent Permission
	parentFolder := result.ParentFolder
	if parentFolder != nil {
		if err := inheritPageRoles(tx.DB(), InheritPageRolesParams{
			ParentFolder: *parentFolder,
			NewPagePkID:  newPage.Pkid,
		}); err != nil {
			return nil, doneTx(err)
		}
	}

	doneTx(nil)

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: newPage,
			PageBody: pageutils.PageBodyParams{
				Document: pageutils.TransformDocModelToDomain(&document),
			},
		},
	), nil
}

func (r *PageRepository) UpdateContent(
	ctx context.Context,
	pagePkID int64,
	content domain.DocumentInput,
) (*domain.Page, *domain.Error) {
	var page = model.Page{}
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}

	var doc model.Document
	if dbErr := r.store.DB().Where("page_pkid = ?", pagePkID).First(&doc).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}
	if content.JsonContent == "" {
		content.JsonContent = "{}"
	}
	doc.JSONContent = &content.JsonContent
	if dbErr := r.store.DB().Clauses(clause.Returning{}).Select("*").Save(&doc).Error; dbErr != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: &page,
			PageBody: pageutils.PageBodyParams{
				Document: pageutils.TransformDocModelToDomain(&doc),
			},
		},
	), nil
}
