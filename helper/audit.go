package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// CreateEventRequest payload POST {EVENT_AUDIT_HOST}/events.
// ProjectId dan ServiceName diisi otomatis dari env jika kosong.
type CreateEventRequest struct {
	ProjectId string `json:"project_id"`
	EventType string `json:"event_type"`
	Action    string `json:"action"`

	ServiceName string `json:"service_name"`

	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`

	ActorID       *string `json:"actor_id,omitempty"`
	ActorUsername *string `json:"actor_username,omitempty"`

	RequestID *string `json:"request_id,omitempty"`
	SessionID *string `json:"session_id,omitempty"`

	OccurredAt *time.Time `json:"occurred_at,omitempty"`

	Before   any            `json:"before,omitempty"`
	After    any            `json:"after,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

var (
	auditHTTPClient = &http.Client{Timeout: 10 * time.Second}
	auditSkipLogged sync.Once
)

func auditHost() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("EVENT_AUDIT_HOST")), "/")
}

func auditEnabled() bool {
	return auditHost() != ""
}

func ptrIfNotEmpty(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// EnrichAuditFromContext mengisi project/service dari env dan actor dari JWT di context.
func EnrichAuditFromContext(ctx context.Context, event *CreateEventRequest) {
	if event == nil {
		return
	}
	if event.ProjectId == "" {
		event.ProjectId = os.Getenv("PROJECT_ID")
	}
	if event.ServiceName == "" {
		event.ServiceName = os.Getenv("SERVICE_NAME")
	}
	if event.OccurredAt == nil {
		now := time.Now()
		event.OccurredAt = &now
	}
	if event.EventType == "" {
		event.EventType = "DATA_CHANGED"
	}

	claims := GetUserInfo(ctx)
	if event.ActorID == nil {
		username := strings.TrimSpace(claims.Nip)
		if username == "" {
			username = strings.TrimSpace(claims.Email)
		}
		event.ActorID = ptrIfNotEmpty(username)
	}

	if event.ActorUsername == nil {
		username := strings.TrimSpace(claims.Nip)
		if username == "" {
			username = strings.TrimSpace(claims.Email)
		}
		event.ActorUsername = ptrIfNotEmpty(username)
	}
}

// PublishAuditEvent mengirim event ke audit API secara sinkron.
// Panggil HANYA setelah transaksi database berhasil di-commit.
// Jika EVENT_AUDIT_HOST kosong, pemanggilan diabaikan.
func PublishAuditEvent(ctx context.Context, event CreateEventRequest) error {
	if !auditEnabled() {
		auditSkipLogged.Do(func() {
			log.Println("audit event dilewati: EVENT_AUDIT_HOST belum di-set")
		})
		return nil
	}

	EnrichAuditFromContext(ctx, &event)

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	url := auditHost() + "/events"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := auditHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("audit event gagal: status %d url=%s", resp.StatusCode, url)
	}
	return nil
}

// PublishAuditAfterCommit mengirim audit di background setelah data lokal tersimpan.
// Kegagalan audit tidak mengubah hasil operasi database.
func PublishAuditAfterCommit(parent context.Context, event CreateEventRequest) {
	if !auditEnabled() {
		auditSkipLogged.Do(func() {
			log.Println("audit event dilewati: EVENT_AUDIT_HOST belum di-set")
		})
		return
	}

	EnrichAuditFromContext(parent, &event)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := PublishAuditEvent(ctx, event); err != nil {
			log.Printf("audit event gagal dikirim: %v", err)
		}
	}()
}
