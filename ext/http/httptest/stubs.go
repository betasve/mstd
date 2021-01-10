package httptest

import (
	"io"
	"io/ioutil"
	"net/http"
	httptest "net/http/httptest"
	"strings"
)

type ClientMock struct{}
type BodyMock struct{}

var HandlerStubFn = func(in string, handler func(w http.ResponseWriter, r *http.Request)) {}
var ListenAndServeStubFn = func(addr string, handler http.Handler) error { return nil }
var DefaultMockFn = func(req *http.Request) (*http.Response, error) {
	res := &http.Response{}
	res.Body = StubbedBody()
	return res, nil
}

var MockFn = DefaultMockFn

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return MockFn(req)
}

func (c *ClientMock) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	HandlerStubFn(path, handler)
}

func (c *ClientMock) ListenAndServe(addr string, handler http.Handler) error {
	return ListenAndServeStubFn(addr, handler)
}

func StubbedBody() io.ReadCloser {
	return ioutil.NopCloser(
		strings.NewReader(`{"token_type":"Bearer","scope":"Tasks.ReadWrite.Shared Tasks.ReadWrite User.Read Mail.Read","expires_in":3600,"ext_expires_in":3600,"access_token":"EwB","refresh_token":"M.R3_BAY"}`),
	)
}

func StubbedResponseWriter() http.ResponseWriter {
	return httptest.NewRecorder()
}
