package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters for password hashing
const (
	hashSaltLen = 16
	hashKeyLen  = 32
	hashTime    = 3
	hashMemory  = 64 * 1024
	hashThreads = 4
)

// HashPassword hashes a password using Argon2id and returns a string containing all parameters, salt, and hash.
func HashPassword(password string) (string, error) {
	salt := make([]byte, hashSaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, hashTime, hashMemory, uint8(hashThreads), hashKeyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := "$argon2id$v=19"
	encoded += fmt.Sprintf("$m=%d,t=%d,p=%d", hashMemory, hashTime, hashThreads)
	encoded += fmt.Sprintf("$%s$%s", b64Salt, b64Hash)
	return encoded, nil
}

// VerifyPassword checks if the provided password matches the encoded Argon2id hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}
	var memory, time, threads uint32
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	calculated := argon2.IDKey([]byte(password), salt, time, memory, uint8(threads), uint32(len(hash)))
	return subtleCompare(hash, calculated), nil
}

// subtleCompare compares two byte slices for equality without leaking timing information.
func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
