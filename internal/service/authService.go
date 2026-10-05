package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/request"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/pkg/emailservice"
	"github.com/boginskiy/psychologistAI/pkg/security"
	"github.com/boginskiy/psychologistAI/pkg/validservice"
	"github.com/google/uuid"

	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
	"github.com/boginskiy/psychologistAI/internal/repository"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/boginskiy/psychologistAI/pkg/jwtservice"
)

const Attempts = 5               // Количество попыток для верификации пользователя
const Retry = 3                  // Количество попыток для верификации пользователя
const OffSet = 3                 // Глубина удаления исторических сессией пользователя
const NameProject = "'A FRIEND'" // Имя проекта

type AuthServ struct {
	Validater   validservice.Validater
	Postman     emailservice.Postman
	JWTManager  jwtservice.JWTManager
	GeoSecurity security.GeoSecurity

	UserRepo    repository.UserRepo
	SessionRepo repository.SessionRepo
}

func NewAuthServ(
	ctx context.Context,
	validater validservice.Validater,
	postman emailservice.Postman,
	jwtManager jwtservice.JWTManager,
	geoSecurity security.GeoSecurity,

	userRepo repository.UserRepo,
	sessionRepo repository.SessionRepo,
) *AuthServ {
	return &AuthServ{
		Validater:   validater,
		Postman:     postman,
		JWTManager:  jwtManager,
		GeoSecurity: geoSecurity,

		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}

func (s *AuthServ) Logout(ctx context.Context) error {
	infoUser, ok := request.GetInfoUserFromContext(ctx)
	if !ok {
		return errs.ErrInfoContext
	}
	// Delete all sessions. Revoked == time.Now()
	s.SessionRepo.CancelSessions(infoUser.UserID)
	return nil
}

func (s *AuthServ) Auth(ctx context.Context, accessTokenReq *dto.TokenReq) (*dto.InfoUser, error) {
	accessClaim := &jwtservice.AccessTokenClaim{}

	_, err := s.JWTManager.CheckAndParseToken(config.SECRET_KEY_JWT_ACCESS_TOKEN, accessTokenReq.Token, accessClaim)
	if err != nil {
		// + logger
		return nil, fmt.Errorf("%w: %w", errs.ErrAuth, err)
	}

	// Check active session
	currentSession, err := s.SessionRepo.Read(accessClaim.SessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrSession, err)
	}

	if currentSession.RevokedAt != nil {
		return nil, errs.ErrActiveSession
	}

	// Deep check user. Если пользователь не прошел проверку лигитимности, кидаем ошибку.
	if !s.isItRealUser(
		currentSession.AddressIP, accessTokenReq.IP,
		currentSession.DeviceInfo, accessTokenReq.Device) {
		// + logger
		return nil, errs.ErrLegitimacyUser
	}

	infoUser := &dto.InfoUser{
		UserID:    accessClaim.UserID,
		UserRoles: accessClaim.UserRoles,
		TokenType: accessClaim.TokenType,
		SessionID: accessClaim.SessionID,
	}
	return infoUser, nil
}

