package entities

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	IsAdm    bool   `json:"is_adm"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	IsAdm    bool   `json:"is_adm"`
}

type UserResponse struct {
	ID int64 `json:"id"`
}
