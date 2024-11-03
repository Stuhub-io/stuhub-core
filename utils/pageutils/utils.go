package pageutils

import (
	"strconv"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/gin-gonic/gin"
)

const PagePkIDParam = "pagePkID"
const PageIDParam = "pageID"

func GetPageIDParam(c *gin.Context) (string, bool) {
	pageID := c.Params.ByName(PageIDParam)
	if pageID == "" {
		return "", false
	}
	return pageID, true
}

func GetPagePkIDParam(c *gin.Context) (int64, bool) {
	pagePkID := c.Params.ByName(PagePkIDParam)
	if pagePkID == "" {
		return int64(-1), false
	}
	pkID, cErr := strconv.Atoi(pagePkID)
	return int64(pkID), cErr == nil
}

func MapPageModelToDomain(model model.Page, ChildPages []domain.Page) *domain.Page {
	archivedAt := ""
	if model.ArchivedAt != nil {
		archivedAt = model.ArchivedAt.String()
	}
	nodeID := ""
	if model.NodeID != nil {
		nodeID = *model.NodeID
	}
	return &domain.Page{
		PkID:           model.Pkid,
		ID:             model.ID,
		SpacePkID:      model.SpacePkid,
		Name:           model.Name,
		ParentPagePkID: model.ParentPagePkid,
		CreatedAt:      model.CreatedAt.String(),
		UpdatedAt:      model.UpdatedAt.String(),
		ViewType:       model.ViewType,
		CoverImage:     model.CoverImage,
		ArchivedAt:     archivedAt,
		NodeID:         nodeID,
		ChildPages:     ChildPages,
	}
}
