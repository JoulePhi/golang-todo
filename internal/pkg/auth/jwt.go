package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var secretKey = []byte("381c64015c0073c9a112cd4f0ff395ad9ac957ebd67b7f3ac2bcbef5a0f2ec99") 

type Claims struct {
    UserID    int64  `json:"user_id"`
    Username  string `json:"username"`
    ExpiresAt int64  `json:"exp"`
}

func GenerateToken(userID int64, username string) (string, error) {
    // Create header
    header := map[string]string{
        "alg": "HS256",
        "typ": "JWT",
    }
    
    headerJSON, err := json.Marshal(header)
    if err != nil {
        return "", err
    }
    
    // Create claims
    claims := Claims{
        UserID:    userID,
        Username:  username,
        ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
    }
    
    claimsJSON, err := json.Marshal(claims)
    if err != nil {
        return "", err
    }
    
    // Create signature
    encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
    encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
    
    signatureInput := encodedHeader + "." + encodedClaims
    signature := generateSignature([]byte(signatureInput))
    
    // Combine to create token
    token := signatureInput + "." + base64.RawURLEncoding.EncodeToString(signature)
    
    return token, nil
}

func ValidateToken(token string) (*Claims, error) {
    // Split token into parts
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid token format")
    }
    
    // Verify signature
    signatureInput := parts[0] + "." + parts[1]
    expectedSignature := generateSignature([]byte(signatureInput))
    
    actualSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
    if err != nil {
        return nil, fmt.Errorf("invalid signature encoding")
    }
    
    if !hmac.Equal(expectedSignature, actualSignature) {
        return nil, fmt.Errorf("invalid signature")
    }
    
    // Decode claims
    claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return nil, fmt.Errorf("invalid claims encoding")
    }
    
    var claims Claims
    if err := json.Unmarshal(claimsJSON, &claims); err != nil {
        return nil, fmt.Errorf("invalid claims format")
    }
    
    // Verify expiration
    if claims.ExpiresAt < time.Now().Unix() {
        return nil, fmt.Errorf("token has expired")
    }
    
    return &claims, nil
}

func generateSignature(data []byte) []byte {
    h := hmac.New(sha256.New, secretKey)
    h.Write(data)
    return h.Sum(nil)
}