package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-gomail/gomail"
)

// GenerateHash creates a bcrypt hash from the provided password
// Returns the hashed password as a string or an error if hashing fails
func GenerateHash(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CompareHash compares a password with its hash
// Returns true if the password matches the hash, false otherwise
func CompareHash(pwd, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err != nil {
		return false
	}
	return true
}

func GenerateJWT(secret string, userID uint, expirationHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expirationHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func SendEmail(to, subject, body, smtpUser, smtpPass, smtpHost string, smtpPort int) error {
    m := gomail.NewMessage()
    m.SetHeader("From", smtpUser)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
    return d.DialAndSend(m)
}
