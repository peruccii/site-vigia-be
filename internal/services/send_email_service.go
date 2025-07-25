package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"

	"github.com/resend/resend-go/v2"
)

type EmailData struct {
	Link string
}

func SendEmail(email, token string) {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	link := fmt.Sprintf("http://site-vigia.com.br/recuperar-senha?token=%s", url.QueryEscape(token))

	subject := "Recuperação de Senha"

	htmlTemplate := `<!DOCTYPE html>
        <html lang="pt-BR">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Recuperação de Senha</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    margin: 0;
                    padding: 0;
                }
                .container {
                    max-width: 600px;
                    margin: 50px auto;
                    background: #ffffff;
                    padding: 20px;
                    border-radius: 5px;
                    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
                }
                h1 {
                    color: #333;
                }
                p {
                    color: #555;
                }
                .btn {
                    display: inline-block;
                    padding: 10px 20px;
                    background-color: #007bff;
                    color: #ffffff;
                    text-decoration: none;
                    border-radius: 5px;
                    margin: 10px 0;
                }
                .btn:hover {
                    background-color: #0056b3;
                }
                .footer {
                    margin-top: 20px;
                    font-size: 12px;
                    color: #888;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>Recuperação de Senha</h1>
                <p>Olá,</p>
                <p>Recebemos um pedido de recuperação de senha para a sua conta. Para redefinir sua senha, clique no botão abaixo:</p>
                <a href="{{.Link}}" class="btn">Redefinir Senha</a>
                <p>Ou copie e cole este link no seu navegador:</p>
                <p style="word-break: break-all; color: #007bff;">{{.Link}}</p>
                <p>Se você não solicitou a recuperação de senha, ignore este e-mail.</p>
                <div class="footer">
                    <p>Atenciosamente,<br>Equipe de Suporte</p>
                </div>
            </div>
        </body>
        </html>`

	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		fmt.Println("Erro ao criar template:", err.Error())
		return
	}

	data := EmailData{
		Link: link,
	}

	var htmlBuffer bytes.Buffer
	err = tmpl.Execute(&htmlBuffer, data)
	if err != nil {
		fmt.Println("Erro ao processar template:", err.Error())
		return
	}

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{"delivered@resend.dev"},
		Html:    htmlBuffer.String(),
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
