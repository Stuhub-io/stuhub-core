package request

import "github.com/Stuhub-io/core/domain"

type SignUrlRequestBody struct {
	PublicID        string                    `binding:"required"      json:"public_id"`
	ResourceType    domain.UploadResourceType `binding:"required"      json:"resource_type"`
	AdditionalQuery string                    `json:"additional_query"`
}
