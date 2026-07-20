package localservice

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	baseURL, stop, err := Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := stop(ctx); err != nil {
			t.Errorf("shutdown: %v", err)
		}
	})
	client := &http.Client{Timeout: time.Second}
	response, err := client.Get(baseURL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || string(body) != "ok\n" {
		t.Fatalf("status=%d body=%q", response.StatusCode, body)
	}
}
