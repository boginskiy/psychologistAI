package emailservice

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

	Retry   = 3
	TimeOut = 5000 * time.Millisecond
	Delay   = 50 * time.Millisecond
)

type EmailService struct {
	dialer *gomail.Dialer
	errCh  chan error
}

func NewEmailService(ctx context.Context) *EmailService {
	d := gomail.NewDialer(
		HostEmail,
		PortEmail,
		FromEmail,
		Password,
	)

	tmpServ := &EmailService{
		dialer: d,
		errCh:  make(chan error, 1),
	}

	go tmpServ.catchingErrors(ctx)
	return tmpServ
}

func (s *EmailService) catchingErrors(ctx context.Context) {
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

func (s *EmailService) Send(MSGer Messager, retry int) {
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

		if retry <= 0 {
			s.sendError(fmt.Errorf("retry must be > 0"))
			return
		}

		// Messager
		msg := MSGer.GetMsg()
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

			for i := 0; i < retry; i++ {
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
				if i == retry-1 {
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
				s.sendError(fmt.Errorf("failed to send email after %d attempts: %w", retry, lastErr))
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

func (s *EmailService) sendError(err error) {
	select {
	case s.errCh <- err:
	default:
		// Игнорируем ошибки, если переполнен канал.
	}
}
