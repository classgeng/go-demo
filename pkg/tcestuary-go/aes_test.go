package tcestuary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 密钥/密文不配对时, 代码不能 panic
// go test -test.count 100000 -test.parallel 100 --run TestAESNotPanic
func TestAESNotPanic(t *testing.T) {
	// 加密
	aeskey := newEncoder().RandomSalt(32)
	plaintext := newEncoder().RandomSalt(32)

	ciphertext, err := _AESEncrypt([]byte(aeskey), []byte(plaintext))
	assert.NoError(t, err)

	// 用错误密钥解密, 代码不能 panic
	fakekey := newEncoder().RandomSalt(32)
	assert.NotPanics(t, func() { _AESDecrypt([]byte(fakekey), ciphertext) })
}
