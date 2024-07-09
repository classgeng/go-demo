package tcestuary

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func _PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
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

func _AESEncrypt(key, origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//AES分组长度为 128 位，所以 blockSize=16 字节
	blockSize := block.BlockSize()
	origData = _PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
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
