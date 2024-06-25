package domain

const (
	SmallPageSize  = 10
	MediumPageSize = 20
	LargePageSize  = 50
	SuperLargeSize = 100
	GiantPageSize  = 200
)

type Pagination struct {
	PageSize   int `json:"size"`
	PageOffset int `json:"offset"`
}
