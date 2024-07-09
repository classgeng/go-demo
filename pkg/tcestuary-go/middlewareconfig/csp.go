package middlewareconfig

type CSPConfig struct {
	Host      string
	IP        string `json:"ip"`
	IPV4      string `json:"ipv4"`
	Port      int
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key" encrypted:"true"`
	UserName  string `json:"userName"`
	User      string `json:"user"`
}

const ConfigCSP = "csp"

func (m *CSPConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigCSP, name, m)
}
