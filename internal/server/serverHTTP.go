package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type ServerHTTP struct {
	Cfg    Config
	Server *http.Server
}

func NewServerHTTP(ctx context.Context, cfg Config) *ServerHTTP {
	return &ServerHTTP{
		Cfg: cfg,
		Server: &http.Server{
			Addr:         cfg.Port,         // Порт.
			ReadTimeout:  cfg.ReadTimeout,  // Если клиент медленно отправляет данные, соединение будет закрыто.
			WriteTimeout: cfg.WriteTimeout, // Ограничивает время время выполнения обработчика + отправку данных клиенту.
			IdleTimeout:  cfg.IdleTimeout,  // Время, в течение которого соединение может оставаться открытым без новых запросов (Keep-Alive)

			// Extra settings:
			// ReadHeaderTimeout: 2 * time.Second, // Ограничивает время на чтение HTTP-заголовков. Полезно для защиты от медленных атак (Slowloris)
			// MaxHeaderBytes: 1 << 20,            // Ограничивает размер заголовков запроса. Защита от DoS-атак с огромными заголовками.
		},
	}
}

func (s *ServerHTTP) Run(ctx context.Context, handler http.Handler) error {
	if handler == nil {
		return fmt.Errorf("Handler cannot be nil")
	}
	s.Server.Handler = handler

	errCh := make(chan error, 1)
	go func() {
		// s.logger.Debug().Msg("Server started")
		errCh <- s.Server.ListenAndServe()
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.Cfg.ShutdownTimeout)
		defer cancel()

		// s.logger.Debug().Msg("Shutdown started")

		shutdownErr := s.Server.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			return shutdownErr
		}

		serverErr := <-errCh
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			return serverErr
		}

		// s.logger.Debug().Msg("Server stoped")
		return nil

	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// TODO
// Сервер со всеми возможными фичами.

// Graceful shutdown — корректное завершение с ожиданием активных соединений
// Middleware — логирование и обработка паник
// REST API — CRUD операции с пользователями
// Health/Ready probes — для Kubernetes
// Отслеживание соединений — для graceful shutdown
// Роутинг — с использованием новых возможностей Go 1.22+ (PathValue)
// Health check
