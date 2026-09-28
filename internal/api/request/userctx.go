package request

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
)

type infoUserKey struct{}

func SetInfoUserToContext(ctx context.Context, infoUser dto.InfoUser) context.Context {
	return context.WithValue(ctx, infoUserKey{}, infoUser)
}

func GetInfoUserFromContext(ctx context.Context) (dto.InfoUser, bool) {
	infoUser, ok := ctx.Value(infoUserKey{}).(dto.InfoUser)
	return infoUser, ok
}
