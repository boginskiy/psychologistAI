package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/pkg/security"
	"github.com/google/uuid"

	models "github.com/boginskiy/psychologistAI/internal/models/user"
	"github.com/boginskiy/psychologistAI/internal/repository"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/boginskiy/psychologistAI/pkg/jwtservice"
)

const AttemptsCnt = 5 // Количество попыток для верификации пользователя
const OffSet = 3      // Глубина удаления исторических сессией пользователя

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
	refreshTokenClaim := jwtservice.RefreshTokenClaim{}

	// =================================================================================

	// Take last active session for current token
	lastActiveSession, err := s.SessionRepo.Read(refreshTokenClaim.ID)
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

	// Create new Tokens and Session
	tokenPair, err := s.createTokenPair(
		refreshTokenClaim.UserID,
		refreshTokenClaim.UserName,
		refreshTokenClaim.UserRoles)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
	}

	newSession := models.NewSession(
		tokenPair.SessionID,
		refreshTokenReq.IP,
		tokenPair.RefreshToken,
		refreshTokenReq.UserAgent,
		refreshTokenReq.OS,
		refreshTokenReq.Browser,
		refreshTokenReq.Device,
		tokenPair.SessionExp,
		refreshTokenClaim.UserID,
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
	tokenPair, err := s.createTokenPair(userDomain.ID, userDomain.Name, userDomain.Role)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
	}

	newSession := models.NewSession(
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

	tokenUser := jwtservice.NewTokenUser(userID, userName, userRole)
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
