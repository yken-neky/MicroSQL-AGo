package services

// EncryptionService provides simple encrypt/decrypt methods used by usecases
type EncryptionService interface {
	Encrypt(plain string) (string, error)
	Decrypt(encrypted string) (string, error)
}
