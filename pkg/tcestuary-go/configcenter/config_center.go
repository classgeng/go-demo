package configcenter

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/jinzhu/copier"
)

// 错误码定义
var (
	ErrNotFound = errors.New("not found")
)

// Scope TCE 定义的资源级别. 与部署架构有关
type Scope int

// 资源级别
const (
	ScopeGlobal Scope = 1 << iota
	ScopeRegion
	ScopeZone
	ScopeGaia
	ScopeAllRegion
	ScopeAllZone
	ScopeAllGaia
	ScopeUnknown // 未知资源级别

	ScopeFlat = ScopeGlobal | ScopeRegion | ScopeZone | ScopeGaia // 表示扁平资源描述结构
)

// 配置中心基础数据结构
type (
	// PasswdSecret 客户纬度私密信息
	PasswdSecret struct {
		V1Aeskey string `json:"aeskey"`
	}

	// SecretConfig: Secret config for transport and storage
	SecretConfig struct {
		Method     string `json:"method,omitempty"`
		PublicKey  string `json:"public_key,omitempty"`
		PrivateKey string `json:"private_key,omitempty"`
		AesKey     string `json:"aes_key,omitempty"`
		V1Aeskey   string `json:"aeskey,omitempty"`
		Sm4Key     string `json:"sm4_key,omitempty"`
		KeyId      string `json:"key_id,omitempty"`
		KMSServer  string `json:"kms_server,omitempty"`
		SecretId   string `json:"secret_id,omitempty"`
		SecretKey  string `json:"secret_key,omitempty"`
	}

	// TSMConfig
	TSMConfig struct {
		PemAppid           string `json:"pem_appid"`
		PemInitLibWithCert string `json:"pem_initlibwithcert"`
	}

	// HashConfig
	HashConfig struct {
		Method string `json:"method"`
	}

	// Scope 的附属信息
	ScopeExtInfo struct {
		MainRegionName string `json:"main_region_name"`
	}

	// Region 地域属性信息
	Region struct {
		RegionID        int    `json:"region_id"`
		RegionArea      string `json:"region_area"`
		RegionCity      string `json:"region_city"`
		RegionName      string `json:"region_name"`
		RegionNameLong  string `json:"region_name_long"`
		RegionNameUpper string `json:"region_name_upper"`
		RegionNameZh    string `json:"region_name_zh"`
		MainZoneName    string `json:"main_zone_name"`
	}

	// Zone 可用区属性信息
	Zone struct {
		ZoneNameZh string `json:"zone_name_zh"`
		RegionID   int    `json:"region_id"`
		ZoneName   string `json:"zone_name"`
		ZoneID     int    `json:"zone_id"`
	}

	// Mysql 映射基础支撑 dbsql 的资源描述
	// 参考文档: http://tapd.oa.com/OneBank/markdown_wikis/#1010130691010773371@toc15
	Mysql struct {
		Host      string   `json:"host"`
		IP        string   `json:"ipv4"`
		Port      int      `json:"port"`
		User      string   `json:"user"`
		Password  string   `json:"pass"`
		Databases []string `json:"db_name_list"`
	}
)

type (
	// OriginConfigCenter 原始配置结构.
	// 配置中心早期设计时, 没有考虑 SDK 实现场景,
	// 相同的key下, 根据场景不通, 写入不通的 JSON 对象
	// SDK 实现过程中注意规避
	OriginConfigCenter struct {
		Base struct {
			Local   json.RawMessage `json:"local"`
			Regions []Region        `json:"region_list"`
			Zones   []Zone          `json:"zone_list"`
		} `json:"base"`
		SDK struct {
			PasswdSecret    SecretConfig `json:"passwd-secret"`
			StorageSecret   SecretConfig `json:"storage-secret"`
			TransportSecret SecretConfig `json:"transport-secret"`
			SignSecret      SecretConfig `json:"sign-secret"`
			TSMSecret       TSMConfig    `json:"tsm"`
			HashSecret      HashConfig   `json:"hash-secret"`
		} `json:"sdk"`
		Mysqls map[string]json.RawMessage `json:"mysql"`
	}

	// ConfigCenter sdk.json 结构化表示
	ConfigCenter struct {
		Base struct {
			Region       Region
			Zone         Zone
			Regions      []Region
			Zones        []Zone
			ScopeExtInfo ScopeExtInfo
		}
		SDK struct {
			PasswdSecret    SecretConfig
			StorageSecret   SecretConfig
			TransportSecret SecretConfig
			SignSecret      SecretConfig
			TSMSecret       TSMConfig
			HashSecret      HashConfig
		}
		Mysqls map[string]*MysqlWrapper
	}

	// MysqlWrapper 当 scope = ALL_REGION / ALL_ZONE 时, 资源描述结构是 JSON 数组
	// 其它 scope 时, 是 JSON 对象
	// Scope 可能取值: ScopeAllRegion / ScopeAllZone / ScopeFlat
	MysqlWrapper struct {
		Scope  Scope
		Object interface{}
	}

	// MysqlWithRegion dbsql.scope 属性为 all_region, 资源描述结构中有 Region 信息
	MysqlWithRegion struct {
		Base    Region `json:"_base"`
		Service Mysql  `json:"_service"`
	}

	// MysqlWithZone dbsql.scope 属性为 all_zone, 资源描述结构中有 Zone 信息
	MysqlWithZone struct {
		Base    Zone  `json:"_base"`
		Service Mysql `json:"_service"`
	}
)

