package restclient

import "github.com/go-resty/resty/v2"

type Config struct {
	HttpClient     *resty.Client
	URL            string
	RequestTimeout int64
	Verbose        bool
}

type ResponseBody[T any] struct {
	Data T `json:"data"`
}

type ResponseBodyList[T any] struct {
	Data []T `json:"data"`
}
