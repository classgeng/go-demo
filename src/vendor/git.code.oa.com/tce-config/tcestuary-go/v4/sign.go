package tcestuary

import (
	"fmt"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity"
)

type Signer interface {
	Sign(msg string) (string, error)            // 签名
	Verify(msg, signValue string) (bool, error) // 验证签名
}

func parseSignConfig() (configcenter.SecretConfig, error) {
	var secretConf configcenter.SecretConfig
	if err := std.Load(); err != nil {
		return secretConf, err
	}
	return std.ConfigCenter.SDK.SignSecret, nil
}

// 签名、验签组件
func NewSigner() (Signer, error) {
	//  判断是否需要初始化TSM
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
	secretConf, err := parseSignConfig()
	if err != nil {
		return nil, err
	}

	f, ok := tcesecurity.SupportSignFunc[secretConf.Method]
	if !ok {
		return nil, fmt.Errorf("not support algorithm: %s", secretConf.Method)
	}
	return f(tcesecurity.SignOpts{
		Method:     secretConf.Method,
		KeyId:      secretConf.KeyId,
		SecretId:   secretConf.SecretId,
		SecretKey:  secretConf.SecretKey,
		KMSServer:  secretConf.KMSServer,
		PublicKey:  secretConf.PublicKey,
		PrivateKey: secretConf.PrivateKey,
	})
}
