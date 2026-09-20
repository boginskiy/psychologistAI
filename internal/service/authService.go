package service

import (
	"context"
	"fmt"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/errs/server"
	"github.com/boginskiy/psychologistAI/internal/errs/users"
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
}

func NewAuthServ(
	ctx context.Context,
	validater Validater,
	notifier Notifier,
	userRepo repository.UserRepo,
	SessionRepo repository.SessionRepo,
	jwtManager jwtservice.JWTManager,
) *AuthServ {
	return &AuthServ{
		Validater:  validater,
		Notifier:   notifier,
		UserRepo:   userRepo,
		JWTManager: jwtManager,
	}
}

func (s *AuthServ) Refresh(ctx context.Context, refreshTokenReq *dto.RefreshTokenRequest) (*dto.Token, error) {

	return nil, nil
}

func (s *AuthServ) Login(ctx context.Context, loginUser *dto.LoginUser) (*dto.TokenPair, error) {
	// Check user
	userDomain, err := s.UserRepo.ReadByEmail(loginUser.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrInvalidCredentials, err)
	}

	// Check password
	err = hashpass.CheckBcryptPassword(userDomain.HashPassword, loginUser.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrInvalidCredentials, err)
	}

	// Verification
	if !userDomain.CheckVerification() {
		if userDomain.Attempts >= AttemptsCnt {
			return nil, users.ErrAttemptsVerification
		}
		userDomain.Attempts += 1
		verificToken, err := userDomain.UpdateVerificationToken()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", server.ErrServer, err)
		}
		s.UserRepo.UpdateItem(userDomain)               // Update user
		s.Notifier.Send(userDomain.Email, verificToken) // Send email to user
		return nil, users.ErrVerification
	}

	// Create Tokens and Session
	tokenPair, err := s.createTokenPair(userDomain)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", server.ErrServer, err)
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
		return nil, fmt.Errorf("%w:%w:%w", server.ErrServer, err1, err2)
	}

	return &dto.TokenPair{
		AccessToken:  accToken,
		RefreshToken: refToken,
		SessionID:    refreshClaim.ID,
		SessionExp:   refreshClaim.ExpiresAt.Time,
	}, nil
}

// HOST_SITE = "psychologistAI.com"

// 	// Tokens-JWT
// 	TIME_LIVE_JWT_ACCESS_TOKEN  = 15 // Minute
// 	SECRET_KEY_JWT_ACCESS_TOKEN = "cjlsjdc3r483ucdhcyeruf9ehrc"

// 	TIME_LIVE_JWT_REFRESH_TOKEN  = 10080 // Minute
// 	SECRET_KEY_JWT_REFRESH_TOKEN = "acsdcfwecwX2343FCSDBJHXSDJ"

// Также при каждой выдаче новой пары токенов старый access token аннулируется временем жизни,
// а refresh token должен ротироваться (выпускать новый jti).

// Ротация Refresh Token (Обязательно): При каждом использовании старого RT выдавайте совершенно новый
// RT, а старый немедленно аннулируйте (заносите его jti в базу отозванных токенов). Это защищает от кражи:
// если злоумышленник перехватил ваш RT и попытался им воспользоваться, настоящий пользователь в этот же
// момент получит ошибку «токен недействителен» при следующем автоматическом рефреше и поймет, что аккаунт
// под угрозой.

// Поле jti: Добавляйте в оба токена уникальный идентификатор (JWT ID). Для AT он помогает отслеживать
// активные сессии, а для RT является первичным ключом в таблице отзыва.

// Разные секреты: Никогда не подписывайте AT и RT одним и тем же ключом. Утечка ключа для AT (который
// гулятирует по всем клиентам) не должна позволять выпускать вечные Refresh-токены.

// Revocation list (Denylist): Поскольку JWT самодостаточны, единственный способ их отозвать —
// хранить идентификаторы (jti) отозванных refresh-токенов в быстром хранилище вроде Redis с TTL до
// даты истечения самого токена. Перед выдачей новой пары всегда проверяйте приходящий RT по этому
// списку.

// !!!
// SessionID (или LinkedJTI): Это ваше главное оружие против кражи токенов. Здесь хранится
// значение jti (уникальный ID) того самого refresh token, которым был выдан данный access token.

// Зачем: Когда приходит запрос с Access Token, вы проверяете: «А существует ли еще в
// моей таблице сессий активный Refresh Token с таким jti?». Если злоумышленник украл
// Access Token, но вы тем временем по кнопке «Выйти на всех устройствах» удалили
// соответствующую строку сессии из Redis/БД, проверка SessionID провалится, и доступ
// будет закрыт немедленно.

// // Внутри вашего BFF при логине:
// fingerprint := r.Header.Get("User-Agent") + ":" + getClientIP(r)
// saveSessionToDB(Session{
//     UserID:      user.ID,
//     RefreshHash: hashSHA256(refreshTokenString),
//     Fingerprint: fingerprint,
//     ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
// })
