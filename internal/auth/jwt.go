package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKey = "your_secret_key"
	defaultTokenDuration = 24 * time.Hour
)

type JwtClaims struct{
	UserID uint `json:"user_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type JwtService interface{
	GenerateToken(userId uint, name string, email string) (string, error);
	ValidateToken(tokenStr string) (*JwtClaims, error)
}

type jwtService struct{
secretKey string
tokenDuration time.Duration
}

func NewJWTService(secretKey string) JwtService {

	if secretKey == "" {
		secretKey=jwtSecretKey;
	}

	return &jwtService{
		secretKey:     secretKey,
		tokenDuration: defaultTokenDuration,
	}
}

func (j *jwtService) GenerateToken(userId uint, name string, email string) (string, error) {
	claims := JwtClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenStr string) (*JwtClaims, error) {
	fmt.Println("tokenStr is ", tokenStr);
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}