package dto

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ActivateUserRequest struct {
	TokenPlaintext string `json:"token"`
}
