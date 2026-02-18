package questions

type Question struct {
	ID        int64  `json:"id" db:"id"`
	RequestID int64  `json:"request_id" db:"request_id"`
	Question  string `json:"question" db:"question"`
	Answer    string `json:"answer" db:"answer"`
	UserID    int64  `json:"user_id" db:"user_id"`
	Level     string `json:"level" db:"level"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type QuestionInput struct {
	RequestID int64  `json:"request_id" db:"request_id"`
	Question  string `json:"question" db:"question"`
	Answer    string `json:"answer" db:"answer"`
	UserID    int64  `json:"user_id" db:"user_id"`
	Level     string `json:"level" db:"level"`
}

type QuestionResponse struct {
	Text  string `json:"text"`
	Topic string `json:"topic"`
	Level string `json:"level"`
}

type QuestionsResponse struct {
	Questions []QuestionResponse `json:"questions"`
}
