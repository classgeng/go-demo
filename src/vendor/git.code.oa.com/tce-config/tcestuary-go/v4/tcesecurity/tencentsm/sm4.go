package sm

/*
#cgo CFLAGS: -g -O2 -I${SRCDIR}/include/
#cgo windows,arm64 LDFLAGS: ${SRCDIR}/lib/windows_arm/libTencentSM.dll ${SRCDIR}/lib/windows_arm/libgmp.a
#cgo windows,!arm64 LDFLAGS: ${SRCDIR}/lib/windows/libTencentSM.dll ${SRCDIR}/lib/windows/libgmp.a
#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/darwin_arm/libTencentSM.a ${SRCDIR}/lib/darwin_arm/libgmp.a
#cgo darwin,!arm64 LDFLAGS: ${SRCDIR}/lib/darwin/libTencentSM.a ${SRCDIR}/lib/darwin/libgmp.a
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/linux_arm/libTencentSM.a ${SRCDIR}/lib/linux_arm/libgmp.a
#cgo linux,!arm64 LDFLAGS: ${SRCDIR}/lib/linux/libTencentSM.a ${SRCDIR}/lib/linux/libgmp.a

#include "sm.h"
*/
import "C"
import (
	"unsafe"
)

/**
 *@brief 生成16字节128bit的SM4 Key，也可调用该接口生成SM4 CBC模式的初始化向量iv，iv长度和key长度一致
 *@param outKey  函数出参 - 16字节密钥。
 *@return 0 -- OK
 */
func GenerateSM4Key(outkey []byte) int {
	if outkey == nil || len(outkey) < 16 {
		panic("invalid parameter")
	}
	return int(C.generateSM4Key((*C.uchar)(unsafe.Pointer(&outkey[0]))))
}

/**
 *@brief SM4 CBC模式对称加解密。加密，使用PKCS#7填充标准
 *@param in  函数入参 - 明文
 *@param inlen  函数入参 - 明文长度
 *@param out  函数出参 - 密文
 *@param outlen  函数出参 - 密文长度
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@return 0 -- OK
 */
func SM4_CBC_Encrypt(in []byte, inlen int, out []byte, outlen *int, key []byte, iv []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	return int(C.SM4_CBC_Encrypt(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0]))))
}

/**
 *@brief SM4 CBC模式对称加解密。解密，使用PKCS#7填充标准
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@return 0 -- OK
 */
func SM4_CBC_Decrypt(in []byte, inlen int, out []byte, outlen *int, key []byte, iv []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	return int(C.SM4_CBC_Decrypt(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0]))))
}

/**
 *@brief SM4 CBC模式对称加解密。加密，无填充。请保证明文为16字节整数倍，否则加密会失败，即出参outlen为0
 *@param in 函数入参 - 明文
 *@param inlen 函数入参 - 明文长度
 *@param out 函数出参 - 密文
 *@param outlen 函数出参 - 密文长度
 *@param key 函数入参 - 秘钥（128bit）
 *@param iv 函数入参 - 初始化向量
 *@return 0 -- OK
 */
func SM4_CBC_Encrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, key []byte, iv []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	return int(
		C.SM4_CBC_Encrypt_NoPadding(
			(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
			(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0]))))
}

/**
 *@brief SM4 CBC模式对称解密，无填充。请保证密文为16字节整数倍，否则解密会失败，即出参outlen为0
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@return 0 -- OK
 */
func SM4_CBC_Decrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, key []byte, iv []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	return int(C.SM4_CBC_Decrypt_NoPadding(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0]))))
}

/**
 *@brief SM4 ECB模式对称加解密。加密
 *@param in  函数入参 - 明文
 *@param inlen  函数入参 - 明文长度
 *@param out  函数出参 - 密文
 *@param outlen  函数出参 - 密文长度
 *@param key  函数入参 - 秘钥（128bit）
 *@return 0 -- OK
 */
func SM4_ECB_Encrypt(in []byte, inlen int, out []byte, outlen *int, key []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	return int(
		C.SM4_ECB_Encrypt(
			(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
			(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0]))))
}

/**
 *@brief SM4 ECB模式对称加解密。解密
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@param key  函数入参 - 秘钥（128bit）
 *@return 0 -- OK
 */
func SM4_ECB_Decrypt(in []byte, inlen int, out []byte, outlen *int, key []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	return int(
		C.SM4_ECB_Decrypt(
			(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
			(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0]))))
}

/**
 *@brief SM4 ECB模式对称加密，无填充。请保证明文为16字节整数倍，否则加密会失败，即出参outlen为0
 *@param in  函数入参 - 明文
 *@param inlen  函数入参 - 明文长度
 *@param out  函数出参 - 密文
 *@param outlen  函数出参 - 密文长度
 *@param key  函数入参 - 秘钥（128bit）
 *@return 0 -- OK
 */
