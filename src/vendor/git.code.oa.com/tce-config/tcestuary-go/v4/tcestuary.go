package tcestuary

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.code.oa.com/tce-config/tcestuary-go/v4/configcenter"
	"git.code.oa.com/tce-config/tcestuary-go/v4/logger"
	"github.com/jinzhu/copier"
)

var std *manager = newManager()

var (
	// ErrNotFound 配置项不存在
	ErrNotFound = errors.New("not found")
	// ErrConfigInValid 配置内容不合法
	ErrConfigInValid = errors.New("not valid")
	// ErrDecryptFail 无法解密密码字段
	ErrDecryptFail = errors.New("decrypt fail")
	// ErrCopyFail 无法创建对象副本
	ErrCopyFail = errors.New("copy fail")
	// ErrKeyFormat 配置项的key格式不满足要求
	ErrKeyFormat = errors.New("key format error")
	// ErrUsageInvalid 接口调用不满足条件
	ErrUsageInvalid = errors.New("usage invalid")
)

type (
	// Mysql 数据库资源描述文件
	Mysql struct {
		Host     string
		IP       string
		Port     int
		User     string
		Password string
		Database string
	}

	// MysqlWithRegion 相比 Mysql 增加 ReginID 字段
	MysqlWithRegion struct {
		Mysql
		RegionID   int
		RegionName string
	}

	// MysqlWithZone 相比 Mysql 增加 ReginID / ZoneID 字段
	MysqlWithZone struct {
		Mysql
		RegionID int
		ZoneID   int
	}
)

// SetConfigDirectory 调试阶段和特殊场景时, 临时修改配置路径. 业务代码中请勿使用.
// 配置文件路径变更, 会影响 SDK 版本 文件的输出. 原则上: SDK 版本文件与配置文件在相同路径下
func SetConfigDirectory(dir string) error {
	// 检查路径是否存在
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("dir not exist, %s", dir)
		}
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("need directory, %s", dir)
	}

	// dir 换成绝对路径, 用于调试信息输出
	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}

	// 检测关键配置, 是否存在
	exist := func(file string) error {
		_, err := os.Stat(file)
		return err
	}
	sdkFile := filepath.Join(dir, "sdk.json")
	if err := exist(sdkFile); err != nil {
		return err
	}

	// 参数验证通过后, 一次性赋值. 防止局部赋值
	std.Directory = dir
	std.ConfigCenterFile = sdkFile

	return nil
}

// GetConfigDirectory 获取config 目录
func GetConfigDirectory() string {
	return std.ConfigCenterFile
}

// GetMysqlConfig 包装配置文件读取 和 解密动作, 向业务提供密码明文
// 使用条件:
// 1. dbsql.scope 声明为 Global / Region / Zone
// 错误码:
// 1. 配置项不存在, 返回 ErrNotFound
// 2. 配置字段不合法, 返回 ErrConfigInValid
// 3. Scope不支持, 返回 ErrUsageInvalid
// 输入参数:
// key: 格式满足 ocloud_api3.hp_api_formal.
//      其中, ocloud_api3 对应 dbsql 实例, hp_api_formal 对应数据库表示标识;
//
// 备注:
// 一期实现: 兼容 password 密文 和 明文. 支持客户升级 和 向历史版本合并代码
// 二期实现: 增加网络请求, 从密码库拉取配置
// 默认配置优先级: 1.密码库; 2.本地密文密码; 3.本地明文密码;
func GetMysqlConfig(key string) (*Mysql, error) {
	// Load 内部逻辑保证仅加载一次配置
	err := std.Load()
	if err != nil {
		return nil, err
	}

	// 解析输入参数, 获取 dbsql 实例 和 数据库名称
	s := strings.Split(key, ".")
	if len(s) != 2 {
		return nil, ErrKeyFormat
	}
	dbsql, database := s[0], s[1]

	// 检查资源等级是否匹配
	scope := std.ConfigCenter.FindMysqlScope(dbsql)
	if scope == configcenter.ScopeUnknown {
		return nil, ErrNotFound
	} else if scope != configcenter.ScopeFlat {
		return nil, ErrUsageInvalid
	}

	// 返回指向全局配置项的指针
	m := std.ConfigCenter.FindMysql(dbsql, database)
	if m == nil {
		return nil, ErrNotFound
	}

	// 配置项检查
	if err := m.Valid(); err != nil {
		return nil, ErrConfigInValid
	}

	// 创建临时对象, 防止暴露全局配置项指针.
	mysql := new(Mysql)
	if err := copier.Copy(mysql, m); err != nil {
		return nil, ErrConfigInValid
	}
	mysql.Database = database // 填充数据库名称
	passwd, err := Decrypt(std.ConfigCenter.SDK.PasswdSecret.V1Aeskey, mysql.Password)
	if err != nil {
		return nil, ErrDecryptFail
	}
	mysql.Password = passwd

	return mysql, nil
}

// SetLogger 允许业务指定日志输出, 用于问题调试;
// 默认: 不输出任何日志
func SetLogger(log logger.Logger) {
	logger.SetLogger(log)
}

// Debug 向终端输出配置信息, 用于异常问题排查
func Debug() {
	std.Debug()
}

