package taskapi

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
	"unicode/utf8"
)

type createInput struct {
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}
type problem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, problem{"METHOD_NOT_ALLOWED", "use POST"})
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType, problem{"UNSUPPORTED_MEDIA_TYPE", "use application/json"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input createInput
	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, problem{"INVALID_JSON", "request JSON is invalid"})
		return
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, problem{"INVALID_JSON", "one JSON value is required"})
		return
	}
	input.Title = strings.TrimSpace(input.Title)
	length := utf8.RuneCountInString(input.Title)
	if length < 1 || length > 100 || input.Priority < 1 || input.Priority > 3 {
		writeJSON(w, http.StatusUnprocessableEntity, problem{"INVALID_TASK", "title or priority is invalid"})
		return
	}
	output := struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Priority int    `json:"priority"`
	}{"task-1", input.Title, input.Priority}
	w.Header().Set("Location", "/tasks/task-1")
	writeJSON(w, http.StatusCreated, output)
}
