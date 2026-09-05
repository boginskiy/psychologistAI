package main

import (
	"context"
	"log"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// TODO
// На практике многие проекты идут по пути: OAuth 2.0 + JWT (как access token) + OIDC (для пользовательской информации).
// Это золотой стандарт для современных распределённых систем.

func main() {
	ctx := context.Background()

	// 1. Подключаемся к OIDC провайдеру через discovery
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Настраиваем OAuth2 клиент
	oauth2Config := oauth2.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"}, // "openid" обязателен[reference:6]
	}

	// 3. Создаём верификатор ID Token
	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: "your-client-id"})
}
