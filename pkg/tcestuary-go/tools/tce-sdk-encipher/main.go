package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"git.code.oa.com/tce-config/tcestuary-go/v4"
	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"github.com/urfave/cli/v2"
)

// 定义命令行错误码, 便于自动化集成识别异常
const encryptError int = 100
const decryptError int = 200

func main() {
	app := &cli.App{
		Name:  "tce-sdk-encipher",
		Usage: "encrypt password field in sdk.json, for test only",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "torwang",
				Email: "torwang@tencent.com",
			},
		},
		Action: encrypt,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func encrypt(c *cli.Context) error {
	// for test
	// tcestuary.SetConfigDirectory(".")

	// 触发配置加载
	_, err := tcestuary.GetMysqlConfig("not-exist.not-exist")
	if err != tcestuary.ErrNotFound {
		return cli.NewExitError(err.Error(), 100)
	}

	// 加密配置
	aeskey := tcestuary.GetConfigCenterPtr().SDK.PasswdSecret.V1Aeskey
	if len(aeskey) != 16 && len(aeskey) != 24 && len(aeskey) != 32 {
		return cli.NewExitError("aeskey length must be 16/24/32 byte", 100)
	}

	// 加密配置项
	center := tcestuary.GetConfigCenterPtr()
	for _, wrap := range center.Mysqls {
		switch wrap.Scope {
		case configcenter.ScopeFlat:
			mysql := wrap.Object.(*configcenter.Mysql)
			if !strings.HasPrefix(mysql.Password, "AES+") {
				p, err := tcestuary.Encrypt(aeskey, mysql.Password)
				if err != nil {
					return cli.NewExitError("encode error, "+err.Error(), 100)
				}
				mysql.Password = p
			}
		case configcenter.ScopeAllRegion:
			mysqlR := wrap.Object.([]*configcenter.MysqlWithRegion)
			for _, mysql := range mysqlR {
				if !strings.HasPrefix(mysql.Service.Password, "AES+") {
					p, err := tcestuary.Encrypt(aeskey, mysql.Service.Password)
					if err != nil {
						return cli.NewExitError("encode error, "+err.Error(), 100)
					}
					mysql.Service.Password = p
				}
			}
		case configcenter.ScopeAllZone:
			mysqlZ := wrap.Object.([]*configcenter.MysqlWithZone)
			for _, mysql := range mysqlZ {
				if !strings.HasPrefix(mysql.Service.Password, "AES+") {
					p, err := tcestuary.Encrypt(aeskey, mysql.Service.Password)
					if err != nil {
						return cli.NewExitError("encode error, "+err.Error(), 100)
					}
					mysql.Service.Password = p
				}
			}
		}

	}

	// 还原成原始配置结构
	origin, err := PackIntoOriginConfigCenter(tcestuary.GetConfigCenterPtr())
	if err != nil {
		return cli.NewExitError(err, 100)
	}
	// 写出配置
	err = write(origin, "/tce/conf/config/tce.config.center/sdk.json")
	if err != nil {
		return cli.NewExitError("overwrite sdk.json error", 100)
	}

	log.Println("encrypt sdk.json done")

	return nil
}

// write 覆盖写 sdk.json
func write(cc *configcenter.OriginConfigCenter, filename string) error {
	data, err := json.MarshalIndent(cc, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// PackIntoOriginConfigCenter 从内存结构化格式还原成原始配置格式
func PackIntoOriginConfigCenter(cc *configcenter.ConfigCenter) (*configcenter.OriginConfigCenter, error) {
	origin := new(configcenter.OriginConfigCenter)

	// Base 拷贝
	origin.Base.Regions = cc.Base.Regions
	origin.Base.Zones = cc.Base.Zones

	// SDK 拷贝
	origin.SDK.PasswdSecret = cc.SDK.PasswdSecret

	// Mysql 拷贝
	origin.Mysqls = make(map[string]json.RawMessage, len(cc.Mysqls))
	for dbsql, wrap := range cc.Mysqls {
		raw, err := json.MarshalIndent(wrap.Object, "", "    ")
		if err != nil {
			return nil, err
		}
		origin.Mysqls[dbsql] = raw
	}

	return origin, nil
}
