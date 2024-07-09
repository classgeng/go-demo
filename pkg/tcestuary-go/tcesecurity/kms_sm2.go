package tcesecurity

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	sm "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm"

	kms "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tseckms/v20190118"
)

const (
	DefaultRegion   = ""
	KMSSm2Algorithm = "kms-sm2"
)

func init() {
	f := func(opts CryptoOpts) (Crypto, error) {
		return NewKMSSm2Crypto(opts.Method, opts.KeyId, opts.SecretId, opts.SecretKey, opts.KMSServer)
	}
	registerCryptoFunc(KMSSm2Algorithm, f)
}

type KMSSm2Crypto struct {
	Client    *kms.Client
	Prefix    string
	Method    string
	Version   string
	KeyId     string
	KMSServer string
}

// 生成摘要
func (c *KMSSm2Crypto) genDigest(body string) (string, error) {
	out := make([]byte, 32)
	data := []byte(body + TceSecurity)
	code := sm.SM3(data[:], len(data), out)
	if code != 0 {
		return "", fmt.Errorf("generate new plaintext failed")
	}
	return hex.EncodeToString(out), nil
}

// Encrypt 加密
func (c *KMSSm2Crypto) Encrypt(plaintext string) (string, error) {
	// 构造请求
	req := kms.NewAsymmetricSm2EncryptRequest()
	req.SetDomain(c.KMSServer)
	req.KeyId = &c.KeyId
	digest, err := c.genDigest(plaintext)
	if err != nil {
		return "", err
	}
	newPlaintext := base64.StdEncoding.EncodeToString(
		[]byte(hex.EncodeToString([]byte(plaintext)) + "|" + digest),
	)
	// plaintext: 待加密数据，长度不能超过160字节
	req.Plaintext = &newPlaintext
	resp, err := c.Client.AsymmetricSm2Encrypt(req)
	if err != nil {
		return "", err
	}
	ciphertext := *resp.Response.Ciphertext
	return c.Prefix + ":" + c.Method + ":" + c.Version + ":" + ciphertext, nil
}

// Decrypt 解密
func (c *KMSSm2Crypto) Decrypt(ciphertext string) (string, error) {
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
	req := kms.NewAsymmetricSm2DecryptRequest()
	req.SetDomain(c.KMSServer)
	req.KeyId = &c.KeyId
	req.Ciphertext = &realCiphertext

	resp, err := c.Client.AsymmetricSm2Decrypt(req)
	if err != nil {
		return "", err
	}
	// 4，验证摘要
	plaintext := *resp.Response.Plaintext
	plaintextBytes, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		return "", err
	}

	pItems := strings.Split(string(plaintextBytes), "|")
	if len(pItems) != 2 {
		return "", fmt.Errorf("invalid ciphertext-data format")
	}
	originPlaintext, err := hex.DecodeString(pItems[0])
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext-data format")
	}

	digest, err := c.genDigest(string(originPlaintext))
	if err != nil {
		return "", nil
	}
	if digest != pItems[1] {
		return "", fmt.Errorf("integrity check failed")
	}

	return string(originPlaintext), nil
}

// NewKMSSm2Crypto
func NewKMSSm2Crypto(method, keyId, secretId, secretKey, KMSServer string) (*KMSSm2Crypto, error) {
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

	return &KMSSm2Crypto{
		Prefix:    AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		Method:    hex.EncodeToString([]byte(method)),
		Version:   hex.EncodeToString([]byte(VERSION)),
		KeyId:     keyId,
		Client:    cli,
		KMSServer: KMSServer,
	}, nil
}
