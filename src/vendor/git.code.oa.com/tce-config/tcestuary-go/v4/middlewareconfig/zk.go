package middlewareconfig

type ZKConfig struct {
	Host        string
	IP          string `json:"ip"`
	IPV4        string `json:"ipv4"`
	Port        int
	URL         string `json:"url"`
	User        string `json:"user"`
	Password    string `encrypted:"true" json:"password"`
	AuthEnabled bool   `json:"auth_enabled"`
}

const ConfigZK = "zk"

func (m *ZKConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigZK, name, m)
}
