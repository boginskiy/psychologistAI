package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/security"

	models "github.com/boginskiy/psychologistAI/internal/models/user"
	"github.com/boginskiy/psychologistAI/internal/repository"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/boginskiy/psychologistAI/pkg/jwtservice"
)

const AttemptsCnt = 5

type AuthServ struct {
	Validater   Validater
	Notifier    Notifier
	UserRepo    repository.UserRepo
	SessionRepo repository.SessionRepo
	JWTManager  jwtservice.JWTManager
	GeoSecurity security.GeoSecurity
}

func NewAuthServ(
	ctx context.Context,
	validater Validater,
	notifier Notifier,
	userRepo repository.UserRepo,
	SessionRepo repository.SessionRepo,
	jwtManager jwtservice.JWTManager,
	geoSecurity security.GeoSecurity,
) *AuthServ {
	return &AuthServ{
		Validater:   validater,
		Notifier:    notifier,
		UserRepo:    userRepo,
		JWTManager:  jwtManager,
		GeoSecurity: geoSecurity,
	}
}

func (s *AuthServ) Refresh(ctx context.Context, refreshTokenReq *dto.RefreshTokenRequest) (*dto.TokenPair, error) {
	// TODO
	// В мидлваре проверим подпись токена, и сделаем парсинг данных.
	// Через контекст можно передать эти данные сюда на дальнейшую обработку

	// Допустим тут у нас есть данные распарсенные с Refresh токена
	dataFromRefresh := jwtservice.RefreshTokenClaim{}

	// =================================================================================

	// Take last active session for current token
	lastActiveSession, err := s.SessionRepo.Read(dataFromRefresh.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrSession, err)
	}

	// Check RefreshTokenHash.
	if !hashpass.CheckHashedTokenWithTokenSHA256(lastActiveSession.RefreshTokenHash, refreshTokenReq.Token) {
		return nil, errs.ErrCompareToken
	}

	// Check ExpiresAt. Сессия валидна до N время после просрочена.
	if lastActiveSession.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errs.ErrTimeLiveSession
	}

	// Check RevokedAt. Если nil — сессия активна. Если заполнена — токен уже был использован.
	// И возможно, тот кто использовал токен - это хакер, или текущий пользователь это хакер
	// Данная атака называется Token Reuse Attack - хакерская атака повторного использования токена

	if lastActiveSession.RevokedAt != nil {
		// Прерываем все сессии у этого пользователя. Revoked == time.Now()
		// Данные сохраняем для анализа. Срок хранения 1 сутки.
		// + logger
		s.SessionRepo.CancelSessions(lastActiveSession.UserID)
		return nil, errs.ErrUsingToken
	}

	if !s.isItRealUser(lastActiveSession, refreshTokenReq) {
		// Если пользователь не прошел проверку лигитимности, прерываем все сессии.
		// + logger
		s.SessionRepo.CancelSessions(lastActiveSession.UserID)
		return nil, errs.ErrLegitimacyUser
	}

	// Выпуск новой пары токенов

	// Меняем статус сессии

	return nil, nil
}

// isItRealUser - Финальная проверка пользователя на лигитимность.
func (s *AuthServ) isItRealUser(session *models.Session, refreshTokenReq *dto.RefreshTokenRequest) bool {
	isTheSameDevice := session.DeviceInfo == refreshTokenReq.Device

	// Разные IP и устройства.
	if session.AddressIP != refreshTokenReq.IP && !isTheSameDevice {
		return false
	}
	// Другая страна (например RU -> BR)
	if s.GeoSecurity.IsCountryChanged(session.AddressIP, refreshTokenReq.IP) && !isTheSameDevice {
		return false
	}
	// Смена провайдера
	if s.GeoSecurity.IsProviderChanged(session.AddressIP, refreshTokenReq.IP) && !isTheSameDevice {
		return false
	}
	return true
}

// TODO...
// 1. Если есть только 1 историческая сессия и она скомпроментирована, т.е. от этой сессии уже выдан новый refresh token хакеру,
//    то нам не с чем будет сопоставить текущие данные пользователя для подтверждения лигитимности.
//
// Действия:
// 	  Проверяем время создания предыдущей сессии с time.Now,
//    если ( time.Now - timeCreatedSession ) < 15 sec, тогда с высокой вероятностью, перед нами хакер
//    дополнительно можно проверить его 'IP' и 'UserAgent'

// SELECT * FROM sessions WHERE user_id = X AND is_active = true;

func (s *AuthServ) Login(ctx context.Context, loginUser *dto.LoginUser) (*dto.TokenPair, error) {
	// Check user
	userDomain, err := s.UserRepo.ReadByEmail(loginUser.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidCredentials, err)
	}

	// Check password
	err = hashpass.CheckBcryptPassword(userDomain.HashPassword, loginUser.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidCredentials, err)
	}

	// Verification
	if !userDomain.CheckVerification() {
		if userDomain.Attempts >= AttemptsCnt {
			return nil, errs.ErrAttemptsVerification
		}
		userDomain.Attempts += 1
		verificToken, err := userDomain.UpdateVerificationToken()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
		}
		s.UserRepo.UpdateItem(userDomain)               // Update user
		s.Notifier.Send(userDomain.Email, verificToken) // Send email to user
		return nil, errs.ErrVerification
	}

	// Create Tokens and Session
	tokenPair, err := s.createTokenPair(userDomain)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
	}

	newSession := models.NewSession(
		userDomain.ID,
		tokenPair.SessionID,
		tokenPair.RefreshToken,
		loginUser.IP,
		loginUser.UserAgent,
		tokenPair.SessionExp,
	)

	s.SessionRepo.Create(newSession)

	return tokenPair, nil
}

func (s *AuthServ) createTokenPair(user *models.User) (*dto.TokenPair, error) {
	configRefresh := jwtservice.NewJWTConfig(
		config.TIME_LIVE_JWT_REFRESH_TOKEN,
		config.SECRET_KEY_JWT_REFRESH_TOKEN,
		config.HOST_NAME,
	)

	configAccess := jwtservice.NewJWTConfig(
		config.TIME_LIVE_JWT_ACCESS_TOKEN,
		config.SECRET_KEY_JWT_ACCESS_TOKEN,
		config.HOST_NAME,
	)

	tokenUser := jwtservice.NewTokenUser(user.ID, user.Name, user.Role)
	refreshClaim := jwtservice.NewRefreshTokenClaim(configRefresh, tokenUser)
	accessClaim := jwtservice.NewAccessTokenClaim(configAccess, refreshClaim)

	refToken, err1 := s.JWTManager.GenerateToken(configRefresh, refreshClaim)
	accToken, err2 := s.JWTManager.GenerateToken(configAccess, accessClaim)

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("%w:%w:%w", errs.ErrServer, err1, err2)
	}

	return &dto.TokenPair{
		AccessToken:  accToken,
		RefreshToken: refToken,
		SessionID:    refreshClaim.ID,
		SessionExp:   refreshClaim.ExpiresAt.Time,
	}, nil
}
