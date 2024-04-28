package response

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}
