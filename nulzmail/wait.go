package nulzmail

import (
	"context"
	"fmt"
	"os"
	"time"
)

const defaultPoll = 2 * time.Second
const defaultWaitTimeout = 60 * time.Second

// WaitForCode polls inbox until at least one message exists, then returns it.
// code is always "" (no extraction from raw/encoded body).
func (c *Client) WaitForCode(ctx context.Context, inboxID string, opts WaitOpts) (code string, msg Message, err error) {
	poll := opts.PollInterval
	if poll <= 0 {
		poll = defaultPoll
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = defaultWaitTimeout
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		page, err := c.ListMessages(ctx, inboxID, "")
		if err != nil {
			return "", Message{}, err
		}
		if os.Getenv("NULZ_DEBUG") != "" && len(page.Messages) > 0 {
			fmt.Fprintf(os.Stderr, "[nulz] poll: %d messages (e.g. subject %q)\n", len(page.Messages), page.Messages[0].Subject)
		}
		if len(page.Messages) == 0 {
			select {
			case <-ctx.Done():
				return "", Message{}, ctx.Err()
			case <-time.After(poll):
				continue
			}
		}
		// Return first (newest) message; caller can parse code from subject/body if needed
		m := page.Messages[0]
		full, err := c.GetMessage(ctx, m.ID)
		if err == nil {
			return "", full, nil
		}
		return "", summaryToMessage(m), nil
	}
	return "", Message{}, &APIError{StatusCode: 408, Message: "timeout waiting for code"}
}

func summaryToMessage(m MessageSummary) Message {
	return Message{
		ID: m.ID, Subject: m.Subject, FromAddr: m.FromAddr, ToAddr: m.ToAddr, Received: m.Received,
	}
}
