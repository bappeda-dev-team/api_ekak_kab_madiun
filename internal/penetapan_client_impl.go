package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type PenetapanClientImpl struct {
	BaseClient
}

func NewPenetapanClient(httpClient *http.Client) *PenetapanClientImpl {
	penetapanHost := os.Getenv("PENETAPAN_SERVICE_HOST")
	return &PenetapanClientImpl{
		BaseClient: newBaseClient(
			penetapanHost,
			"",
			httpClient,
		),
	}
}

type SyncPenetapanPkPegawaiRequest struct {
	PegawaiID string `json:"pegawai_id"`
	KodeOpd   string `json:"kode_opd"`
	Tahun     int    `json:"tahun"`
}

func (c *PenetapanClientImpl) SyncPenetapanPkPegawai(
	ctx context.Context,
	pegawaiId string,
	kodeOpd string,
	tahun int,
) error {
	url := fmt.Sprintf("%s/individu/rekin/sync", c.host)

	payload := SyncPenetapanPkPegawaiRequest{
		PegawaiID: pegawaiId,
		KodeOpd:   kodeOpd,
		Tahun:     tahun,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("gagal membuat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Session Id tidak ditemukan, mungkin akan 401")
	}

	log.Printf("POST %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gagal request sync penetapan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"unexpected status: %d, response: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}
