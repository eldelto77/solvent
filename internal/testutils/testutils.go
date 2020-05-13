package testutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func AssertEquals(t *testing.T, expected interface{}, actual interface{}, title string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v should be '%v' but was '%v'", title, expected, actual)
	}
}

func AssertNotEquals(t *testing.T, expected interface{}, actual interface{}, title string) {
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("%v should not be '%v' but was '%v'", title, expected, actual)
	}
}

type Response struct {
	response   *http.Response
	Body       map[string]interface{}
	StatusCode int
}

func NewResponse(response *http.Response) (Response, error) {
	var body map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&body)
	if err != nil && !errors.Is(err, io.EOF) {
		return Response{}, err
	}

	parsedResponse := Response{
		response:   response,
		Body:       body,
		StatusCode: response.StatusCode,
	}

	return parsedResponse, nil
}

type TestServer struct {
	*httptest.Server
	T      *testing.T
	Client *http.Client
}

func NewTestServer(t *testing.T, handler http.Handler) *TestServer {
	ts := httptest.NewServer(handler)
	return &TestServer{
		Server: ts,
		T:      t,
		Client: http.DefaultClient,
	}
}

func (ts *TestServer) GET(path string) Response {
	return ts.request("GET", path, "")
}

func (ts *TestServer) POST(path string, body string) Response {
	return ts.request("POST", path, body)
}

func (ts *TestServer) PUT(path string, body string) Response {
	return ts.request("PUT", path, body)
}

func (ts *TestServer) request(verb, path string, body string) Response {
	url := ts.URL + path
	bodyData := bytes.NewBufferString(body)

	req, err := http.NewRequest(verb, url, bodyData)
	if err != nil {
		ts.T.Fatalf("http.NewRequest error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := ts.Client.Do(req)
	if err != nil {
		ts.T.Fatalf("ts.Client.Do error: %v", err)
	}
	defer response.Body.Close()

	parsedResponse, err := NewResponse(response)
	if err != nil {
		ts.T.Fatalf("NewResponse error: %v", err)
	}

	return parsedResponse
}
