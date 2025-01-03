package pageutils

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/userutils"
)

func TransformDocModelToDomain(doc *model.Document) *domain.Document {
	if doc == nil {
		return nil
	}
	jsonContent := ""
	if doc.JSONContent != nil {
		jsonContent = *doc.JSONContent
	}
	return &domain.Document{
		PkID:        doc.Pkid,
		PagePkID:    doc.PagePkid,
		Content:     doc.Content,
		JsonContent: jsonContent,
		CreatedAt:   doc.CreatedAt.String(),
		UpdatedAt:   doc.UpdatedAt.String(),
	}
}

type PageBodyParams struct {
	Document *domain.Document
	Asset    *domain.Asset
	Author   *domain.User
}

type PageModelToDomainParams struct {
	Page            *model.Page
	ChildPages      []domain.Page
	PageBody        PageBodyParams
	InheritFromPage *domain.Page
	Permissions     *domain.PageRolePermissions
}

func TransformPageModelToDomain(
	params PageModelToDomainParams,
) *domain.Page {
	model := params.Page
	childPages := params.ChildPages
	pageBody := params.PageBody
	inheritFromPage := params.InheritFromPage
	Permissions := params.Permissions

	if model == nil {
		return nil
	}
	archivedAt := ""
	if model.ArchivedAt != nil {
		archivedAt = model.ArchivedAt.String()
	}
	nodeID := ""
	if model.NodeID != nil {
		nodeID = *model.NodeID
	}

	return &domain.Page{
		PkID:             model.Pkid,
		ID:               model.ID,
		OrganizationPkID: *model.OrgPkid,
		AuthorPkID:       model.AuthorPkid,
		Name:             model.Name,
		ParentPagePkID:   model.ParentPagePkid,
		CreatedAt:        model.CreatedAt.String(),
		UpdatedAt:        model.UpdatedAt.String(),
		ViewType:         domain.PageViewFromString(model.ViewType),
		CoverImage:       model.CoverImage,
		ArchivedAt:       archivedAt,
		NodeID:           nodeID,
		ChildPages:       childPages,
		Document:         pageBody.Document,
		Asset:            pageBody.Asset,
		Path:             model.Path,
		GeneralRole:      domain.PageRoleFromString(model.GeneralRole),
		Author:           pageBody.Author,
		InheritFromPage:  inheritFromPage,
		Permissions:      Permissions,
	}
}

type PageRoleWithUser struct {
	model.PageRole
	User            *model.User  `gorm:"foreignKey:user_pkid" json:"user"` // Define foreign key relationship
	InheritFromPage *domain.Page `gorm:"-"                    json:"inherit_from_page"`
	Page            *model.Page  `gorm:"foreignKey:page_pkid" json:"page"` // Define foreign key relationship
}

func TransformPageRoleModelToDomain(
	model PageRoleWithUser,
) *domain.PageRoleUser {
	return &domain.PageRoleUser{
		PkID:            model.Pkid,
		PagePkID:        model.PagePkid,
		User:            userutils.TransformUserModelToDomain(model.User),
		Email:           model.Email,
		Role:            domain.PageRoleFromString(model.Role),
		CreatedAt:       model.CreatedAt.String(),
		UpdatedAt:       model.UpdatedAt.String(),
		InheritFromPage: model.InheritFromPage,
	}
}

func TransformAssetModalToDomain(asset *model.Asset) *domain.Asset {
	if asset == nil {
		return nil
	}
	extension := ""
	if asset.Extension != nil {
		extension = *asset.Extension
	}

	size := int64(0)
	if asset.Size != nil {
		size = *asset.Size
	}

	return &domain.Asset{
		PkID:       asset.Pkid,
		PagePkID:   asset.PagePkid,
		URL:        asset.URL,
		CreatedAt:  asset.CreatedAt.String(),
		UpdatedAt:  asset.UpdatedAt.String(),
		Thumbnails: domain.AssetThumbnailFromString(asset.Thumbnails),
		Extension:  extension,
		Size:       size,
	}
}

func TransformPagePublicTokenModelToDomain(model model.PublicToken) *domain.PagePublicToken {
	archivedAt := ""
	if model.ArchivedAt != nil {
		archivedAt = model.ArchivedAt.String()
	}
	return &domain.PagePublicToken{
		PkID:       model.Pkid,
		ArchivedAt: archivedAt,
		PagePkID:   model.PagePkid,
		ID:         model.ID,
		CreatedAt:  model.CreatedAt.String(),
	}
}
