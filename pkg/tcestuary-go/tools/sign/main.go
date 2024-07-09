package main

import (
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
		Name:  "sign",
		Usage: "tce config center shell tool",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "wentaoyin",
				Email: "wentaoyin@tencent.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "Sign",
				Usage:     "签名",
				ArgsUsage: "",
				Action:    Sign,
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
			{
				Name:      "Verify",
				Usage:     "验签",
				ArgsUsage: "",
				Action:    Verify,
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
					&cli.StringFlag{
						Name:     "sign",
						Usage:    "签名",
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

// 签名
func Sign(ctx *cli.Context) error {

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}
	s, err := tcestuary.NewSigner()
	if err != nil {
		return err
	}
	ret, err := s.Sign(ctx.String("msg"))
	if err != nil {
		return err
	}
	fmt.Println("输入消息： ", ctx.String("msg"))
	fmt.Println("签名：")
	fmt.Println(ret)
	return nil
}

// 验签
func Verify(ctx *cli.Context) error {

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}
	s, err := tcestuary.NewSigner()
	if err != nil {
		return err
	}
	ret, err := s.Verify(ctx.String("msg"), ctx.String("sign"))
	if err != nil {
		return err
	}
	fmt.Println("输入签名：", ctx.String("msg"))
	fmt.Println("验证签名：", ret)
	return nil
}
