package tcesecurity

// 加密、解密

// Crypto配置参数
type CryptoOpts struct {
	Method string

	// For AES
	AesKey string

	// For tsm-sm4
	Sm4Key string
	// For RSA, tsm-sm2
	PrivateKey string
	PublicKey  string
	// For kms-sm2, kms-4
	KeyId     string
	SecretId  string
	SecretKey string
	KMSServer string
}

// Crypto 加密、解密接口
type Crypto interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// 加密解密算法
type AlgorithmFunc func(CryptoOpts) (Crypto, error)

// SupportAlgorithm 支持的加密算法
var SupportAlgorithm = make(map[string]AlgorithmFunc, 0)

// 注册支持加密算法
func registerCryptoFunc(k string, f AlgorithmFunc) {
	SupportAlgorithm[k] = f
}
