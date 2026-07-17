package taskapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"read RFC 9110","priority":2}`))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response := httptest.NewRecorder()
	Create(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status=%d", response.Code)
	}
	if got := response.Header().Get("Location"); got != "/tasks/task-1" {
		t.Fatalf("Location=%q", got)
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["id"] != "task-1" {
		t.Fatalf("body=%v", body)
	}
}

func TestCreateFailures(t *testing.T) {
	tests := []struct {
		name, method, contentType, body string
		status                          int
	}{
		{"method", http.MethodGet, "application/json", `{}`, 405},
		{"media", http.MethodPost, "text/plain", `{}`, 415},
		{"syntax", http.MethodPost, "application/json", `{"title":`, 400},
		{"unknown", http.MethodPost, "application/json", `{"title":"x","priority":2,"prio":1}`, 400},
		{"semantic", http.MethodPost, "application/json", `{"title":" ","priority":4}`, 422},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/tasks", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			res := httptest.NewRecorder()
			Create(res, req)
			if res.Code != tc.status {
				t.Fatalf("status=%d, want=%d body=%s", res.Code, tc.status, res.Body.String())
			}
		})
	}
}
