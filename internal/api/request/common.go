package request

type PaginationRequest struct {
	Page int64 `form:"page,default=0"      json:"page"`
	Size int64 `form:"size,default=100000" json:"size"`
}
