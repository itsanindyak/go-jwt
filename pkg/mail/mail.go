package mail

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/itsanindyak/go-jwt/config"
	"github.com/itsanindyak/go-jwt/pkg/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailData struct {
	Name string
	OTP  string
}

func Send(m MailData, toMail string) error {
	subject := "Authentication otp"
	tmpl, err := template.ParseFiles("pkg/mail/template/otp.html")
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var body bytes.Buffer

	if err = tmpl.Execute(&body, m); err != nil {
		logger.Error(err.Error())
		return err
	}

	from := mail.NewEmail(config.SMTP_NAME, config.SMTP_SENDER)
	to := mail.NewEmail(m.Name, toMail)

	message := mail.NewSingleEmail(from, subject, to, "", body.String())

	client := sendgrid.NewSendClient(config.SENDGRID_API_KEY)

	response, err := client.Send(message)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if response.StatusCode >= 400 {
		logger.Error(fmt.Sprintf("SendGrid error: status %d, body: %s", response.StatusCode, response.Body))
		return fmt.Errorf("SendGrid error: %s", response.Body)
	}

	return nil
}
