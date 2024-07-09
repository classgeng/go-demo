package middlewareconfig

type CMQConfig struct {
	Host     string
	IPV4     string `json:"ipv4"`
	IP       string `json:"ip"`
	Port     int
	User     string `json:"user"`
	Pass     string `encrypted:"true" json:"pass"`
	Username string `json:"userName"`
	Password string `encrypted:"true" json:"password"`
}

const ConfigCMQ = "cmq"

func (m *CMQConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigCMQ, name, m)
}
