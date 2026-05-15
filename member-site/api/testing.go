package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// NewTestServer provides an interface for end-to-end testing an [http.Handler].
// It sets up an end-to-end server on a random port, provides methods for generating
// requests, and automatically serializes and deserializes JSON.
func NewTestServer(t *testing.T, ctx context.Context, handler http.Handler) *TestServer {
	t.Helper()

	s := httptest.NewServer(handler)
	t.Cleanup(func() {
		s.Close()
	})
	if ctx == nil {
		ctx = t.Context()
	}
	ts := TestServer{Server: s, Header: http.Header{}, Context: ctx}
	return &ts
}

type TestServer struct {
	*httptest.Server
	Header http.Header
	context.Context
	client http.Client
}

type TestResponse struct {
	*http.Response
	Data       map[string]any
	Pagination map[string]any
	Error      map[string]any
}

func (ts *TestServer) SetHeader(header, value string) {
	ts.Header[header] = []string{value}
}

func (ts *TestServer) UnsetHeader(header string) {
	ts.Header[header] = []string{}
}

func (ts *TestServer) Request(method, path string, body any) (*TestResponse, error) {
	url := fmt.Sprintf("%v%v", ts.Server.URL, path)
	ts.SetHeader("Content-Type", "application/json")
	ts.SetHeader("Accept", "application/json")

	var br io.Reader
	if body == nil {
		br = bytes.NewReader([]byte{}) // empty
	} else if rbody, ok := body.(io.Reader); ok {
		br = rbody
	} else if strbody, ok := body.(string); ok {
		br = strings.NewReader(strbody)
	} else {
		// serialize as json
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		br = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ts, method, url, br)
	if err != nil {
		return nil, err
	}
	maps.Copy(req.Header, ts.Header)
	res, err := ts.client.Do(req)
	if err != nil {
		return nil, err
	}

	tr := TestResponse{Response: res, Data: map[string]any{}}
	if res.Header["Content-Type"][0] == "application/json" {
		jd := map[string]any{}
		if err := json.NewDecoder(res.Body).Decode(&jd); err != nil {
			return nil, err
		}
		if jd["status"] == "ERROR" {
			maps.Copy(tr.Error, jd["error"].(map[string]any))
		} else {
			maps.Copy(tr.Data, jd["data"].(map[string]any))
		}
		if jd["pagination"] != nil {
			maps.Copy(tr.Pagination, jd["pagination"].(map[string]any))
		}
	}
	return &tr, nil
}