// NewConfigCenter 初始化结构内部资源
func NewConfigCenter() *ConfigCenter {
	return &ConfigCenter{
		Mysqls: make(map[string]*MysqlWrapper, 0),
	}
}

// Parse 解析配置文件
func (c *ConfigCenter) Parse(data []byte) error {

	sourceCC := new(OriginConfigCenter)
	if err := json.Unmarshal(data, sourceCC); err != nil {
		return err
	}

	// Base 解析
	if err := copier.Copy(&c.Base, &sourceCC.Base); err != nil {
		return err
	}
	if sourceCC.Base.Local != nil { // 兼容 conf.base 未声明的场景
		if region, zone, scopeExtInfo, err := sourceCC.parseBaseLocal(sourceCC.Base.Local); err == nil {
			c.Base.Region = *region
			c.Base.Zone = *zone
			c.Base.ScopeExtInfo = *scopeExtInfo
		} else {
			return err
		}
	}

	// SDK 复制
	if err := copier.Copy(&c.SDK, &sourceCC.SDK); err != nil {
		return err
	}

	// Mysql 解析
	for dbsql, message := range sourceCC.Mysqls {
		if w, err := sourceCC.parseMysqls(message); err == nil {
			c.Mysqls[dbsql] = w
			//} else {
			//	return err
		}
	}

	return nil
}

// Local 字段表示: 当前地域/可用区/Gaia. 配置中将信息揉合在一个 JSON 对象中传递, 逻辑上拆解开, 方便后续使用
func (c *OriginConfigCenter) parseBaseLocal(message json.RawMessage) (*Region, *Zone, *ScopeExtInfo, error) {
	region, zone, scopeExtInfo := new(Region), new(Zone), new(ScopeExtInfo)

	if err := json.Unmarshal(message, region); err != nil {
		return nil, nil, nil, err
	}

	if err := json.Unmarshal(message, zone); err != nil {
		return nil, nil, nil, err
	}

	if err := json.Unmarshal(message, scopeExtInfo); err != nil {
		return nil, nil, nil, err
	}

	return region, zone, scopeExtInfo, nil
}

// 不同的 Scope 对应的资源描述结构不同, 此处根据描述结构类型, 将配置转化成格式化内存结构, 方便后续使用
func (c *OriginConfigCenter) parseMysqls(message json.RawMessage) (*MysqlWrapper, error) {
	w := new(MysqlWrapper)

	// 判断是否为扁平资源描述
	mysql := &Mysql{}
	if err := json.Unmarshal(message, mysql); err == nil {
		if mysqlErr := mysql.Valid(); mysqlErr == nil {
			w.Scope = ScopeFlat
			w.Object = mysql
			return w, nil
		}
	}

	// 判断是否为 ALL_REGION 数组类型
	mysqlR := make([]*MysqlWithRegion, 0)
	if err := json.Unmarshal(message, &mysqlR); err == nil && len(mysqlR) > 0 {
		mr := mysqlR[0]
		serviceErr, baseErr := mr.Service.Valid(), mr.Base.Valid()
		if serviceErr == nil && baseErr == nil {
			w.Scope = ScopeAllRegion
			w.Object = mysqlR
			return w, nil
		}
	}

	// 判断是否为 ALL_ZONE 数组类型
	mysqlZ := make([]*MysqlWithZone, 0)
	if err := json.Unmarshal(message, &mysqlZ); err == nil && len(mysqlZ) > 0 {
		mz := mysqlZ[0]
		serviceErr, baseErr := mz.Service.Valid(), mz.Base.Valid()
		if serviceErr == nil && baseErr == nil {
			w.Scope = ScopeAllZone
			w.Object = mysqlZ
			return w, nil
		}
	}

	return nil, errors.New("unknown msyql config type")
}

