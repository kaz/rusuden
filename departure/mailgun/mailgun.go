package mailgun

import (
	"context"
	"fmt"
	"strings"

	"github.com/kaz/rusuden/departure"
	"github.com/mailgun/mailgun-go/v4"
)

type (
	MailgunDeparture struct {
		client mailgun.Mailgun

		sender    string
		recipient string
	}
)

func New(key, sender, recipient string) departure.Departure {
	domain := ""
	if fragments := strings.Split(sender, "@"); len(fragments) == 2 {
		domain = fragments[1]
	}
	return &MailgunDeparture{mailgun.NewMailgun(domain, key), sender, recipient}
}

func (md *MailgunDeparture) Depart(ctx context.Context, text string) error {
	msg := md.client.NewMessage(md.sender, "【rusuden】留守番電話", text, md.recipient)

	if _, _, err := md.client.Send(ctx, msg); err != nil {
		return fmt.Errorf("md.client.Send failed: %w", err)
	}
	return nil
}
