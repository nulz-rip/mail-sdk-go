// Example: create inbox and wait for verification code.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nulz-rip/mail-sdk-go/nulzmail"
)

func main() {
	if os.Getenv("NULZ_API_KEY") == "" {
		log.Fatal("set NULZ_API_KEY")
	}
	client := nulzmail.New()
	ctx := context.Background()

	inbox, err := client.CreateInbox(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inbox:", inbox.Address)

	// Quick check: do we see any messages? (same account as dashboard)
	page, _ := client.ListMessages(ctx, inbox.ID, "")
	fmt.Println("Messages in inbox:", len(page.Messages))
	if len(page.Messages) > 0 {
		fmt.Println("First subject:", page.Messages[0].Subject)
	}

	_, msg, err := client.WaitForCode(ctx, inbox.ID, nulzmail.WaitOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("From:", msg.FromAddr, "Subject:", msg.Subject)
}
