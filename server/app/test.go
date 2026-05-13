package app

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

	"rabidaudio.com/coven-door/server/database"
)

func NewTest(t *testing.T, routers ...RouteBuilder) *TestApp {
	t.Helper()

	ctx := database.PrepareForTest(t)
	app := New(database.Get(ctx), routers...)

	s := httptest.NewServer(app)
	t.Cleanup(func() {
		s.Close()
	})
	tapp := TestApp{Server: s, Header: http.Header{}, Context: ctx}
	return &tapp
}

type TestApp struct {
	*httptest.Server
	http.Header
	context.Context
	http.Client
}

type TestResponse struct {
	*http.Response
	Data  map[string]any
	Error map[string]any
}

func (ta *TestApp) DB() database.DB {
	return database.Get(ta.Context)
}

func (ta *TestApp) SetHeader(header, value string) {
	ta.Header[header] = []string{value}
}

func (ta *TestApp) UnsetHeader(header string) {
	ta.Header[header] = []string{}
}

func (ta *TestApp) Request(method, path string, body any) (*TestResponse, error) {
	url := fmt.Sprintf("%v%v", ta.Server.URL, path)
	ta.SetHeader("Content-Type", "application/json")
	ta.SetHeader("Accept", "application/json")

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

	req, err := http.NewRequestWithContext(ta, method, url, br)
	if err != nil {
		return nil, err
	}
	maps.Copy(req.Header, ta.Header)
	res, err := ta.Do(req)
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
		// TODO: pagination
	}
	return &tr, nil
}
