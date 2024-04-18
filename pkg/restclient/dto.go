package restclient

import "github.com/go-resty/resty/v2"

type Config struct {
	HttpClient     *resty.Client
	URL            string
	RequestTimeout int64
	Verbose        bool
}
