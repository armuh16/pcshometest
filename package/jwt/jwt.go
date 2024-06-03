package jwt

import (
	"context"
	"fmt"
	"pcstakehometest/enum"
	"time"

	"github.com/golang-jwt/jwt"
)

type ClaimData struct {
	UserID int    `json:"user_id,omitempty"`
	UUID   string `json:"uuid,omitempty"`
}

type InternalClaimData struct {
	UserID int           `json:"user_id,omitempty"`
	Role   enum.RoleType `json:"role,omitempty"`
}

// Claim struct
type Claim struct {
	Data ClaimData `json:"data"`
	jwt.StandardClaims
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

// RequestToken
func RequestToken(ctx context.Context, data ClaimData, secret string, accessExpiredAt, refreshExpiredAt int64) (*Token, error) {
	// Generate access Token JWT
	accessToken, err := GenerateToken(Claim{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: accessExpiredAt,
		},
	}, secret)
	if err != nil {
		return nil, err
	}

	// Generate refresh Token JWT
	refreshToken, err := GenerateToken(Claim{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: refreshExpiredAt,
		},
	}, secret)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

// GenerateToken
func GenerateToken(c Claim, secretString string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(secretString))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// ParseToken
func ParseToken(tokenString string, secretString string) (*jwt.Token, error) {
	t, err := jwt.Parse(tokenString, func(jt *jwt.Token) (interface{}, error) {
		if _, ok := jt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jt.Header["alg"])
		}
		return []byte(secretString), nil
	})
	return t, err
}

// IsValidToken validate JWT Token
func IsValidToken(tokenString string, secretString string) (bool, error) {
	token, err := ParseToken(tokenString, secretString)
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

// ParseClaim func
func ParseClaim(tokenString string, secretString string) (*Claim, error) {
	claims := Claim{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretString), nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