func SM4_ECB_Encrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, key []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	return int(C.SM4_ECB_Encrypt_NoPadding(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0]))))
}

/**
 *@brief SM4 ECB模式对称解密，无填充。请保证密文为16字节整数倍，否则解密会失败，即出参outlen为0
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@param key  函数入参 - 秘钥（128bit）
 *@return 0 -- OK
 */
func SM4_ECB_Decrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, key []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	return int(
		C.SM4_ECB_Decrypt_NoPadding(
			(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
			(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0]))))
}

/**
 *@brief SM4 GCM模式对称加解密。加密，使用PKCS7填充，实际上GCM模式可不填充，非短明文加密推荐使用SM4_GCM_Encrypt_NoPadding替代。
*@param in  函数入参 - 明文
*@param inlen  函数入参 - 明文长度
*@param out  函数出参 - 密文
*@param outlen  函数出参 - 密文长度
*@param tag  函数出参 - GMAC值，即消息验证码
*@param taglen  既作函数入参也作为函数出参 - GMAC长度，通常取16字节
*@param key  函数入参 - 秘钥（128bit）
*@param iv  函数入参 - 初始化向量
*@param aad  函数入参 - 附加验证消息
*@param aadlen  函数入参 - 附加验证消息长度
*@return 成功为0，一般加密失败是由参数错误导致
*/
func SM4_GCM_Encrypt(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen *int, key []byte, iv []byte, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	return int(C.SM4_GCM_Encrypt(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (*C.size_t)(unsafe.Pointer(taglen)),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
 *@brief SM4 GCM模式对称加解密。解密，使用PKCS7填充，实际上GCM模式可不填充。
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@param tag  函数入参 - GMAC值，即消息验证码
 *@param taglen  函数入参 - GMAC长度，通常取16字节
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@param aad  函数入参 - 附加验证消息
 *@param aadlen  函数入参 - 附加验证消息长度
 *@return 成功为0，GCM的解密失败主要是tag校验失败
 */
func SM4_GCM_Decrypt(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen int, key []byte, iv []byte, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	tag = append(tag, 0)
	return int(C.SM4_GCM_Decrypt(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (C.size_t)(taglen),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
 *@brief SM4 GCM模式对称加解密。加密，无填充，明文长度无要求。
 *@param in  函数入参 - 明文
 *@param inlen  函数入参 - 明文长度
 *@param out  函数出参 - 密文
 *@param outlen  函数出参 - 密文长度(GCM NOPADDING模式密文长度与明文长度一致)
 *@param tag  函数出参 - GMAC值，即消息验证码
 *@param taglen  既作函数入参也作为函数出参 - GMAC长度，通常取16字节
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@param aad  函数入参 - 附加验证消息
 *@param aadlen  函数入参 - 附加验证消息长度
 *@return 成功为0，一般加密失败是由参数错误导致
 */
func SM4_GCM_Encrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen *int, key []byte, iv []byte, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	return int(C.SM4_GCM_Encrypt_NoPadding(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (*C.size_t)(unsafe.Pointer(taglen)),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
 *@brief SM4 GCM模式对称加解密。解密，无填充，密文长度无要求。
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度(GCM NOPADDING模式密文长度与明文长度一致)
 *@param tag  函数入参 - GMAC值，即消息验证码
 *@param taglen  函数入参 - GMAC长度，通常取16字节
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@param aad  函数入参 - 附加验证消息
 *@param aadlen  函数入参 - 附加验证消息长度
 *@return 返回解密是否失败，GCM的解密失败主要是tag校验失败
 */
func SM4_GCM_Decrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen int, key []byte, iv []byte, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	tag = append(tag, 0)
	return int(C.SM4_GCM_Decrypt_NoPadding(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (C.size_t)(taglen),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
 *@brief SM4 GCM模式对称加解密。加密，使用PKCS7填充，实际上GCM模式可不填充，非短明文加密推荐使用SM4_GCM_Encrypt_NoPadding替代。
*@param in  函数入参 - 明文
*@param inlen  函数入参 - 明文长度
*@param out  函数出参 - 密文
*@param outlen  函数出参 - 密文长度
*@param tag  函数出参 - GMAC值，即消息验证码
*@param taglen  既作函数入参也作为函数出参 - GMAC长度，通常取16字节
*@param key  函数入参 - 秘钥（128bit）
*@param iv  函数入参 - 初始化向量
*@param ivlen 按照NIST SP800-38D标准实现GCM部分算法,RFC5647标准iv推荐使用12字节，96bit
*@param aad  函数入参 - 附加验证消息
*@param aadlen  函数入参 - 附加验证消息长度
*@return 成功为0，一般加密失败是由参数错误导致
*/
func SM4_GCM_Encrypt_NIST_SP800_38D(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen *int, key []byte, iv []byte, ivlen int, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	return int(C.SM4_GCM_Encrypt_NIST_SP800_38D(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (*C.size_t)(unsafe.Pointer(taglen)),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])), (C.size_t)(ivlen),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
 *@brief SM4 GCM模式对称加解密。解密，使用PKCS7填充，实际上GCM模式可不填充。
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@param tag  函数入参 - GMAC值，即消息验证码
 *@param taglen  函数入参 - GMAC长度，通常取16字节
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@param ivlen 按照NIST SP800-38D标准实现GCM部分算法,RFC5647标准iv推荐使用12字节，96bit
 *@param aad  函数入参 - 附加验证消息
 *@param aadlen  函数入参 - 附加验证消息长度
 *@return 成功为0，GCM的解密失败主要是tag校验失败
 */
func SM4_GCM_Decrypt_NIST_SP800_38D(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen int, key []byte, iv []byte, ivlen int, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	tag = append(tag, 0)
	return int(C.SM4_GCM_Decrypt_NIST_SP800_38D(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (C.size_t)(taglen),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])), (C.size_t)(ivlen),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
 *@brief SM4 GCM模式对称加解密。加密，无填充，明文长度无要求。
 *@param in  函数入参 - 明文
 *@param inlen  函数入参 - 明文长度
 *@param out  函数出参 - 密文
 *@param outlen  函数出参 - 密文长度(GCM NOPADDING模式密文长度与明文长度一致)
 *@param tag  函数出参 - GMAC值，即消息验证码
 *@param taglen  既作函数入参也作为函数出参 - GMAC长度，通常取16字节
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@param ivlen 按照NIST SP800-38D标准实现GCM部分算法,RFC5647标准iv推荐使用12字节，96bit
 *@param aad  函数入参 - 附加验证消息
 *@param aadlen  函数入参 - 附加验证消息长度
 *@return 成功为0，一般加密失败是由参数错误导致
 */
func SM4_GCM_Encrypt_NoPadding_NIST_SP800_38D(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen *int, key []byte, iv []byte, ivlen int, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	return int(C.SM4_GCM_Encrypt_NoPadding_NIST_SP800_38D(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (*C.size_t)(unsafe.Pointer(taglen)),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])), (C.size_t)(ivlen),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
 *@brief SM4 GCM模式对称加解密。解密，无填充，密文长度无要求。
 *@param in  函数入参 - 密文
 *@param inlen  函数入参 - 密文长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度(GCM NOPADDING模式密文长度与明文长度一致)
 *@param tag  函数入参 - GMAC值，即消息验证码
 *@param taglen  函数入参 - GMAC长度，通常取16字节
 *@param key  函数入参 - 秘钥（128bit）
 *@param iv  函数入参 - 初始化向量
 *@param ivlen 按照NIST SP800-38D标准实现GCM部分算法,RFC5647标准iv推荐使用12字节，96bit
 *@param aad  函数入参 - 附加验证消息
 *@param aadlen  函数入参 - 附加验证消息长度
 *@return 返回解密是否失败，GCM的解密失败主要是tag校验失败
 */
func SM4_GCM_Decrypt_NoPadding_NIST_SP800_38D(in []byte, inlen int, out []byte, outlen *int, tag []byte, taglen int, key []byte, iv []byte, ivlen int, aad []byte, aadlen int) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	aad = append(aad, 0)
	tag = append(tag, 0)
	return int(C.SM4_GCM_Decrypt_NoPadding_NIST_SP800_38D(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&tag[0])), (C.size_t)(taglen),
		(*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0])), (C.size_t)(ivlen),
		(*C.uchar)(unsafe.Pointer(&aad[0])), (C.size_t)(aadlen)))

}

/**
SM4 CTR模式对称加解密。加密，CTR模式不需要填充。
@param in 函数入参 - 明文
@param inlen 函数入参 - 明文长度
@param out  函数出参 - 密文 - 为out分配的内存长度为明文长度可保证安全
@param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
@param key 函数入参 - 秘钥（128bit）
@param iv 函数入参 - 初始化向量
*/
func SM4_CTR_Encrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, key []byte, iv []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	return int(C.SM4_CTR_Encrypt_NoPadding(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0]))))
}

/**
SM4 CTR模式对称加解密。解密，CTR模式不需要填充。
@param in  函数入参 - 密文
@param inlen  函数入参 - 密文长度
@param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
@param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
@param key  函数入参 - 秘钥（128bit）
@param iv  函数入参 - 初始化向量
*/
func SM4_CTR_Decrypt_NoPadding(in []byte, inlen int, out []byte, outlen *int, key []byte, iv []byte) int {
	if in == nil || inlen <= 0 || out == nil || outlen == nil || key == nil || iv == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	key = append(key, 0)
	iv = append(iv, 0)
	return int(C.SM4_CTR_Decrypt_NoPadding(
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen), (*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.size_t)(unsafe.Pointer(outlen)), (*C.uchar)(unsafe.Pointer(&key[0])), (*C.uchar)(unsafe.Pointer(&iv[0]))))
}
