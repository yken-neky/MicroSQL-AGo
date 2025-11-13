package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims representa las claims en el JWT con información adicional del usuario
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret      string
	expiryHours int
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret:      secret,
		expiryHours: 24, // default 24 horas
	}
}

// NewJWTServiceWithExpiry crea un nuevo JWTService con expiración personalizada
func NewJWTServiceWithExpiry(secret string, expiryHours int) *JWTService {
	if expiryHours <= 0 {
		expiryHours = 24
	}
	return &JWTService{
		secret:      secret,
		expiryHours: expiryHours,
	}
}

// Generate genera un JWT con solo userID (método legacy, mantenido para compatibilidad)
func (s *JWTService) Generate(userID uint, expHours int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(expHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// GenerateToken genera un JWT completo con userID, username y role
func (s *JWTService) GenerateToken(userID uint, username string, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.expiryHours) * time.Hour)

	claims := TokenClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Validate valida un JWT y retorna el userID (método legacy)
func (s *JWTService) Validate(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(float64)
	return uint(sub), nil
}

// ValidateToken valida y parsea un JWT token completo, retornando las claims
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

// ExtractBearerToken extrae el token del header "Authorization: Bearer <token>"
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("authorization header must start with 'Bearer '")
	}

	return authHeader[7:], nil
}
