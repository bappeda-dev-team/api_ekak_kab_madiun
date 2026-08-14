package helper

import (
	"ekak_kabupaten_madiun/model/web"
	"fmt"
	"os"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwks keyfunc.Keyfunc
)

func InitJWT() error {
	keycloakCerts := os.Getenv("KEYCLOAK_CERTS_URL")

	var err error
	jwks, err = keyfunc.NewDefault([]string{keycloakCerts})
	if err != nil {
		return fmt.Errorf("failed to initialize Keycloak JWKS: %w", err)
	}

	return nil
}

func ValidateJWT(tokenString string) (web.JWTClaim, error) {
	if jwks == nil {
		return web.JWTClaim{}, fmt.Errorf("JWT is not initialized")
	}

	token, err := jwt.Parse(
		tokenString,
		jwks.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if err != nil {
		return web.JWTClaim{}, fmt.Errorf("invalid JWT: %w", err)
	}

	if !token.Valid {
		return web.JWTClaim{}, fmt.Errorf("invalid JWT")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return web.JWTClaim{}, fmt.Errorf("invalid JWT claims")
	}

	// mapping claims sementara
	_ = claims

	return web.JWTClaim{}, nil
}
