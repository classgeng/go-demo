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
#include <stdlib.h>
#include<string.h>

typedef struct stTCSMItemWrapper {
    int itemType;
    void* value;
    int valuelen;
} TstTCSMItemWrapper;

void createAndCopy(void** p, unsigned char* src, int size){
    if(size > 0 && p != NULL && src != NULL){
		*p = (unsigned char*)malloc(size);
		memcpy(*p, src, size);
	}
}

void crelease(void* p){
    if(p != NULL){
		free(p);
	}
}

*/
import "C"
import (
	"unsafe"
)

/**
 *@brief SM2上下文结构体的大小
 */
func SM2CtxSize() int {
	return int(C.SM2CtxSize())
}

/**
 *@brief 使用SM2获取公私钥或加解密之前，必须调用SM2InitCtx或者SM2InitCtxWithPubKey函数
 *@param ctx  函数出参 - 上下文
 *@return 0-成功，其他见错误码
 */
func SM2InitCtx(ctx *SM2_ctx_t) int {
	if ctx == nil {
		panic("invalid parameter")
	}
	var ret = -1
	if IsContextPoolEnable {
		sm2Ctx := ContextPool.Get().(*SM2_ctx_t)
		*ctx = *sm2Ctx
		ret = 0
	} else {
		ret = int(C.SM2InitCtx(&ctx.Context))
	}
	return ret
}

/**
 * @brief 使用SM2获取公私钥或加解密之前，必须调用SM2InitCtx或者SM2InitCtxWithPubKey函数.如果使用固定公钥加密，可调用SM2InitCtxWithPubKey，将获得较大性能提升
 * @param ctx  函数出参 - 上下文
 * @param pubkey  函数入参 - 公钥
 * @return int --0-- successful
 */
func SM2InitCtxWithPubKey(ctx *SM2_ctx_t, pubkey []byte) int {
	if ctx == nil || pubkey == nil {
		panic("invalid parameter")
	}
	if len(pubkey) < 130 {
		panic("memory len is too small")
	}
	pubkey = append(pubkey, 0)
	return int(C.SM2InitCtxWithPubKey(&ctx.Context, (*C.char)(unsafe.Pointer(&pubkey[0]))))
}

/**
 *@brief 使用完SM2算法后，必须调用free函数释放
 *@param ctx  函数入参 - 上下文
 *@return int --0-- successful
 */
func SM2FreeCtx(ctx *SM2_ctx_t) int {
	if ctx == nil {
		panic("invalid parameter")
	}
	var ret int = -1
	if IsContextPoolEnable {
		ContextPool.Put(ctx)
		ret = 0
	} else {
		ret = int(C.SM2FreeCtx(&ctx.Context))
	}
	return ret

}

/**
 *@brief 生成私钥
 *@param ctx  函数入参 - 上下文
 *@param out  函数出参 - 私钥，私钥实际上为256bit的大整数，这里输出的为256bit二进制内容Hex后的ASCII编码的可见字符串，长度为64字节，为保证字符串的结束符0，out至少需分配65字节空间。
 *@return  0表示成功，其他值为错误码
 */
func GeneratePrivateKey(ctx *SM2_ctx_t, out []byte) int {
	if ctx == nil || out == nil {
		panic("invalid parameter")
	}
	return int(C.generatePrivateKey(&ctx.Context, (*C.char)(unsafe.Pointer(&out[0]))))
}

/**
 *@brief根据私钥生成对应公钥，
 *@param ctx 函数入参 - 上下文
 *@param privateKey 函数入参 - 私钥，私钥实际上为256bit的大整数，这里输出的为256bit二进制内容Hex后的ASCII编码的可见字符串，长度为64字节
 *@param outPubKey 函数出参 - 公钥，公钥格式为04 | X | Y，其中X和Y为256bit大整数，这里输出的为04 | X | Y的二进制内容Hex后的ASCII编码的可见字符串，长度为130字节，为保证字符串的结束符0，outPubKey至少需分配131字节空间。
 *@return  0表示成功，其他值为错误码
 */
func GeneratePublicKey(ctx *SM2_ctx_t, privateKey []byte, outPubKey []byte) int {
	if ctx == nil || privateKey == nil || outPubKey == nil {
		panic("invalid parameter")
	}
	privateKey = append(privateKey, 0)
	return int(C.generatePublicKey(&ctx.Context, (*C.char)(unsafe.Pointer(&privateKey[0])), (*C.char)(unsafe.Pointer(&outPubKey[0]))))
}

/**
 *@brief生成公私钥对
 *@param ctx 函数入参 - 上下文
 *@param outPriKey 函数出参 - 私钥，私钥实际上为256bit的大整数，这里输出的为256bit二进制内容Hex后的ASCII编码的可见字符串，长度为64字节，为保证字符串的结束符0，outPubKey至少需分配65字节空间。
 *@param outPubKey 函数出参 - 公钥，公钥格式为04 | X | Y，其中X和Y为256bit大整数，这里输出的为04 | X | Y的二进制内容Hex后的ASCII编码的可见字符串，长度为130字节，为保证字符串的结束符0，outPubKey至少需分配131字节空间。
 *@return  0表示成功，其他值为错误码
 */
func GenerateKeyPair(ctx *SM2_ctx_t, outPriKey []byte, outPubKey []byte) int {
	if ctx == nil || outPriKey == nil || outPubKey == nil {
		panic("invalid parameter")
	}
	return int(C.generateKeyPair(&ctx.Context, (*C.char)(unsafe.Pointer(&outPriKey[0])), (*C.char)(unsafe.Pointer(&outPubKey[0]))))
}

/**
 *@briefSM2非对称加解密算法，加密
 *@param ctx 函数入参 - 上下文
 *@param in  函数入参 - 待加密消息
 *@param inlen  函数入参 - 消息长度(字节单位)
 *@param strPubKey  函数入参 - 公钥
 *@param pubkeyLen  函数入参 - 公钥长度
 *@param out  函数出参 - 密文
 *@param outlen  函数出参 - 密文长度
 *@return  0表示成功，其他值为错误码
 */
