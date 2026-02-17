package auth

type Auth struct {
	UserID       int64  `json:"user_id" db:"user_id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

type AuthInput struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
