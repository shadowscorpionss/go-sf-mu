package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	censorsstorage "sf-mu/pkg/storage/censors"

	"testing"
	"time"
)

func TestCommentHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	psgr, err := censorsstorage.New(ctx, "postgres://postgres:postgres@localhost:5432/comm")
	if err != nil {
		t.Fatal(err)
	}

	api := New(psgr)

	var testBody = []byte(`{"newsID": 1,"content": "улий"}`)
	var testBody2 = []byte(`{"newsID": 1,"content": "Тест анус "}`)
	var testBody3 = []byte(`{"newsID": 1,"content": "Тест блядво "}`)
	var testBody4 = []byte(`{"newsID": 1,"content": "Тест въебать "}`)
	var testBody5 = []byte(`{"newsID": 1,"content": "Тест ups "}`)

	req := httptest.NewRequest(http.MethodPost, "/comments/check", bytes.NewBuffer(testBody))
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/check", bytes.NewBuffer(testBody2))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/check", bytes.NewBuffer(testBody3))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/check", bytes.NewBuffer(testBody4))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/stop", bytes.NewBuffer(testBody5))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}
}
