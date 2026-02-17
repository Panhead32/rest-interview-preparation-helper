package model

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewSuccessResponse(message string, data any) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, errMsg string) *Response {
	return &Response{
		Success: false,
		Message: message,
		Error:   errMsg,
	}
}
