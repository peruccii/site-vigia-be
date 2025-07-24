package services

import (
	"fmt"
	"net/url"

	"github.com/resend/resend-go/v2"
)

func SendEmail(email, token string) {
	apiKey := "re_xxxxxxxxx"

	client := resend.NewClient(apiKey)

	link := fmt.Sprintf("https://site-vigia.com.br/recuperar-senha?token=%s", url.QueryEscape(token))

	subject := fmt.Sprintf("Link para recuperar senha: %s", link)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{"delivered@resend.dev"},
		Html:    "<strong>hello world</strong>",
		Subject: subject,
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(sent.Id)
}
