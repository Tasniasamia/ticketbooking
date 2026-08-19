package payment

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "errors"
    "io"
    "os"
)

func getEncryptionKey() []byte {
    key := os.Getenv("ENCRYPTION_KEY") // .env তে 32 character এর key রাখো
    if len(key) != 32 {
        // development এর জন্য fallback (production এ অবশ্যই .env থেকে নাও)
        key = "12345678901234567890123456789012"
    }
    return []byte(key)
}

func EncryptCredentials(creds map[string]string) (string, error) {
    plain, err := json.Marshal(creds)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(getEncryptionKey())
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, aesGCM.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := aesGCM.Seal(nonce, nonce, plain, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptCredentials(encrypted string) (map[string]string, error) {
    data, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(getEncryptionKey())
    if err != nil {
        return nil, err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := aesGCM.NonceSize()
    if len(data) < nonceSize {
        return nil, errors.New("invalid ciphertext")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plain, err := aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    var creds map[string]string
    if err := json.Unmarshal(plain, &creds); err != nil {
        return nil, err
    }
    return creds, nil
}