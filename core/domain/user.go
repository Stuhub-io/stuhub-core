package domain

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
