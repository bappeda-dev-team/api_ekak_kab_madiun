package helper

import (
	"context"
	"fmt"
	"log"
	"os"

	"ekak_kabupaten_madiun/model/web"

	"github.com/coreos/go-oidc/v3/oidc"
)

var keycloakVerifier *oidc.IDTokenVerifier

func InitKeycloakJWT() error {
	log.Printf("VALIDATE JWT USING KEYCLOAK")
	issuer := os.Getenv("KEYCLOAK_ISSUER")

	if issuer == "" {
		return fmt.Errorf("KEYCLOAK_ISSUER is not configured")
	}

	log.Printf(
		"Mencoba koneksi ke Keycloak issuer: %s",
		issuer,
	)

	provider, err := oidc.NewProvider(
		context.Background(),
		issuer,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to initialize Keycloak OIDC provider: %w",
			err,
		)
	}

	keycloakVerifier = provider.Verifier(
		&oidc.Config{
			SkipClientIDCheck: true,
		},
	)

	jwtProvider = &KeycloakJWTProvider{
		verifier: keycloakVerifier,
	}

	return nil
}

type KeycloakJWTProvider struct {
	verifier *oidc.IDTokenVerifier
}

func (p *KeycloakJWTProvider) Validate(
	tokenString string,
) (web.JWTClaim, error) {

	idToken, err := p.verifier.Verify(
		context.Background(),
		tokenString,
	)
	if err != nil {
		return web.JWTClaim{}, fmt.Errorf(
			"invalid JWT: %w",
			err,
		)
	}

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
