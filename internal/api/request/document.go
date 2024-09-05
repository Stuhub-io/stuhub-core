package request

type CreateDocumentBody struct {
	Page        CreatePageBody `json:"page" binding:"required"`
	JsonContent string         `json:"json_content" binding:"omitempty"`
}