// GetMysqlConfigAllRegion 获取数据库 region 实例列表.
// 使用条件:
// 1. dbsql 组件声明为 region 级别;
// 2. cc.declear.json 中声明为 all_region 级别引用;
// 条件不满足的情况下, 调用接口返回: ErrUsageInvalid
//
// key 规则: 参考 GetMysqlConfig 说明
func GetMysqlConfigAllRegion(key string) ([]*MysqlWithRegion, error) {
	// Load 内部逻辑保证仅加载一次配置
	err := std.Load()
	if err != nil {
		return nil, err
	}

	// 解析输入参数, 获取 dbsql 实例 和 数据库名称
	s := strings.Split(key, ".")
	if len(s) != 2 {
		return nil, ErrKeyFormat
	}
	dbsql, database := s[0], s[1]

	// 检查资源等级是否匹配
	scope := std.ConfigCenter.FindMysqlScope(dbsql)
	if scope == configcenter.ScopeUnknown {
		return nil, ErrNotFound
	} else if scope != configcenter.ScopeAllRegion {
		return nil, ErrUsageInvalid
	}

	// 查找配置
	mysqlR := make([]*MysqlWithRegion, 0)
	mysqls := std.ConfigCenter.FindMysqlAllRegion(dbsql, database)
	for _, m := range mysqls {
		r := new(MysqlWithRegion)

		if err := copier.Copy(r, m.Service); err != nil {
			return nil, ErrConfigInValid
		}
		passwd, err := Decrypt(std.ConfigCenter.SDK.PasswdSecret.V1Aeskey, r.Password)
		if err != nil {
			return nil, ErrDecryptFail
		}

		r.Database = database // 填充数据库名称
		r.Password = passwd
		r.RegionID = m.Base.RegionID
		r.RegionName = m.Base.RegionName

		mysqlR = append(mysqlR, r)
	}

	return mysqlR, nil
}

// GetMysqlConfigAllZone 获取数据库 zone 实例列表.
// 使用条件:
// 1. dbsql 组件声明为 zone 级别;
// 2. cc.declear.json 中声明为 all_zone 级别引用;
// 条件不满足的情况下, 调用接口返回: ErrUsageInvalid
//
// key 规则: 参考 GetMysqlConfig 说明
func GetMysqlConfigAllZone(key string) ([]*MysqlWithZone, error) {
	// Load 内部逻辑保证仅加载一次配置
	err := std.Load()
	if err != nil {
		return nil, err
	}

	// 解析输入参数, 获取 dbsql 实例 和 数据库名称
	s := strings.Split(key, ".")
	if len(s) != 2 {
		return nil, ErrKeyFormat
	}
	dbsql, database := s[0], s[1]

	// 检查资源等级是否匹配
	scope := std.ConfigCenter.FindMysqlScope(dbsql)
	if scope == configcenter.ScopeUnknown {
		return nil, ErrNotFound
	} else if scope != configcenter.ScopeAllZone {
		return nil, ErrUsageInvalid
	}

	// 查找配置
	mysqlZ := make([]*MysqlWithZone, 0)
	mysqls := std.ConfigCenter.FindMysqlAllZone(dbsql, database)
	for _, m := range mysqls {
		z := new(MysqlWithZone)

		if err := copier.Copy(z, m.Service); err != nil {
			return nil, ErrConfigInValid
		}
		passwd, err := Decrypt(std.ConfigCenter.SDK.PasswdSecret.V1Aeskey, z.Password)
		if err != nil {
			return nil, ErrDecryptFail
		}

		z.Database = database // 填充数据库名称
		z.Password = passwd
		z.RegionID = m.Base.RegionID
		z.ZoneID = m.Base.ZoneID

		mysqlZ = append(mysqlZ, z)
	}

	return mysqlZ, nil
}

func GetMainRegionName() (string, error) {
	if len(std.ConfigCenter.Base.ScopeExtInfo.MainRegionName) == 0 {
		err := std.Load()
		if err != nil {
			return "", err
		}
	}
	return std.ConfigCenter.Base.ScopeExtInfo.MainRegionName, nil
}

// Region 地域信息
type Region struct {
	RegionID        int    `json:"region_id"`
	RegionArea      string `json:"region_area"`
	RegionCity      string `json:"region_city"`
	RegionName      string `json:"region_name"`
	RegionNameLong  string `json:"region_name_long"`
	RegionNameUpper string `json:"region_name_upper"`
	RegionNameZh    string `json:"region_name_zh"`
	MainZoneName    string `json:"main_zone_name"`
}

// GetRegion 当 regionID 不满足业务需求, 调用接口获取完整描述信息
func GetRegion(regionID int) (*Region, error) {
	// Load 内部逻辑保证仅加载一次配置
	err := std.Load()
	if err != nil {
		return nil, err
	}

	r := std.ConfigCenter.FindRegion(regionID)
	if r == nil {
		return nil, ErrNotFound
	}

	region := new(Region)
	copier.Copy(region, r)

	return region, nil
}

// Zone 可用区信息
type Zone struct {
	RegionID   int    `json:"region_id"`
	ZoneID     int    `json:"zone_id"`
	ZoneName   string `json:"zone_name"`
	ZoneNameZh string `json:"zone_name_zh"`
}

// GetZone 当 ZoneID 不满足业务需求, 调用接口获取完整描述信息
func GetZone(regionID int, zoneID int) (*Zone, error) {
	// Load 内部逻辑保证仅加载一次配置
	err := std.Load()
	if err != nil {
		return nil, err
	}

	z := std.ConfigCenter.FindZone(regionID, zoneID)
	if z == nil {
		return nil, ErrNotFound
	}

	zone := new(Zone)
	copier.Copy(zone, z)

	return zone, nil
}

// GetConfigCenterPtr 支持 sdk.json 加密工具, 请勿调用
func GetConfigCenterPtr() *configcenter.ConfigCenter {
	return std.ConfigCenter
}
