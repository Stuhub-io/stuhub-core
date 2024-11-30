package domain

import (
	"encoding/json"
)

type Asset struct {
	PkID       int64          `json:"pkid"`
	PagePkID   int64          `json:"page_pkid"`
	UpdatedAt  string         `json:"updated_at"`
	CreatedAt  string         `json:"created_at"`
	URL        string         `json:"url"`
	Size       int64          `json:"size"`
	Extension  string         `json:"extension"`
	Thumbnails AssetThumbnail `json:"thumbnails"`
}

type AssetThumbnail struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

func AssetThumbnailFromString(val string) AssetThumbnail {
	var asset AssetThumbnail
	json.Unmarshal([]byte(val), &asset)
	return asset
}

func (a *AssetThumbnail) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

type AssetInput struct {
	URL        string         `json:"url"`
	Size       int64          `json:"size"`
	Extension  string         `json:"extension"`
	Thumbnails AssetThumbnail `json:"thumbnails"`
}

type AssetPageInput struct {
	PageInput
	Asset AssetInput `json:"asset"`
}
