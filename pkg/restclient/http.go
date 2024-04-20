package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func DoHTTPRequest(ctx context.Context, method string, data any, url string) (*http.Response, error) {
	client := &http.Client{}

	var requestBody []byte
	if method == http.MethodPost {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		requestBody = body
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	log.Info().Msgf("making %s request to: %s", method, req.URL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return resp, nil
}
