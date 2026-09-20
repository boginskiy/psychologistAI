package response

type InfoResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (r *InfoResponse) GetStatus() int {
	return r.Status
}

func (r *InfoResponse) ErrorUpdate(err error, status int) {
	r.Status = status
	r.Message = err.Error()
}

func (r *InfoResponse) InfoUpdate(msg string, status int) {
	r.Status = status
	r.Message = msg
}
