package nulzmail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError from JSON response
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("nulz mail api: %s (status %d)", e.Message, e.StatusCode)
	}
	if e.Body != "" {
		return fmt.Sprintf("nulz mail api: status %d: %s", e.StatusCode, truncate(e.Body, 200))
	}
	return fmt.Sprintf("nulz mail api: status %d", e.StatusCode)
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// api error JSON: "error" string or { "error": { "code", "message" } }

func parseError(resp *http.Response) *APIError {
	e := &APIError{StatusCode: resp.StatusCode}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	e.Body = strings.TrimSpace(string(body))
	var p struct {
		Error interface{} `json:"error"`
	}
	if json.Unmarshal(body, &p) == nil && p.Error != nil {
		switch v := p.Error.(type) {
		case string:
			e.Message = v
		case map[string]interface{}:
			if c, ok := v["code"].(string); ok {
				e.Code = c
			}
			if m, ok := v["message"].(string); ok {
				e.Message = m
			}
		}
	}
	return e
}
