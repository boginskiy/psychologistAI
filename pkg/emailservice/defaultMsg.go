package emailservice

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type VerifLetter struct {
	Email       string
	Token       string
	NameProject string
}

func NewVerifLetter(email, token, nameProject string) *VerifLetter {
	return &VerifLetter{
		Email:       email,
		Token:       token,
		NameProject: nameProject,
	}
}

func (v *VerifLetter) GetMsg() *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", FromEmail)
	msg.SetHeader("To", v.Email)
	msg.SetHeader("Subject", Subject)

	// Ссылка
	verifyLink := Link + v.Token

	// Текст письма с ссылкой для подтверждения
	body := fmt.Sprintf(`
    Здравствуйте!

    Для подтверждения email перейдите по ссылке: %s

    Если вы не регистрировались на нашем сайте, проигнорируйте это письмо.

    С уважением,
    Команда проекта '%s'
	`, verifyLink, v.NameProject)

	msg.SetBody("text/plain", body)
	return msg
}
