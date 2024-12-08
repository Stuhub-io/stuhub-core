package domain

import (
	"encoding/json"
	"errors"
)

type SignUrlInput struct {
	PublicID        string             `json:"public_id"`
	ResourceType    UploadResourceType `json:"resource_type"`
	AdditionalQuery string             `json:"additional_query"`
}

type UploadResourceType string

func (t UploadResourceType) String() string {
	return string(t)
}

func (r *UploadResourceType) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch UploadResourceType(value) {
	case "image", "video", "raw", "auto":
		*r = UploadResourceType(value)
		return nil
	default:
		return errors.New("invalid resource_type, must be image | video | raw")
	}
}

type SignedUrl struct {
	Url       string `json:"url"`
	Signature string `json:"signature"`
	ApiKey    string `json:"api_key"`
	Timestamp string `json:"timestamp"`
}
