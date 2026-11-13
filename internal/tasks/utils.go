package tasks

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	// hash := sha256.Sum256([]byte(token))
	// tokenHash := base64.URLEncoding.EncodeToString(hash[:])

	return token, nil
}
