package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var hmacSecret = []byte("replace-with-a-long-random-secret-at-least-32-bytes")

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 接口：产生token
func GenerateToken(UserID, Role string, ttl time.Duration) (string, error) {
	now := time.Now()

	// 创建jwt的负荷
	claims := CustomClaims{
		UserID: UserID,
		Role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.example.com",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{"api.example.com"},
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// 创建jwt对象
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	// 返回签名后的token
	return token.SignedString(hmacSecret)
}

// 接口：解析并验证token
func ParseToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSecret, nil
		},
		jwt.WithIssuer("auth.example.com"),
		jwt.WithAudience("api.example.com"),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// func main() {
// 	tokenString, err := GenerateToken("user_123", "admin", 2*time.Hour)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("JWT:")
// 	fmt.Println(tokenString)

// 	claims, err := ParseToken(tokenString)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("UserID:", claims.UserID)
// 	fmt.Println("Role:", claims.Role)
// 	fmt.Println("Subject:", claims.Subject)
// }
