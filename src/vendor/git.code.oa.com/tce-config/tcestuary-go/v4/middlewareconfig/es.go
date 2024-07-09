package middlewareconfig

type ESConfig struct {
	Host        string
	IPV4        string `json:"ipv4"`
	IP          string `json:"ip"`
	Port        int
	AdminUser   string `json:"admin_user"`
	AdminPass   string `encrypted:"true" json:"admin_pass"`
	ClusterType string `json:"cluster_type"`
	User        string `json:"user"`
	Password    string `encrypted:"true" json:"password"`
}

const ConfigES = "es"

func (m *ESConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigES, name, m)
}
