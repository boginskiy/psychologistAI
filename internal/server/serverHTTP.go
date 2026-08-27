package server

import (
	"context"
	"net/http"
)

type ServerHTTP struct {
	Server *http.Server
}

func NewServerHTTP(ctx context.Context, cfg Config) *ServerHTTP {
	return &ServerHTTP{
		Server: &http.Server{
			Addr: cfg.Port},
	}
}

func (s *ServerHTTP) Run(ctx context.Context, handler http.Handler) error {
	s.Server.Handler = handler
	return s.Server.ListenAndServe()
}

// TODO
// Сервер со всеми возможными фичами.

// Graceful shutdown — корректное завершение с ожиданием активных соединений
// Middleware — логирование и обработка паник
// REST API — CRUD операции с пользователями
// Health/Ready probes — для Kubernetes
// Таймауты — Read/Write/Idle
// Отслеживание соединений — для graceful shutdown
// Роутинг — с использованием новых возможностей Go 1.22+ (PathValue)
// Health check
