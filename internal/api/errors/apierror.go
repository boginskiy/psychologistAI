package errors

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(err error, status int) *APIError {
	return &APIError{
		Status:  status,
		Message: err.Error(),
	}
}
