package tcesecurity

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	sm "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm"
)

// tsm2 非对称加密
const (
	Tsm2Algorithm = "tsm-sm2"
)

func init() {
	f := func(opts CryptoOpts) (Crypto, error) {
		privateKeyBytes, err := base64.StdEncoding.DecodeString(opts.PrivateKey)
		if err != nil {
			return nil, err
		}

		publicKeyBytes, err := base64.StdEncoding.DecodeString(opts.PublicKey)
		if err != nil {
			return nil, err
		}
		return NewTSM2Crypto(opts.Method, publicKeyBytes, privateKeyBytes)
	}
	registerCryptoFunc(Tsm2Algorithm, f)
}

// TSM-SM2算法
type TSM2Crypto struct {
	Version    string
	Prefix     string
	Method     string
	Ctx        *sm.SM2_ctx_t
	PrivateKey []byte
	PublicKey  []byte
}

// Encrypt 加密
func (c *TSM2Crypto) Encrypt(plaintext string) (string, error) {
	// 1，如果加密数据带有已加密前缀信息，则直接返回
	if strings.HasPrefix(plaintext, c.Prefix) {
		return plaintext, nil
	}
	plaintextByte := []byte(plaintext)
	ciphertext := make([]byte, len(plaintextByte)+200)
	var ciphertextLen int
	// 2，加密
	if code := sm.SM2Encrypt(
		c.Ctx, plaintextByte, len(plaintextByte), c.PublicKey, len(c.PublicKey),
		ciphertext, &ciphertextLen); code != 0 {
		return "", fmt.Errorf("encrypt failed, code: %d", code)
	}

	// 3，构造返回
	body := base64.StdEncoding.EncodeToString(ciphertext[:ciphertextLen])
	return c.Prefix + ":" + c.Method + ":" + c.Version + ":" + body, nil
}

// Decrypt 解密
func (c *TSM2Crypto) Decrypt(ciphertext string) (string, error) {
	// 1，如果解密数据前缀错误，直接返回
	if !strings.HasPrefix(ciphertext, c.Prefix) {
		return ciphertext, nil
	}
	// 2，验证method
	items := strings.Split(ciphertext, ":")
	if len(items) != 4 {
		return "", fmt.Errorf("invalid ciphertext-data format")
	}
	if items[1] != c.Method {
		return "", fmt.Errorf("invalid entrypted-data method")
	}

	// 3，解码
	ciphertextByte, err := base64.StdEncoding.DecodeString(items[3])
	if err != nil {
		return "", fmt.Errorf("invalid entrypted-data format")
	}
	// 4，解密
	plaintext := make([]byte, len(ciphertextByte)-96)
	var plaintextLen int
	if code := sm.SM2Decrypt(
		c.Ctx, ciphertextByte, len(ciphertextByte), c.PrivateKey, len(c.PrivateKey),
		plaintext, &plaintextLen); code != 0 {
		return "", fmt.Errorf("decrypt failed, code: %d", code)
	}
	return string(plaintext), nil
}

func NewTSM2Crypto(method string, publicKey, privateKey []byte) (*TSM2Crypto, error) {
	var ctx sm.SM2_ctx_t
	if code := sm.SM2InitCtx(&ctx); code != 0 {
		return nil, fmt.Errorf("init sm2 failed, code: %d", code)
	}
	return &TSM2Crypto{
		Version:    hex.EncodeToString([]byte(VERSION)),
		Prefix:     AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		Method:     hex.EncodeToString([]byte(method)),
		Ctx:        &ctx,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
