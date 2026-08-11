package handler

import (
	"encoding/json"
	"net/http"
	"notification-payment-service/service"
)

type Handler struct {
	mailService    *service.MailService
	smsService     *service.SMSService
	paymentService *service.PaymentService
}

func NewHandler(mailService *service.MailService, smsService *service.SMSService, paymentService *service.PaymentService) *Handler {
	return &Handler{
		mailService:    mailService,
		smsService:     smsService,
		paymentService: paymentService,
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func (h *Handler) SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Sadece POST isteği atılabilir."})
		return
	}

	var dto service.SendEmailDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz JSON formatı"})
		return
	}

	if err := h.mailService.SendEmail(dto); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "E-posta başarıyla gönderildi."})
}

func (h *Handler) SendSMSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Sadece POST isteği atılabilir."})
		return
	}

	var dto service.SendSMSDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz JSON formatı"})
		return
	}

	if err := h.smsService.SendSMS(dto); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "SMS başarıyla gönderildi."})
}

func (h *Handler) ProcessPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Sadece POST isteği atılabilir."})
		return
	}

	var dto service.ProcessPaymentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz JSON formatı"})
		return
	}

	paymentID, err := h.paymentService.ProcessPayment(dto)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message":    "Ödeme başarıyla gerçekleşti.",
		"payment_id": paymentID,
	})
}
