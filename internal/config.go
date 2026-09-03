package internal

import (
	"context"
	"net/http"
	"os"
)

type BaseClient struct {
	host       string
	path       string
	httpClient *http.Client
}

// constructor
func newBaseClient(host, path string, httpClient *http.Client) BaseClient {
	return BaseClient{
		host:       host,
		path:       path,
		httpClient: httpClient,
	}
}

// key untuk context session id
type ctxKey string

const SessionIDKey ctxKey = "X-Session-Id"

// Inject session ID ke context (opsional)
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

// Ambil session ID dari context
func getSessionID(ctx context.Context) string {
	if v := ctx.Value(SessionIDKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return os.Getenv("DEV_SESSION_ID") // fallback
}