func SM2Encrypt(ctx *SM2_ctx_t, in []byte, inlen int, strPubKey []byte, pubkeyLen int, out []byte, outlen *int) int {
	if ctx == nil || in == nil || inlen <= 0 || strPubKey == nil || pubkeyLen <= 0 || out == nil || outlen == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	strPubKey = append(strPubKey, 0)
	return int(C.SM2Encrypt(&ctx.Context, (*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen),
		(*C.char)(unsafe.Pointer(&strPubKey[0])), (C.size_t)(pubkeyLen),
		(*C.uchar)(unsafe.Pointer(&out[0])), (*C.size_t)(unsafe.Pointer(outlen))))
}

/**
 *@briefSM2非对称加解密算法，解密
 *@param ctx  函数入参 - 上下文
 *@param in  函数入参 - 待解密密文
 *@param inlen  函数入参 - 密文长度(字节单位)
 *@param strPriKey  函数入参 - 私钥
 *@param prikeyLen  函数入参 - 私钥长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@return  0表示成功，其他值为错误码
 */
func SM2Decrypt(ctx *SM2_ctx_t, in []byte, inlen int, strPriKey []byte, prikeyLen int, out []byte, outlen *int) int {
	if ctx == nil || in == nil || inlen <= 0 || strPriKey == nil || prikeyLen <= 0 || out == nil || outlen == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	strPriKey = append(strPriKey, 0)
	return int(C.SM2Decrypt(&ctx.Context, (*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen),
		(*C.char)(unsafe.Pointer(&strPriKey[0])), (C.size_t)(prikeyLen),
		(*C.uchar)(unsafe.Pointer(&out[0])), (*C.size_t)(unsafe.Pointer(outlen))))
}

/**
 *@briefSM2签名验签算法，签名
 *@param ctx 函数入参 - 上下文
 *@param msg 函数入参 - 待签名消息
 *@param msglen 函数入参 - 待签名消息长度
 *@param id 函数入参 - 用户ID(作用是加入到签名hash中，对于传入值无特殊要求)
 *@param idlen 函数入参 - 用户ID长度
 *@param strPubKey 函数入参 - 公钥(作用是加入到签名hash中)
 *@param pubkeyLen 函数入参 - 公钥长度
 *@param strPriKey 函数入参 - 私钥
 *@param prikeyLen 函数入参 - 私钥长度
 *@param sig 函数出参 - 签名结果
 *@param siglen 函数出参 - 签名结果长度
 */
func SM2Sign(ctx *SM2_ctx_t, msg []byte, msglen int, id []byte, idlen int, strPubKey []byte, pubkeyLen int, strPriKey []byte, prikeyLen int, sig []byte, siglen *int) int {
	if ctx == nil || msg == nil || msglen <= 0 || id == nil || idlen <= 0 || strPubKey == nil || pubkeyLen <= 0 ||
		strPriKey == nil || prikeyLen <= 0 || sig == nil || siglen == nil {
		panic("invalid parameter")
	}
	msg = append(msg, 0)
	id = append(id, 0)
	strPubKey = append(strPubKey, 0)
	strPriKey = append(strPriKey, 0)
	return int(C.SM2Sign(&ctx.Context,
		(*C.uchar)(unsafe.Pointer(&msg[0])), (C.size_t)(msglen),
		(*C.char)(unsafe.Pointer(&id[0])), (C.size_t)(idlen),
		(*C.char)(unsafe.Pointer(&strPubKey[0])), (C.size_t)(pubkeyLen),
		(*C.char)(unsafe.Pointer(&strPriKey[0])), (C.size_t)(prikeyLen),
		(*C.uchar)(unsafe.Pointer(&sig[0])), (*C.size_t)(unsafe.Pointer(siglen))))
}

/**
 *@briefSM2签名验签算法，验签
 *@param ctx 函数入参 - 上下文
 *@param msg 函数入参 - 待验签内容
 *@param msglen 函数入参 - 待验签内容长度
 *@param id 函数入参 - 用户ID
 *@param idlen 函数入参 - 用户ID长度
 *@param sig 函数入参 - 签名结果
 *@param siglen 函数入参 - 签名结果长度
 *@param strPubKey 函数入参 - 公钥
 *@param pubkeyLen 函数入参 - 公钥长度
 *@return 0表示成功，其他值为错误码
 */
func SM2Verify(ctx *SM2_ctx_t, msg []byte, msglen int, id []byte, idlen int, sig []byte, siglen int, strPubKey []byte, pubkeyLen int) int {
	if ctx == nil || msg == nil || msglen <= 0 || id == nil || idlen <= 0 || sig == nil || siglen <= 0 || strPubKey == nil || pubkeyLen <= 0 {
		panic("invalid parameter")
	}
	msg = append(msg, 0)
	id = append(id, 0)
	sig = append(sig, 0)
	strPubKey = append(strPubKey, 0)
	return int(C.SM2Verify(&ctx.Context,
		(*C.uchar)(unsafe.Pointer(&msg[0])), (C.size_t)(msglen),
		(*C.char)(unsafe.Pointer(&id[0])), (C.size_t)(idlen),
		(*C.uchar)(unsafe.Pointer(&sig[0])), (C.size_t)(siglen),
		(*C.char)(unsafe.Pointer(&strPubKey[0])), (C.size_t)(pubkeyLen)))

}

