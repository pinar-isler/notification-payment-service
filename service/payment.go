package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"notification-payment-service/config"
)

type PaymentService struct {
	cfg *config.Config
}

func NewPaymentService(cfg *config.Config) *PaymentService {
	return &PaymentService{cfg: cfg}
}

type ProcessPaymentDTO struct {
	CardHolderName string `json:"card_holder_name"`
	CardNumber     string `json:"card_number"`
	ExpireMonth    string `json:"expire_month"`
	ExpireYear     string `json:"expire_year"`
	Cvc            string `json:"cvc"`
	Price          string `json:"price"`
}

type iyzicoPaymentResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	PaymentID    string `json:"paymentId"`
}

func (s *PaymentService) ProcessPayment(dto ProcessPaymentDTO) (string, error) {
	payload := map[string]interface{}{
		"locale":         "tr",
		"conversationId": "123456789",
		"price":          dto.Price,
		"paidPrice":      dto.Price,
		"currency":       "TRY",
		"installment":    "1",
		"basketId":       "B67832",
		"paymentChannel": "WEB",
		"paymentGroup":   "PRODUCT",
		"paymentCard": map[string]string{
			"cardHolderName": dto.CardHolderName,
			"cardNumber":     dto.CardNumber,
			"expireMonth":    dto.ExpireMonth,
			"expireYear":     dto.ExpireYear,
			"cvc":            dto.Cvc,
			"registerCard":   "0",
		},
		"buyer": map[string]string{
			"id":                  "BY789",
			"name":                "Pinar",
			"surname":             "Isler",
			"gsmNumber":           "+905350000000",
			"email":               "email@example.com",
			"identityNumber":      "74300864791",
			"registrationAddress": "Nisantasi",
			"ip":                  "85.95.255.255",
			"city":                "Istanbul",
			"country":             "Turkey",
		},
		"shippingAddress": map[string]string{
			"contactName": "Pinar Isler",
			"city":        "Istanbul",
			"country":     "Turkey",
			"address":     "Nisantasi",
		},
		"billingAddress": map[string]string{
			"contactName": "Pinar Isler",
			"city":        "Istanbul",
			"country":     "Turkey",
			"address":     "Nisantasi",
		},
		"basketItems": []map[string]string{
			{
				"id":        "BI101",
				"name":      "Test Urun",
				"category1": "Yazilim",
				"itemType":  "VIRTUAL",
				"price":     dto.Price,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/payment/auth", s.cfg.IyzicoBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-iyzi-rnd", "234567")
	req.Header.Set("Authorization", fmt.Sprintf("IYZWS %s:%s", s.cfg.IyzicoAPIKey, s.cfg.IyzicoSecretKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("iyzico istek hatasi: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var iyzicoResp iyzicoPaymentResponse
	if err := json.Unmarshal(bodyBytes, &iyzicoResp); err != nil {
		return "", fmt.Errorf("yanit ayristirilamadi: %v", err)
	}

	if iyzicoResp.Status != "success" {
		return "", fmt.Errorf("odeme basarisiz: %s (Hata Kodu: %s)", iyzicoResp.ErrorMessage, iyzicoResp.ErrorCode)
	}

	return iyzicoResp.PaymentID, nil
}
