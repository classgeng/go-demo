package main

import (
	"fmt"
	"log"
	"os"

	"git.code.oa.com/tce-config/tcestuary-go/v4"
	"github.com/urfave/cli/v2"
)

// 定义命令行错误码, 便于自动化集成识别异常
const (
	ErrParamNotExist int = 100
	ErrGetConfig     int = 200
	decryptError     int = 300
)

func main() {
	app := &cli.App{
		Name:  "tce-config-sdk",
		Usage: "tce config tool for shell",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "torwang",
				Email: "torwang@tencent.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "getMysqlConfig",
				Usage:     "get database config (host, ip, port, user, passwd)",
				ArgsUsage: "configkey",
				Action:    getMysqlConfig,
				Flags:     []cli.Flag{},
			},
			{
				Name:      "debug",
				Usage:     "print sdk config",
				ArgsUsage: " ",
				Action:    debug,
				Flags:     []cli.Flag{},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getMysqlConfig(c *cli.Context) error {
	if c.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(c, "getMysqlConfig", ErrParamNotExist)
	}

	mysql, err := tcestuary.GetMysqlConfig(c.Args().First())
	if err != nil {
		return cli.NewExitError(err.Error(), ErrGetConfig)
	}
	fmt.Println(mysql.Host, mysql.IP, mysql.Port, mysql.User, mysql.Password)

	return nil
}

func debug(c *cli.Context) error {
	// 触发一次加载
	tcestuary.GetMysqlConfig("dbsql.not-exist")

	tcestuary.Debug()
	return nil
}