/**
 *@brief SM2密钥交换算法接口，计算共享密钥。SM2密钥交换算法需要AB两方参与，其中一方为发起方，一方为响应方。我方为A，持有一对自身的公私钥和一对自身的临时公私钥。对方为B，持有对方的公钥和临时公钥。
 *@param ctx 函数入参 - 上下文
 *@param strPriKey_A 函数入参 -  A的私钥
 *@param strPubKey_A 函数入参 -  A的公钥
 *@param strTempPriKey_A 函数入参 -  A的临时私钥
 *@param strTempPubKey_A 函数入参 -  A的临时公钥
 *@param strPubKey_B 函数入参 -  B的公钥
 *@param strTempPubKey_B 函数入参 -  B的临时公钥
 *@param id_A 函数入参 -  A的ID，可为NULL，当使用NULL时，将使用国标中的默认值，字符串“1234567812345678”
 *@param idlen_A 函数入参 -  A的ID长度，当ID为空时，填0
 *@param id_B 函数入参 -  B的ID，可为NULL，当使用NULL时，将使用国标中的默认值，字符串“1234567812345678”
 *@param idlen_B 函数入参 -  B的ID长度，当ID为空时，填0
 *@param klen 函数入参 -  需要计算的共享密钥长度
 *@param sharedKey 函数出参 - 共享密钥
 *@param A_is_initiator 函数出参 - A是否为发起方，即我方是否为发起方，如果我方为发起方则填1，否则为0
 *@return 0表示成功，其他值为错误码
 */
func SM2CalculateSharedKey(ctx *SM2_ctx_t, strPriKey_A []byte, strPubKey_A []byte, strTempPriKey_A []byte, strTempPubKey_A []byte, strPubKey_B []byte, strTempPubKey_B []byte, id_A []byte, idlen_A int, id_B []byte, idlen_B int, klen int, sharedKey []byte, A_is_initiator int) int {
	if ctx == nil || strPriKey_A == nil || strPubKey_A == nil || strTempPriKey_A == nil || strTempPubKey_A == nil ||
		strPubKey_B == nil || strTempPubKey_B == nil || klen <= 0 || sharedKey == nil {
		panic("invalid parameter")
	}
	strPriKey_A = append(strPriKey_A, 0)
	strPubKey_A = append(strPubKey_A, 0)
	strTempPriKey_A = append(strTempPriKey_A, 0)
	strTempPubKey_A = append(strTempPubKey_A, 0)
	strPubKey_B = append(strPubKey_B, 0)
	strTempPubKey_B = append(strTempPubKey_B, 0)
	id_A = append(id_A, 0)
	id_B = append(id_B, 0)

	var idA *C.char = nil
	var idB *C.char = nil
	if id_A != nil {
		idA = (*C.char)(unsafe.Pointer(&id_A[0]))
	}
	if id_B != nil {
		idB = (*C.char)(unsafe.Pointer(&id_B[0]))
	}
	return int(C.SM2CalculateSharedKey(&ctx.Context,
		(*C.char)(unsafe.Pointer(&strPriKey_A[0])),
		(*C.char)(unsafe.Pointer(&strPubKey_A[0])),
		(*C.char)(unsafe.Pointer(&strTempPriKey_A[0])),
		(*C.char)(unsafe.Pointer(&strTempPubKey_A[0])),
		(*C.char)(unsafe.Pointer(&strPubKey_B[0])),
		(*C.char)(unsafe.Pointer(&strTempPubKey_B[0])),
		idA, (C.size_t)(idlen_A),
		idB, (C.size_t)(idlen_B),
		(C.size_t)(klen),
		(*C.uchar)(unsafe.Pointer(&sharedKey[0])), C.int(A_is_initiator),
	))

}

/**
 *@brief SM2非对称加解密算法，加密的兼容接口
 *@param ctx 函数入参 - 上下文
 *@param in 函数入参 - 待加密消息
 *@param inlen 函数入参 - 消息长度(字节单位)
 *@param strPubKey 函数入参 - 公钥
 *@param pubkeyLen 函数入参 - 公钥长度
 *@param out 函数出参 - 密文
 *@param outlen 函数出参 - 密文长度
 *@param mode 密文输出格式
 *@return  0表示成功，其他值为错误码
 */
func SM2EncryptWithMode(ctx *SM2_ctx_t, in []byte, inlen int, strPubKey []byte, pubkeyLen int, out []byte, outlen *int, mode SM2CipherMode) int {
	if ctx == nil || in == nil || inlen <= 0 || strPubKey == nil || pubkeyLen <= 0 || out == nil || outlen == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	strPubKey = append(strPubKey, 0)
	return int(C.SM2EncryptWithMode(&ctx.Context,
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen),
		(*C.char)(unsafe.Pointer(&strPubKey[0])), (C.size_t)(pubkeyLen),
		(*C.uchar)(unsafe.Pointer(&out[0])), (*C.size_t)(unsafe.Pointer(outlen)),
		ConvertSMCipherModeToC(mode)))
}

/**
 *@brief SM2非对称加解密算法，解密的兼容接口
 *@param ctx 函数入参 - 上下文
 *@param in  函数入参 - 待解密密文
 *@param inlen  函数入参 - 密文长度(字节单位)
 *@param strPriKey  函数入参 - 私钥
 *@param prikeyLen  函数入参 - 私钥长度
 *@param out  函数出参 - 明文
 *@param outlen  函数出参 - 明文长度
 *@param mode  密文格式
 *@return  0表示成功，其他值为错误码
 */
func SM2DecryptWithMode(ctx *SM2_ctx_t, in []byte, inlen int, strPriKey []byte, prikeyLen int, out []byte, outlen *int, mode SM2CipherMode) int {
	if ctx == nil || in == nil || inlen <= 0 || strPriKey == nil || prikeyLen <= 0 || out == nil || outlen == nil {
		panic("invalid parameter")
	}
	in = append(in, 0)
	strPriKey = append(strPriKey, 0)
	return int(C.SM2DecryptWithMode(&ctx.Context,
		(*C.uchar)(unsafe.Pointer(&in[0])), (C.size_t)(inlen),
		(*C.char)(unsafe.Pointer(&strPriKey[0])), (C.size_t)(prikeyLen),
		(*C.uchar)(unsafe.Pointer(&out[0])), (*C.size_t)(unsafe.Pointer(outlen)),
		ConvertSMCipherModeToC(mode)))
}

