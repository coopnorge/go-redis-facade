package database

import (
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
)

// Encryption defines interface
type Encryption interface {
	Encrypt(value []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// EncryptionClient handles encryption and decryption tasks
type EncryptionClient struct {
	a   tink.AEAD
	aad []byte
}

// EncryptionConfig is configuration needed to set up the encryption client
type EncryptionConfig struct {
	CredentialPath string
	RedisKeyURI    string
	Aad            []byte
}

// NewEncryptionClient creates a new client for encrypting and decrypting data
func NewEncryptionClient(c EncryptionConfig) (*EncryptionClient, error) {
	kmsClient, err := gcpkms.NewClientWithCredentials(c.RedisKeyURI, c.CredentialPath)
	if err != nil {
		return nil, err
	}
	registry.RegisterKMSClient(kmsClient)

	dek := aead.AES128CTRHMACSHA256KeyTemplate()
	kh, err := keyset.NewHandle(aead.KMSEnvelopeAEADKeyTemplate(c.RedisKeyURI, dek))
	if err != nil {
		return nil, err
	}

	a, err := aead.New(kh)
	if err != nil {
		return nil, err
	}

	return &EncryptionClient{a: a, aad: c.Aad}, nil
}

// Encrypt encrypts data
func (e *EncryptionClient) Encrypt(value []byte) ([]byte, error) {
	return e.a.Encrypt(value, e.aad)
}

// Decrypt decrypts data
func (e *EncryptionClient) Decrypt(ciphertext []byte) ([]byte, error) {
	return e.a.Decrypt(ciphertext, e.aad)
}
