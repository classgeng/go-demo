package tcestuary

import (
	"encoding/json"
	"fmt"
	"os"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity"
)

type StorageSecurity interface {
	Encrypt(string) (string, error) // 加密，明文输入长度限制与算法相关
	Decrypt(string) (string, error) // 解密，密文输入长度限制与算法相关
}

func parseStorageSecretConfig() (configcenter.SecretConfig, error) {
	storageEnv := os.Getenv("STORAGE_SECRET")
	var secretConf configcenter.SecretConfig
	if storageEnv == "" {
		if err := std.Load(); err != nil {
			return secretConf, err
		}
		return std.ConfigCenter.SDK.StorageSecret, nil
	}
	err := json.Unmarshal([]byte(storageEnv), &secretConf)
	return secretConf, err
}

// NewStorageSecurity 存储安全组件
func NewStorageSecurity() (StorageSecurity, error) {
	// 判断是否需要初始化TSM
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
	secretConf, err := parseStorageSecretConfig()
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

// NewPasswdSecret 存储安全组件
func NewPasswdSecret() (StorageSecurity, error) {
	// 判断是否需要初始化TSM
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
	secretConf, err := parsePasswdSecretConfig()
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

func parsePasswdSecretConfig() (configcenter.SecretConfig, error) {
	var secretConf configcenter.SecretConfig
	if err := std.Load(); err != nil {
		return secretConf, err
	}
	// 兼容 method 为空场景，历史版本，走默认 aes
	if std.ConfigCenter.SDK.PasswdSecret.Method == "" {
		std.ConfigCenter.SDK.PasswdSecret.Method = tcesecurity.Aes256CbcAlgorithm
		std.ConfigCenter.SDK.PasswdSecret.AesKey = std.ConfigCenter.SDK.PasswdSecret.V1Aeskey
	}
	return std.ConfigCenter.SDK.PasswdSecret, nil
}