/**
 *@brief SM2签名验签算法，签名的兼容接口
 *@param ctx 函数入参 - 上下文
 *@param msg 函数入参 - 待签名消息
 *@param msglen 函数入参 - 待签名消息长度
 *@param id 函数入参 - 用户ID(作用是加入到签名hash中，对于传入值无特殊要求)
 *@param idlen 函数入参 - 用户ID长度
 *@param strPubKey 函数入参 - 公钥(作用是加入到签名hash中)
 *@param pubkeyLen 函数入参 - 公钥长度
 *@param strPriKey 函数入参 - 私钥
 *@param prikeyLen 函数入参 - 私钥长度
 *@param sig 函数出参 - 签名结果
 *@param siglen 函数出参 - 签名结果长度
 *@param mode 签名格式
 */
func SM2SignWithMode(ctx *SM2_ctx_t, msg []byte, msglen int, id []byte, idlen int, strPubKey []byte, pubkeyLen int, strPriKey []byte, prikeyLen int, sig []byte, siglen *int, signMode SM2SignMode) int {
	if ctx == nil || msg == nil || msglen <= 0 || id == nil || idlen <= 0 || strPubKey == nil || pubkeyLen <= 0 ||
		strPriKey == nil || prikeyLen <= 0 || sig == nil || siglen == nil {
		panic("invalid parameter")
	}
	msg = append(msg, 0)
	id = append(id, 0)
	strPubKey = append(strPubKey, 0)
	strPriKey = append(strPriKey, 0)
	return int(C.SM2SignWithMode(&ctx.Context,
		(*C.uchar)(unsafe.Pointer(&msg[0])), (C.size_t)(msglen),
		(*C.char)(unsafe.Pointer(&id[0])), (C.size_t)(idlen),
		(*C.char)(unsafe.Pointer(&strPubKey[0])), (C.size_t)(pubkeyLen),
		(*C.char)(unsafe.Pointer(&strPriKey[0])), (C.size_t)(prikeyLen),
		(*C.uchar)(unsafe.Pointer(&sig[0])), (*C.size_t)(unsafe.Pointer(siglen)), ConvertSMSignModeToC(signMode)))
}

/**
 *@brief SM2签名验签算法，验签的兼容接口
 *@param ctx 函数入参 - 上下文
 *@param msg 函数入参 - 待验签内容
 *@param msglen 函数入参 - 待验签内容长度
 *@param id 函数入参 - 用户ID
 *@param idlen 函数入参 - 用户ID长度
 *@param sig 函数入参 - 签名结果
 *@param siglen 函数入参 - 签名结果长度
 *@param strPubKey 函数入参 - 公钥
 *@param pubkeyLen 函数入参 - 公钥长度
 *@param mode 签名格式
 *@return 0表示成功，其他值为错误码
 */
func SM2VerifyWithMode(ctx *SM2_ctx_t, msg []byte, msglen int, id []byte, idlen int, sig []byte, siglen int, strPubKey []byte, pubkeyLen int, signMode SM2SignMode) int {
	if ctx == nil || msg == nil || msglen <= 0 || id == nil || idlen <= 0 || sig == nil || siglen <= 0 || strPubKey == nil || pubkeyLen <= 0 {
		panic("invalid parameter")
	}
	msg = append(msg, 0)
	id = append(id, 0)
	strPubKey = append(strPubKey, 0)
	sig = append(sig, 0)
	return int(C.SM2VerifyWithMode(&ctx.Context,
		(*C.uchar)(unsafe.Pointer(&msg[0])), (C.size_t)(msglen),
		(*C.char)(unsafe.Pointer(&id[0])), (C.size_t)(idlen),
		(*C.uchar)(unsafe.Pointer(&sig[0])), (C.size_t)(siglen),
		(*C.char)(unsafe.Pointer(&strPubKey[0])), (C.size_t)(pubkeyLen),
		ConvertSMSignModeToC(signMode)))
}

/**
 *@brief  签名时使用的随机数由外部设置
 *@param ctx
 *@param sign_random
 *@return 0 -- successful
 */
func SM2SetRandomDataCtx(ctx *SM2_ctx_t, sign_random []byte) int {
	if ctx == nil || sign_random == nil {
		panic("invalid parameter")
	}
	return int(C.SM2SetRandomDataCtx(&ctx.Context, (*C.char)(unsafe.Pointer(&sign_random[0]))))
}

/**
 *@brief 判断外部设置的签名时随机数是否有效
 */
func IsSM2CtxRandomDataVaild(ctx *SM2_ctx_t) int {
	if ctx == nil {
		panic("invalid parameter")
	}
	return int(C.IsSM2CtxRandomDataVaild(&ctx.Context))
}

/**
 *@brief 为SM2增加外部熵源，该接口主要用于输入外部随机熵，一般情况下无需调用该接口，模块内部使用的随机熵通常情况下可保证熵值足够。
 *@param ctx 函数入参 - 上下文
 *@param buf 函数入参 - 熵buf
 *@param buflen 函数入参 - 熵buf长度
 */
func SM2ReSeed(ctx *SM2_ctx_t, buf []byte, buflen int) int {
	if ctx == nil || buf == nil || buflen <= 0 {
		panic("invalid parameter")
	}
	buf = append(buf, 0)
	return int(C.SM2ReSeed(&ctx.Context, (*C.uchar)(unsafe.Pointer(&buf[0])), (C.size_t)(buflen)))
}

func SM2CertCtxSize() int {
	return int(C.SM2CertCtxSize())
}

/**
 *@brief 初始化证书管理上下文sm2_cert_ctx_t
 *@param ctx 函数入参 - 上下文
 *@param dir 函数入参 - 目录，指定存储证书信息的目录,不导入证书时可以传NULL
 *@return 0表示成功，其他值为错误码
 */
func SM2CertInitCtx(ctx *SM2_cert_ctx_t, dir []byte) int {
	if ctx == nil {
		panic("invalid parameter")
	}
	dir = append(dir, 0)
	if dir == nil {
		return int(C.SM2CertInitCtx(&ctx.Context, (*C.char)(unsafe.Pointer(nil))))
	}
	return int(C.SM2CertInitCtx(&ctx.Context, (*C.char)(unsafe.Pointer(&dir[0]))))
}

