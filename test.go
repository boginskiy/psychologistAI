package main

import (
	"fmt"
	"time"

	gomail "gopkg.in/mail.v2"
)

const (
	HostEmail = "smtp.gmail.com"
	PortEmail = 587
	FromEmail = "gophkeeper@gmail.com"
	Password  = "upiplnvviujgnevc"

	Subject = "Подтверждение email"
	Retry   = 3
	TimeOut = 200 * time.Millisecond
)

func CreateMess(ToEmail, token string) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", FromEmail)
	msg.SetHeader("To", ToEmail)
	msg.SetHeader("Subject", Subject)

	// Ссылка
	verifyLink := fmt.Sprintf("localhost:8080/api/v1/user/verify/%s", token)

	// Текст письма с ссылкой для подтверждения
	body := fmt.Sprintf(`
    Здравствуйте!

    Для подтверждения email перейдите по ссылке:
    %s

    Если вы не регистрировались на нашем сайте, проигнорируйте это письмо.

    С уважением,
    Команда проекта 'Psychologist AI'
	`, verifyLink)

	msg.SetBody("text/plain", body)
	return msg
}

func main() {

	dialer := gomail.NewDialer(
		HostEmail,
		PortEmail,
		FromEmail,
		Password,
	)

	msg := CreateMess("1.boginskiy@mail.ru", "<<1323456789>>")
	fmt.Println(msg)

	err := dialer.DialAndSend(msg)
	fmt.Println(err)

}
