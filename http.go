package forget

import (
	"bytes"
	"net/http"
)

// HTTPClient represents a type that can make HTTP requests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTPDefaultClient is the default Requester implementation.
type HTTPDefaultClient struct {
	c *http.Client
}

// Do sends an HTTP request and returns an HTTP response.
func (self *HTTPDefaultClient) Do(req *http.Request) (*http.Response, error) {
	return self.c.Do(req)
}

// MockBody implements io.ReadCloser.
type MockBody struct {
	*bytes.Buffer
}

// Closer mocks the Close() method.
func (self *MockBody) Close() error {
	return nil
}

// HTTPMockClient is a HTTPClient implementation used for testing.
type HTTPMockClient struct {
	Body       *MockBody
	StatusCode int
	Error      error
}

// NewHTTPMockClient returns a *HTTPMockClient instance.
func NewHTTPMockClient(body *bytes.Buffer, sc int) *HTTPMockClient {
	return &HTTPMockClient{
		Body:       &MockBody{body},
		StatusCode: sc,
	}
}

// Do sends a mock request and returns a response.
func (self *HTTPMockClient) Do(_ *http.Request) (*http.Response, error) {
	if self.Error != nil {
		return nil, self.Error
	}

	return &http.Response{
		Body:       self.Body,
		StatusCode: self.StatusCode,
	}, nil
}
