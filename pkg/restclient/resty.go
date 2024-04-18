package restclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func GetResty(client *resty.Client, data any, config Config, path string) (*resty.Response, error) {
	if config.Verbose {
		client.SetDebug(true)
	}

	// Make the GET request
	resp, err := client.R().
		SetHeaders(map[string]string{
			"accept":       "application/json",
			"content-type": "application/json",
		}).
		SetResult(data).
		Get(config.URL + path)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP error: %s", resp.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status())
	}

	return resp, nil
}
