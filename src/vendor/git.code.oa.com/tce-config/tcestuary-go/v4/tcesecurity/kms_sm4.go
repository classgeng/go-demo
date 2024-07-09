package tcesecurity

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	kms "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tseckms/v20190118"
)

const (
	KMSSm4Algorithm = "kms-sm4-128-gcm"
	RemoteAlgorithm = "SM4_CBC_PKCS7PADDING"
)

func init() {
	f := func(opts CryptoOpts) (Crypto, error) {
		return NewKMSSm4Crypto(opts.Method, opts.KeyId, opts.SecretId, opts.SecretKey, opts.KMSServer)
	}
	registerCryptoFunc(KMSSm4Algorithm, f)
}

type KMSSm4Crypto struct {
	Client          *kms.Client
	Prefix          string
	Method          string
	RemoteAlgorithm string // kms中的加密算法
	Version         string
	KeyId           string
	KMSServer       string
}

// Encrypt 加密
func (c *KMSSm4Crypto) Encrypt(plaintext string) (string, error) {
	// 构造请求
	req := kms.NewEncryptRequest()
	req.SetDomain(c.KMSServer)
	req.KeyId = &c.KeyId
	req.Algorithm = &c.RemoteAlgorithm
	newPlaintext := base64.StdEncoding.EncodeToString([]byte(plaintext))
	req.Plaintext = &newPlaintext
	resp, err := c.Client.Encrypt(req)
	if err != nil {
		return "", err
	}
	ciphertext := *resp.Response.CiphertextBlob
	return c.Prefix + ":" + c.Method + ":" + c.Version + ":" + ciphertext, nil
}

// Decrypt 解密
func (c *KMSSm4Crypto) Decrypt(ciphertext string) (string, error) {
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
	realCiphertext := items[3]
	// 3，构造请求
	req := kms.NewDecryptRequest()
	req.SetDomain(c.KMSServer)
	req.CiphertextBlob = &realCiphertext
	resp, err := c.Client.Decrypt(req)
	if err != nil {
		return "", err
	}
	// 4，验证摘要
	plaintext := *resp.Response.Plaintext
	oriPlaintext, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		return "", err
	}

	return string(oriPlaintext), nil
}

// NewKMSSm4Crypto
func NewKMSSm4Crypto(method, keyId, secretId, secretKey, KMSServer string) (*KMSSm4Crypto, error) {
	if KMSServer == "" {
		return nil, fmt.Errorf("invalid kmsServer config")
	}
	cli, err := kms.NewClientWithSecretId(secretId, secretKey, DefaultRegion)

	if err != nil {
		return nil, err
	}
	// config transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli.WithHttpTransport(transport)
	return &KMSSm4Crypto{
		Prefix:          AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		RemoteAlgorithm: RemoteAlgorithm,
		Method:          hex.EncodeToString([]byte(method)),
		Version:         hex.EncodeToString([]byte(VERSION)),
		KeyId:           keyId,
		Client:          cli,
		KMSServer:       KMSServer,
	}, nil
}
