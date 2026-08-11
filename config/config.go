package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	SendGridAPIKey    string
	SendGridFromEmail string
	NetgsmUsercode    string
	NetgsmPassword    string
	NetgsmHeader      string
	IyzicoAPIKey      string
	IyzicoSecretKey   string
	IyzicoBaseURL     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Uyarı: .env dosyası bulunamadı, sistem çevre değişkenleri kullanılacak.")
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		SendGridAPIKey:    os.Getenv("SENDGRID_API_KEY"),
		SendGridFromEmail: os.Getenv("SENDGRID_FROM_EMAIL"),
		NetgsmUsercode:    os.Getenv("NETGSM_USERCODE"),
		NetgsmPassword:    os.Getenv("NETGSM_PASSWORD"),
		NetgsmHeader:      os.Getenv("NETGSM_HEADER"),
		IyzicoAPIKey:      os.Getenv("IYZICO_API_KEY"),
		IyzicoSecretKey:   os.Getenv("IYZICO_SECRET_KEY"),
		IyzicoBaseURL:     getEnv("IYZICO_BASE_URL", "https://sandbox-api.iyzipay.com"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
