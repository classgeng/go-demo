package tcestuary

import (
	"encoding/base64"
	"fmt"
	"sync"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	sm "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm"
)

var initOnce sync.Once

func parseTSMSecretConfig() (configcenter.TSMConfig, error) {
	var tsmConf configcenter.TSMConfig
	if err := std.Load(); err != nil {
		return tsmConf, nil
	}
	return std.ConfigCenter.SDK.TSMSecret, nil
}

// 初始化TSM
func initTencentSM(appid, bundle, cert []byte) error {
	var code int
	initOnce.Do(func() {
		code = sm.InitTencentSMWithCert(appid, bundle, cert)
	})
	if code == 0 {
		return nil
	}
	return fmt.Errorf("Init TSM failed")
}

// 通过TSMConfig初始化TSM
func InitTencentSMWithConfig(tsmConf configcenter.TSMConfig) error {
	cert, err := base64.StdEncoding.DecodeString(tsmConf.PemInitLibWithCert)
	if err != nil {
		return fmt.Errorf("invalid tsm cert: %s", err.Error())
	}
	return initTencentSM([]byte(tsmConf.PemAppid), nil, cert)
}
