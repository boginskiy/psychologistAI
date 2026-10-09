package adapters

import (
	"fmt"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/pkg/cooker"
)

func FromTokenPair(cooker cooker.Cooker, tokenPair *dto.TokenPair) (*dto.CookiePair, error) {
	access, err := cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, tokenPair.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("access cookie: %w", err)
	}

	refresh, err := cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, tokenPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh cookie: %w", err)
	}

	return &dto.CookiePair{Access: access, Refresh: refresh}, nil
}
