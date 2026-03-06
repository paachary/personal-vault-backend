package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"go-mongo-project/config"
	"io"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	cachedKey    []byte
	cachedKeyErr error
	keyOnce      sync.Once
)

func getEncryptionKey() ([]byte, error) {
	keyOnce.Do(func() {
		keyHex := os.Getenv(config.CREDENTIAL_ENCRYPTION_KEY)

		key, err := hex.DecodeString(keyHex)
		if err != nil {
			cachedKeyErr = errors.New("CREDENTIAL_ENCRYPTION_KEY must be a valid hex string")
			return
		}

		if len(key) != 32 {
			cachedKeyErr = errors.New("CREDENTIAL_ENCRYPTION_KEY must be exactly 32 bytes for AES-256")
			return
		}
		cachedKey = []byte(key)
	})
	return cachedKey, cachedKeyErr
}

// EncryptCredential encrypts a plaintext credential using AES-256-GCM and returns
// a base64-encoded ciphertext. Returns an empty string unchanged.
func EncryptCredential(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptCredential decrypts a base64-encoded AES-256-GCM ciphertext and returns
// the original plaintext. Returns an empty string unchanged.
func DecryptCredential(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errors.New("invalid credential format")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid credential format")
	}
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", errors.New("failed to decrypt credential")
	}
	return string(plaintext), nil
}

// EncryptStructFields encrypts all string fields in a bson.M map
func EncryptStructFields(data bson.M, excludeFields ...string) (bson.M, error) {
	encrypted := make(bson.M)
	excludeMap := make(map[string]bool)

	for _, field := range excludeFields {
		excludeMap[field] = true
	}

	for key, value := range data {
		if excludeMap[key] {
			encrypted[key] = value
			continue
		}

		if strValue, ok := value.(string); ok && strValue != "" {
			encValue, err := EncryptCredential(strValue)
			if err != nil {
				return nil, err
			}
			encrypted[key] = encValue
		} else {
			encrypted[key] = value
		}
	}

	return encrypted, nil
}

// GenerateDeterministicHash creates a consistent hash from multiple input strings
// This is useful for uniqueness checks with encrypted data
func GenerateDeterministicHash(inputs ...string) string {
	// Combine all inputs into a single string
	combined := ""
	for _, input := range inputs {
		combined += input + "|" // Use delimiter to prevent collision
	}

	// Create SHA-256 hash
	hash := sha256.Sum256([]byte(combined))

	// Return as hex string
	return hex.EncodeToString(hash[:])
}

// GenerateUniqueKey creates a shorter unique key (first 16 chars of hash)
func GenerateUniqueKey(inputs ...string) string {
	fullHash := GenerateDeterministicHash(inputs...)
	// Return first 16 characters for shorter keys
	if len(fullHash) > 16 {
		return fullHash[:16]
	}
	return fullHash
}
