package pageutils

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
)

func MapPageModelToDomain(model model.Page) *domain.Page {
	return &domain.Page{
		PkId:           model.Pkid,
		ID:             model.ID,
		SpacePkID:      model.SpacePkid,
		Name:           model.Name,
		ParentPagePkID: model.ParentPagePkid,
		CreatedAt:      model.CreatedAt.String(),
		UpdatedAt:      model.UpdatedAt.String(),
		ViewType:       model.ViewType,
	}
}
