package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	JwtSecretKey         string `env:"JWT_SECRET_KEY"`
	JwtAccessTokenHours  int64  `env:"JWT_ACCESS_TOKEN_HOURS,default=24"`
	JwtRefreshTokenHours int64  `env:"JWT_REFRESH_TOKEN_HOURS,default=168"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Authorizer struct {
	secretKey       []byte
	accessDuration  int64
	refreshDuration int64
}

func NewAuthorizer(config Config) Authorizer {
	return Authorizer{
		secretKey:       []byte(config.JwtSecretKey),
		accessDuration:  config.JwtAccessTokenHours,
		refreshDuration: config.JwtRefreshTokenHours,
	}
}

func (a Authorizer) Validate(token string) (bool, string, error) {
	ok, identity, err := a.verifyToken(token)
	if err != nil {
		return false, "", fmt.Errorf("token verification: %w")
	}
	return ok, identity, nil
}

func (a Authorizer) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * time.Duration(a.accessDuration)).Unix(),
		})

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", fmt.Errorf("token sign: %w", err)
	}
	return tokenString, nil
}

func (a Authorizer) verifyToken(tokenString string) (bool, string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return a.secretKey, nil
	})

	identityData, ok := claims["username"]
	if !ok {
		return false, "", fmt.Errorf("claim not found")
	}
	identity, ok := identityData.(string)
	if !ok {
		return false, "", fmt.Errorf("claim invalid")

	}
	if err != nil {
		return false, identity, fmt.Errorf("parse token: %w", err)
	}
	return token.Valid, identity, nil
}
