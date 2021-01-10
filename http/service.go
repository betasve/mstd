package http

import (
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	HandleFunc(path string, handler func(http.ResponseWriter, *http.Request))
	ListenAndServe(addr string, handler http.Handler) error
}

type Client struct{}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	httpClient := http.Client{}
	return httpClient.Do(req)
}

func (c *Client) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, handler)
}

func (c *Client) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
