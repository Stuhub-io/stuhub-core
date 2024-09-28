package docutils

import (
	"strconv"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/gin-gonic/gin"
)

const (
	DocumentPkIDParam = "documentPkID"
)

func GetDocumentParams(c *gin.Context) (int64, bool) {
	documentPkID := c.Params.ByName(DocumentPkIDParam)
	if documentPkID == "" {
		return int64(-1), false
	}
	docPkID, cErr := strconv.Atoi(documentPkID)
	return int64(docPkID), cErr == nil
}

func TransformDocModalToDomain(doc model.Document) *domain.Document {
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
