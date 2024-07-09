package tcesecurity

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

//兼容历史版本

var ErrorFormat error = errors.New("format error")

const (
	Aes256CbcAlgorithm = "aes-256-cbc"
)

func init() {
	f := func(opts CryptoOpts) (Crypto, error) {
		return NewAesCbcCrypto(opts.Method, []byte(opts.AesKey))
	}
	registerCryptoFunc(Aes256CbcAlgorithm, f)
}

// AES-CBC算法
type AesCbcCrypto struct {
	Method string
	aesKey []byte
	Prefix string

	format   string
	size     int
	salt     string
	saltSize int // 随机盐的长度
	version  int
}

// NewAesCbcCrypto return AES-CBC Crypto
func NewAesCbcCrypto(method string, aesKey []byte) (*AesCbcCrypto, error) {
	s := &AesCbcCrypto{
		Prefix:   "AES+",
		format:   "AES+V%d+%s",
		size:     2,
		saltSize: rand.Intn(8) + 1,
		version:  1,
		Method:   hex.EncodeToString([]byte(method)),
		aesKey:   []byte(aesKey),
	}
	s.RandomSalt(s.saltSize)
	return s, nil
}

type encoder struct {
	prefix   string
	format   string
	size     int
	salt     string
	saltSize int // 随机盐的长度
	version  int
}

// RandomSalt 随机盐生成器. string 返回值用于构建测试用例
func (s *AesCbcCrypto) RandomSalt(n int) string {
	letters := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := uint32(len(letters))

	var bb bytes.Buffer
	bb.Grow(n)
	for i := 0; i < n; i++ {
		bb.WriteByte(letters[rand.Uint32()%length])
	}
	s.salt = bb.String()

	return s.salt
}

// 兼容 未填充
func (s *AesCbcCrypto) Encrypt(string) (string, error) {
	return "", nil
}

// Decrypt 提供给业务代码调用, 从配置密文中获取明文配置. 另外, 构建解密工具
func (s *AesCbcCrypto) Decrypt(crypted string) (string, error) {
	// 产品逻辑上需要兼容 明文密码 和 AES密码
	// 如果不是 AES 格式密文, 则认为是明文密码
	if !s.WithPrefix(crypted) {
		return crypted, nil
	}

	crypted, err := s.Unwrap(crypted)
	if err != nil {
		return "", err
	}

	bytes, err := hex.DecodeString(crypted)
	if err != nil {
		return "", err
	}

	originWithSalt, err := _AESDecrypt(s.aesKey, bytes)
	if err != nil {
		return "", fmt.Errorf("decode. %s", err)
	}

	origin, err := s.Unsalt(originWithSalt)
	if err != nil {
		return "", fmt.Errorf("decode. %s", err)
	}

	return origin, err
}

func (s *AesCbcCrypto) WithPrefix(str string) bool {
	return strings.HasPrefix(str, s.Prefix)
}

// 密文格式的生成/解析工具
func (s *AesCbcCrypto) Wrap(str string) string {
	return fmt.Sprintf(s.format, s.version, str)
}
func (s *AesCbcCrypto) Unwrap(str string) (string, error) {
	var crypted string

	n, err := fmt.Sscanf(str, s.format, &s.version, &crypted)
	if err != nil {
		return "", err
	}
	if n != s.size {
		return "", ErrorFormat
	}
	return crypted, nil
}

func _AESDecrypt(key, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为 128 位，所以 blockSize=16 字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = _PKCS5UnPadding(origData, blockSize)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

func _PKCS5UnPadding(origData []byte, blockSize int) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	// padding 的取值范围: [1, blockSize]
	// fix: aeskey 和 密文 不匹配时, 潜在的 slice 操作越界
	if unpadding > blockSize || unpadding < 1 {
		return nil, errors.New("aes unpadding error. aeskey and ciphertext may not match")
	}
	return origData[:(length - unpadding)], nil
}

func (s *AesCbcCrypto) Unsalt(source []byte) (string, error) {
	buff := bytes.NewBuffer(source)

	// 取盐长度
	c, err := buff.ReadByte()
	if err != nil {
		return "", err
	}
	s.saltSize = int(c)

	// 取盐
	salt := make([]byte, s.saltSize)
	n, err := buff.Read(salt)
	if err != nil {
		return "", err
	}
	if n != s.saltSize {
		return "", fmt.Errorf("read salt error")
	}
	s.salt = string(salt)

	// 原始明文
	return string(buff.Bytes()), nil

}
