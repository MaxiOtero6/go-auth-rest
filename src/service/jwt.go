package service

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/repository"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	key        []byte
	repository repository.IJWTRepository
}

var jwtKey = os.Getenv("JWT_SECRET_KEY")

func NewJWTService(jwtRepository *repository.IJWTRepository) JWTService {
	var key = []byte("secret")

	var repo repository.IJWTRepository = &repository.JWTRepository{}

	if jwtKey != "" {
		key = []byte(jwtKey)
	}

	if jwtRepository != nil {
		repo = *jwtRepository
	}

	return JWTService{key: []byte(key), repository: repo}
}

func (service *JWTService) Sign(payload model.User) (string, error) {
	nowUtc := time.Now().UTC()

	claim := model.JWTPayload{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(nowUtc.Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(nowUtc),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(service.key)

	return tokenString, err
}

func (service *JWTService) isJWT(tokenString string) bool {
	jwtRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)$`)

	return jwtRegex.MatchString(tokenString)
}

func (service *JWTService) Verify(tokenString string) (bool, error) {
	if !service.isJWT(tokenString) {
		return false, &model.ValidationError{Title: "Invalid JWT", Detail: "The provided token doesn't have JWT format"}
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return service.key, nil
	})

	return token.Valid, err
}

func (service *JWTService) Decode(tokenString string) (model.JWTPayload, error) {
	if !service.isJWT(tokenString) {
		return model.JWTPayload{}, &model.ValidationError{
			Title:  "Invalid JWT",
			Detail: "The provided token doesn't have JWT format",
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return service.key, nil
	})

	if claims, ok := token.Claims.(*model.JWTPayload); ok && token.Valid {
		return *claims, nil
	} else {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return model.JWTPayload{}, &model.AuthenticationError{
				Title:  "Expired token",
				Detail: "Your token has expired, please try logging in again",
			}
		}

		return model.JWTPayload{}, err
	}
}

func (service *JWTService) Blacklist(tokenString string) error {
	lastDotIndex := strings.LastIndex(tokenString, ".")
	signature := tokenString[lastDotIndex:]

	decoded, err := service.Decode(tokenString)

	if err != nil {
		return err
	}

	exp := decoded.ExpiresAt.Time

	service.repository.Blacklist(signature, exp)

	return nil
}

func (service *JWTService) IsBlacklisted(tokenString string) bool {
	lastDotIndex := strings.LastIndex(tokenString, ".")
	signature := tokenString[lastDotIndex:]

	return service.repository.IsBlacklisted(signature)
}
