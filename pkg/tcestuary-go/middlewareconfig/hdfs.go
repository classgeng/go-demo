package middlewareconfig

type HdfsConfig struct {
	Host       string
	IP         string `json:"ip"`
	IPV4       string `json:"ipv4"`
	Port       int
	Principal  string `encrypted:"true" json:"principal"`
	KeytabFile string `json:"keytab_file"`
}

const ConfigHdfs = "hdfs"

func (m *HdfsConfig) GetConfig(name string) error {
	return unmarshallConfig(ConfigHdfs, name, m)
}
