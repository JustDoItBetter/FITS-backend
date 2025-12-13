package crypto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRSAKeyPair(t *testing.T) {
	t.Run("generates key pair successfully", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()

		require.NoError(t, err)
		assert.NotNil(t, keyPair)
		assert.NotNil(t, keyPair.PrivateKey)
		assert.NotNil(t, keyPair.PublicKey)
	})

	t.Run("generates 4096-bit keys", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		assert.Equal(t, 4096, keyPair.PrivateKey.N.BitLen())
	})

	t.Run("public key matches private key", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		assert.Equal(t, keyPair.PublicKey.N, keyPair.PrivateKey.N)
		assert.Equal(t, keyPair.PublicKey.E, keyPair.PrivateKey.E)
	})

	t.Run("generates different keys each time", func(t *testing.T) {
		keyPair1, err1 := GenerateRSAKeyPair()
		keyPair2, err2 := GenerateRSAKeyPair()

		require.NoError(t, err1)
		require.NoError(t, err2)

		assert.NotEqual(t, keyPair1.PrivateKey.D, keyPair2.PrivateKey.D)
		assert.NotEqual(t, keyPair1.PublicKey.N, keyPair2.PublicKey.N)
	})
}

func TestPEMConversion(t *testing.T) {
	t.Run("converts private key to PEM and back", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		pemString, err := PrivateKeyToPEM(keyPair.PrivateKey)
		require.NoError(t, err)
		assert.NotEmpty(t, pemString)

		parsedKey, err := ParsePrivateKeyFromPEM(pemString)
		require.NoError(t, err)
		assert.Equal(t, keyPair.PrivateKey.D, parsedKey.D)
	})

	t.Run("converts public key to PEM and back", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		pemString, err := PublicKeyToPEM(keyPair.PublicKey)
		require.NoError(t, err)
		assert.NotEmpty(t, pemString)

		parsedKey, err := ParsePublicKeyFromPEM(pemString)
		require.NoError(t, err)
		assert.Equal(t, keyPair.PublicKey.N, parsedKey.N)
		assert.Equal(t, keyPair.PublicKey.E, parsedKey.E)
	})
}

func TestSaveAndLoadKeys(t *testing.T) {
	t.Run("saves and loads private key", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		tempDir := t.TempDir()
		privPath := filepath.Join(tempDir, "test.key")

		err = SavePrivateKeyToFile(keyPair.PrivateKey, privPath)
		require.NoError(t, err)

		loadedKey, err := LoadPrivateKeyFromFile(privPath)
		require.NoError(t, err)
		assert.Equal(t, keyPair.PrivateKey.D, loadedKey.D)
	})

	t.Run("saves and loads public key", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		tempDir := t.TempDir()
		pubPath := filepath.Join(tempDir, "test.pub")

		err = SavePublicKeyToFile(keyPair.PublicKey, pubPath)
		require.NoError(t, err)

		loadedKey, err := LoadPublicKeyFromFile(pubPath)
		require.NoError(t, err)
		assert.Equal(t, keyPair.PublicKey.N, loadedKey.N)
	})

	t.Run("private key file has correct permissions", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		tempDir := t.TempDir()
		privPath := filepath.Join(tempDir, "secure.key")

		err = SavePrivateKeyToFile(keyPair.PrivateKey, privPath)
		require.NoError(t, err)

		info, err := os.Stat(privPath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
	})

	t.Run("public key file has correct permissions", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		tempDir := t.TempDir()
		pubPath := filepath.Join(tempDir, "test.pub")

		err = SavePublicKeyToFile(keyPair.PublicKey, pubPath)
		require.NoError(t, err)

		info, err := os.Stat(pubPath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
	})
}

func TestSignData(t *testing.T) {
	t.Run("signs data successfully", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		data := []byte("test data to sign")
		signature, err := SignData(data, keyPair.PrivateKey)

		require.NoError(t, err)
		assert.NotEmpty(t, signature)
	})

	t.Run("generates different signatures for different data", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		sig1, err1 := SignData([]byte("data1"), keyPair.PrivateKey)
		sig2, err2 := SignData([]byte("data2"), keyPair.PrivateKey)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, sig1, sig2)
	})

	t.Run("handles empty data", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		signature, err := SignData([]byte(""), keyPair.PrivateKey)

		require.NoError(t, err)
		assert.NotEmpty(t, signature)
	})
}

func TestVerifySignature(t *testing.T) {
	t.Run("verifies valid signature", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		data := []byte("verify this data")
		signature, err := SignData(data, keyPair.PrivateKey)
		require.NoError(t, err)

		err = VerifySignature(data, signature, keyPair.PublicKey)
		assert.NoError(t, err)
	})

	t.Run("rejects signature with wrong data", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		originalData := []byte("original data")
		signature, err := SignData(originalData, keyPair.PrivateKey)
		require.NoError(t, err)

		differentData := []byte("different data")
		err = VerifySignature(differentData, signature, keyPair.PublicKey)
		assert.Error(t, err)
	})

	t.Run("rejects signature with wrong key", func(t *testing.T) {
		keyPair1, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		keyPair2, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		data := []byte("signed with key 1")
		signature, err := SignData(data, keyPair1.PrivateKey)
		require.NoError(t, err)

		err = VerifySignature(data, signature, keyPair2.PublicKey)
		assert.Error(t, err)
	})

	t.Run("rejects invalid signature", func(t *testing.T) {
		keyPair, err := GenerateRSAKeyPair()
		require.NoError(t, err)

		data := []byte("some data")
		invalidSignature := "invalid-base64-signature"

		err = VerifySignature(data, invalidSignature, keyPair.PublicKey)
		assert.Error(t, err)
	})
}

// Benchmark tests
func BenchmarkGenerateRSAKeyPair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateRSAKeyPair()
	}
}

func BenchmarkSignData(b *testing.B) {
	keyPair, _ := GenerateRSAKeyPair()
	data := []byte("benchmark data to sign")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SignData(data, keyPair.PrivateKey)
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	keyPair, _ := GenerateRSAKeyPair()
	data := []byte("benchmark data")
	signature, _ := SignData(data, keyPair.PrivateKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifySignature(data, signature, keyPair.PublicKey)
	}
}
