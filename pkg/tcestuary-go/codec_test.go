package tcestuary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -test.count 100000 -test.parallel 100 --run TestAESFunc

// AES密码: 加解密函数测试
func TestAESFunc(t *testing.T) {
	aeskey := newEncoder().RandomSalt(32)
	origin := newEncoder().RandomSalt(16)

	pass, err := Encrypt(aeskey, origin)
	assert.NoError(t, err)
	t.Logf("encrypted: %s\n", pass)

	tpass, err := Decrypt(aeskey, pass)
	assert.NoError(t, err)
	t.Logf("decrypted: %s\n", tpass)

	assert.Equal(t, origin, tpass)
}

// 产品需求: 配置文件同时支持明文密码和AES密码
// 如果是明文密码, 则原样返回
func TestSourceCipher(t *testing.T) {
	aeskey := newEncoder().RandomSalt(32)
	password := newEncoder().RandomSalt(16)

	tpass, err := Decrypt(aeskey, password)
	assert.NoError(t, err)

	assert.Equal(t, password, tpass)
}
