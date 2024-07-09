package tcesecurity

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
)

const (
	Rsa2048Algorithm = "rsa-2048"
	Rsa1024Algorithm = "rsa-1024"
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
		return NewRsaCrypto(opts.Method, publicKeyBytes, privateKeyBytes)
	}
	registerCryptoFunc(Rsa2048Algorithm, f)
	registerCryptoFunc(Rsa1024Algorithm, f)
}

// Rsa加密、解密
type RsaCrypto struct {
	Prefix     string
	Method     string
	Version    string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// Encrypt 加密
func (r *RsaCrypto) Encrypt(plaintext string) (string, error) {
	// 1，如果加密数据带有已加密前缀信息，则直接返回
	if strings.HasPrefix(plaintext, r.Prefix) {
		return plaintext, nil
	}

	// 2，加密
	bts, err := rsa.EncryptPKCS1v15(rand.Reader, r.PublicKey, []byte(plaintext))
	if err != nil {
		return "", err
	}
	// 3，构造返回
	body := base64.StdEncoding.EncodeToString(bts)
	return r.Prefix + ":" + r.Method + ":" + r.Version + ":" + body, nil

}

// Decrypt 解密
func (r *RsaCrypto) Decrypt(ciphertext string) (string, error) {
	// 1，如果解密数据前缀错误，直接返回
	if !strings.HasPrefix(ciphertext, r.Prefix) {
		return ciphertext, nil
	}
	// 2，验证method
	items := strings.Split(ciphertext, ":")
	if len(items) != 4 {
		return "", fmt.Errorf("invalid ciphertext-data format")
	}
	if items[1] != r.Method {
		return "", fmt.Errorf("invalid entrypted-data method")
	}
	// 3，解码
	rawEncrypted, err := base64.StdEncoding.DecodeString(items[3])
	if err != nil {
		fmt.Println(err.Error())
		return "", fmt.Errorf("invalid entrypted-data format 2")
	}
	// 4，解密
	bts, err := rsa.DecryptPKCS1v15(rand.Reader, r.PrivateKey, rawEncrypted)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

// NewResCrypto
func NewRsaCrypto(method string, publicKey, privateKey []byte) (*RsaCrypto, error) {
	// 构造publicKey
	pubBlock, _ := pem.Decode(publicKey)
	if pubBlock == nil {
		return nil, fmt.Errorf("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}
	// 构造privateKey
	priBlock, _ := pem.Decode(privateKey)
	if priBlock == nil {
		return nil, fmt.Errorf("private key error")
	}
	validPrivateKey, err := x509.ParsePKCS1PrivateKey(priBlock.Bytes)
	if err != nil {
		return nil, err
	}
	// 判断privateKey是否可用
	err = validPrivateKey.Validate()
	if err != nil {
		return nil, err
	}

	return &RsaCrypto{
		Prefix:     AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		Method:     hex.EncodeToString([]byte(method)),
		Version:    hex.EncodeToString([]byte(VERSION)),
		PrivateKey: validPrivateKey,
		PublicKey:  pubInterface.(*rsa.PublicKey),
	}, nil
}
