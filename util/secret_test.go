package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	password := "hashed_password_mock"
	secret := "my_secret_api_key"

	// 1. Round trip (V2)
	encrypted, err := EncryptSecretWithPassword(password, secret)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	if !strings.HasPrefix(encrypted, "ENCv2:") {
		t.Errorf("Encrypted string should have prefix 'ENCv2:', got %s", encrypted)
	}

	decrypted, err := DecryptSecretWithPassword(password, encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	if decrypted != secret {
		t.Errorf("Decrypted secret mismatch. Want %s, got %s", secret, decrypted)
	}

	// 2. Already encrypted
	enc2, err := EncryptSecretWithPassword(password, encrypted)
	if err != nil {
		t.Fatalf("Re-encryption failed: %v", err)
	}
	if enc2 != encrypted {
		t.Errorf("Should not re-encrypt already encrypted string")
	}

	// 3. Wrong password
	wrongPass := "wrong_password"
	_, err = DecryptSecretWithPassword(wrongPass, encrypted)
	if err == nil {
		t.Error("Decryption should fail with wrong password")
	}

	// 4. Corrupted ciphertext
	corrupted := encrypted + "junk"
	_, err = DecryptSecretWithPassword(password, corrupted)
	if err == nil {
		t.Error("Decryption should fail with corrupted ciphertext")
	}

	// 5. Empty secret
	emptyEnc, err := EncryptSecretWithPassword(password, "")
	if err != nil {
		t.Fatalf("Empty secret encryption failed: %v", err)
	}
	if emptyEnc != "" {
		t.Errorf("Empty secret should result in empty string, got %s", emptyEnc)
	}

	// 6. Empty password
	noPassEnc, err := EncryptSecretWithPassword("", secret)
	if err != nil {
		t.Fatalf("No password encryption failed: %v", err)
	}
	if noPassEnc != secret {
		t.Errorf("No password should return original secret")
	}

	// 7. Legacy V1 Decryption
	v1Key := deriveSecretKey(password)
	block, _ := aes.NewCipher(v1Key)
	aead, _ := cipher.NewGCM(block)
	nonce := make([]byte, aead.NonceSize())
	rand.Read(nonce)
	ciphertext := aead.Seal(nil, nonce, []byte(secret), nil)
	v1Enc := "ENC:" + base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))

	decV1, err := DecryptSecretWithPassword(password, v1Enc)
	if err != nil {
		t.Fatalf("V1 Decryption failed: %v", err)
	}
	if decV1 != secret {
		t.Errorf("V1 Decrypted secret mismatch. Want %s, got %s", secret, decV1)
	}
}
