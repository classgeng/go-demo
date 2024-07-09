package middlewareconfig

type MysqlConfig struct {
	Host       string
	Port       int
	User       string
	Pass       string   `encrypted:"true"`
	DBNameList []string `json:"db_name_list"`
	IPV4       string   `json:"ipv4"`
}

const ConfigMysql = "mysql"

func (m *MysqlConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigMysql, name, m)
}
