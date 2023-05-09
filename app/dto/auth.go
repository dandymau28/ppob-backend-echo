package dto

import "github.com/golang-jwt/jwt"

type (
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Uuid     string `json:"user_id"`
	}

	RegisterRequest struct {
		Username string `json:"username" validate:"required,min=8"`
		Password string `json:"password" validate:"required,min=6,max=6,numeric"`
		Name     string `json:"name" validate:"required"`
	}

	RegisterResponse struct {
		Username  string `json:"username" validate:"required,min=8"`
		Name      string `json:"name" validate:"required"`
		CreatedAt string `json:"created_at"`
	}

	UserCredential struct {
		Uuid     string
		Name     string
		Username string
		Password string
	}

	JwtCustomClaims struct {
		Username string `json:"username"`
		UserID   string `json:"user_id"`
		jwt.StandardClaims
	}
)
