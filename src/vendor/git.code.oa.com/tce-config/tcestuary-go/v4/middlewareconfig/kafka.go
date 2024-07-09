package middlewareconfig

type KafkaConfig struct {
	Host     string
	IPV4     string `json:"ipv4"`
	IP       string `json:"ip"`
	Port     int
	Username string `json:"username"`
	User     string `json:"user"`
	Password string `encrypted:"true" json:"password"`
}

const ConfigKafka = "kafka"

func (m *KafkaConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigKafka, name, m)
}
