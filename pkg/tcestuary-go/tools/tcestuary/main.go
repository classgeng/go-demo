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
		Name:  "tcestuary",
		Usage: "tce config center shell tool",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "torwang",
				Email: "torwang@tencent.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "getMysqlConfig",
				Usage:     "获取数据库配置信息",
				ArgsUsage: "dbsql.database",
				Action:    GetMysqlConfig,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "configDirectory",
						Usage:    "sdk.json 配置文件路径 `DIR`",
						Required: false,
					},
				},
			},
			{
				Name:      "getMysqlConfigAllRegion",
				Usage:     "获取数据库配置信息(所有 Region 实例)",
				ArgsUsage: "dbsql.database",
				Action:    GetMysqlConfigAllRegion,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "configDirectory",
						Usage:    "sdk.json 配置文件路径 `DIR`",
						Required: false,
					},
				},
			},
			{
				Name:      "getMysqlConfigMainRegion",
				Usage:     "获取数据库配置信息(Main Region 实例)",
				ArgsUsage: "dbsql.database",
				Action:    GetMysqlConfigMainRegion,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "configDirectory",
						Usage:    "sdk.json 配置文件路径 `DIR`",
						Required: false,
					},
				},
			},
			{
				Name:      "getMysqlConfigAllZone",
				Usage:     "获取数据库配置信息(所有 Zone 实例)",
				ArgsUsage: "dbsql.database",
				Action:    GetMysqlConfigAllZone,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "configDirectory",
						Usage:    "sdk.json 配置文件路径 `DIR`",
						Required: false,
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

// GetMysqlConfig 成功时, 向 STDOUT 顺序写出: host / ip / port / user / passwd
func GetMysqlConfig(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(ctx, "getMysqlConfig", ToolError)
	}

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}

	mysql, err := tcestuary.GetMysqlConfig(ctx.Args().First())
	if err != nil {
		return cli.NewExitError(err, ToolError)
	}

	fmt.Println(mysql.Host, mysql.IP, mysql.Port, mysql.User, mysql.Password)

	return nil
}

// GetMysqlConfigAllRegion 成功时, 向 STDOUT 顺序写出多行: host / ip / port / user / passwd / regionid
func GetMysqlConfigAllRegion(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(ctx, "getMysqlConfigAllRegion", ToolError)
	}

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}

	mysql, err := tcestuary.GetMysqlConfigAllRegion(ctx.Args().First())
	if err != nil {
		return cli.NewExitError(err, ToolError)
	}

	for _, item := range mysql {
		fmt.Println(item.Host, item.IP, item.Port, item.User, item.Password, item.RegionID)
	}

	return nil
}

// GetMysqlConfigMainRegion 成功时, 向 STDOUT 顺序写出多行: host / ip / port / user / passwd / regionid
func GetMysqlConfigMainRegion(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(ctx, "getMysqlConfigMainRegion", ToolError)
	}

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}

	mysql, err := tcestuary.GetMysqlConfigAllRegion(ctx.Args().First())
	if err != nil {
		return cli.NewExitError(err, ToolError)
	}

	mainRegionName, err := tcestuary.GetMainRegionName()
	if err != nil {
		return cli.NewExitError(err, ToolError)
	}

	for _, item := range mysql {
		if mainRegionName == item.RegionName {
			fmt.Println(item.Host, item.IP, item.Port, item.User, item.Password, item.RegionID)
			break
		}
	}

	return nil
}

// GetMysqlConfigAllZone 成功时, 向 STDOUT 顺序写出多行: host / ip / port / user / passwd / regionid / zoneid
func GetMysqlConfigAllZone(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(ctx, "getMysqlConfigAllZone", ToolError)
	}

	if ctx.IsSet("configDirectory") {
		tcestuary.SetConfigDirectory(ctx.String("configDirectory"))
	}

	mysql, err := tcestuary.GetMysqlConfigAllZone(ctx.Args().First())
	if err != nil {
		return cli.NewExitError(err, ToolError)
	}

	for _, item := range mysql {
		fmt.Println(item.Host, item.IP, item.Port, item.User, item.Password, item.RegionID, item.ZoneID)
	}

	return nil
}
