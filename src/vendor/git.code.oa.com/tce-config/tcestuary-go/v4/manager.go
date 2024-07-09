package tcestuary

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"git.code.oa.com/tce-config/tcestuary-go/v4/logger"
)

type manager struct {
	// 配置目录
	Directory string

	// 配置目录下的 sdk.json
	ConfigCenter     *configcenter.ConfigCenter
	ConfigCenterFile string

	loadError        error // 加载错误信息
	loadCounter      int32 // 单元测试, 验证加载次数
	loadOnce         sync.Once
	writeCounter     int32 // 单元测试, 验证写出次数
	writeVersionOnce sync.Once
}

func newManager() *manager {
	c := &manager{
		Directory:        "/tce/conf/config/tce.config.center",
		ConfigCenterFile: "/tce/conf/config/tce.config.center/sdk.json",
		ConfigCenter:     configcenter.NewConfigCenter(),
	}
	return c
}

// Load 整个生命周期中, 仅在启动时加载一次配置文件
// 返回值表示“首次加载配置”的错误信息, 通常为 nil
func (c *manager) Load() error {

	c.writeVersionOnce.Do(c.WriteSDKVersion)

	// 通用加载逻辑. 其中, object 必须是对象指针
	load := func(file string, object interface{}) error {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			logger.Printf("read %s error, %s", file, err)
			return err
		}

		if err := c.ConfigCenter.Parse(b); err != nil {
			return err
		}

		return nil
	}

	c.loadOnce.Do(func() {
		atomic.AddInt32(&c.loadCounter, 1)
		c.loadError = load(c.ConfigCenterFile, c.ConfigCenter)
		if c.loadError != nil {
			logger.Printf("for the first time, load cofig file error")
		}
	})

	return c.loadError
}

// Debug 向终端输出已加载配置信息, 支持异常调试
func (c *manager) Debug() {
	log.Printf("config directory: %s\n", c.Directory)
	log.Printf("config file: %s\n", c.ConfigCenterFile)
	log.Printf("config info:\n")
	c.ConfigCenter.Debug()
}

// WriteSDKVersion 向配置目录下输出 SDK 版本号,追踪使用情况, 支持SDK升级
func (c *manager) WriteSDKVersion() {
	atomic.AddInt32(&c.writeCounter, 1)

	filename := filepath.Join(c.Directory, "sdk.golang.version")

	err := ioutil.WriteFile(filename, []byte(Version), 0644)
	if err != nil {
		logger.Printf("write version file error, %s", err)
	}
}
