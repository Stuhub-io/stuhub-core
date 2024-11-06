package domain

const (
	SmallPageSize  = 10
	MediumPageSize = 20
	LargePageSize  = 50
	SuperLargeSize = 100
	GiantPageSize  = 200
)

type Pagination struct {
	Size int64 `json:"size"`
	Page int64 `json:"page"`
}
