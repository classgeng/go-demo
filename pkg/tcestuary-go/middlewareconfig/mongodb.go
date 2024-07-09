package middlewareconfig

type MongodbConfig struct {
	Host      string
	IP        string `json:"ip"`
	IPV4      string `json:"ipv4"`
	Port      int
	URL       string `json:"url"`
	User      string `json:"user"`
	Password  string `encrypted:"true" json:"password"`
	AdminUser string `json:"admin_user"`
	AdminPass string `encrypted:"true" json:"admin_pass"`
	Role      string
	DBName    string `json:"db_name"`
}

const ConfigMongodb = "mongodb"

func (m *MongodbConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigMongodb, name, m)
}
