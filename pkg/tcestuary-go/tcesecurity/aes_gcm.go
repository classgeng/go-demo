package tcesecurity

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	Aes256GcmAlgorithm = "aes-256-gcm"
)

func init() {
	f := func(opts CryptoOpts) (Crypto, error) {
		return NewAesGcmCrypto(opts.Method, []byte(opts.AesKey))
	}
	registerCryptoFunc(Aes256GcmAlgorithm, f)
}

// AES-GCM算法
type AesGcmCrypto struct {
	Version string
	Prefix  string
	Method  string
	aesKey  []byte
	iv      []byte
	aad     []byte
}

// Encrypt加密
func (a *AesGcmCrypto) Encrypt(plaintext string) (string, error) {
	// 1，如果加密数据带有已加密前缀信息，则直接返回
	if strings.HasPrefix(plaintext, a.Prefix) {
		return plaintext, nil
	}
	// 2，加密
	// 2.1 cipher
	c, err := aes.NewCipher(a.aesKey)
	if err != nil {
		return "", err
	}

	// 2.2 gcm
	gcm, err := cipher.NewGCMWithNonceSize(c, 16)
	if err != nil {
		return "", err
	}
	// 2.3 seal
	bts := gcm.Seal(nil, a.iv, []byte(plaintext), a.aad)

	// cipher 将tag追加到密文后，16位
	tag := bts[len(bts)-16:]

	// 3，构造返回
	return a.Prefix + ":" + a.Method + ":" + a.Version + ":" +
		hex.EncodeToString(tag) + ":" + hex.EncodeToString(bts[:len(bts)-16]), nil
}

// Decrypt解密
func (a *AesGcmCrypto) Decrypt(ciphertext string) (string, error) {
	// 1，如果解密数据前缀错误，直接返回
	if !strings.HasPrefix(ciphertext, a.Prefix) {
		return ciphertext, nil
	}

	// 2，验证method
	items := strings.Split(ciphertext, ":")
	if len(items) != 5 {
		return "", fmt.Errorf("invalid ciphertext-data format")
	}
	if items[1] != a.Method {
		return "", fmt.Errorf("invalid entrypted-data method")
	}
	// 3，解码
	tag, err := hex.DecodeString(items[3])
	if err != nil {
		return "", fmt.Errorf("invalid entrypted-data tag")
	}

	rawEncrypted, err := hex.DecodeString(items[4])
	if err != nil {
		return "", fmt.Errorf("invalid entrypted-data format")
	}

	c, err := aes.NewCipher(a.aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCMWithNonceSize(c, 16)
	if err != nil {
		return "", err
	}
	bts, err := gcm.Open(nil, a.iv, append(rawEncrypted, tag...), a.aad)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

// NewAesGcmCrypto return AES-GCM Crypto
func NewAesGcmCrypto(method string, aesKey []byte) (*AesGcmCrypto, error) {
	return &AesGcmCrypto{
		Version: hex.EncodeToString([]byte(VERSION)),
		Prefix:  AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		Method:  hex.EncodeToString([]byte(method)),
		aesKey:  []byte(aesKey),
		iv:      []byte(aesKey[:16]),
		aad:     []byte(aesKey[:8]),
	}, nil
}
