package forget

import (
	"bytes"
	"testing"
)

var distResponse = `{"status_code":200,"status_txt":"","data":{"distribution":"colors","Z":148235,"T":1425056403,"data":[{"bin":"red","count":1,"p":6.746045131041927e-06},{"bin":"blue","count":1,"p":6.746045131041927e-06}]}}`
var distErrResponse = `{"status_code":500,"status_txt":"MISSING_ARG_DISTRIBUTION","data":null}`

func TestDistribution(t *testing.T) {
	client := mockClient(distResponse, 200)
	res, err := client.Distribution("colors")
	if err != nil {
		t.Fail()
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected 200 got %d", res.StatusCode)
	}
	if res.Distribution.Name != "colors" {
		t.Errorf("Expected 'colors' got '%s'", res.Distribution.Name)
	}
}

func TestMostProbable(t *testing.T) {
	client := mockClient(distResponse, 200)
	res, err := client.MostProbable("colors", 10)
	if err != nil {
		t.Fail()
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected 200 got %d", res.StatusCode)
	}
	if res.Distribution.Name != "colors" {
		t.Errorf("Expected 'colors' got '%s'", res.Distribution.Name)
	}
}

func TestField(t *testing.T) {
	client := mockClient(distResponse, 200)
	res, err := client.Field("colors", "925")
	if err != nil {
		t.Fail()
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected 200 got %d", res.StatusCode)
	}
	if res.Distribution.Name != "colors" {
		t.Errorf("Expected 'colors' got '%s'", res.Distribution.Name)
	}
}

func TestIncrement(t *testing.T) {
	client := mockClient("OK", 200)
	err := client.Increment("colors", "red")
	if err != nil {
		t.Fail()
	}
}

func TestIncrementByN(t *testing.T) {
	client := mockClient("OK", 200)
	err := client.IncrementByN("colors", "red", 3)
	if err != nil {
		t.Fail()
	}
}

func TestDatabaseSize(t *testing.T) {
	client := mockClient(`{"status_code":200,"status_txt":"","data":42}`, 200)
	res, err := client.DatabaseSize()
	if err != nil {
		t.Fail()
	}

	if res != 42 {
		t.Errorf("Expected 42 got %d", res)
	}
}

func TestError(t *testing.T) {
	client := mockClient(distResponse, 500)
	_, err := client.Distribution("colors")
	if err == nil {
		t.Fail()
	}

	client = mockClient(distErrResponse, 200)
	_, err = client.Distribution("colors")
	if err == nil {
		t.Fail()
	}
}

func mockClient(body string, status_code int) *Client {
	client := NewClient("http://forgettable.io:51000")
	client.C = NewHTTPMockClient(bytes.NewBufferString(body), status_code)
	return client
}
