package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
)

const (
    saltSize    = 16
    hashSize    = 32
)

func HashPassword(password string) (string, error) {
    salt := make([]byte, saltSize)
    _, err := rand.Read(salt)
    if err != nil {
        return "", fmt.Errorf("error generating salt: %v", err)
    }

    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    hashedPassword := hash.Sum(nil)

    // Combine salt and hash
    final := make([]byte, saltSize+hashSize)
    copy(final[:saltSize], salt)
    copy(final[saltSize:], hashedPassword)

    return base64.StdEncoding.EncodeToString(final), nil
}

func VerifyPassword(hashedPassword, password string) (bool, error) {
    decoded, err := base64.StdEncoding.DecodeString(hashedPassword)
    if err != nil {
        return false, fmt.Errorf("error decoding hash: %v", err)
    }

    if len(decoded) != saltSize+hashSize {
        return false, fmt.Errorf("invalid hash length")
    }

    salt := decoded[:saltSize]
    expectedHash := decoded[saltSize:]

    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    actualHash := hash.Sum(nil)

    return subtle.ConstantTimeCompare(expectedHash, actualHash) == 1, nil
}