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
	KMSSignAlgorithm    = "kms-sign" // 配置文件Method
	KMSSm2SignAlgorithm = "SM2DSA"   // 远程签名算法
	RawMessageType      = "RAW"
)

func init() {
	f := func(opts SignOpts) (Signer, error) {
		return NewKMSSign(opts.Method, opts.KeyId, opts.SecretId, opts.SecretKey, opts.KMSServer)
	}
	registerSignFunc(KMSSignAlgorithm, f)
}

type KMSSign struct {
	Client          *kms.Client
	Prefix          string
	Method          string
	Version         string
	KeyId           string
	KMSServer       string
	RemoteAlgorithm string
	MessageType     string
}

// Sign 生成签名
func (s *KMSSign) Sign(msg string) (string, error) {
	req := kms.NewSignByAsymmetricKeyRequest()
	req.SetDomain(s.KMSServer)
	req.KeyId = &s.KeyId
	req.Algorithm = &s.RemoteAlgorithm
	req.MessageType = &s.MessageType
	b64Msg := base64.StdEncoding.EncodeToString([]byte(msg))
	req.Message = &b64Msg
	resp, err := s.Client.SignByAsymmetricKey(req)
	if err != nil {
		return "", err
	}
	signature := *resp.Response.Signature
	return s.Prefix + ":" + s.Method + ":" + s.Version + ":" + signature, nil
}

// Verify 验证签名
func (s *KMSSign) Verify(msg, signValue string) (bool, error) {
	// 1，如果解密数据前缀错误，直接返回
	if !strings.HasPrefix(signValue, s.Prefix) {
		return false, nil
	}
	// 2，验证method
	items := strings.Split(signValue, ":")
	if len(items) != 4 {
		return false, fmt.Errorf("invalid signvalue-data format")
	}
	if items[1] != s.Method {
		return false, fmt.Errorf("invalid signvalue-data method")
	}

	realSign := items[3]
	// 3，构造请求
	req := kms.NewVerifyByAsymmetricKeyRequest()
	req.SetDomain(s.KMSServer)
	req.KeyId = &s.KeyId
	req.Algorithm = &s.RemoteAlgorithm
	req.MessageType = &s.MessageType
	req.SignatureValue = &realSign
	b64Msg := base64.StdEncoding.EncodeToString([]byte(msg))
	req.Message = &b64Msg
	resp, err := s.Client.VerifyByAsymmetricKey(req)
	if err != nil {
		return false, err
	}
	return *resp.Response.SignatureValid, nil
}

// NewKMSSign
func NewKMSSign(method, keyId, secretId, secretKey, KMSServer string) (*KMSSign, error) {
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
	return &KMSSign{
		Prefix:          AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		Method:          hex.EncodeToString([]byte(method)),
		Version:         hex.EncodeToString([]byte(VERSION)),
		RemoteAlgorithm: KMSSm2SignAlgorithm,
		MessageType:     RawMessageType,
		KeyId:           keyId,
		Client:          cli,
		KMSServer:       KMSServer,
	}, nil
}
