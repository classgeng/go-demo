package tcesecurity

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	sm "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm"
)

// sm4 对称加密

const (
	TSM4Algorithm = "tsm-sm4-128-gcm"
)

func init() {
	f := func(opts CryptoOpts) (Crypto, error) {
		return NewTSM4Crypto(opts.Method, []byte(opts.Sm4Key))
	}
	registerCryptoFunc(TSM4Algorithm, f)
}

// TSM-SM4算法
type TSM4Crypto struct {
	Version string
	Prefix  string
	Method  string
	sm4Key  []byte
	iv      []byte
	aad     []byte
}

// Encrypt加密
func (c *TSM4Crypto) Encrypt(plaintext string) (string, error) {
	// 1，如果加密数据带有已加密前缀信息，则直接返回
	if strings.HasPrefix(plaintext, c.Prefix) {
		return plaintext, nil
	}
	plaintextByte := []byte(plaintext)
	tag := make([]byte, 16)
	ciphertext := make([]byte, len(plaintextByte)+16-len(plaintextByte)%16)
	var ciphertextLen int
	tagLen := 16

	code := sm.SM4_GCM_Encrypt_NIST_SP800_38D(
		plaintextByte, len(plaintextByte), ciphertext, &ciphertextLen,
		tag, &tagLen, c.sm4Key, c.iv, len(c.iv), c.aad, len(c.aad))
	if code != 0 {
		return "", fmt.Errorf("encrypt failed, code: %d", code)
	}
	// 3，构造返回
	return c.Prefix + ":" + c.Method + ":" + c.Version + ":" +
		base64.StdEncoding.EncodeToString(tag[:tagLen]) + ":" +
		base64.StdEncoding.EncodeToString(ciphertext[:ciphertextLen]), nil
}

// Decrypt解密
func (c *TSM4Crypto) Decrypt(ciphertext string) (string, error) {
	// 1，如果解密数据前缀错误，直接返回
	if !strings.HasPrefix(ciphertext, c.Prefix) {
		return ciphertext, nil
	}

	// 2，验证method
	items := strings.Split(ciphertext, ":")
	if len(items) != 5 {
		return "", fmt.Errorf("invalid ciphertext-data format")
	}
	if items[1] != c.Method {
		return "", fmt.Errorf("invalid ciphertext-data method")
	}

	tag, err := base64.StdEncoding.DecodeString(items[3])
	if err != nil {
		return "", err
	}

	realCiphertext, err := base64.StdEncoding.DecodeString(items[4])
	if err != nil {
		return "", err
	}
	// 3，解码
	plaintext := make([]byte, len(realCiphertext))
	var plaintextLen int
	code := sm.SM4_GCM_Decrypt_NIST_SP800_38D(
		realCiphertext, len(realCiphertext), plaintext, &plaintextLen,
		tag, len(tag), c.sm4Key, c.iv, len(c.iv), c.aad, len(c.aad))
	if code != 0 {
		return "", fmt.Errorf("decrypt failed, code: %d", code)
	}
	return string(plaintext[:plaintextLen]), nil
}

// NewTSM4Crypto return TSM-SM4 Crypto
func NewTSM4Crypto(method string, sm4Key []byte) (*TSM4Crypto, error) {
	return &TSM4Crypto{
		Version: hex.EncodeToString([]byte(VERSION)),
		Prefix:  AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		Method:  hex.EncodeToString([]byte(method)),
		sm4Key:  []byte(sm4Key),
		iv:      []byte(sm4Key[:16]),
		aad:     []byte(sm4Key[:8]),
	}, nil
}
