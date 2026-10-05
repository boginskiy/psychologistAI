package emailservice

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

const Subject = "Подтверждение email"

type VerifLetter struct {
	Email       string
	VerifyPath  string
	NameProject string
}

func NewVerifLetter(email, verifyPath, nameProject string) *VerifLetter {
	return &VerifLetter{
		Email:       email,
		VerifyPath:  verifyPath,
		NameProject: nameProject,
	}
}

func (v *VerifLetter) GetMsg() *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", FromEmail)
	msg.SetHeader("To", v.Email)
	msg.SetHeader("Subject", Subject)

	// Текст письма с ссылкой для подтверждения
	body := fmt.Sprintf(`
    Здравствуйте!

    Для подтверждения email перейдите по ссылке: %s

    Если вы не регистрировались на нашем сайте, проигнорируйте это письмо.

    С уважением,
    Команда проекта '%s'
	`, v.VerifyPath, v.NameProject)

	msg.SetBody("text/plain", body)
	return msg
}
