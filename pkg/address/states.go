package address

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
)

const (
	getStatePath = "/states/"
)

type Address interface {
	GetState(name string) (state *State, err error)
	GetAllStates() (states *[]State, err error)
}

type Config struct {
	httpClient     *resty.Client
	URL            string
	RequestTimeout int64
	Verbose        bool
}

func (a Config) GetState(name string) (state *State, err error) {
	if a.Verbose {
		a.httpClient.SetDebug(true)
	}
	resp, err := a.httpClient.R().
		SetHeaders(map[string]string{
			"accept":       "application/json",
			"content-type": "application/json",
		}).
		SetResult(&State{}).
		Get(a.URL + getStatePath + name)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %s", resp.Error())
	}

	data := resp.Result().(*State)
	return data, nil
}

func (a Config) GetAllStates() (states *[]State, err error) {
	if a.Verbose {
		a.httpClient.SetDebug(true)
	}

	resp, err := a.httpClient.R().
		SetHeaders(map[string]string{
			"accept":       "application/json",
			"content-type": "application/json",
		}).
		SetResult([]State{}).
		Get(a.URL + getStatePath)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %s", resp.Error())
	}

	data := resp.Result().(*[]State)
	return data, nil
}

func New(cfg *Config) (Address, error) {
	_, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	// defaults request timeouts to 60secs
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 60
	}

	return &Config{
		httpClient:     resty.New(),
		URL:            cfg.URL,
		RequestTimeout: cfg.RequestTimeout,
		Verbose:        cfg.Verbose,
	}, nil
}
