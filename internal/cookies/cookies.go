package cookies

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/pkg/cooker"
)

type AppCookies struct {
	Cfg    config.Config
	Cooker cooker.Cooker
}

func NewAppCookies(cfg config.Config, cookie cooker.Cooker) *AppCookies {
	return &AppCookies{
		Cfg:    cfg,
		Cooker: cookie,
	}
}

func (c *AppCookies) GetCookiesWithTokens(tokenPair *dto.TokenPair) (access *http.Cookie, refresh *http.Cookie, err error) {
	access, err = c.Cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, tokenPair.AccessToken)
	if err != nil {
		return nil, nil, err
	}

	refresh, err = c.Cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, tokenPair.RefreshToken)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}
