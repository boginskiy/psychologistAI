package response

type MessRes struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func NewMessRes(msg string, status int) *MessRes {
	return &MessRes{Status: status, Message: msg}
}

func (r *MessRes) GetStatus() int {
	return r.Status
}

func (r *MessRes) UpdateErr(err error, status int) {
	r.Message = err.Error()
	r.Status = status
}

func (r *MessRes) UpdateInfo(info string, status int) {
	r.Message = info
	r.Status = status
}
