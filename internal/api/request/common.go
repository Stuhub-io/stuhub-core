package request

type PaginationRequest struct {
	Page int64 `form:"page,default=0"       json:"page"`
	Size int64 `form:"size,default=1000000" json:"size"`
}
