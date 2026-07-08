package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type IsustrategicClient struct {
	BaseClient
}

func NewIsuStrategicClient(host string, httpClient *http.Client) *IsustrategicClient {
	return &IsustrategicClient{
		BaseClient: newBaseClient("https://isu-strategis-dev.zeabur.app", "", httpClient),
	}
}

func (c *IsustrategicClient) GetDataIsuStrategic(ctx context.Context, kodeOpd string, tahun string) ([]IsuStrategisResponse, error) {
	url := fmt.Sprintf("%s/isu_strategis/kebelakang/%s/%s", c.host, kodeOpd, tahun)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("isu strategis di panggil")

	// log.Printf("url:%s ", url)

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Session Id ditemukan, mungkin akan 401")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper IsuStrategicWrapper
    if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
        return nil, fmt.Errorf("gagal decode response: %w", err)
    }

    if len(wrapper.Data) == 0 {
        return nil, nil
    }

    return wrapper.Data, nil
}
func (c *IsustrategicClient) GetDataPermasalahan(ctx context.Context, kodeOpd string, tahun string) ([]PermasalahanResp, error) {
	url := fmt.Sprintf("%s/isu_strategis/kebelakang/%s/%s", c.host, kodeOpd, tahun)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("isu strategis di panggil")

	// log.Printf("url:%s ", url)

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Session Id ditemukan, mungkin akan 401")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper PermasalahanWrapper
    if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
        return nil, fmt.Errorf("gagal decode response: %w", err)
    }

    if len(wrapper.Data) == 0 {
        return nil, nil
    }

    return wrapper.Data, nil
}
