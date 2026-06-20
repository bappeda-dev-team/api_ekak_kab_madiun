package helper

import (
	"ekak_kabupaten_madiun/model/web"

	"context"
)

// see contextkey.go
func GetUserInfo(ctx context.Context) web.JWTClaim {
	if ctx == nil {
		return web.JWTClaim{}
	}
	val := ctx.Value(UserInfoKey)
	if val == nil {
		return web.JWTClaim{}
	}

	claims, ok := val.(web.JWTClaim)
	if !ok {
		return web.JWTClaim{}
	}

	return claims
}
