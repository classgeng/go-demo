#### 代码引用参考

安装：
```shell script
% go get git.code.oa.com/tce-config/tcestuary-go/v4
```
注意：如果需要执行go mod vendor, 须确认路径vendor下内容，如果缺少库文件，可以从代码仓库复制补齐。

可能缺少tencentsm路径下的include/* 和lib/*，补齐后的tencentsm路径如下：
```shell script
% tree vendor/git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm
vendor/git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm
├── base.go
├── include
│   └── sm.h
├── lib
│   ├── darwin
│   │   ├── libTencentSM.a
│   │   └── libgmp.a
│   ├── darwin_arm
│   │   ├── libTencentSM.a
│   │   └── libgmp.a
│   ├── linux
│   │   ├── libTencentSM.a
│   │   └── libgmp.a
│   ├── linux_arm
│   │   ├── libTencentSM.a
│   │   └── libgmp.a
│   ├── windows
│   │   ├── libTencentSM.dll
│   │   └── libgmp.a
│   └── windows_arm
│       ├── libTencentSM.dll
│       └── libgmp.a
├── sm2.go
├── sm3.go
└── sm4.go
```

补充说明: 

git.code.oa.com 是内部域名，GOMODULE 会遇到私有仓库的问题:

1. 设置private repository

    ```shell
    go env -w GONOPROXY="git.code.oa.com"
    go env -w GONOSUMDB="git.code.oa.com"
    go env -w GOPRIVATE="git.code.oa.com"
    ```

2. 配置git replace，~/.gitconfig中增加如下内容

   ```tex
   [url "git@git.code.oa.com:"]
       insteadOf= https://git.code.oa.com/
   ```

#### 使用案例:
```
import (
	"log"

	"git.code.oa.com/tce-config/tcestuary-go/v4"
)

func main() {
	// 测试期间, 临时设置配置路径. 业务代码中请勿调用
	tcestuary.SetConfigDirectory(".")

	// key 采用两段式:
	// ocloud_api3: 对应 dbsql 别名;
	// api_sync: 对应数据库名称;
	m, err := tcestuary.GetMysqlConfig("ocloud_api3.api_sync")
	if err != nil {
		log.Println("return error, ", err)
	}

	// 从返回的 Mysql 配置对象中取配置
	// 参考文档: http://tapd.oa.com/OneBank/markdown_wikis/#1010130691010773371@toc1
	log.Println(m.Host)
	log.Println(m.IP)
	log.Println(m.Port)
	log.Println(m.User)
	log.Println(m.Password)
	log.Println(m.Database)

	// 配置加载异常, 调用 Debug 打印配置加载情况. 输出到终端
	tcestuary.Debug()
}
```

#### SDK 接口说明

##### 地域相关接口

|  接口名称   | 描述  |
|  ----  | ----  |
| GetRegion	 | 获取地域属性信息 |
| GetZone  | 获取可用区属性信息 |

##### Mysql 配置接口

|  接口名称   | 描述  |
|  ----  | ----  |
| GetMysqlConfig	 | 数据库五元组, dbsql.scope = GLOBAL/REGION/ZONE 时使用 |
| GetMysqlConfigAllRegon  | 所有 Region 的数据库五元组列表, dbsql.scope = ALL_REGION 时使用|
| GetMysqlConfigAllZone  | 所有 Zone 的数据库五元组列表, dbsql.scope = ALL_ZONE 时使用|

===

#### 加解密 Demo

通常不需要调用加解密接口, 如果有需要请通知 torwang.

```
import (
	"fmt"

	"git.code.oa.com/tce-config/tcestuary-go/v4"
)

func main() {
	//key的长度必须是16、24或者32字节，分别用于选择AES-128, AES-192, or AES-256
	aeskey := "f13c3f40c60db7f32ce6a5e0143f09ea"
	fmt.Println("密钥:", aeskey)
	origin := "a98a62d62#a+3@a2ed"
	fmt.Println("明文:", origin)

	pass, err := tcestuary.Encrypt(aeskey, origin)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("加密后:", pass)

	tpass, err := tcestuary.Decrypt(aeskey, pass)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("解密后:", tpass)

	if origin != tpass {
		fmt.Println("fatal, not equal")
	}
}

// 密钥: f13c3f40c60db7f32ce6a5e0143f09ea
// 明文: a98a62d62#a+3@a2ed
// 加密后: AES+oPHDfoh4+34a6dda167861f008235acf656a09d3546b0801083fece707fea6982c5c9f114
// 解密后: a98a62d62#a+3@a2ed
```

#### 加解密工具
1. 安装工具
```
go get git.code.oa.com/tce-config/encipher/tce-encipher
```

2. 工具使用说明
```
./tce-encipher encrypt --aeskey "f13c3f40c60db7f32ce6a5e0143f09ea" "a98a62d62#a+3@a2e"

./tce-encipher decrypt --aeskey "f13c3f40c60db7f32ce6a5e0143f09ea" AES+V1+0fc61e78b71f288f160b1797bc0a02ba162956bc87d2c14e5bb696068942c4a64
```

3. 说明

    - 密码盐随机生成, 每次加密输出可能不同

#### 密钥格式说明
```
AES+oPHDfoh4+34a6dda167861f008235acf656a09d3546b0801083fece707fea6982c5c9f114
```
- AES 标识加密方法
- oPHDfoh4 随机密码盐
- 34a6dda167861f008235acf656a09d3546b0801083fece707fea6982c5c9f114 加盐后的密文

#### 兼容性说明
背景:
1. 升级过程, 需要先升级组件, 再刷新现场配置包.
2. 部分组件, 需要向历史版本合入, 历史版本只有明文密码.

因此, 加解密工具需要兼容“明文密码”, 即: 如果输入密文非 “AES+” 开头, 则原样返回

#### 国密加密、解密
该版本支持的国密加密算法包括：kms-sm2, kms-sm4, tsm-sm2, tsm-sm4, 另外，还支持aes-256-gcm，rsa-1024, rsa-2048算法。
算法的选择由配置文件决定，不需要在代码里指明：生产环境默认读取/tce/conf/config/tce.config.center/sdk.json配置文件，该文件在渲染时写入了必要的秘钥信息。
国密加解密的应用场景包括：安全存储和安全传输。

##### 相关接口
|  接口名称   | 描述  |
|  ----  | ----  |
| Encrypt  | 加密 |
| Decrypt  | 解密 |

```go
// 安全存储
type StorageSecurity interface {
	Encrypt(string) (string, error) // 加密，明文输入长度限制与算法相关
	Decrypt(string) (string, error) // 解密，密文输入长度限制与算法相关
}

// 安全传输
type TransportSecurity interface {
	Encrypt(string) (string, error) // 加密，明文输入长度限制与算法相关
	Decrypt(string) (string, error) // 解密，密文输入长度限制与算法相关
}
```

##### Demo
```go
package main

import (
    "git.code.oa.com/tce-config/tcestuary-go/v4"
    "fmt"
)

func main() {
    // 1，安全存储 
    st, err := tcestuary.NewStorageSecurity()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    // 1.1 加密
    plaintext := "明文"
    ciphertext, _ := st.Encrypt(plaintext)

    // 1.2 解密
    plaintext2, _ := st.Decrypt(ciphertext)
    fmt.Println(plaintext == plaintext2)

    // 2, 安全传输
    tr, err := tcestuary.NewTransportSecurity()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    // 2.1 加密
    ciphertext2, _ := tr.Encrypt(plaintext)
    
    // 2.2 解密
    plaintext3, _ := tr.Decrypt(ciphertext2)
    fmt.Println(plaintext == plaintext3)

}
```

#### 国密签名、验签
算法的选择由配置文件决定，不需要在代码里指明：生产环境默认读取/tce/conf/config/tce.config.center/sdk.json配置文件，该文件在渲染时写入了必要的秘钥信息。
该版本支持的国密签名算法包括：kms-sign（基于kms-sm2）, tsm-sign（基于tsm-sm2）

##### 相关接口
|  接口名称   | 描述  |
|  ----  | ----  |
| Sign  | 签名 |
| Verify  | 验签 |

```go
type Signer interface {
	Sign(msg string) (string, error)            // 签名
	Verify(msg, signValue string) (bool, error) // 验证签名
}
```

##### Demo

```go
package main

import (
    "git.code.oa.com/tce-config/tcestuary-go/v4"
    "fmt"
)

func main() {
    // 1，安全存储
    signer, err := tcestuary.NewSigner()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    // 1.1 签名
    msg := "消息"
    sign, _ := signer.Sign(msg)

    // 1.2 验签
    ok, _ := signer.Verify(msg, sign)
    fmt.Println(ok)
}
```

#### 安全Hash

基于配置实现散列算法，支持SHA256和tsm-sm3

##### Demo
```go
package main

import (
    "encoding/hex"
"git.code.oa.com/tce-config/tcestuary-go/v4"
    "fmt"
)

func main() {
    // 1，安全存储
    hasher, err := tcestuary.NewTHasher()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    // 1.1 初始化
    thash, _ := hasher.New()
   
    // 1.2 更新
    _ = thash.Update([]byte("test msg"))
    _ = thash.Update([]byte("append msg"))
    // 1.3 散列
    ret, _ := thash.Digest()
    fmt.Println(hex.EncodeToString(ret))
}
```

#### 支撑组件认证配置

|  支持支撑组件类型   | 定义结构Struct  |
|  ----  | ----  |
| cmq	 | CMQConfig |
| csp  | CSPConfig|
| es  | ESConfig|
| hdfs	 | HdfsConfig |
| kafka  | KafkaConfig|
| mongodb  | MongodbConfig|
| mysql	 | MysqlConfig |
| redis  | RedisConfig|
| zk  | ZkConfig|

##### Demo

```go
package main

import (
   "fmt"
   "git.code.oa.com/tce-config/tcestuary-go/v4"
   "git.code.oa.com/tce-config/tcestuary-go/v4/middlewareconfig"
)

func main() { 
   // mysql
   mysqlConfig := &middlewareconfig.MysqlConfig{}
   mysqlConfig.GetConfig("dbsql_billing")
   fmt.Printf("%v", mysqlConfig)
   
   // redis
   redisConfig := &middlewareconfig.RedisConfig{}
   redisConfig.GetConfig("ckv_cas")
   fmt.Printf("%v", redisConfig)

   // kafka
   kafkaConfig := &middlewareconfig.KafkaConfig{}
   kafkaConfig.GetConfig("wtag")
   fmt.Printf("%v", kafkaConfig)

   // cmq
   cmqConfig := &middlewareconfig.CMQConfig{}
   cmqConfig.GetConfig("mq_waccount")
   fmt.Printf("%v", cmqConfig)

   // csp
   cspConfig := &middlewareconfig.CSPConfig{}
   cspConfig.GetConfig("yehe_file")
   fmt.Printf("%v", cspConfig)

   // mongodb
   mongodbConfig := &middlewareconfig.MongodbConfig{}
   mongodbConfig.GetConfig("bsp_document")
   fmt.Printf("%v", mongodbConfig)

   // es
   esConfig := &middlewareconfig.ESConfig{}
   esConfig.GetConfig("es_audit")
   fmt.Printf("%v", esConfig)
    
   // zk
   zkConfig := &middlewareconfig.ZKConfig{}
   zkConfig.GetConfig("yunapi3_zk")
   fmt.Printf("%v", zkConfig)

   // hdfs
   hdfsConfig := &middlewareconfig.HdfsConfig{}
   hdfsConfig.GetConfig("file_server")
   fmt.Printf("%v", hdfsConfig)
   
}
```