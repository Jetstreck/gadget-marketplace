package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"gadget-marketplace/config"
)

type EmailService interface {
	SendWelcomeEmail(toEmail, username string) error
	SendOrderInvoiceEmail(toEmail, username, productName string, quantity int, totalPrice float64, newDeposit float64) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) sendMail(toEmail, subject, bodyHTML string) error {
	if s.cfg.SMTPUser == "" || s.cfg.SMTPPassword == "" {
		log.Printf("[EMAIL NOTIFICATION SIMULATED] To: %s | Subject: %s | Message: Email notification triggered successfully! (Set SMTP_USER and SMTP_PASSWORD in .env for Mailjet/Mailtrap/Inboxes delivery)", toEmail, subject)
		return nil
	}

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	headers := make(map[string]string)
	headers["From"] = s.cfg.SenderEmail
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + bodyHTML

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	err := smtp.SendMail(addr, auth, s.cfg.SenderEmail, []string{toEmail}, []byte(message))
	if err != nil {
		log.Printf("Warning: Failed to send email via SMTP (%v)", err)
		return err
	}

	log.Printf("Email successfully sent to %s via Mailjet/Mailtrap/Inboxes SMTP", toEmail)
	return nil
}

func (s *emailService) SendWelcomeEmail(toEmail, username string) error {
	subject := "Welcome to Gadget Marketplace!"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome to Gadget Marketplace, %s! 🎉</h2>
			<p>Thank you for registering your account with email: <strong>%s</strong>.</p>
			<p>Explore our latest gadget catalog and enjoy seamless shopping!</p>
			<hr>
			<p><small>Gadget Marketplace Team</small></p>
		</body>
		</html>
	`, username, toEmail)

	return s.sendMail(toEmail, subject, strings.TrimSpace(body))
}

func (s *emailService) SendOrderInvoiceEmail(toEmail, username, productName string, quantity int, totalPrice float64, newDeposit float64) error {
	subject := "Order Invoice - Gadget Marketplace Purchase"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Purchase Invoice 🛒</h2>
			<p>Dear <strong>%s</strong>,</p>
			<p>Your purchase checkout was successful!</p>
			<table border="1" cellpadding="8" cellspacing="0">
				<tr><th>Item</th><td>%s</td></tr>
				<tr><th>Quantity</th><td>%d</td></tr>
				<tr><th>Total Price</th><td>Rp %.2f</td></tr>
				<tr><th>Remaining Deposit</th><td>Rp %.2f</td></tr>
			</table>
			<br>
			<p>Thank you for shopping with Gadget Marketplace!</p>
		</body>
		</html>
	`, username, productName, quantity, totalPrice, newDeposit)

	return s.sendMail(toEmail, subject, strings.TrimSpace(body))
}
