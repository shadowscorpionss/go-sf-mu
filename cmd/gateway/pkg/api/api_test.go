package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	config "sf-mu/pkg/configs"
	"testing"
)

func TestAPI_endpoints(t *testing.T) {

	cfg := config.New()
	api := New(cfg, cfg.News.AdrPort, cfg.Censor.AdrPort, cfg.Comments.AdrPort)

	var testBody1 = []byte(`{"newsID": 3,"content": "Тест qwerty "}`)
	var testBody2 = []byte(`{"newsID": 3,"content": "Тест ups "}`)
	var testBody3 = []byte(`{"id": 3}`)

	req := httptest.NewRequest(http.MethodGet, "/news", nil)
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/search?id=2", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody1))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody2))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusCreated)
	}

	req = httptest.NewRequest(http.MethodDelete, "/comments/del", bytes.NewBuffer(testBody3))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("error code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

}
