package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gomail "gopkg.in/mail.v2"
)

// TODO. Нужно структурированное логирование.
// Сделать профилирование сервиса. Есть узкие места - проверить.

const (
	HostEmail = "smtp.gmail.com"
	PortEmail = 587
	FromEmail = "gophkeeper@gmail.com"
	Password  = "upiplnvviujgnevc"

	Subject = "Подтверждение email"
	Retry   = 3
	TimeOut = 5000 * time.Millisecond
	Delay   = 50 * time.Millisecond
	Link    = "http://localhost:8080/api/v1/user/verify/"
)

type EmailServ struct {
	dialer *gomail.Dialer
	retry  int
	errCh  chan error
}

func NewEmailServ(ctx context.Context) *EmailServ {
	d := gomail.NewDialer(
		HostEmail,
		PortEmail,
		FromEmail,
		Password,
	)

	tmpServ := &EmailServ{
		dialer: d,
		retry:  Retry,
		errCh:  make(chan error, 1),
	}

	go tmpServ.CatchingErrors(ctx)
	return tmpServ
}

func (s *EmailServ) CatchingErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-s.errCh:
			if !ok {
				return
			}
			// TODO. Записать событие в инфраструктурное логирование.
			log.Printf("email service error: %s", err.Error())
		}
	}
}

func (s *EmailServ) Send(email, token string) {
	go func() {
		// Context
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		ctxTimeOut, cancel := context.WithTimeout(ctx, TimeOut)

		// Recover
		defer func() {
			if r := recover(); r != nil {
				s.sendError(fmt.Errorf("panic recovered: %v", r))
			}
			cancel()
			stop()
		}()

		if s.retry <= 0 {
			s.sendError(fmt.Errorf("retry must be > 0"))
			return
		}

		msg := s.CreateMess(email, token)
		done := make(chan struct{}, 1)

		go func() {
			defer close(done)
			defer func() {
				if r := recover(); r != nil {
					s.sendError(fmt.Errorf("send panic in retry loop: %v", r))
				}
			}()

			var lastErr error
			delay := Delay

			for i := 0; i < s.retry; i++ {
				select {
				case <-ctxTimeOut.Done():
					s.sendError(fmt.Errorf("context cancelled before attempt %d", i+1))
					return
				default:
				}

				lastErr = s.dialer.DialAndSend(msg)
				if lastErr == nil {
					return
				}
				if i == s.retry-1 {
					break
				}

				// Ожидание с учётом контекста
				timer := time.NewTimer(delay)
				select {
				case <-ctxTimeOut.Done():
					timer.Stop()
					s.sendError(fmt.Errorf("context cancelled while waiting"))
					return
				case <-timer.C:
				}
				delay *= 2
			}

			// Если дошли сюда – все попытки провалились
			if lastErr != nil {
				s.sendError(fmt.Errorf("failed to send email after %d attempts: %w", s.retry, lastErr))
			}
		}()

		// Ожидаем завершения либо отмены контекста
		select {
		case <-ctxTimeOut.Done():
			s.sendError(fmt.Errorf("overall timeout exceeded"))
		case <-done:
		}
	}()
}

func (s *EmailServ) sendError(err error) {
	select {
	case s.errCh <- err:
	default:
	}
}

func (s *EmailServ) CreateMess(ToEmail, token string) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", FromEmail)
	msg.SetHeader("To", ToEmail)
	msg.SetHeader("Subject", Subject)

	// Ссылка
	verifyLink := Link + token

	// Текст письма с ссылкой для подтверждения
	body := fmt.Sprintf(`
    Здравствуйте!

    Для подтверждения email перейдите по ссылке: %s

    Если вы не регистрировались на нашем сайте, проигнорируйте это письмо.

    С уважением,
    Команда проекта 'Psychologist AI'
	`, verifyLink)

	msg.SetBody("text/plain", body)
	return msg
}
