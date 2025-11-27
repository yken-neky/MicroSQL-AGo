package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// AESGCMService implements a simple AES-GCM based encryption service
// Key must be 16, 24 or 32 bytes long (AES-128/192/256)
type AESGCMService struct {
	key []byte
}

func NewAESGCMService(key string) *AESGCMService {
	return &AESGCMService{key: []byte(key)}
}

// Encrypt encrypts plain text and returns base64 encoded ciphertext (nonce + ciphertext)
func (s *AESGCMService) Encrypt(plain string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce gen: %w", err)
	}

	ct := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

// Decrypt accepts base64 encoded ciphertext (nonce + ct)
func (s *AESGCMService) Decrypt(encrypted string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}

	if len(raw) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := raw[:gcm.NonceSize()]
	ct := raw[gcm.NonceSize():]

	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("open: %w", err)
	}
	return string(pt), nil
}
