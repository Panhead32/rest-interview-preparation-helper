package requests

type Request struct {
	ID          int64  `json:"id" db:"id"`
	RequestText string `json:"request_text" db:"request_text"`
	UserID      int64  `json:"user_id" db:"user_id"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

type RequestInput struct {
	RequestText string `json:"request_text" db:"request_text"`
	UserID      int64  `json:"user_id" db:"user_id"`
}
