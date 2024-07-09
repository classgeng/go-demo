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
 *@brief SM3上下文结构体的大小
 */
func SM3CtxSize() int {
	return int(C.SM3CtxSize())
}

func SM3Init(ctx *SM3_ctx_t) int {
	if ctx == nil {
		panic("invalid parameter")
	}
	return int(C.SM3Init(&ctx.Context))
}

func SM3Update(ctx *SM3_ctx_t, data []byte, datalen int) int {
	if ctx == nil || data == nil || datalen <= 0 {
		panic("invalid parameter")
	}
	data = append(data, 0)
	return int(C.SM3Update(&ctx.Context, (*C.uchar)(unsafe.Pointer(&data[0])), (C.size_t)(datalen)))
}

func SM3Final(ctx *SM3_ctx_t, digest []byte) int {
	if ctx == nil || digest == nil {
		panic("invalid parameter")
	}
	return int(C.SM3Final(&ctx.Context, (*C.uchar)(unsafe.Pointer(&digest[0]))))
}

/**
 *@brief SM3 hash算法， 内部依次调用了init update和final三个接口
 */
func SM3(data []byte, datalen int, digest []byte) int {
	if data == nil || digest == nil || datalen <= 0 {
		panic("invalid parameter")
	}
	data = append(data, 0)
	return int(C.SM3((*C.uchar)(unsafe.Pointer(&data[0])), (C.size_t)(datalen), (*C.uchar)(unsafe.Pointer(&digest[0]))))
}

/**
 * @brief 基于sm3算法计算HMAC值 ctx init
 * @param key HMAC用的秘钥
 * @param key_len 秘钥长度
 * @return 0 -- OK
 */
func SM3HMACInit(key []byte, keyLen int) *HmacSm3Ctx {
	if key == nil || keyLen <= 0 {
		panic("invalid parameter")
	}
	key = append(key, 0)
	var ret HmacSm3Ctx
	cRet := C.SM3_HMAC_Init((*C.uchar)(unsafe.Pointer(&key[0])), (C.size_t)(keyLen))
	ret.Context = cRet
	return &ret
}

/**
 * @brief 基于sm3算法计算HMAC值 update数据
 * @param ctx hmac上下文结构指针
 * @param data 做HMAC计算的数据
 * @param data_len 数据长度
 * @return 0 -- OK
 */
func SM3HmacUpdate(ctx *HmacSm3Ctx, data []byte, dataLen int) int {
	if data == nil || ctx == nil || dataLen <= 0 {
		panic("invalid parameter")
	}
	data = append(data, 0)
	return int(C.SM3_HMAC_Update(ctx.Context, (*C.uchar)(unsafe.Pointer(&data[0])), (C.size_t)(dataLen)))
}

/**
 * @brief 基于sm3算法计算HMAC值 最终计算HMAC值
 * @param ctx hmac上下文结构指针
 * @param mac 输出的HMAC字节码
 * @return 0 -- OK
 */
func SM3HmacFinal(ctx *HmacSm3Ctx, mac []byte, macLen int) int {
	if mac == nil || len(mac) != macLen || macLen != SM3_HMAC_SIZE {
		panic("invalid parameter")
	}
	return int(C.SM3_HMAC_Final(ctx.Context, (*C.uchar)(unsafe.Pointer(&mac[0]))))
}

/**
 * @brief 基于sm3算法计算HMAC值
 * @param data 做HMAC计算的数据
 * @param data_len 数据长度
 * @param key HMAC用的秘钥
 * @param key_len 秘钥长度
 * @param mac 输出的HMAC字节码
 * @return 0 -- OK
 */
func SM3_HMAC(ctx *HmacSm3Ctx, data []byte, dataLen int, key []byte, keyLen int, mac []byte, macLen int) int {
	if data == nil || key == nil || mac == nil || len(mac) != macLen || macLen != SM3_HMAC_SIZE {
		panic("invalid parameter")
	}
	data = append(data, 0)
	key = append(key, 0)
	return int(C.SM3_HMAC((*C.uchar)(unsafe.Pointer(&data[0])), (C.size_t)(dataLen), (*C.uchar)(unsafe.Pointer(&key[0])), (C.size_t)(keyLen), (*C.uchar)(unsafe.Pointer(&mac[0]))))
}

/**
 *@brief 密钥导出函数。sharelen不可大于1024字节。
 */
func SM3KDF(share []byte, shareLen int, outkey []byte, keyLen int) int {
	if share == nil || shareLen <= 0 || outkey == nil || keyLen <= 0 {
		panic("invalid parameter")
	}
	return int(C.SM3KDF((*C.uchar)(unsafe.Pointer(&share[0])), (C.size_t)(shareLen), (*C.uchar)(unsafe.Pointer(&outkey[0])), (C.size_t)(keyLen)))
}

/**
基于SM3的PBKDF2算法接口，可用于口令哈希。

@param msg 函数入参 - 待哈希的原始消息
@param msglen 函数入参 - 待哈希的原始消息长度
@param salt 函数入参 - 外部添加的盐值，可为空，建议添加动态盐以增强安全性，该函数内部也会添加静态盐
@param saltlen 函数入参 - 外部添加的盐值长度
@param roundtimes 函数入参 - 算法的轮次，数值越大，构建彩虹表攻击的难度越大，但数值越大其性能开销越大
@param outbuf 函数输出- 输出长度为32字节数据
*/
func SM3BasedPBKDF2(msg []byte, msglen int, salt []byte, saltlen int, roundtimes int, outbuf []byte) int {
	if msg == nil || msglen <= 0 || salt == nil || saltlen <= 0 || outbuf == nil {
		panic("invalid parameter")
	}
	return int(C.SM3BasedPBKDF2((*C.uchar)(unsafe.Pointer(&msg[0])), C.int(msglen), (*C.uchar)(unsafe.Pointer(&salt[0])), C.int(saltlen), C.int(roundtimes), (*C.uchar)(unsafe.Pointer(&outbuf[0]))))

}
