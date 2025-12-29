package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		3,       // time
		64*1024, // memory
		4,       // threads
		32,      // key length (AES-256)
	)
}

func GenerateDEK() ([]byte, error) {
	dek := make([]byte, 32) // AES-256
	_, err := rand.Read(dek)
	return dek, err
}

func EncryptDEKWithKEK(masterPassword string, dek []byte) (string, error) {
	// per-user KEK salt
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// derive KEK from master password
	kek := DeriveKey(masterPassword, salt)

	block, err := aes.NewCipher(kek)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nil, nonce, dek, nil)

	// store salt + nonce + encrypted DEK
	out := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return base64.RawStdEncoding.EncodeToString(out), nil
}

func DecryptDEKWithKEK(masterPassword, encrypted_DEK string) ([]byte, error) {
	raw, err := base64.RawStdEncoding.DecodeString(encrypted_DEK)
	if err != nil {
		return nil, err
	}

	salt := raw[:saltLen]
	nonce := raw[saltLen : saltLen+12]
	ciphertext := raw[saltLen+12:]

	kek := DeriveKey(masterPassword, salt)

	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, nonce, ciphertext, nil)
}

func EncryptWithDEK(dek []byte, plaintext string) (string, error) {
	if len(dek) != 32 {
		return "", errors.New("invalid DEK length")
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)

	out := append(nonce, ciphertext...)
	return base64.RawStdEncoding.EncodeToString(out), nil
}

func DecryptWithDEK(dek []byte, encoded string) (string, error) {
	raw, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(raw) < aead.NonceSize() {
		return "", errors.New("invalid ciphertext")
	}

	nonce := raw[:aead.NonceSize()]
	ciphertext := raw[aead.NonceSize():]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
