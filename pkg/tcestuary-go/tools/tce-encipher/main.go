package main

import (
	"fmt"
	"log"
	"os"

	"git.code.oa.com/tce-config/tcestuary-go/v4"
	"github.com/urfave/cli/v2"
)

// 定义命令行错误码, 便于自动化集成识别异常
const encryptError int = 100
const decryptError int = 200

func main() {
	app := &cli.App{
		Name:  "encipher",
		Usage: "tce tool for password management",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "torwang",
				Email: "torwang@tencent.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "encrypt",
				Aliases:   []string{"e"},
				Usage:     "encrypt cleartext with aeskey",
				ArgsUsage: "cleartext",
				Action:    encrypt,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "aeskey",
						Aliases:  []string{"k"},
						Usage:    "AES key",
						Required: true,
					},
				},
			},
			{
				Name:      "decrypt",
				Aliases:   []string{"d"},
				Usage:     "decrypt ciphertext with aeskey",
				ArgsUsage: "ciphertext",
				Action:    decrypt,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "aeskey",
						Aliases:  []string{"k"},
						Usage:    "AES key",
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

func encrypt(c *cli.Context) error {
	if c.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(c, "encrypt", encryptError)
	}

	aeskey, cleartext := c.String("aeskey"), c.Args().First()
	if length := len([]byte(aeskey)); length != 16 && length != 24 && length != 32 {
		return cli.NewExitError("aeskey should be of 16/24/32 Byte", encryptError)
	}

	ciphertext, err := tcestuary.Encrypt(aeskey, cleartext)
	if err != nil {
		return cli.NewExitError(err, encryptError)
	}
	fmt.Println(ciphertext)

	return nil
}

func decrypt(c *cli.Context) error {
	if c.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(c, "decrypt", decryptError)
	}

	aeskey, ciphertext := c.String("aeskey"), c.Args().First()
	if length := len([]byte(aeskey)); length != 16 && length != 24 && length != 32 {
		return cli.NewExitError("aeskey should be of 16/24/32 Byte", decryptError)
	}

	cleartext, err := tcestuary.Decrypt(aeskey, ciphertext)
	if err != nil {
		return cli.NewExitError(err, decryptError)
	}
	fmt.Println(cleartext)

	return nil
}
