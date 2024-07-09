package tcestuary

import (
	"encoding/json"
	"fmt"
	"os"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity"
)

type TransportSecurity interface {
	Encrypt(string) (string, error) // 加密，明文输入长度限制与算法相关
	Decrypt(string) (string, error) // 解密，密文输入长度限制与算法相关
}

func parseTransportSecretConfig() (configcenter.SecretConfig, error) {
	transportEnv := os.Getenv("TRANSPORT_SECRET")
	var secretConf configcenter.SecretConfig
	if transportEnv == "" {
		if err := std.Load(); err != nil {
			return secretConf, err
		}
		return std.ConfigCenter.SDK.TransportSecret, nil
	}
	err := json.Unmarshal([]byte(transportEnv), &secretConf)
	return secretConf, err
}

// NewTransportSecurity 传输安全组件
func NewTransportSecurity() (TransportSecurity, error) {
	// 判断是否加载TSM
	tsmConf, err := parseTSMSecretConfig()
	if err != nil {
		return nil, err
	}
	if tsmConf.PemAppid != "" {
		if err := InitTencentSMWithConfig(tsmConf); err != nil {
			return nil, err
		}
	}
	// 密钥配置
	secretConf, err := parseTransportSecretConfig()
	if err != nil {
		return nil, err
	}

	f, ok := tcesecurity.SupportAlgorithm[secretConf.Method]
	if !ok {
		return nil, fmt.Errorf("not support algorithm: %s", secretConf.Method)
	}
	return f(tcesecurity.CryptoOpts{
		Method:     secretConf.Method,
		AesKey:     secretConf.AesKey,
		Sm4Key:     secretConf.Sm4Key,
		PrivateKey: secretConf.PrivateKey,
		PublicKey:  secretConf.PublicKey,
		KeyId:      secretConf.KeyId,
		SecretId:   secretConf.SecretId,
		SecretKey:  secretConf.SecretKey,
		KMSServer:  secretConf.KMSServer,
	})
}
