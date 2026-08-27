package server

import (
	"context"
	"net/http"
)

type Server interface {
	Run(ctx context.Context, handler http.Handler) error
}
