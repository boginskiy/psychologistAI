package response

type InfoRes struct {
	Status  int    `json:"-"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func NewInfoRes(status int, msg string) *InfoRes {
	return &InfoRes{
		Status:  status,
		Message: msg,
	}
}
