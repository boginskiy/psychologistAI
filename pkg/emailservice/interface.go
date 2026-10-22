package emailservice

import gomail "gopkg.in/mail.v2"

type Messager interface {
	GetMsg() *gomail.Message
}

type Postman interface {
	Send(MSGer Messager, retry int)
}
