package authService

import (
	"errors"
	"ppob-backend/app/dto"
	"ppob-backend/app/repository/authRepository"
	"ppob-backend/config"
	"ppob-backend/model"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	authService struct {
		Config         *config.SystemConfig
		authRepository authRepository.AuthRepository
	}

	AuthService interface {
		Login(request dto.LoginRequest) (dto.LoginResponse, error)
		Register(request dto.RegisterRequest) (dto.RegisterResponse, error)
	}
)

func NewAuthService(config *config.SystemConfig, authRepository authRepository.AuthRepository) AuthService {
	return &authService{
		Config:         config,
		authRepository: authRepository,
	}
}

func (s *authService) Login(request dto.LoginRequest) (dto.LoginResponse, error) {
	s.Config.Logger.Infof("Login request for %s", request.Username)
	credential := s.authRepository.GetUserCredentialByUsername(request.Username)
	emptyResponse := dto.LoginResponse{}

	if (dto.UserCredential{}) == credential {
		s.Config.Logger.Warnf("Failed login for %s: no username found", request.Username)
		return emptyResponse, errors.New("username/password mismatch")
	}

	err := bcrypt.CompareHashAndPassword([]byte(credential.Password), []byte(request.Password))

	if err != nil {
		s.Config.Logger.Warnf("Failed login for %s: failed compare password", request.Username)
		return emptyResponse, errors.New("username/password mismatch")
	}

	claims := &dto.JwtCustomClaims{
		Username: request.Username,
		UserID:   credential.Uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.Config.JwtSecret))

	if err != nil {
		s.Config.Logger.Warnf("Failed login for %s: failed to sign token %+v", request.Username, err)
		return emptyResponse, errors.New("failed to sign token")
	}

	// userID, _ := strconv.ParseUint(credential.ID, 10, 64)

	response := dto.LoginResponse{
		Username: request.Username,
		Name:     credential.Name,
		Token:    t,
		Uuid:     credential.Uuid,
	}
	return response, err
}

func (s *authService) Register(request dto.RegisterRequest) (dto.RegisterResponse, error) {
	s.Config.Logger.Infof("Registering for %s", request.Username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	uuid := uuid.New()

	user := model.User{
		Name:     request.Name,
		Username: request.Username,
		Password: string(hashedPassword),
		Uuid:     uuid.String(),
	}
	response := dto.RegisterResponse{}

	if err != nil {
		s.Config.Logger.Warnf("Error registering for %s: error on hashing password", request.Username)
		return response, errors.New("failed to register")
	}

	user.CreatedAt = time.Now()

	err = s.authRepository.SaveUser(&user)

	if err != nil {
		s.Config.Logger.Warnf("Error registering for %s: error on saving user", request.Username)
		return response, errors.New("failed to register")
	}

	err = s.authRepository.CreateWallet(&user)

	if err != nil {
		s.Config.Logger.Warnf("Error registering for %s: error on creating wallet", request.Username)
		return response, errors.New("failed to register")
	}

	response.Name = user.Name
	response.Username = user.Username
	response.CreatedAt = user.CreatedAt.String()

	return response, nil
}
