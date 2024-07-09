package sm

/*
#cgo windows,arm64 LDFLAGS: ${SRCDIR}/lib/windows_arm/libTencentSM.dll ${SRCDIR}/lib/windows_arm/libgmp.a
#cgo windows,!arm64 LDFLAGS: ${SRCDIR}/lib/windows/libTencentSM.dll ${SRCDIR}/lib/windows/libgmp.a
#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/darwin_arm/libTencentSM.a ${SRCDIR}/lib/darwin_arm/libgmp.a
#cgo darwin,!arm64 LDFLAGS: ${SRCDIR}/lib/darwin/libTencentSM.a ${SRCDIR}/lib/darwin/libgmp.a
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/linux_arm/libTencentSM.a ${SRCDIR}/lib/linux_arm/libgmp.a
#cgo linux,!arm64 LDFLAGS: ${SRCDIR}/lib/linux/libTencentSM.a ${SRCDIR}/lib/linux/libgmp.a

#cgo CFLAGS: -g -O2 -I${SRCDIR}/include/

#include "sm.h"
#include <stdlib.h>
*/
import "C"
import (
	"sync"
	"unsafe"
)

const SM3_BLOCK_SIZE int = 64

const SM3_DIGEST_LENGTH int = 32
const SM3_HMAC_SIZE int = SM3_DIGEST_LENGTH

//-----------------------------------证书生成项及类型定义----------------------------------
const TYPE_CERTITEM_NECESSARY_MAX int = 10        //必要参数type最大值
const TYPE_CERTITEM_NECESSARY_COUNT int = 4       //必要参数type数量
const TYPE_CERTITEM_CSR_DER int = 1               //(unsigned char*)   //cert request der
const TYPE_CERTITEM_SERIALNUMBER int = 2          //(char*)            //序列号
const TYPE_CERTITEM_SIGN_PRIKEY int = 3           //(char*)            //私钥
const TYPE_CERTITEM_VALID_DAYS int = 4            //(int*)            //有效天数
const TYPE_CERTITEM_USEAGE int = 20               //(int*)            //签名证书(1), 加解密证书(2), 两者都有(3)
const TYPE_CERTITEM_CA_DER int = 21               //(unsigned char*)   //ca cert der,自签名时不需要
const TYPE_CERTITEM_SUBJECT_BEGIN int = 100       //subject下项目type begin
const TYPE_CERTITEM_SUBJECT_END int = 199         //subject下项目type end
const TYPE_CERTITEM_SUBJECT_DESCRIPTION int = 101 //(char*)            //subject下description项目
const TYPE_CERTITEM_SUBJECT_CN int = 102          //(char*)            //subject下commonname项目

//---------------------------证书读取项目及定义-----------------------------------
const TYPE_READCERTITEM_PUBKEY int = 2002                                           //(char*)            //公钥
const TYPE_READCERTITEM_SERIALNUMBER int = TYPE_CERTITEM_SERIALNUMBER               //(char*)            //序列号
const TYPE_READCERTITEM_SUBJECT_DESCRIPTION int = TYPE_CERTITEM_SUBJECT_DESCRIPTION //(char*)            //subject下description项目
const TYPE_READCERTITEM_SUBJECT_CN int = TYPE_CERTITEM_SUBJECT_CN                   //(char*)            //subject下commonname项目

var ContextPool *sync.Pool

var IsContextPoolEnable bool = false

//SM2上下文
type SM2_ctx_t struct {
	Context C.sm2_ctx_t
}

//SM3上下文
type SM3_ctx_t struct {
	Context C.sm3_ctx_t
}

//SM3 HMAC上下文
type HmacSm3Ctx struct {
	Context *C.TstHmacSm3Ctx
}

//SM2证书上下文
type SM2_cert_ctx_t struct {
	Context C.sm2_cert_ctx_t
}

//SM2证书时间结构体
type SM2_valid_time_t struct {
	Context C.sm2_valid_time_t
}

//Sm2签名模式
type SM2SignMode int

//
type TstTCSMItem struct {
	ItemType int
	Value    []byte
	Valuelen int
}

const (
	SM2SignMode_RS_ASN1 SM2SignMode = iota
	SM2SignMode_RS
)

