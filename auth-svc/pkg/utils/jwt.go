package utils

import (
	"auth-svc/pkg/models"
	"errors"
	"time"

	"github.com/goforj/godump"
	"github.com/golang-jwt/jwt"
)

type JWTWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type jwtClaims struct {
	jwt.StandardClaims
	Id    string
	Email string
	Role  string
}

func (w *JWTWrapper) GenerateToken(user models.User, role string) (signedToken string, err error) {
	claims := &jwtClaims{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			//Audience:  "",
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(w.ExpirationHours)).Unix(),
			//Id:        "",
			//IssuedAt:  0,
			Issuer: w.Issuer,
			//NotBefore: 0,
			//Subject:   "",
		},
	}

	godump.Dump(claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(w.SecretKey))

	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (w *JWTWrapper) ValidateToken(singedToken string, role string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		singedToken,
		&jwtClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(w.SecretKey), nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.Role != role {
		return nil, errors.New("role not matching")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}
	return claims, nil
}
