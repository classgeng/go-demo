package main

import (
	"fmt"
	"git.code.oa.com/tce-config/tcestuary-go/v4"
	"git.code.oa.com/tce-config/tcestuary-go/v4/middlewareconfig"
	"log"
)

func main() {

	fmt.Printf("%v", "start...")

	// 测试期间, 临时设置配置路径. 业务代码中请勿调用
	// tcestuary.SetConfigDirectory("/Library/GolandProjects/godemo/src/config")

	// mysql
	mysqlConfig := &middlewareconfig.MysqlConfig{}
	mysqlConfig.GetConfig("ocloud_api3")

	fmt.Printf("mysqlConfig.Pass: %v", mysqlConfig.Pass + " ")

	// key 采用两段式:
	// ocloud_api3: 对应 dbsql 别名;
	// api_sync: 对应数据库名称;
	m, err := tcestuary.GetMysqlConfig("ocloud_api3.api_sync")
	if err != nil {
		log.Println("return error, ", err)
	} else {
		log.Println(m.Host)
		log.Println(m.IP)
		log.Println(m.Port)
		log.Println(m.User)
		log.Println(m.Password)
		log.Println(m.Database)
	}

	// 配置加载异常, 调用 Debug 打印配置加载情况. 输出到终端
	tcestuary.Debug()
}
