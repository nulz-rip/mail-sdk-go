package nulzmail

import "time"

// Inbox from API
type Inbox struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

// Message from API
type Message struct {
	ID       string  `json:"id"`
	Subject  string  `json:"subject"`
	FromAddr string  `json:"from_addr"`
	ToAddr   string  `json:"to_addr"`
	BodyText *string `json:"body_text"`
	BodyHTML *string `json:"body_html"`
	Received string  `json:"received_at"`
}

// MessageSummary in list
type MessageSummary struct {
	ID       string `json:"id"`
	Subject  string `json:"subject"`
	FromAddr string `json:"from_addr"`
	ToAddr   string `json:"to_addr"`
	Received string `json:"received_at"`
}

// MessagesPage list response (API returns messages, total)
type MessagesPage struct {
	Messages []MessageSummary `json:"messages"`
	Total    int64            `json:"total"`
	Cursor   string           `json:"cursor,omitempty"` // optional; use as offset for next page
}

// WaitOpts for WaitForCode
type WaitOpts struct {
	PollInterval time.Duration // default 2s
	Timeout      time.Duration // default 60s
}
