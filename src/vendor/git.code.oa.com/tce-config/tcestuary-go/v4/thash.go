package tcestuary

import (
	"fmt"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity"
)

type THash interface {
	Update([]byte) error
	Digest() ([]byte, error)
}

type THasher interface {
	New() (THash, error)
}

type thasher struct {
	f tcesecurity.HashFunc
}

func (h *thasher) New() (THash, error) {
	return h.f()
}

func parseHashSecretConfig() (configcenter.HashConfig, error) {
	var hashConf configcenter.HashConfig
	if err := std.Load(); err != nil {
		return hashConf, err
	}
	return std.ConfigCenter.SDK.HashSecret, nil
}

// NewHasher 返回Hash生成器
func NewTHasher() (THasher, error) {
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
	hashConf, err := parseHashSecretConfig()
	if err != nil {
		return nil, err
	}
	f, ok := tcesecurity.SupportHashFunc[hashConf.Method]
	if !ok {
		return nil, fmt.Errorf("not support algorithm: %s", hashConf.Method)
	}
	return &thasher{f}, nil
}
