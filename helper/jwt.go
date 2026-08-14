package helper

import (
	"context"
	"fmt"
	"os"
	"log"
	"ekak_kabupaten_madiun/model/web"

	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	jwtVerifier *oidc.IDTokenVerifier
)

func InitJWT() error {
	jwtIssuer := os.Getenv("KEYCLOAK_ISSUER")
	if jwtIssuer == "" {
		return fmt.Errorf("KEYCLOAK_ISSUER is not configured")
	}

	log.Printf("Mencoba koneksi ke keycloak issuer: %s", jwtIssuer)
	provider, err := oidc.NewProvider(
		context.Background(),
		jwtIssuer,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Keycloak OIDC provider: %w", err)
	}

	jwtVerifier = provider.Verifier(
		&oidc.Config{
			SkipClientIDCheck: true,
		},
	)

	return nil
}

func ValidateJWT(tokenString string) (web.JWTClaim, error) {
	if jwtVerifier == nil {
		return web.JWTClaim{}, fmt.Errorf("JWT is not initialized")
	}

	idToken, err := jwtVerifier.Verify(
		context.Background(),
		tokenString,
	)

	if err != nil {
		return web.JWTClaim{}, fmt.Errorf("invalid JWT: %w", err)
	}

	// if !token.Valid {
	// 	return web.JWTClaim{}, fmt.Errorf("invalid JWT")
	// }
	var claims struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
		Issuer  string `json:"iss"`
		Iat     int64  `json:"iat"`
		Exp     int64  `json:"exp"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return web.JWTClaim{}, fmt.Errorf(
			"failed to parse JWT claims: %w",
			err,
		)
	}

	return web.JWTClaim{
		Issuer: claims.Issuer,
		Email:  claims.Email,
		Iat:    claims.Iat,
		Exp:    claims.Exp,
	}, nil
}
