package nulzmail

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const defaultBase = "https://v1.nulz.lol/v1"
const defaultTimeout = 30 * time.Second

// Client for nulz mail API
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// New client. Key from arg or NULZ_API_KEY env.
func New(apiKey ...string) *Client {
	key := ""
	if len(apiKey) > 0 && apiKey[0] != "" {
		key = apiKey[0]
	} else {
		key = os.Getenv("NULZ_API_KEY")
	}
	return &Client{
		baseURL: strings.TrimSuffix(defaultBase, "/"),
		apiKey:  key,
		http: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// SetBaseURL override
func (c *Client) SetBaseURL(u string) {
	c.baseURL = strings.TrimSuffix(u, "/")
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, u, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "ApiKey "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseError(resp)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// CreateInbox creates one alias (inbox). Sends {} so backend gets valid JSON.
func (c *Client) CreateInbox(ctx context.Context) (Inbox, error) {
	return c.CreateInboxWithPrefix(ctx, "")
}

// CreateInboxWithPrefix creates alias with optional prefix (e.g. "femboy").
func (c *Client) CreateInboxWithPrefix(ctx context.Context, prefix string) (Inbox, error) {
	var body any = struct{}{}
	if prefix != "" {
		body = struct{ Prefix string `json:"prefix"` }{Prefix: prefix}
	}
	var out Inbox
	err := c.do(ctx, "POST", "/aliases", body, &out)
	return out, err
}

// ListInboxes returns all aliases
func (c *Client) ListInboxes(ctx context.Context) ([]Inbox, error) {
	var out struct {
		Aliases []Inbox `json:"aliases"`
	}
	err := c.do(ctx, "GET", "/aliases", nil, &out)
	if err != nil {
		return nil, err
	}
	if out.Aliases == nil {
		return []Inbox{}, nil
	}
	return out.Aliases, nil
}

// DeleteInbox by id
func (c *Client) DeleteInbox(ctx context.Context, inboxID string) error {
	path := "/aliases/" + url.PathEscape(inboxID)
	return c.do(ctx, "DELETE", path, nil, nil)
}

// ListMessages with optional cursor (API uses limit/offset; cursor ignored)
func (c *Client) ListMessages(ctx context.Context, inboxID, cursor string) (MessagesPage, error) {
	path := "/aliases/" + url.PathEscape(inboxID) + "/messages?limit=50"
	if cursor != "" {
		path = "/aliases/" + url.PathEscape(inboxID) + "/messages?limit=50&offset=" + url.QueryEscape(cursor)
	}
	var out MessagesPage
	err := c.do(ctx, "GET", path, nil, &out)
	return out, err
}

// GetMessage by id
func (c *Client) GetMessage(ctx context.Context, messageID string) (Message, error) {
	path := "/messages/" + url.PathEscape(messageID)
	var out Message
	err := c.do(ctx, "GET", path, nil, &out)
	return out, err
}