func ConvertSMSignModeToC(mode SM2SignMode) C.SM2SignMode {
	var ret C.SM2SignMode
	switch mode {
	case SM2SignMode_RS_ASN1:
		ret = C.SM2SignMode_RS_ASN1
	case SM2SignMode_RS:
		ret = C.SM2SignMode_RS
	default:
		ret = C.SM2SignMode_RS_ASN1
	}
	return ret
}

//sm2加密模式
type SM2CipherMode int

const (
	SM2CipherMode_C1C3C2_ASN1 SM2CipherMode = iota
	SM2CipherMode_C1C3C2
	SM2CipherMode_C1C2C3_ASN1
	SM2CipherMode_C1C2C3
	SM2CipherMode_04C1C3C2
	SM2CipherMode_04C1C2C3
)

func ConvertSMCipherModeToC(mode SM2CipherMode) C.SM2CipherMode {
	var ret C.SM2CipherMode
	switch mode {
	case SM2CipherMode_C1C3C2_ASN1:
		ret = C.SM2CipherMode_C1C3C2_ASN1
	case SM2CipherMode_C1C3C2:
		ret = C.SM2CipherMode_C1C3C2
	case SM2CipherMode_C1C2C3_ASN1:
		ret = C.SM2CipherMode_C1C2C3_ASN1
	case SM2CipherMode_C1C2C3:
		ret = C.SM2CipherMode_C1C2C3
	default:
		ret = C.SM2SignMode_RS_ASN1
	}
	return ret
}

//证书模式
type SM2CSRMode int

const (
	SM2CSRMode_Single SM2CSRMode = iota
	SM2CSRMode_Double
)

func ConvertSM2CsrModeToC(mode SM2CSRMode) C.SM2CSRMode {
	var ret C.SM2CSRMode
	switch mode {
	case SM2CSRMode_Single:
		ret = C.SM2CSRMode_Single
	case SM2CSRMode_Double:
		ret = C.SM2CSRMode_Double
	default:
		ret = C.SM2CSRMode_Single
	}
	return ret
}

/**
 *@brief 获取当前sdk版本
 */
func Version() string {
	return C.GoString(C.version())
}

/**
 *@brief appid 认证
 *@param appid appid or module id
 *@param token from tencentsm.oa.com
 *@return int --0-- successful
 */
func InitTencentSM(appid []byte, token []byte) int {
	tempAppid := append(appid, 0)
	tempToken := append(token, 0)
	// fmt.Printf("final app=0x%x", string(tempAppid))
	// fmt.Printf("final token=%s", tempToken)
	return int(C.initTencentSM((*C.char)(unsafe.Pointer(&tempAppid[0])), (*C.char)(unsafe.Pointer(&tempToken[0]))))
}

/**
 * @brief 用证书和appid进行初始化库操作
 *
 * @param appid 应用id，在官网上输入的同一个
 * @param cert 官网下发的证书
 * @return int 0 -- successful
 */
func InitTencentSMWithCert(appid []byte, bundleid []byte, cert []byte) int {
	if appid == nil || len(appid) <= 0 || cert == nil || len(cert) <= 0 {
		panic("invalid parameter")
	}
	if bundleid != nil {
		return int(C.initTencentSMWithCert((*C.char)(unsafe.Pointer(&appid[0])), (*C.char)(unsafe.Pointer(&bundleid[0])), (*C.char)(unsafe.Pointer(&cert[0]))))

	}
	return int(C.initTencentSMWithCert((*C.char)(unsafe.Pointer(&appid[0])), (*C.char)(unsafe.Pointer(nil)), (*C.char)(unsafe.Pointer(&cert[0]))))
}

func SetContextPoolEnable(size int) {
	if size < 3 {
		size = 3
	}
	IsContextPoolEnable = true

	ContextPool = &sync.Pool{
		New: func() interface{} {
			ret := new(SM2_ctx_t)
			C.SM2InitCtx(&ret.Context)
			return ret
		},
	}
	for i := 0; i < size; i++ {
		var ctx = new(SM2_ctx_t)
		C.SM2InitCtx(&ctx.Context)
		ContextPool.Put(ctx)
	}

}

func ClearContextPool() {
	IsContextPoolEnable = false
}