/**
 *@brief 销毁证书管理上下文sm2_cert_ctx_t
 *@param ctx 函数入参 - 上下文
 *@return 0表示成功，其他值为错误码
 */
func SM2CertFreeCtx(ctx *SM2_cert_ctx_t) int {
	if ctx == nil {
		panic("invalid parameter")
	}
	return int(C.SM2CertFreeCtx(&ctx.Context))
}

/**
 *@brief 生成SM2的证书请求串，生成的CSR以DER格式编码，支持普通单证书和国密双证书
 *@param ctx 函数入参 - 上下文
 *@param country_name 函数入参 - CountryName
 *@param province 函数入参 - StateOrProvinceName
 *@param locality_name 函数入参 - LocalityName
 *@param organization_name 函数入参 - OrganizationName
 *@param organization_unit_name 函数入参 - OrganizationalUnitName
 *@param common_name 函数入参 - CommonName
 *@param email 函数入参 - Email
 *@param challenge_password 函数入参 - ChallengePassword
 *@param public_key 函数入参 - 公钥
 *@param private_key 函数入参 - 私钥，用作CSR签名
 *@param temp_public_key 函数入参 - 临时公钥，国密双证书需要填写该字段，单证书填NULL，双证书使用该公钥用作交换公私钥，用于解密CA下发的加密密钥对
 *@param out 函数出参 - CSR输出，请保证out指向的内存空间足够，例如可以分配8K空间，der编码
 *@param outlen 函数出参 - CSR输出的长度
 *@param mode 模式 -  选择生成SM2单证书或双证书
 *@return 0表示成功，其他值为错误码
 */
func SM2GenerateCSR(ctx *SM2_cert_ctx_t, country_name []byte, province []byte, locality_name []byte, organization_name []byte,
	organization_unit_name []byte, common_name []byte, email []byte, challenge_password []byte, public_key []byte,
	private_key []byte, temp_public_key []byte, out []byte, outlen *int, mode SM2CSRMode) int {
	if ctx == nil || out == nil || outlen == nil {
		panic("invalid parameter")
	}
	locality_name = append(locality_name, 0)
	organization_name = append(organization_name, 0)
	country_name = append(country_name, 0)
	province = append(province, 0)
	organization_unit_name = append(organization_unit_name, 0)
	common_name = append(common_name, 0)
	email = append(email, 0)
	public_key = append(public_key, 0)
	challenge_password = append(challenge_password, 0)
	private_key = append(private_key, 0)
	temp_public_key = append(temp_public_key, 0)
	var temp_pubkey *C.char = nil
	if temp_public_key != nil {
		temp_pubkey = (*C.char)(unsafe.Pointer(&temp_public_key[0]))
	}
	return int(C.SM2GenerateCSR(&ctx.Context, (*C.char)(unsafe.Pointer(&country_name[0])),
		(*C.char)(unsafe.Pointer(&province[0])), (*C.char)(unsafe.Pointer(&locality_name[0])),
		(*C.char)(unsafe.Pointer(&organization_name[0])), (*C.char)(unsafe.Pointer(&organization_unit_name[0])),
		(*C.char)(unsafe.Pointer(&common_name[0])), (*C.char)(unsafe.Pointer(&email[0])),
		(*C.char)(unsafe.Pointer(&challenge_password[0])),
		(*C.char)(unsafe.Pointer(&public_key[0])),
		(*C.char)(unsafe.Pointer(&private_key[0])), temp_pubkey,
		(*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.int)(unsafe.Pointer(outlen)), ConvertSM2CsrModeToC(mode)))
}

/**
 *@brief 导入证书
 *@param ctx 函数入参 - 上下文
 *@param cert_data 函数入参 - 证书数据，为der编码二进制数据
 *@param cert_data_len 函数入参 - 证书数据长度
 *@param cert_id 函数出参 - 导入后生成的证书ID，后续可通过证书ID进行访问
 *@return 0表示成功，其他值为错误码
 */
func SM2CertImport(ctx *SM2_cert_ctx_t, cert_data []byte, cert_data_len int, cert_id []byte) int {
	if ctx == nil || cert_data == nil || cert_data_len <= 0 || cert_id == nil {
		panic("invalid parameter")
	}
	cert_data = append(cert_data, 0)
	return int(C.SM2CertImport(&ctx.Context, (*C.uchar)(unsafe.Pointer(&cert_data[0])), (C.int)(cert_data_len), (*C.char)(unsafe.Pointer(&cert_id[0]))))
}

/**
 *@brief 导出证书
 *@param ctx 函数入参 - 上下文
 *@param cert_data 函数出参 - 证书数据，为der编码二进制数据
 *@param cert_data_len 函数出参 - 证书数据长度
 *@param cert_id 函数入参 - 证书ID
 *@return 0表示成功，其他值为错误码
 */
func SM2CertExport(ctx *SM2_cert_ctx_t, cert_data []byte, cert_data_len *int, cert_id []byte) int {
	if ctx == nil || cert_data == nil || cert_data_len == nil || cert_id == nil {
		panic("invalid parameter")
	}
	cert_id = append(cert_id, 0)
	return int(C.SM2CertExport(&ctx.Context, (*C.uchar)(unsafe.Pointer(&cert_data[0])),
		(*C.int)(unsafe.Pointer(cert_data_len)), (*C.char)(unsafe.Pointer(&cert_id[0]))))
}

/**
 *@brief 删除证书
 *@param ctx 函数入参 - 上下文
 *@param cert_id 函数入参 - 证书ID
 *@return 0表示成功，其他值为错误码
 */
func SM2CertDelete(ctx *SM2_cert_ctx_t, cert_id []byte) int {
	if ctx == nil || cert_id == nil {
		panic("invalid parameter")
	}
	cert_id = append(cert_id, 0)
	return int(C.SM2CertDelete(&ctx.Context, (*C.char)(unsafe.Pointer(&cert_id[0]))))
}

