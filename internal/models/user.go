package models

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	IsAdm    bool   `json:"is_adm"`
}

type AddUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	IsAdm    bool   `json:"is_adm"`
}

type AddUserResponse struct {
	ID int64 `json:"id"`
}
