package main

import (
	"fmt"
	"log"
	"net/http"

	"notification-payment-service/config"
	"notification-payment-service/handler"
	"notification-payment-service/service"
)

func main() {
	cfg := config.LoadConfig()

	mailService := service.NewMailService(cfg)
	smsService := service.NewSMSService(cfg)
	paymentService := service.NewPaymentService(cfg)

	h := handler.NewHandler(mailService, smsService, paymentService)

	// Ana dizine (http://localhost:8080) girildiğinde index.html'i sunar
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "index.html")
	})

	// API Endpoint'leri
	http.HandleFunc("/api/v1/mail/send", h.SendEmailHandler)
	http.HandleFunc("/api/v1/sms/send", h.SendSMSHandler)
	http.HandleFunc("/api/v1/payment/process", h.ProcessPaymentHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Uygulama sunucusu http://localhost%s üzerinde çalışıyor...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
