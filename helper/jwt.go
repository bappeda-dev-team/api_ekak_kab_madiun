package helper

import (
	"fmt"
	"os"

	"ekak_kabupaten_madiun/model/web"
)

type JWTProvider interface {
	Validate(tokenString string) (web.JWTClaim, error)
}

var jwtProvider JWTProvider

func InitJWT() error {
	provider := os.Getenv("AUTH_PROVIDER")

	switch provider {
	case "keycloak":
		return InitKeycloakJWT()

	default:
		jwtProvider = NewInternalJWTProvider()
		return nil
	}
}

func ValidateJWT(tokenString string) (web.JWTClaim, error) {
	if jwtProvider == nil {
		return web.JWTClaim{}, fmt.Errorf("JWT provider is not initialized")
	}

	return jwtProvider.Validate(tokenString)
}