/**
 *@brief 根据证书ID读取公钥
 *@param ctx 函数入参 - 上下文
 *@param cert_id 函数入参 - 证书ID
 *@param public_key 函数出参 - 公钥
#return 0表示成功，其他值为错误码
*/
func SM2CertReadPublicKey(ctx *SM2_cert_ctx_t, cert_id []byte, public_key []byte) int {
	if ctx == nil || public_key == nil || cert_id == nil {
		panic("invalid parameter")
	}
	cert_id = append(cert_id, 0)
	return int(C.SM2CertReadPublicKey(&ctx.Context, (*C.char)(unsafe.Pointer(&cert_id[0])), (*C.char)(unsafe.Pointer(&public_key[0]))))
}

/**
 *@brief 根据证书ID读取有效期
 *@param ctx 函数入参 - 上下文
 *@param cert_id 函数入参 - 证书ID
 *@param valid_time 函数出参 - 有效期
#return 0表示成功，其他值为错误码
*/
func SM2CertReadValidTime(ctx *SM2_cert_ctx_t, cert_id []byte, valid_time *SM2_valid_time_t) int {
	if ctx == nil || valid_time == nil || cert_id == nil {
		panic("invalid parameter")
	}
	cert_id = append(cert_id, 0)
	return int(C.SM2CertReadValidTime(&ctx.Context, (*C.char)(unsafe.Pointer(&cert_id[0])), &valid_time.Context))
}

/**
 *@brief 判断证书是否为根证书，即是否为自签发证书
 *@param ctx 函数入参 - 上下文
 *@param cert_id 函数入参 - 证书ID
 *@param b_root_cert 函数出参 - 是否根证书，1为根证书，0为不是根证书
 *@return 0表示成功，其他值为错误码
 */
func SM2CertIsRoot(ctx *SM2_cert_ctx_t, cert_id []byte, b_root_cert *int) int {
	if ctx == nil || b_root_cert == nil || cert_id == nil {
		panic("invalid parameter")
	}
	cert_id = append(cert_id, 0)
	return int(C.SM2CertIsRoot(&ctx.Context, (*C.char)(unsafe.Pointer(&cert_id[0])), (*C.int)(unsafe.Pointer(b_root_cert))))
}

/**
*@brief 进行本地证书链的验证，请确保你所验证的证书的所有上级证书均已通过导入接口事先导入，否则会由于证书链不完整验证失败
        证书链的验证这里未验证是否过期，原因为获取的本地时间戳可能不准。
        如果调用方需要验证是否过期。可通过SM2CertReadValidTime接口读取有效期并与可信的时间戳进行对比。
        证书链的验证这里未验证是否吊销。
        如果验证方需要验证是否吊销。可自行前往CA请求证书状态。
*@param ctx 函数入参 - 上下文
*@param cert_id 函数入参 - 证书ID
*@return 0表示成功，其他值为错误码
*/
func SM2VerifyCertChain(ctx *SM2_cert_ctx_t, cert_id []byte) int {
	if ctx == nil || cert_id == nil {
		panic("invalid parameter")
	}
	cert_id = append(cert_id, 0)
	return int(C.SM2VerifyCertChain(&ctx.Context, (*C.char)(unsafe.Pointer(&cert_id[0]))))
}

/**
 * @brief 根据csr生成证书
 * @param ctx input 上下文
 * @param der_csr input cert request der
 * @param der_csr_len input cert request der len
 * @param der_ca input ca cert der
 * @param der_ca_len imput  ca len
 * @param serialNumber input 序列号
 * @param prikey input ca 私钥
 * @param valid_days input 有效天数
 * @param isSign input 是否签名用证书
 * @param pcertder output cert der
 * @param derlen input&output input buf len and output cert der len
 * @return int --0-- successful
 */

func SM2CertGenerate(ctx *SM2_cert_ctx_t, der_csr []byte, der_csr_len int,
	der_ca []byte, der_ca_len int, serialNumber []byte, prikey []byte, valid_days int,
	isSign int, pcertder []byte, derlen *int) int {
	if ctx == nil || der_csr == nil || der_csr_len <= 0 || der_ca == nil ||
		der_ca_len <= 0 || prikey == nil || pcertder == nil || derlen == nil {
		panic("invalid parameter")
	}
	prikey = append(prikey, 0)
	der_csr = append(der_csr, 0)
	der_ca = append(der_ca, 0)
	serialNumber = append(serialNumber, 0)
	prikey = append(prikey, 0)
	return int(C.SM2CertGenerate(&ctx.Context, (*C.uchar)(unsafe.Pointer(&der_csr[0])), (C.int)(der_csr_len),
		(*C.uchar)(unsafe.Pointer(&der_ca[0])), (C.int)(der_ca_len), (*C.char)(unsafe.Pointer(&serialNumber[0])),
		(*C.char)(unsafe.Pointer(&prikey[0])), (C.int)(valid_days), (C.int)(isSign), (*C.uchar)(unsafe.Pointer(&pcertder[0])), (*C.int)(unsafe.Pointer(derlen))))
}

/**
 *@brief 将DER编码格式的CSR转换为PEM格式
 *@param der 函数入参 - der CSR数据
 *@param derlen 函数入参 - der CSR数据的长度
 *@param pem 函数出参 - pem格式的CSR数据
 *@return 0表示成功，其他值为错误码
 */
func SM2CSRConvertDER2PEM(der []byte, derlen int, pem []byte) int {
	if der == nil || pem == nil || derlen <= 0 {
		panic("invalid parameter")
	}
	der = append(der, 0)
	return int(C.SM2CSRConvertDER2PEM((*C.uchar)(unsafe.Pointer(&der[0])), (C.size_t)(derlen), (*C.char)(unsafe.Pointer(&pem[0]))))
}

/**
 *@brief 将PEM编码格式的CSR转换为DER格式
 *@param pem 函数入参 - der CSR数据
 *@param der 函数入参 - der CSR数据的长度
 *@param derlen 函数出参 - pem格式的CSR数据
 *@return 0表示成功，其他值为错误码
 */
