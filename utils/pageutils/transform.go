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
}

func TransformPageModelToDomain(
	model model.Page,
	ChildPages []domain.Page,
	pageBody PageBodyParams,
) *domain.Page {
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
		ChildPages:       ChildPages,
		Document:         pageBody.Document,
		Asset:            pageBody.Asset,
		Path:             model.Path,
		IsGeneralAccess:  model.IsGeneralAccess,
		GeneralRole:      model.GeneralRole,
	}
}

type PageRoleWithUser struct {
	model.PageRole
	User model.User `gorm:"foreignKey:user_pkid" json:"user"` // Define foreign key relationship
}

func TransformPageRoleModelToDomain(
	model PageRoleWithUser,
) *domain.PageRoleUser {
	return &domain.PageRoleUser{
		PkID:      model.Pkid,
		PagePkID:  model.PagePkid,
		User:      *userutils.TransformUserModelToDomain(model.User),
		Role:      model.Role,
		CreatedAt: model.CreatedAt.String(),
		UpdatedAt: model.UpdatedAt.String(),
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