func (s *AuthServ) Refresh(ctx context.Context, refreshTokenReq *dto.TokenReq) (*dto.TokenPair, error) {
	refreshClaim := &jwtservice.RefreshTokenClaim{}

	_, err := s.JWTManager.CheckAndParseToken(config.SECRET_KEY_JWT_REFRESH_TOKEN, refreshTokenReq.Token, refreshClaim)
	if err != nil {
		// + logger
		return nil, fmt.Errorf("%w: %w", errs.ErrAuth, err)
	}

	// Take last active session for current token
	lastActiveSession, err := s.SessionRepo.Read(refreshClaim.ID)
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

	// Deep check user. Если пользователь не прошел проверку лигитимности, прерываем все сессии.
	if !s.isItRealUser(
		lastActiveSession.AddressIP, refreshTokenReq.IP,
		lastActiveSession.DeviceInfo, refreshTokenReq.Device) {
		// + logger
		s.SessionRepo.CancelSessions(lastActiveSession.UserID)
		return nil, errs.ErrLegitimacyUser
	}

	// Create new Tokens and Session
	tokenPair, err := s.createTokenPair(
		refreshClaim.UserID,
		refreshClaim.UserName,
		refreshClaim.UserRoles)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
	}

	newSession := domain.NewSession(
		tokenPair.SessionID,
		refreshTokenReq.IP,
		tokenPair.RefreshToken,
		refreshTokenReq.UserAgent,
		refreshTokenReq.OS,
		refreshTokenReq.Browser,
		refreshTokenReq.Device,
		tokenPair.SessionExp,
		refreshClaim.UserID,
	)

	// Привязываем новую сессию к старой. Цепочка сессий.
	newSession.PreviousID = lastActiveSession.ID

	// Контроль количеств исторических сессий. Сохраняем в БД 3 крайних сессии.
	// 1 - newSession                   - текущая новая сессия.
	// 2 - lastActiveSession            - последняя активная сессия.
	// 3 - lastActiveSession.PreviousID - самая старая сессия. После нее не должно быть иных сессий.

	// Создаем запись с новой сессией.
	// Cancel предыдущую сессию
	// Удаляем самую старую сессию
	err = s.SessionRepo.UpdateAfterRefresh(newSession, OffSet)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrUpdateDBAfterRefresh, err)
	}

	return tokenPair, nil
}

// isItRealUser - Финальная проверка пользователя на лигитимность.
func (s *AuthServ) isItRealUser(IPFromSession, IPFromReq, DeviceFromSession, DeviceFromReq string) bool {
	isTheSameDevice := DeviceFromSession == DeviceFromReq

	// Разные IP и устройства.
	if IPFromSession != IPFromReq && !isTheSameDevice {
		return false
	}
	// Другая страна (например RU -> BR)
	if s.GeoSecurity.IsCountryChanged(IPFromSession, IPFromReq) && !isTheSameDevice {
		return false
	}
	// Смена провайдера
	if s.GeoSecurity.IsProviderChanged(IPFromSession, IPFromReq) && !isTheSameDevice {
		return false
	}
	return true
}

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
		if userDomain.Attempts >= Attempts {
			return nil, errs.ErrAttemptsVerification
		}
		userDomain.Attempts += 1
		verificToken, err := userDomain.UpdateVerificationToken()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
		}

		// Update user
		s.UserRepo.UpdateItem(userDomain)

		// Create message and Send email to user
		verifLetter := emailservice.NewVerifLetter(userDomain.Email, verificToken, NameProject)
		s.Postman.Send(verifLetter, Retry)

		return nil, errs.ErrVerification
	}

	// Create Tokens and Session
	tokenPair, err := s.createTokenPair(userDomain.ID, userDomain.Name, userDomain.Role)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
	}

	newSession := domain.NewSession(
		tokenPair.SessionID,
		loginUser.IP,
		tokenPair.RefreshToken,
		loginUser.UserAgent,
		loginUser.OS,
		loginUser.Browser,
		loginUser.Device,
		tokenPair.SessionExp,
		userDomain.ID,
	)

	s.SessionRepo.Create(newSession)
	return tokenPair, nil
}

func (s *AuthServ) createTokenPair(userID uuid.UUID, userName string, userRole []string) (*dto.TokenPair, error) {
	configRefresh := jwtservice.NewClaimConfig(
		config.TIME_LIVE_JWT_REFRESH_TOKEN,
		config.HOST_NAME,
	)

	configAccess := jwtservice.NewClaimConfig(
		config.TIME_LIVE_JWT_ACCESS_TOKEN,
		config.HOST_NAME,
	)

	tokenUser := jwtservice.NewTokenUser(userID, userName, userRole)
	refreshClaim := jwtservice.NewRefreshTokenClaim(configRefresh, tokenUser)
	accessClaim := jwtservice.NewAccessTokenClaim(configAccess, refreshClaim)

	refToken, err1 := s.JWTManager.GenerateToken(config.SECRET_KEY_JWT_REFRESH_TOKEN, refreshClaim)
	accToken, err2 := s.JWTManager.GenerateToken(config.SECRET_KEY_JWT_ACCESS_TOKEN, accessClaim)

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
