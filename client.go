package forget

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var statusCodeExpectedValue = 200
var incrExpectedValue = []byte("OK")

// Value stores the values for a distribution field.
type Value struct {
	Field       string  `json:"bin"`
	Count       int     `json:"count"`
	Probability float64 `json:"p"`
}

// Distribution stores the values for a distribution.
type Distribution struct {
	Name       string  `json:"distribution"`
	Values     []Value `json:"data"`
	Z          int     `json:"Z"`
	Time       int     `json:"T"`
	Rate       float64 `json:"rate"`
	Prune      bool    `json:"prune"`
	isFull     bool
	hasDecayed bool
}

// Response holds the values returned from a Forgettable API request.
type Response struct {
	StatusCode   int    `json:"status_code"`
	StatusTxt    string `json:"status_txt"`
	Distribution `json:"data"`
}

// databaseSizeResponse holds the values returned from a Forgettable "dbsize" request.
type databaseSizeResponse struct {
	StatusCode int    `json:"status_code"`
	StatusTxt  string `json:"status_txt"`
	Size       int    `json:"data"`
}

// Client makes requests to Forgettable servers.
type Client struct {
	RootURL string
	C       HTTPClient
}

// NewClient creates and returns a *Client instance.
func NewClient(ru string) *Client {
	return &Client{ru, &HTTPDefaultClient{http.DefaultClient}}
}

// Distribution returns the values for the given distribution.
func (self *Client) Distribution(distribution string) (*Response, error) {
	vals := url.Values{}
	vals.Add("distribution", distribution)

	return self.send("/dist", vals)
}

// MostProbable returns the N most probable values from the given distribution.
func (self *Client) MostProbable(distribution string, n int) (*Response, error) {
	vals := url.Values{}
	vals.Add("distribution", distribution)
	vals.Add("N", strconv.Itoa(n))

	return self.send("/nmostprobable", vals)
}

// Field returns the values for a single distribution field.
func (self *Client) Field(distribution, field string) (*Response, error) {
	vals := url.Values{}
	vals.Add("distribution", distribution)
	vals.Add("field", field)

	return self.send("/get", vals)
}

// Increment increments the value of a distribution field by a single point.
// The call was a success when the method does not return an error.
func (self *Client) Increment(distribution, field string) error {
	vals := url.Values{}
	vals.Add("distribution", distribution)
	vals.Add("field", field)
	body, err := self.request("/incr", vals)
	if err != nil {
		return err
	}

	return makeIncrementResponse(body)
}

// IncrementByN increments the value of a distribution field by the value of n.
// The call was a success when the method does not return an error.
func (self *Client) IncrementByN(distribution, field string, n int) error {
	vals := url.Values{}
	vals.Add("distribution", distribution)
	vals.Add("field", field)
	vals.Add("N", strconv.Itoa(n))
	body, err := self.request("/incr", vals)
	if err != nil {
		return err
	}

	return makeIncrementResponse(body)
}

// DatabaseSize returns the size of the Forgettable database.
func (self *Client) DatabaseSize() (int, error) {
	body, err := self.request("/dbsize", nil)
	if err != nil {
		return 0, err
	}

	res, err := makeDatabaseSizeResponse(body)
	if err != nil {
		return 0, err
	}

	return res.Size, nil
}

// send sends a request to the Forgettable server and returns a *Response.
func (self *Client) send(endpoint string, vals url.Values) (*Response, error) {
	body, err := self.request(endpoint, vals)
	if err != nil {
		return nil, err
	}

	dist, err := makeDistributionResponse(body)
	if err != nil {
		return nil, err
	}
	if dist.StatusCode != 200 {
		return nil, errors.New(dist.StatusTxt)
	}

	return dist, nil
}

// request sends a request to the Forgettable server and returns the response body.
func (self *Client) request(endpoint string, vals url.Values) ([]byte, error) {
	var url string
	if vals != nil {
		url = fmt.Sprintf("%s/%s?%s", self.RootURL, endpoint, vals.Encode())
	} else {
		url = fmt.Sprintf("%s/%s", self.RootURL, endpoint)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	res, err := self.C.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return []byte{}, errors.New(
			fmt.Sprintf("Got response status code %d.", res.StatusCode))
	}

	return ioutil.ReadAll(res.Body)
}

// makeDistributionResponse turns a response body into a *Response instance.
func makeDistributionResponse(body []byte) (*Response, error) {
	dist := &Response{}
	err := json.Unmarshal(body, dist)
	if err != nil {
		return nil, err
	}

	return dist, nil
}

// makeIncrementResponse turns a response body into a increment error.
func makeIncrementResponse(body []byte) error {
	if !bytes.Equal(body, incrExpectedValue) {
		res, err := makeDistributionResponse(body)
		if err != nil {
			return err
		}

		return errors.New(res.StatusTxt)
	}

	return nil
}

// makeDatabaseSizeResponse turns a response body into a *databaseSizeResponse.
func makeDatabaseSizeResponse(body []byte) (*databaseSizeResponse, error) {
	res := &databaseSizeResponse{}
	err := json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != statusCodeExpectedValue {
		return nil, errors.New(res.StatusTxt)
	}

	return res, nil
}
