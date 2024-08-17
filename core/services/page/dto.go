package page

import "github.com/Stuhub-io/core/domain"

type CreatePageDto struct {
	SpacePkID      int64
	ViewType       domain.PageViewType
	ParentPagePkID *int64
	Name           string
}
