package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"notification-payment-service/config"
)

type SMSService struct {
	cfg *config.Config
}

func NewSMSService(cfg *config.Config) *SMSService {
	return &SMSService{cfg: cfg}
}

type SendSMSDTO struct {
	Phone   []string `json:"phone"`
	Message string   `json:"message"`
}

type netgsmRequestPayload struct {
	Usercode string   `json:"usercode"`
	Password string   `json:"password"`
	Header   string   `json:"msgheader"`
	Message  string   `json:"message"`
	Gsm      []string `json:"gsm"`
}

func (s *SMSService) SendSMS(dto SendSMSDTO) error {
	payload := netgsmRequestPayload{
		Usercode: s.cfg.NetgsmUsercode,
		Password: s.cfg.NetgsmPassword,
		Header:   s.cfg.NetgsmHeader,
		Message:  dto.Message,
		Gsm:      dto.Phone,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://api.netgsm.com.tr/sms/send/post", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("netgsm istek hatasi: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sms gönderilemedi, HTTP Status: %d", resp.StatusCode)
	}

	return nil
}
