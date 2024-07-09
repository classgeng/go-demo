package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"git.code.oa.com/tce-config/tcestuary-go/v4"

	"github.com/urfave/cli/v2"
)

// ToolError 定义命令行错误码, 便于自动化集成识别异常
const ToolError int = 200

func main() {
	app := &cli.App{
		Name:  "Tencent Hash Client",
		Usage: "tce config center shell tool",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "wentaoyin",
				Email: "wentaoyin@tencent.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "Hash",
				Usage:     "Hash",
				ArgsUsage: "",
				Action:    Hash,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "configDirectory",
						Usage:    "sdk.json 配置文件路径 `DIR`",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "msg",
						Usage:    "消息",
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// 加密
func Hash(ctx *cli.Context) error {

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}
	hasher, err := tcestuary.NewTHasher()
	if err != nil {
		return err
	}
	c, err := hasher.New()
	if err != nil {
		return err
	}
	_ = c.Update([]byte(ctx.String("msg")))
	out, _ := c.Digest()
	fmt.Println("输入消息： ", ctx.String("msg"))
	fmt.Println("散列值：")
	fmt.Println(hex.EncodeToString(out))
	return nil
}
