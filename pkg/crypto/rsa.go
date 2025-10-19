package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	rsaKeyBits = 4096
)

// RSAKeyPair represents an RSA public/private key pair
type RSAKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// GenerateRSAKeyPair generates a new RSA key pair
func GenerateRSAKeyPair() (*RSAKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	return &RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// PrivateKeyToPEM converts a private key to PEM format
func PrivateKeyToPEM(key *rsa.PrivateKey) (string, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	if privateKeyPEM == nil {
		return "", fmt.Errorf("failed to encode private key to PEM")
	}

	return string(privateKeyPEM), nil
}

// PublicKeyToPEM converts a public key to PEM format
func PublicKeyToPEM(key *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	if publicKeyPEM == nil {
		return "", fmt.Errorf("failed to encode public key to PEM")
	}

	return string(publicKeyPEM), nil
}

// ParsePrivateKeyFromPEM parses a private key from PEM format
func ParsePrivateKeyFromPEM(pemString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// ParsePublicKeyFromPEM parses a public key from PEM format
func ParsePublicKeyFromPEM(pemString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return publicKey, nil
}

// SavePrivateKeyToFile saves a private key to a file
func SavePrivateKeyToFile(key *rsa.PrivateKey, filepath string) error {
	pemString, err := PrivateKeyToPEM(key)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, []byte(pemString), 0600); err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}

	return nil
}

// SavePublicKeyToFile saves a public key to a file
func SavePublicKeyToFile(key *rsa.PublicKey, filepath string) error {
	pemString, err := PublicKeyToPEM(key)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, []byte(pemString), 0644); err != nil {
		return fmt.Errorf("failed to write public key to file: %w", err)
	}

	return nil
}

// LoadPrivateKeyFromFile loads a private key from a file
func LoadPrivateKeyFromFile(filepath string) (*rsa.PrivateKey, error) {
	pemBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	return ParsePrivateKeyFromPEM(string(pemBytes))
}

// LoadPublicKeyFromFile loads a public key from a file
func LoadPublicKeyFromFile(filepath string) (*rsa.PublicKey, error) {
	pemBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	return ParsePublicKeyFromPEM(string(pemBytes))
}

// SignData signs data with a private key using RSA-PSS
func SignData(data []byte, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256(data)

	signature, err := rsa.SignPSS(
		rand.Reader,
		privateKey,
		crypto.SHA256,
		hash[:],
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies a signature with a public key using RSA-PSS
func VerifySignature(data []byte, signature string, publicKey *rsa.PublicKey) error {
	hash := sha256.Sum256(data)

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	err = rsa.VerifyPSS(
		publicKey,
		crypto.SHA256,
		hash[:],
		signatureBytes,
		nil,
	)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}
