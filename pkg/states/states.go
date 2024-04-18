package states

import (
	"github.com/go-resty/resty/v2"
	"github.com/leetatech/leeta_backend/pkg/restclient"
	"net/url"
)

const (
	getStatePath = "/states/"
)

type StateMethods interface {
	GetState(name string) (*State, error)
	GetAllStates() (states *[]State, err error)
}

type Config struct {
	httpClient     *resty.Client
	URL            string
	RequestTimeout int64
	Verbose        bool
}

func (a Config) GetState(name string) (*State, error) {
	var state State
	_, err := restclient.GetResty(a.httpClient, &state, restclient.Config{HttpClient: a.httpClient, URL: a.URL, RequestTimeout: a.RequestTimeout, Verbose: a.Verbose}, getStatePath+name)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (a Config) GetAllStates() (*[]State, error) {

	var states []State
	_, err := restclient.GetResty(a.httpClient, &states, restclient.Config{HttpClient: a.httpClient, URL: a.URL, RequestTimeout: a.RequestTimeout, Verbose: a.Verbose}, getStatePath)
	if err != nil {
		return nil, err
	}

	return &states, nil
}

func New(cfg *Config) (StateMethods, error) {
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
