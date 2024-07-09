package tcesecurity

// 签名、验签

// Sign配置参数
type SignOpts struct {
	Method string
	// For kms-sign
	KeyId     string
	SecretId  string
	SecretKey string
	KMSServer string
	// For tsm
	PublicKey  string
	PrivateKey string
}

// Signer签名、验签
type Signer interface {
	Sign(string) (string, error)
	Verify(string, string) (bool, error)
}

// 签名算法
type SignFunc func(opts SignOpts) (Signer, error)

// 支持的签名算法
var SupportSignFunc = make(map[string]SignFunc, 0)

// 注册支持的签名算法
func registerSignFunc(k string, f SignFunc) {
	SupportSignFunc[k] = f
}
