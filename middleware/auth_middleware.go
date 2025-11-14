package middleware

import (
	"context"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

var JWKS *keyfunc.JWKS

func InitJWKS(jwksURL string) error {
	var err error
	JWKS, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	return err
}

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	publicPaths := []struct {
		path    string
		pattern string
	}{
		{"/user/login", "^/user/login$"},
		{"/api/pokin_opd/findall/", "^/api/pokin_opd/findall/[^/]+/[^/]+$"},
		{"/api/pokin_pemda/subtematik/", "^/api/pokin_pemda/subtematik/[^/]+$"},
		{"/pohon_kinerja/pokin_atasan/", "^/pohon_kinerja/pokin_atasan/[^/]+$"},
		{"/rekin/atasan/", "^/rekin/atasan/[^/]+$"},
		{"/api_internal/rencana_kinerja/findall", "^/api_internal/rencana_kinerja/findall$"},
	}

	currentPath := request.URL.Path
	for _, route := range publicPaths {
		if strings.HasPrefix(currentPath, route.path) || helper.MatchPattern(currentPath, route.pattern) {
			middleware.Handler.ServeHTTP(writer, request)
			return
		}
	}

	tokenString := request.Header.Get("Authorization")
	if tokenString == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak ditemukan",
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	tokenHeader := request.Header.Get("Authorization")
	if tokenHeader == "" || !strings.HasPrefix(tokenHeader, "Bearer ") {
		writeUnauthorized(writer, "Missing or invalid Authorization header")
		return
	}
	rawToken := strings.TrimPrefix(tokenHeader, "Bearer ")

	token, err := jwt.Parse(rawToken, JWKS.Keyfunc)
	if err != nil {
		log.Printf("JWT parse error: %v", err)
		writeUnauthorized(writer, "Invalid token")
		return
	}
	if !token.Valid {
		log.Println("JWT is not valid")
		writeUnauthorized(writer, "Invalid token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		writeUnauthorized(writer, "Invalid claims")
		return
	}

	// Optional: validate audience or issuer
	issuer := os.Getenv("KEYCLOAK_ISSUER")
	if iss, ok := claims["iss"].(string); !ok || iss != issuer {
		writeUnauthorized(writer, "Invalid issuer")
		return
	}

	clientId, hasClientId := claims["client_id"].(string)
	if hasClientId {
		allowedService := os.Getenv("ALLOWED_SERVICE_CLIENT_ID")

		if clientId != allowedService {
			writeForbidden(writer, "client not allowed")
			return
		}
	} else {
		sessionId := request.Header.Get("X-Session-Id")
		if sessionId == "" {
			writeUnauthorized(writer, "User session not found")
			return
		}
	}

	ctx := context.WithValue(request.Context(), helper.UserInfoKey, claims)
	ctx = context.WithValue(ctx, "client_id", clientId)
    ctx = context.WithValue(ctx, "is_service_token", hasClientId)

	ctx = context.WithValue(ctx, "user_session_id", request.Header.Get("X-Session-Id"))

	request = request.WithContext(ctx)

	middleware.Handler.ServeHTTP(writer, request)
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   http.StatusUnauthorized,
		Status: "UNAUTHORIZED",
		Data:   message,
	})
}

func writeForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   http.StatusForbidden,
		Status: "ACCESS NOT ALLOWED",
		Data:   message,
	})
}
