package smtp

import (
	"context"
	"fmt"

	"github.com/wneessen/go-mail"
)

type Client struct {
	MailClient *mail.Client
	FromEmail  string
	FromName   string
}

// EmailMessage represents an email to be sent
type EmailMessage struct {
	To       string
	Subject  string
	Body     string
	HTMLBody string
}

// SendEmail sends a plain text email
func (c *Client) SendEmail(ctx context.Context, msg *EmailMessage) error {
	if c.MailClient == nil {
		return fmt.Errorf("SMTP client not initialized")
	}

	message := mail.NewMsg()
	if err := message.FromFormat(c.FromName, c.FromEmail); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}
	if err := message.To(msg.To); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	message.Subject(msg.Subject)

	if msg.HTMLBody != "" {
		message.SetBodyString(mail.TypeTextHTML, msg.HTMLBody)
		if msg.Body != "" {
			message.AddAlternativeString(mail.TypeTextPlain, msg.Body)
		}
	} else {
		message.SetBodyString(mail.TypeTextPlain, msg.Body)
	}

	if err := c.MailClient.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email with a verification code
func (c *Client) SendPasswordResetEmail(ctx context.Context, to, code, userName string) error {
	subject := "Password Reset Code - Inventory API"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .code { font-size: 32px; font-weight: bold; color: #4CAF50; letter-spacing: 5px; text-align: center; padding: 20px; background: #f5f5f5; border-radius: 8px; margin: 20px 0; }
        .warning { color: #666; font-size: 14px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Password Reset Request</h2>
        <p>Hi %s,</p>
        <p>We received a request to reset your password. Use the following code to reset your password:</p>
        <div class="code">%s</div>
        <p class="warning">This code will expire in 15 minutes. If you didn't request this, please ignore this email.</p>
        <p>Best regards,<br>Inventory API Team</p>
    </div>
</body>
</html>
`, userName, code)

	plainBody := fmt.Sprintf(`Password Reset Request

Hi %s,

We received a request to reset your password. Use the following code to reset your password:

%s

This code will expire in 15 minutes. If you didn't request this, please ignore this email.

Best regards,
Inventory API Team
`, userName, code)

	return c.SendEmail(ctx, &EmailMessage{
		To:       to,
		Subject:  subject,
		Body:     plainBody,
		HTMLBody: htmlBody,
	})
}

// IsEnabled checks if SMTP client is available
func (c *Client) IsEnabled() bool {
	return c != nil && c.MailClient != nil
}