func SM2CSRConvertPEM2DER(pem []byte, der []byte, derlen *int) int {
	if der == nil || pem == nil || derlen == nil {
		panic("invalid parameter")
	}
	pem = append(pem, 0)
	return int(C.SM2CSRConvertPEM2DER((*C.char)(unsafe.Pointer(&pem[0])), (*C.uchar)(unsafe.Pointer(&der[0])), (*C.size_t)(unsafe.Pointer(derlen))))
}

/**
 *@brief 将DER编码格式的证书转换为PEM格式
 *@param der 函数入参 - der CSR数据
 *@param derlen 函数入参 - der CSR数据的长度
 *@param pem 函数出参 - pem格式的CSR数据
 *@return 0表示成功，其他值为错误码
 */
func SM2CRTConvertDER2PEM(der []byte, derlen int, pem []byte) int {
	if der == nil || pem == nil || derlen <= 0 {
		panic("invalid parameter")
	}
	der = append(der, 0)
	return int(C.SM2CRTConvertDER2PEM((*C.uchar)(unsafe.Pointer(&der[0])), (C.size_t)(derlen), (*C.char)(unsafe.Pointer(&pem[0]))))
}

/**
 *@brief 将PEM编码格式的证书转换为DER格式
 *@param pem 函数入参 - der CSR数据
 *@param der 函数入参 - der CSR数据的长度
 *@param derlen 函数出参 - pem格式的CSR数据
 *@return 0表示成功，其他值为错误码
 */
func SM2CRTConvertPEM2DER(pem []byte, der []byte, derlen *int) int {
	if der == nil || pem == nil || derlen == nil {
		panic("invalid parameter")
	}
	pem = append(pem, 0)
	return int(C.SM2CRTConvertPEM2DER((*C.char)(unsafe.Pointer(&pem[0])), (*C.uchar)(unsafe.Pointer(&der[0])), (*C.size_t)(unsafe.Pointer(derlen))))
}

/**
 * @brief der 2 pem for prikey
 * @param der der格式私钥
 * @param derlen der格式私钥长度
 * @param pem pem格式私钥
 * @return int   --0-- successful
 */
func SM2PrikeyConvertDER2PEM(der []byte, derlen int, pem []byte) int {
	if der == nil || pem == nil || derlen <= 0 {
		panic("invalid parameter")
	}
	der = append(der, 0)
	return int(C.SM2PrikeyConvertDER2PEM((*C.uchar)(unsafe.Pointer(&der[0])), (C.size_t)(derlen), (*C.char)(unsafe.Pointer(&pem[0]))))
}

/**
 * @brief pem 2 der for prikeypem
 * @param pem  pem格式私钥
 * @param der der格式私钥
 * @param derlen  der格式私钥长度
 * @return int  --0-- successful
 */
func SM2PrikeyDerGetFromPem(pem []byte, der []byte, derlen *int) int {
	if der == nil || pem == nil || derlen == nil {
		panic("invalid parameter")
	}
	pem = append(pem, 0)
	return int(C.SM2PrikeyDerGetFromPem((*C.char)(unsafe.Pointer(&pem[0])), (*C.uchar)(unsafe.Pointer(&der[0])), (*C.size_t)(unsafe.Pointer(derlen))))
}

/**
 * @brief 从key der 中获取prikey 字符串,未压缩格式
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式私钥
 * @param der_len  der格式私钥长度
 * @param prikeystr  私钥字符串
 * @return int int --0-- successful
 */
func SM2PrikeyStrGetFromDer(ctx *SM2_cert_ctx_t, der []byte, derlen int, prikeystr []byte) int {
	if ctx == nil || der == nil || prikeystr == nil || derlen <= 0 {
		panic("invalid parameter")
	}
	der = append(der, 0)
	return int(C.SM2PrikeyStrGetFromDer(&ctx.Context, (*C.uchar)(unsafe.Pointer(&der[0])), (C.int)(derlen), (*C.char)(unsafe.Pointer(&prikeystr[0]))))
}

/**
 * @brief der to str for prikey
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式私钥
 * @param der_len der格式私钥长度
 * @param prikeystr 私钥字符串
 * @return int  --0-- successful
 */
func SM2PrikeyStr2Der(ctx *SM2_cert_ctx_t, der []byte, derlen *int, prikeystr []byte) int {
	if ctx == nil || der == nil || prikeystr == nil || derlen == nil {
		panic("invalid parameter")
	}
	return int(C.SM2PrikeyStr2Der(&ctx.Context, (*C.uchar)(unsafe.Pointer(&der[0])), (*C.int)(unsafe.Pointer(derlen)), (*C.char)(unsafe.Pointer(&prikeystr[0]))))
}

/**
 * @brief der 2 pem for pubkey
 * @param der der格式公钥
 * @param derlen der格式公钥长度
 * @param pem pem格式公钥
 * @return int   --0-- successful
 */
func SM2PubkeyConvertDER2PEM(der []byte, derlen int, pem []byte) int {
	if der == nil || pem == nil || derlen <= 0 {
		panic("invalid parameter")
	}
	der = append(der, 0)
	return int(C.SM2PubkeyConvertDER2PEM((*C.uchar)(unsafe.Pointer(&der[0])), (C.size_t)(derlen), (*C.char)(unsafe.Pointer(&pem[0]))))
}

/**
 * @brief pem 2 der for pubkeypem
 * @param pem  pem格式公钥
 * @param der der格式公钥
 * @param derlen  der格式公钥长度
 * @return int  --0-- successful
 */
func SM2PubkeyDerGetFromPem(pem []byte, der []byte, derlen *int) int {
	if pem == nil || der == nil || derlen == nil {
		panic("invalid parameter")
	}
	pem = append(pem, 0)
	return int(C.SM2PubkeyDerGetFromPem((*C.char)(unsafe.Pointer(&pem[0])), (*C.uchar)(unsafe.Pointer(&der[0])), (*C.size_t)(unsafe.Pointer(derlen))))
}

/**
 * @brief 从key der 中获取pubkey 字符串,未压缩格式
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式公钥
 * @param der_len  der格式公钥长度
 * @param pubkeystr  公钥字符串
 * @return int int --0-- successful
 */
