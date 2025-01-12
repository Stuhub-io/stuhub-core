package postgres

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
)

// Asset Page.
func (r *PageRepository) CreateAsset(ctx context.Context, assetInput domain.AssetPageInput) (*domain.Page, *domain.Error) {

	initPageResult, iErr := r.initPageModel(preloadPageResult(r.store.DB(), PreloadPageResultParams{
		Author: true, // Init Page with Parent Page Preload
	}), assetInput.PageInput)

	if iErr != nil {
		return nil, iErr
	}

	newPage := initPageResult.Page

	author := &model.User{}
	if err := r.store.DB().Where("pkid = ?", newPage.AuthorPkid).First(author).Error; err != nil {
		return nil, domain.ErrBadRequest
	}

	tx, doneTx := r.store.NewTransaction()

	err := preloadPageResult(r.store.DB(), PreloadPageResultParams{
		Author: true,
	}).Create(&newPage).Error

	if err != nil {
		return nil, doneTx(err)
	}

	asset := model.Asset{
		PagePkid:   newPage.Pkid,
		URL:        assetInput.Asset.URL,
		Size:       &assetInput.Asset.Size,
		Extension:  &assetInput.Asset.Extension,
		Thumbnails: assetInput.Asset.Thumbnails.String(),
	}

	rerr := tx.DB().Create(&asset).Error
	if rerr != nil {
		return nil, doneTx(err)
	}

	parentFolder := initPageResult.ParentFolder
	if parentFolder != nil {
		if err := inheritPageRoles(tx.DB(), InheritPageRolesParams{
			ParentFolder:            parentFolder.Page,
			ParentFolderAuthorEmail: parentFolder.Author.Email,
			NewPagePkID:             newPage.Pkid,
			NewPageAuthorPkID:       author.Pkid,
			NewPageAuthorEmail:      author.Email,
		}); err != nil {
			return nil, doneTx(err)
		}
	}

	doneTx(nil)
	// Commit Tx

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: newPage,
			PageBody: pageutils.PageBodyParams{
				Asset: pageutils.TransformAssetModalToDomain(&asset),
			},
		},
	), nil
}
