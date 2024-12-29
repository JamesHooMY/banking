package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomAPIKey() string {
	// Generate 32 random bytes (256 bits) for the API key
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("Error generating API key:", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func GenerateRandomSecretKey() string {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("Error generating secret key:", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func GenerateHashedSecretKey(secret string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifySecretKey(storedHash string, providedSecret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedSecret))
	return err == nil
}
