package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"sf-mu/pkg/models/news"
	newsstorage "sf-mu/pkg/storage/news"

	"testing"
	"time"
)

func TestAPI_endpoints(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := newsstorage.New(ctx, "postgres://postgres:postgres@localhost:5432/aggregator")
	if err != nil {
		t.Fatalf("could not connect to database: %v", err)
	}
	api := New(db)

	req := httptest.NewRequest(http.MethodGet, "/news?=page=2&s=", nil)
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}
	response := struct {
		Posts      []news.Post
		Pagination news.Pagination
	}{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
	b, err = io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}
	var data []news.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}
	const wantLen = 1
	if len(data) < wantLen {
		t.Fatalf("получено %d записей, ожидалось >= %d", len(data), wantLen)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/search?id=2", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
	b, err = io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}
	var post news.Post
	err = json.Unmarshal(b, &post)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/qwerty", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusNotFound) {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}
}
