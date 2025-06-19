package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// EncryptionService handles encryption and decryption of sensitive data
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(password string) *EncryptionService {
	// Generate key from password using SHA256
	hash := sha256.Sum256([]byte(password))
	return &EncryptionService{
		key: hash[:],
	}
}

// Encrypt encrypts plaintext using AES-GCM
func (e *EncryptionService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func (e *EncryptionService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// APIKeyManager manages encrypted API keys
type APIKeyManager struct {
	encService *EncryptionService
	app        *Application
}

// NewAPIKeyManager creates a new API key manager
func NewAPIKeyManager(app *Application, masterPassword string) *APIKeyManager {
	return &APIKeyManager{
		encService: NewEncryptionService(masterPassword),
		app:        app,
	}
}

// StoreAPIKey stores an encrypted API key
func (m *APIKeyManager) StoreAPIKey(service, apiKey string) error {
	encrypted, err := m.encService.Encrypt(apiKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt API key: %w", err)
	}

	query := `
		INSERT OR REPLACE INTO api_keys (service_name, encrypted_key, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	timestamp := getCurrentTimestamp()
	_, err = m.app.db.Exec(query, service, encrypted, timestamp, timestamp)
	if err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	m.app.logger.Infof("Stored API key for service: %s", service)
	return nil
}

// GetAPIKey retrieves and decrypts an API key
func (m *APIKeyManager) GetAPIKey(service string) (string, error) {
	query := `
		SELECT encrypted_key FROM api_keys 
		WHERE service_name = ?
	`
	var encryptedKey string
	err := m.app.db.QueryRow(query, service).Scan(&encryptedKey)
	if err != nil {
		return "", fmt.Errorf("API key not found for service %s: %w", service, err)
	}

	apiKey, err := m.encService.Decrypt(encryptedKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt API key for service %s: %w", service, err)
	}

	// Update last_used timestamp
	updateQuery := `UPDATE api_keys SET last_used = ? WHERE service_name = ?`
	m.app.db.Exec(updateQuery, getCurrentTimestamp(), service)

	return apiKey, nil
}

// DeleteAPIKey removes an API key
func (m *APIKeyManager) DeleteAPIKey(service string) error {
	query := `DELETE FROM api_keys WHERE service_name = ?`
	_, err := m.app.db.Exec(query, service)
	if err != nil {
		return fmt.Errorf("failed to delete API key for service %s: %w", service, err)
	}

	m.app.logger.Infof("Deleted API key for service: %s", service)
	return nil
}

// ListAPIKeys returns a list of services with stored API keys
func (m *APIKeyManager) ListAPIKeys() ([]APIKeyInfo, error) {
	query := `
		SELECT service_name, created_at, updated_at, last_used 
		FROM api_keys 
		ORDER BY service_name
	`
	rows, err := m.app.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer rows.Close()

	var keys []APIKeyInfo
	for rows.Next() {
		var key APIKeyInfo
		var lastUsed *int64
		err := rows.Scan(&key.Service, &key.CreatedAt, &key.UpdatedAt, &lastUsed)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key row: %w", err)
		}
		if lastUsed != nil {
			key.LastUsed = *lastUsed
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// APIKeyInfo represents information about an API key
type APIKeyInfo struct {
	Service   string `json:"service"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	LastUsed  int64  `json:"last_used,omitempty"`
}
