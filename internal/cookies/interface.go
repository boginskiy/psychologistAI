package cookies

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
)

type AppCooker interface {
	GetCookiesWithTokens(tokenPair *dto.TokenPair) (access *http.Cookie, refresh *http.Cookie, err error)
}
