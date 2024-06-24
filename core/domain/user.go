package domain

type User struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar"`
	IsOAuth     bool   `json:"is_oauth"`
	IsActivated bool   `json:"is_activated"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
