package nulzmail

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"aliases":[]}`))
	}))
	defer srv.Close()

	c := New("secret-key")
	c.SetBaseURL(srv.URL)
	_, err := c.ListInboxes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "ApiKey secret-key" {
		t.Errorf("auth header: got %q", gotAuth)
	}
}

func TestErrorParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{
				"code":    "invalid",
				"message": "bad request",
			},
		})
	}))
	defer srv.Close()

	c := New("key")
	c.SetBaseURL(srv.URL)
	_, err := c.ListInboxes(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("status: got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "bad request" {
		t.Errorf("message: got %q", apiErr.Message)
	}
}

func TestWaitForCodeReturnsFirstMessage(t *testing.T) {
	// WaitForCode returns first message, code always ""
	listCalls := 0
	getCalls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/aliases/x/messages" {
			listCalls++
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"messages":[{"id":"m1","subject":"hi","from_addr":"a@b","to_addr":"c@d","received_at":""}],"total":1}`))
			return
		}
		if r.URL.Path == "/messages/m1" {
			getCalls++
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":"m1","subject":"hi","from_addr":"a@b","to_addr":"c@d","received_at":"","body_text":null,"body_html":null}`))
			return
		}
		t.Fatalf("unexpected path %s", r.URL.Path)
	}))
	defer srv.Close()

	c := New("key")
	c.SetBaseURL(srv.URL)
	code, msg, err := c.WaitForCode(context.Background(), "x", WaitOpts{Timeout: 1, PollInterval: 1})
	if err != nil {
		t.Fatal(err)
	}
	if code != "" {
		t.Errorf("code: got %q (want empty)", code)
	}
	if msg.Subject != "hi" || msg.FromAddr != "a@b" {
		t.Errorf("msg: subject=%q from=%q", msg.Subject, msg.FromAddr)
	}
	if listCalls < 1 || getCalls < 1 {
		t.Errorf("listCalls=%d getCalls=%d", listCalls, getCalls)
	}
}
