package tcestuary

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	// ErrorFormat 格式错误
	ErrorFormat error = errors.New("format error")
)

type encoder struct {
	prefix   string
	format   string
	size     int
	salt     string
	saltSize int // 随机盐的长度
	version  int
}

func newEncoder() *encoder {
	rand.Seed(time.Now().UnixNano())

	e := &encoder{
		prefix:   "AES+",
		format:   "AES+V%d+%s",
		size:     2,
		saltSize: rand.Intn(8) + 1,
		version:  1,
	}
	e.RandomSalt(e.saltSize)

	return e
}

// 检查是否可能是 AES 加密格式
func (s *encoder) WithPrefix(str string) bool {
	return strings.HasPrefix(str, s.prefix)
}

// 密文格式的生成/解析工具
func (s *encoder) Wrap(str string) string {
	return fmt.Sprintf(s.format, s.version, str)
}
func (s *encoder) Unwrap(str string) (string, error) {
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

// 加盐格式: [盐长度(1Byte)][盐][明文]
func (s *encoder) Salt(str string) ([]byte, error) {
	buff := bytes.NewBufferString("")

	buff.Grow(s.saltSize + len(s.salt) + len(str))

	buff.WriteByte(byte(s.saltSize))
	buff.WriteString(s.salt)
	buff.WriteString(str)

	return buff.Bytes(), nil
}
func (s *encoder) Unsalt(source []byte) (string, error) {
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

// RandomSalt 随机盐生成器. string 返回值用于构建测试用例
func (s *encoder) RandomSalt(n int) string {
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

// Encrypt 构建加密工具
func Encrypt(key string, origin string) (string, error) {
	coder := newEncoder()

	salted, _ := coder.Salt(origin)

	crypted, err := _AESEncrypt([]byte(key), salted)
	if err != nil {
		return "", err
	}

	return coder.Wrap(hex.EncodeToString(crypted)), nil
}

// Decrypt 提供给业务代码调用, 从配置密文中获取明文配置. 另外, 构建解密工具
func Decrypt(key string, crypted string) (string, error) {
	coder := newEncoder()

	// 产品逻辑上需要兼容 明文密码 和 AES密码
	// 如果不是 AES 格式密文, 则认为是明文密码
	if !coder.WithPrefix(crypted) {
		return crypted, nil
	}

	crypted, err := coder.Unwrap(crypted)
	if err != nil {
		return "", err
	}

	bytes, err := hex.DecodeString(crypted)
	if err != nil {
		return "", err
	}

	originWithSalt, err := _AESDecrypt([]byte(key), bytes)
	if err != nil {
		return "", fmt.Errorf("decode. %s", err)
	}

	origin, err := coder.Unsalt(originWithSalt)
	if err != nil {
		return "", fmt.Errorf("decode. %s", err)
	}

	return origin, err
}
