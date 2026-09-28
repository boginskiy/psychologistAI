package response

type InfoBody struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func NewInfoBodyWithErr(err error, status int) *InfoBody {
	return &InfoBody{
		Status:  status,
		Message: err.Error(),
	}
}

func NewInfoBody(msg string, status int) *InfoBody {
	return &InfoBody{
		Status:  status,
		Message: msg,
	}
}

func (r *InfoBody) GetStatus() int {
	return r.Status
}

func (r *InfoBody) ErrorUpdate(err error, status int) {
	r.Status = status
	r.Message = err.Error()
}

func (r *InfoBody) InfoUpdate(msg string, status int) {
	r.Status = status
	r.Message = msg
}
