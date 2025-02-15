package domain

type PageStar struct {
	PkID      int64  `json:"pkid"`
	UserPkID  int64  `json:"user_pkid"`
	PagePkID  int64  `json:"page_pkid"`
	Order     int    `json:"order"`
	CreatedAt string `json:"created_at"`
}

type StarPageInput struct {
	ActorUserPkID int64   `json:"user_pkid"`
	PagePkID      int64   `json:"page_pkid"`
	Order         float64 `json:"order"`
}