// FindMysqlScope 判断资源级别
// Mysql 资源级别可能是: ScopeAllRegion / ScopeAllZone / ScopeFlat
//
// 如果, 配置文件中不存在 dbsql, 返回 ErrNotFound
func (c *ConfigCenter) FindMysqlScope(dbsql string) Scope {
	if m, ok := c.Mysqls[dbsql]; ok {
		return m.Scope
	}
	return ScopeUnknown
}

// FindMysql 判断 database 是否属于 dbsql 实例
// nil 表示: database 配置中不存在
func (c *ConfigCenter) FindMysql(dbsql string, database string) *Mysql {
	if mysql, ok := c.Mysqls[dbsql]; ok {
		for _, dbname := range mysql.Object.(*Mysql).Databases {
			if dbname == database {
				return mysql.Object.(*Mysql)
			}
		}
	}
	return nil
}

// FindMysqlAllRegion 判断 database 是否属于 dbsql 实例
// nil 表示: database 配置中不存在
func (c *ConfigCenter) FindMysqlAllRegion(dbsql string, database string) []*MysqlWithRegion {
	if w, ok := c.Mysqls[dbsql]; ok {
		return w.Object.([]*MysqlWithRegion)
	}
	return nil
}

// FindMysqlAllZone 判断 database 是否属于 dbsql 实例
// nil 表示: database 配置中不存在
func (c *ConfigCenter) FindMysqlAllZone(dbsql string, database string) []*MysqlWithZone {
	if w, ok := c.Mysqls[dbsql]; ok {
		return w.Object.([]*MysqlWithZone)
	}
	return nil
}

// FindRegion 查找 Region 属性
func (c *ConfigCenter) FindRegion(regionID int) *Region {
	for i, region := range c.Base.Regions {
		if region.RegionID == regionID {
			return &c.Base.Regions[i]
		}
	}
	return nil
}

// FindZone 查找 Zone 属性
func (c *ConfigCenter) FindZone(regionID int, zoneID int) *Zone {
	for i, zone := range c.Base.Zones {
		if zone.RegionID == regionID && zone.ZoneID == zoneID {
			return &c.Base.Zones[i]
		}
	}
	return nil
}

// Debug 输出已加载配置内容, 支持异常调试
func (c *ConfigCenter) Debug() {
	buff := []byte{}

	buff, _ = json.MarshalIndent(c.Base, "", "  ")
	log.Printf("base: %s\n", string(buff))

	buff, _ = json.MarshalIndent(c.SDK, "", "  ")
	log.Printf("sdk: %s\n", string(buff))

	for dbsql, mysql := range c.Mysqls {
		buff, _ = json.MarshalIndent(*mysql, "", "  ")
		log.Printf("mysql %s %+v\n", dbsql, string(buff))
	}
}

// Valid 检查配置项
func (region *Region) Valid() error {
	if region.RegionID <= 0 {
		return errors.New("region id error")
	}
	if region.RegionName == "" {
		return errors.New("region name empty")
	}
	return nil
}

// Valid 检查配置项
func (zone *Zone) Valid() error {
	if zone.RegionID <= 0 {
		return errors.New("region id error")
	}
	if zone.ZoneID <= 0 {
		return errors.New("zone id error")
	}
	if zone.ZoneName == "" {
		return errors.New("zone name empty")
	}
	return nil
}

// Valid 配置项检查
func (mysql *Mysql) Valid() error {
	if mysql.Host == "" {
		return errors.New("host is empty")
	}
	if mysql.IP == "" {
		return errors.New("ip is empty")
	}
	if mysql.Port < 1 || mysql.Port > 65535 {
		return errors.New("port not valid")
	}
	if mysql.User == "" {
		return errors.New("user is empty")
	}
	if mysql.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}
