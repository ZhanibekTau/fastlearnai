package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
)

const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// EmailService sends transactional emails over Gmail SMTP using an app password.
type EmailService struct {
	login    string
	password string
}

func NewEmailService(login, password string) *EmailService {
	return &EmailService{login: login, password: password}
}

func (s *EmailService) GeneratePassword(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate password: %w", err)
		}
		result[i] = passwordChars[n.Int64()]
	}
	return string(result), nil
}

func (s *EmailService) SendNewPassword(toEmail, name, password string) error {
	subject := "Восстановление пароля — FastLearn"
	body := fmt.Sprintf(
		"Здравствуйте, %s!\n\nВаш новый пароль: %s\n\nПожалуйста, войдите с этим паролем и при желании смените его в настройках профиля.",
		name, password,
	)
	return s.send(toEmail, subject, body)
}

func (s *EmailService) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.login, s.password, "smtp.gmail.com")
	msg := "From: " + s.login + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body
	return smtp.SendMail("smtp.gmail.com:587", auth, s.login, []string{to}, []byte(msg))
}
