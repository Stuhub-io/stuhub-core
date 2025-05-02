package request

type ActivityPaginationRequest struct {
	EndTime string `form:"end_time" json:"end_time"`
	Limit   int    `form:"limit,default=100" json:"limit"`
}
