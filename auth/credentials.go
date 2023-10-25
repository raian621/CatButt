package auth

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegistrationInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
