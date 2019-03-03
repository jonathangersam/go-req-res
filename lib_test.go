package go_req_res

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_CaptureResponse(t *testing.T) {
	// arrange
	handlerResponseBody := sampleResponse{
		Name:  "foo",
		Age:   42,
		Fresh: true,
	}
	raw, e := json.Marshal(handlerResponseBody)
	if e != nil {
		t.Fatal(e)
	}

	mock := mockHandler{
		respStatus: http.StatusOK,
		respBody:   raw,
	}

	log := bytes.NewBuffer([]byte{})
	wrappedHandler := CaptureResponse(log, mock.handler)

	ts := httptest.NewServer(http.HandlerFunc(wrappedHandler))
	defer ts.Close()
	client := ts.Client()

	req, e := http.NewRequest(http.MethodGet, ts.URL, strings.NewReader("foo"))
	if e != nil {
		t.Fatal(e)
	}

	// act
	res, e := client.Do(req)

	// assert
	if e != nil {
		t.Fatal(e)
	}

	body := bytes.NewBufferString("")
	io.Copy(body, res.Body)

	// http response should contain actual response
	if strings.Contains(body.String(), "foo") == false {
		t.Errorf("Expected response body to contain substring 'foo', but didn't. Contents: %s", body.String())
	}

	// http response should contain expected http status
	if res.StatusCode != mock.respStatus {
		t.Errorf("Expected response http status code to be %d, but got %d", mock.respStatus, res.StatusCode)
	}

	// logger should contain http response
	if strings.Contains(log.String(), "foo") == false {
		t.Errorf("Expected log to contain substring 'foo', but didn't. Contents: %s", log.String())
	}
	//t.Log("[log]", log.String())
}

func Test_CaptureRequest(t *testing.T) {
	// arrange
	handlerResponseBody := sampleResponse{
		Name:  "foo",
		Age:   42,
		Fresh: true,
	}
	raw, e := json.Marshal(handlerResponseBody)
	if e != nil {
		t.Fatal(e)
	}

	mock := mockHandler{
		respStatus: http.StatusOK,
		respBody:   raw,
	}

	log := bytes.NewBuffer([]byte{})
	wrappedHandler := CaptureRequest(log, mock.handler)

	ts := httptest.NewServer(http.HandlerFunc(wrappedHandler))
	defer ts.Close()
	client := ts.Client()

	reqBody := sampleRequest{
		Name: "bar",
	}
	rawReqBody, e := json.Marshal(reqBody)
	if e != nil {
		t.Fatal(e)
	}

	req, e := http.NewRequest(http.MethodGet, ts.URL, bytes.NewBuffer(rawReqBody))
	if e != nil {
		t.Fatal(e)
	}

	// act
	res, e := client.Do(req)

	// assert
	if e != nil {
		t.Fatal(e)
	}

	body := bytes.NewBufferString("")
	io.Copy(body, res.Body)

	// http response should contain actual response
	if strings.Contains(body.String(), "foo") == false {
		t.Errorf("Expected response body to contain substring 'foo', but didn't. Contents: %s", body.String())
	}

	// http response should contain expected http status
	if res.StatusCode != mock.respStatus {
		t.Errorf("Expected response http status code to be %d, but got %d", mock.respStatus, res.StatusCode)
	}

	// logger should contain http response
	if strings.Contains(log.String(), "bar") == false {
		t.Errorf("Expected log to contain substring 'bar', but didn't. Contents: %s", log.String())
	}

	// handler should have received the request body.
	if mock.receivedName != "bar" {
		t.Errorf("Expected handler to have received Name 'bar' in request body, but got %s", mock.receivedName)
	}
	t.Log("[log]", log.String())
}

type mockHandler struct {
	receivedName string

	respStatus int
	respBody   []byte
}

func (m *mockHandler) handler(w http.ResponseWriter, r *http.Request) {
	var req sampleRequest
	e := json.NewDecoder(r.Body).Decode(&req)
	if e != nil {
		panic(e)
	}
	m.receivedName = req.Name

	w.WriteHeader(m.respStatus)
	w.Write(m.respBody)
}

type sampleRequest struct {
	Name string
}

type sampleResponse struct {
	Name  string
	Age   int
	Fresh bool
}
