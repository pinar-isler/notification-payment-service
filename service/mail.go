package service

import (
	"fmt"
	"notification-payment-service/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService struct {
	cfg *config.Config
}

func NewMailService(cfg *config.Config) *MailService {
	return &MailService{cfg: cfg}
}

type SendEmailDTO struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func (s *MailService) SendEmail(dto SendEmailDTO) error {
	from := mail.NewEmail("Sistem Bildirimi", s.cfg.SendGridFromEmail)
	// Çoklu alıcı için V3 Mail nesnesi oluşturuyoruz
	m := mail.NewV3Mail()
	m.SetFrom(from)
	m.Subject = dto.Subject

	// Dizi içindeki tüm alıcıları (to) e-posta nesnesine ekliyoruz
	p := mail.NewPersonalization()
	for _, emailAddr := range dto.To {
		p.AddTos(mail.NewEmail("Müşteri", emailAddr))
	}
	m.AddPersonalizations(p)

	// Mail içeriği (hem düz metin hem HTML olarak)
	m.AddContent(mail.NewContent("text/plain", dto.Body))

	client := sendgrid.NewSendClient(s.cfg.SendGridAPIKey)
	response, err := client.Send(m)
	if err != nil {
		return fmt.Errorf("sendgrid baglanti hatasi: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("mail gonderilemedi, status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
