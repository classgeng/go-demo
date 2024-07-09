package tcesecurity

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	sm "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm"
)

// tsm-sm2签名、验签
const (
	Tsm2SignAlgorithm = "tsm-sign"
)

func init() {
	f := func(opts SignOpts) (Signer, error) {
		privateKeyBytes, err := base64.StdEncoding.DecodeString(opts.PrivateKey)
		if err != nil {
			return nil, err
		}

		publicKeyBytes, err := base64.StdEncoding.DecodeString(opts.PublicKey)
		if err != nil {
			return nil, err
		}
		return NewTSMSign(opts.Method, publicKeyBytes, privateKeyBytes)
	}
	registerSignFunc(Tsm2SignAlgorithm, f)
}

// TSM-SM2 签名算法
type TSMSign struct {
	Version    string
	Prefix     string
	Method     string
	Id         []byte
	Ctx        *sm.SM2_ctx_t
	PrivateKey []byte
	PublicKey  []byte
}

// 签名
func (s *TSMSign) Sign(msg string) (string, error) {
	sign := make([]byte, 131)
	var signLen int
	msgByte := []byte(msg)

	if code := sm.SM2Sign(
		s.Ctx, msgByte, len(msgByte), s.Id, len(s.Id), s.PublicKey,
		len(s.PublicKey), s.PrivateKey, len(s.PrivateKey), sign, &signLen); code != 0 {
		return "", fmt.Errorf("sign failed, code: %d", code)
	}
	b64Sign := base64.StdEncoding.EncodeToString(sign[:signLen])
	return s.Prefix + ":" + s.Method + ":" + s.Version + ":" + b64Sign, nil
}

// 验签
func (s *TSMSign) Verify(msg, signValue string) (bool, error) {
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
		return false, fmt.Errorf("invalid entrypted-data method")
	}
	msgByte := []byte(msg)
	b64Sign := items[3]
	realSign, err := base64.StdEncoding.DecodeString(b64Sign)
	if err != nil {
		return false, fmt.Errorf("invalid signvalue-data format")
	}
	if code := sm.SM2Verify(s.Ctx, msgByte, len(msgByte), s.Id, len(s.Id),
		realSign, len(realSign), s.PublicKey, len(s.PublicKey),
	); code == 0 {
		return true, nil
	}
	return false, nil
}

func NewTSMSign(method string, publicKey, privateKey []byte) (*TSMSign, error) {
	var ctx sm.SM2_ctx_t
	if code := sm.SM2InitCtx(&ctx); code != 0 {
		return nil, fmt.Errorf("init sm2 failed, code: %d", code)
	}
	return &TSMSign{
		Version:    hex.EncodeToString([]byte(VERSION)),
		Prefix:     AlreadyEncryptPrefix + hex.EncodeToString([]byte(TceSecurity)),
		Method:     hex.EncodeToString([]byte(method)),
		Ctx:        &ctx,
		Id:         []byte(TceSecurity),
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
