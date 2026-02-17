package users

type User struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	NickName  string `json:"nickname" db:"nickname"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type UserFullInfo struct {
	User
	Email string `json:"email" db:"email"`
}

type UserInput struct {
	Name     string `json:"name" db:"name"`
	Surname  string `json:"surname" db:"surname"`
	NickName string `json:"nickname" db:"nickname"`
}