func SM2PubkeyStrGetFromDer(ctx *SM2_cert_ctx_t, der []byte, derlen int, pubkeystr []byte) int {
	if ctx == nil || der == nil || pubkeystr == nil || derlen <= 0 {
		panic("invalid parameter")
	}
	der = append(der, 0)
	return int(C.SM2PubkeyStrGetFromDer(&ctx.Context, (*C.uchar)(unsafe.Pointer(&der[0])), (C.int)(derlen), (*C.char)(unsafe.Pointer(&pubkeystr[0]))))
}

/**
 * @brief der to str for pubkey
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式公钥
 * @param der_len der格式公钥长度
 * @param pubkeystr 公钥字符串
 * @return int  --0-- successful
 */
func SM2PubkeyStr2Der(ctx *SM2_cert_ctx_t, der []byte, derlen *int, pubkeystr []byte) int {
	if ctx == nil || der == nil || pubkeystr == nil || derlen == nil {
		panic("invalid parameter")
	}
	return int(C.SM2PubkeyStr2Der(&ctx.Context, (*C.uchar)(unsafe.Pointer(&der[0])), (*C.int)(unsafe.Pointer(derlen)), (*C.char)(unsafe.Pointer(&pubkeystr[0]))))
}

/**
 * @brief 检查csr是否sm2 pubkey，并输出pubkey
 *
 * @param ctx input 上下文
 * @param der input csr der data
 * @param derlen input csr der data len
 * @param pubkey output !NULL时输出公钥buf，size > 130; NULL时，不输出
 * @return 0表示成功，其他值为错误码
 */
func SM2CSRCheckSm2Pubkey(ctx *SM2_cert_ctx_t, der []byte, derlen int, pubkey []byte) int {
	if ctx == nil || der == nil || pubkey == nil || derlen <= 0 {
		panic("invalid parameter")
	}
	der = append(der, 0)
	return int(C.SM2CSRCheckSm2Pubkey(&ctx.Context, (*C.uchar)(unsafe.Pointer(&der[0])), (C.int)(derlen), (*C.char)(unsafe.Pointer(&pubkey[0]))))
}

/**
 * @brief 根据csr生成证书，其中der_csr、serialNumber、prikey、valid_days是必须项。
 * @param ctx input 上下文
 * @param items input 生成证书输入的项目;
 *            subject,issue等中的项目，如果已有情况下回增加一个，例如:description
 * @param itemcount input 参数items数组长度
 * @param pcertder output cert der buf；not null
 * @param derlen input&output input buf len and output cert der len
 * @return int --0-- successful
 */
func SM2CertGenerate2(ctx *SM2_cert_ctx_t, items []TstTCSMItem, itemcount int, pcertder []byte, derlen *int) int {
	if ctx == nil || pcertder == nil || derlen == nil || items == nil {
		panic("invalid parameter")
	}
	//wrapper := make([]C.TstTCSMItem, itemcount)
	wrapper := make([]C.TstTCSMItemWrapper, itemcount)
	for i := 0; i < itemcount; i++ {

		wrapper[i].valuelen = (C.int)(items[i].Valuelen)
		wrapper[i].itemType = (C.int)(items[i].ItemType)
		if items[i].Value != nil {
			temp := (*C.uchar)(unsafe.Pointer(&items[i].Value[0]))
			//如何防止内存泄漏
			C.createAndCopy(&wrapper[i].value, temp, (C.int)(items[i].Valuelen))
		}
	}
	ret := int(C.SM2CertGenerate2(&ctx.Context, (*C.TstTCSMItem)(unsafe.Pointer(&wrapper[0])), (C.int)(itemcount), (*C.uchar)(unsafe.Pointer(&pcertder[0])), (*C.int)(unsafe.Pointer(derlen))))
	//释放内存
	for i := 0; i < itemcount; i++ {
		if items[i].Value != nil {
			C.crelease(wrapper[i].value)
		}
	}
	return ret

}

/**
 * @brief 读取证书信息
 *
 * @param ctx input 上下文
 * @param der input 证书der
 * @param certlen input 证书der数据长度
 * @param itemid input 欲获取项目id，例如：TYPE_READCERTITEM_**
 * @param itemcount output 相同项目数量
 * @param outstr 输出数据字符串，格式：*itemcount == 1时为字符串格式数据（byte[]格式转为hex）;*itemcount > 1时为数组格式：["字符串","字符串"]
 * @param len input&output input buf len and output str len
 * @return  int 0-- successful
 */
func SM2GetCertItem(ctx *SM2_cert_ctx_t, der []byte, certlen int, itemid int, itemcount *int, outstr []byte, len *int) int {
	return int(C.SM2GetCertItem(&ctx.Context, (*C.uchar)(unsafe.Pointer(&der[0])), (C.int)(certlen), (C.int)(itemid),
		(*C.int)(unsafe.Pointer(itemcount)), (*C.char)(unsafe.Pointer(&outstr[0])), (*C.int)(unsafe.Pointer(len))))
}

/**
 * @brief 解析prikey pem格式字符串并输出私钥等信息
 *        目前支持格式EC PRIVATE KEY和PRIVATE KEY
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param keypem  [input] pem格式字符串
 * @param outPrikey [output] 输出的私钥字串，hex格式
 * @param extdata  [output] 扩展输出，备用
 * @return int 0-- successful
 */
func SM2ParsePrikeyPem(ctx *SM2_cert_ctx_t, keypem []byte, outPriKey []byte, extdata []string) int {
	arg := make([](*C.char), 0)
	//l := len(extdata)
	for i := range extdata {
		char := C.CString(extdata[i])
		defer C.free(unsafe.Pointer(char))
		strptr := (*C.char)(unsafe.Pointer(char))
		arg = append(arg, strptr)
	}
	return int(C.SM2ParsePrikeyPem(&ctx.Context, (*C.char)(unsafe.Pointer(&keypem[0])),
		(*C.char)(unsafe.Pointer(&outPriKey[0])), (**C.char)(unsafe.Pointer(&arg[0]))))

}
