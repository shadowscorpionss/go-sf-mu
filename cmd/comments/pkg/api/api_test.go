package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	comments "sf-mu/pkg/models/comments"
	commentsstorage "sf-mu/pkg/storage/comments"

	"testing"
	"time"
)

func TestCommentHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	psgr, err := commentsstorage.New(ctx, "postgres://postgres:postgres@localhost:5432/comm")
	if err != nil {
		t.Fatal(err)
	}
	api := New(psgr)

	var testBody = []byte(`{"newsID": 1,"content": "Тест"}`)

	req := httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody))
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/comments?news_id=1", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}
	var data []comments.Comment
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("failed parsing server response: %v", err)
	}
	const wantLen = 1
	if len(data) < wantLen {
		t.Fatalf("got %d rows, wanted %d", len(data), wantLen)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/add", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusConflict)
	}

}
