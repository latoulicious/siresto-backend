package utils

type StandardResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(message string, data interface{}) StandardResponse {
	return StandardResponse{
		Message: message,
		Status:  200,
		Data:    data,
	}
}

func Error(message string, status int) StandardResponse {
	return StandardResponse{
		Message: message,
		Status:  status,
	}
}
