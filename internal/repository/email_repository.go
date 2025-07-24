package repository

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type EmailRepository interface {
	SendEmail(email string) string
}

type Provider struct{}

func (p *Provider) SendEmail(email string) string {
	apiKey := "re_xxxxxxxxx"

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{"delivered@resend.dev"},
		Html:    "<strong>hello world</strong>",
		Subject: "Hello from Golang",
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return fmt.Sprint("Email sent with ID ", sent.Id)
}
