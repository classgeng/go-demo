package middlewareconfig

type RedisConfig struct {
	Host string
	IP   string `json:"ip"`
	IPV4 string `json:"ipv4"`
	Port int
	User string `json:"user"`
	Password string `encrypted:"true" json:"password"`
	Pass string `encrypted:"true" json:"pass"`
}

const ConfigRedis = "redis"

func (m *RedisConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigRedis, name, m)
}
