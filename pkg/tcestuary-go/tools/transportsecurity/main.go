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
		Name:  "transport security",
		Usage: "tce config center shell tool",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "wentaoyin",
				Email: "wentaoyin@tencent.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "Encrypt",
				Usage:     "加密",
				ArgsUsage: "",
				Action:    Encrypt,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "configDirectory",
						Usage:    "sdk.json 配置文件路径 `DIR`",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "plaintext",
						Usage:    "明文",
						Required: true,
					},
				},
			},
			{
				Name:      "Decrypt",
				Usage:     "解密",
				ArgsUsage: "",
				Action:    Decrypt,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "configDirectory",
						Usage:    "sdk.json 配置文件路径 `DIR`",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "ciphertext",
						Usage:    "密文",
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
func Encrypt(ctx *cli.Context) error {

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}
	st, err := tcestuary.NewTransportSecurity()
	if err != nil {
		return err
	}
	ret, err := st.Encrypt(ctx.String("plaintext"))
	if err != nil {
		return err
	}
	fmt.Println("输入明文： ", ctx.String("plaintext"))
	fmt.Println("加密结果：")
	fmt.Println(ret)
	return nil
}

// 解密
func Decrypt(ctx *cli.Context) error {

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}
	st, err := tcestuary.NewTransportSecurity()
	if err != nil {

		return err
	}
	ret, err := st.Decrypt(ctx.String("ciphertext"))
	if err != nil {
		return err
	}
	fmt.Println("输入密文：", ctx.String("ciphertext"))
	fmt.Println("解析明文：", ret)
	return nil
}
