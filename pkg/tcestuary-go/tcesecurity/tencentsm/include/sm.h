/*
Copyright 2019, Tencent Technology (Shenzhen) Co Ltd
Description: This file is part of the Tencent SM (Pro Version) Library.
*/

#ifndef HEADER_SM_H
#define HEADER_SM_H

#ifdef OS_ANDROID
#ifdef DEBUG
#include <android/log.h>
#endif /* DEBUG */
#endif /* OS_ANDROID */

#include <stdint.h>
#include <stddef.h>
#include <time.h>

#ifdef __cplusplus
extern "C" {
#endif
#undef SMLib_EXPORT
#if defined (_WIN32) && !defined (_WIN_STATIC)
#if defined(SMLib_EXPORTS)
#define  SMLib_EXPORT __declspec(dllexport)
#else
#define  SMLib_EXPORT __declspec(dllimport)
#endif
#else /* defined (_WIN32) */
#define SMLib_EXPORT
#endif

#define SM3_DIGEST_LENGTH  32
#define SM3_HMAC_SIZE    (SM3_DIGEST_LENGTH)
  
#define TENCENTSM_VERSION ("Pro_1.7.6")

#define SM2_PUBKEY_LEN 130
#define SM2_PRIKEY_LEN 64

typedef struct {
  void *group;
  void *generator;
  void *jcb_generator;
  void *jcb_compute_var;
  void *bn_vars;
  void *ec_vars;
  void *pre_comp_g;
  void *pre_comp_p;
  void *rand_ctx;
  void *pubkey_x;
  void *pubkey_y;
  void *sign_random;
} sm2_ctx_t;

typedef struct {
  uint32_t digest[8];
  int nblocks;
  unsigned char block[64];
  int num;
} sm3_ctx_t;
typedef struct stHmacSm3Ctx TstHmacSm3Ctx;

/* ---------------------------------------------------------------- 以下为SM2/SM3/SM4通用能力接口 ---------------------------------------------------------------- */

SMLib_EXPORT const char* version(void);

/**
 SDK初始化接口v1：使用token来初始SDK。
 该接口在整个进程生命周期只需调用一次，同时接口也支持重入。
 使用SDK提供的密码功能需首先调用SDK初始化接口(v1或v2版本接口仅需调用其中一个)。
 
 @param appid appid or module id
 @param token from tencentsm
 @return int --0-- successful
 */
SMLib_EXPORT int initTencentSM(const char* appid, const char* token);
/**
 SDK初始化接口v2：使用证书来初始化SDK。
 该接口在整个进程生命周期只需调用一次，同时接口也支持重入。
 使用SDK提供的密码功能需首先调用SDK初始化接口(v1或v2版本接口仅需调用其中一个)。
 
 * @param appid 应用id，在官网上输入的同一个
 * @param cert 官网下发的证书
 * @return int 0 -- successful
 */
SMLib_EXPORT int initTencentSMWithCert(const char* appid, const char* bundleid, const char* cert);

/**
 SM2上下文结构体的大小
 */
SMLib_EXPORT int SM2CtxSize(void);

/**
 SM2上下文ctx初始化接口。
 
 在使用SM2接口进行密钥生成、加密解密、签名验签之前，必须调用该接口或SM2InitCtxWithPubKey( )。
 该接口只需调用一次，在后续的密钥生成、加密解密、签名验签运算中，无需再次调用该接口。
 该接口所涉及ctx不是线程安全的，如需支持多线程，可对涉及ctx参数的接口调用加锁以保证线程安全，或不同线程使用不同的ctx。
 如需支持多线程，推荐在不同线程使用不同的ctx，以防止加锁带来的性能损耗，ctx在SM2FreeCtx之前的整个线程生命周期可复用。
 
 @param ctx  函数出参 - 上下文
 @return int --0-- successful
 */
SMLib_EXPORT int SM2InitCtx(sm2_ctx_t *ctx);

/**
 SM2上下文ctx初始化接口。
 
 该接口内部基于输入的公钥进行了预处理，针对该公钥的密码运算性能将得到较大提升。
 对于某个固定公钥需要进行密集的密码运算场景可使用该接口来替代SM2InitCtx接口进行ctx的初始化，以获得性能提升。
 该接口只需调用一次，在后续的密钥生成、加密解密、签名验签运算中，无需再次调用该接口。
 该接口初始化的ctx仍然可以用于其他公钥相关的密码运算，但是其性能与SM2InitCtx初始化的ctx性能一致。
 该接口所涉及ctx不是线程安全的，如需支持多线程，可对涉及ctx参数的接口调用加锁以保证线程安全，或不同线程使用不同的ctx。
 如需支持多线程，推荐在不同线程使用不同的ctx，以防止加锁带来的性能损耗，ctx在SM2FreeCtx之前的整个线程生命周期可复用。
 
 @param ctx  函数出参 - 上下文
 @param pubkey  函数入参 - 公钥
 @return int --0-- successful
*/
SMLib_EXPORT int SM2InitCtxWithPubKey(sm2_ctx_t *ctx,const char* pubkey);

/**
 使用完SM2算法后，必须调用SM2FreeCtx函数释放相关数据。
 如果ctx需要在整个线程生命周期复用的话，可在线程退出前释放。
 
 @param ctx  函数入参 - 上下文
 @return int --0-- successful
 */
SMLib_EXPORT int SM2FreeCtx(sm2_ctx_t *ctx);

/**
 生成私钥
 @param ctx  函数入参 - 上下文
 @param out  函数出参 - 私钥，私钥实际上为256bit的大整数，这里输出的为256bit二进制内容Hex后的ASCII编码的可见字符串，长度为64字节，为保证字符串的结束符0，out至少需分配65字节空间。
 @return  0表示成功，其他值为错误码
*/
SMLib_EXPORT int generatePrivateKey(sm2_ctx_t *ctx, char *out);

/**
 根据私钥生成对应公钥，
 @param ctx 函数入参 - 上下文
 @param privateKey 函数入参 - 私钥，私钥实际上为256bit的大整数，这里输出的为256bit二进制内容Hex后的ASCII编码的可见字符串，长度为64字节
 @param outPubKey 函数出参 - 公钥，公钥格式为04 | X | Y，其中X和Y为256bit大整数，这里输出的为04 | X | Y的二进制内容Hex后的ASCII编码的可见字符串，长度为130字节，为保证字符串的结束符0，outPubKey至少需分配131字节空间。
 @return  0表示成功，其他值为错误码
*/
SMLib_EXPORT int generatePublicKey(sm2_ctx_t *ctx, const char *privateKey, char *outPubKey);

/**
 生成公私钥对
 @param ctx 函数入参 - 上下文
 @param outPriKey 函数出参 - 私钥，私钥实际上为256bit的大整数，这里输出的为256bit二进制内容Hex后的ASCII编码的可见字符串，长度为64字节，为保证字符串的结束符0，outPubKey至少需分配65字节空间。
 @param outPubKey 函数出参 - 公钥，公钥格式为04 | X | Y，其中X和Y为256bit大整数，这里输出的为04 | X | Y的二进制内容Hex后的ASCII编码的可见字符串，长度为130字节，为保证字符串的结束符0，outPubKey至少需分配131字节空间。
 @return  0表示成功，其他值为错误码
 */
SMLib_EXPORT int generateKeyPair(sm2_ctx_t *ctx, char *outPriKey, char *outPubKey);

/**
 SM2非对称加解密算法，加密
 @param ctx 函数入参 - 上下文
 @param in  函数入参 - 待加密消息
 @param inlen  函数入参 - 消息长度(字节单位)
 @param strPubKey  函数入参 - 公钥
 @param pubkeyLen  函数入参 - 公钥长度
 @param out  函数出参 - 密文 - 应当为out分配的内存长度遵循以下规则：密文长度 = 明文长度 + 96 + ASN1编码增量，其中ASN1编码增量长度不定，为简单起见，可直接分配 密文长度 = 明文长度 + 200，此长度可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为输出密文的实际长度
 @return  0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2Encrypt(sm2_ctx_t *ctx, const unsigned char *in, size_t inlen, const char *strPubKey, size_t pubkeyLen, unsigned char *out, size_t *outlen);

/**
 SM2非对称加解密算法，解密
 @param ctx  函数入参 - 上下文
 @param in  函数入参 - 待解密密文
 @param inlen  函数入参 - 密文长度(字节单位)
 @param strPriKey  函数入参 - 私钥
 @param prikeyLen  函数入参 - 私钥长度
 @param out  函数出参 - 明文  - 为out分配的内存长度与inlen一致可保证安全，明文长度一定小于密文长度
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为输出明文的实际长度
 @return  0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2Decrypt(sm2_ctx_t *ctx, const unsigned char *in, size_t inlen, const char *strPriKey, size_t prikeyLen,unsigned char *out, size_t *outlen);

/**
 SM2签名验签算法，签名
 @param ctx 函数入参 - 上下文
 @param msg 函数入参 - 待签名消息
 @param msglen 函数入参 - 待签名消息长度
 @param id 函数入参 - 用户ID(作用是加入到签名hash中，对于传入值无特殊要求)
 @param idlen 函数入参 - 用户ID长度
 @param strPubKey 函数入参 - 公钥(作用是加入到签名hash中)
 @param pubkeyLen 函数入参 - 公钥长度
 @param strPriKey 函数入参 - 私钥
 @param prikeyLen 函数入参 - 私钥长度
 @param sig 函数出参 - 签名结果 - 为sig分配的内存长度为100可保证内存安全，实际长度为 64 + ASN1编码增量长度
 @param siglen 函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*siglen置为sig指针所指向内存的分配大小，函数返回后*siglen将被置为签名的实际长度
 */
SMLib_EXPORT int SM2Sign(sm2_ctx_t *ctx, const unsigned char *msg, size_t msglen, const char *id, size_t idlen, const char *strPubKey, size_t pubkeyLen, const char *strPriKey, size_t prikeyLen,unsigned char *sig, size_t *siglen);

/**
 SM2签名验签算法，验签
 @param ctx 函数入参 - 上下文
 @param msg 函数入参 - 待验签内容
 @param msglen 函数入参 - 待验签内容长度
 @param id 函数入参 - 用户ID
 @param idlen 函数入参 - 用户ID长度
 @param sig 函数入参 - 签名结果
 @param siglen 函数入参 - 签名结果长度
 @param strPubKey 函数入参 - 公钥
 @param pubkeyLen 函数入参 - 公钥长度
 @return 0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2Verify(sm2_ctx_t *ctx, const unsigned char *msg, size_t msglen,const char *id, size_t idlen, const unsigned char *sig, size_t siglen, const char *strPubKey, size_t pubkeyLen);

/** 
 SM2密钥交换算法接口，计算共享密钥。
 SM2密钥交换算法需要AB两方参与，其中一方为发起方，一方为响应方。
 我方为A，持有一对自身的公私钥和一对自身的临时公私钥。
 对方为B，持有对方的公钥和临时公钥。
 @param ctx 函数入参 - 上下文
 @param strPriKey_A 函数入参 -  A的私钥
 @param strPubKey_A 函数入参 -  A的公钥
 @param strTempPriKey_A 函数入参 -  A的临时私钥
 @param strTempPubKey_A 函数入参 -  A的临时公钥
 @param strPubKey_B 函数入参 -  B的公钥
 @param strTempPubKey_B 函数入参 -  B的临时公钥
 @param id_A 函数入参 -  A的ID，可为NULL，当使用NULL时，将使用国标中的默认值，字符串“1234567812345678”
 @param idlen_A 函数入参 -  A的ID长度，当ID为空时，填0
 @param id_B 函数入参 -  B的ID，可为NULL，当使用NULL时，将使用国标中的默认值，字符串“1234567812345678”
 @param idlen_B 函数入参 -  B的ID长度，当ID为空时，填0
 @param klen 函数入参 -  需要计算的共享密钥长度
 @param sharedKey 函数出参 - 共享密钥
 @param A_is_initiator 函数出参 - A是否为发起方，即我方是否为发起方，如果我方为发起方则填1，否则为0
 @return 0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2CalculateSharedKey(sm2_ctx_t *ctx,const char* strPriKey_A,const char* strPubKey_A,const char* strTempPriKey_A,const char* strTempPubKey_A,const char* strPubKey_B,const char* strTempPubKey_B,const char *id_A, size_t idlen_A,const char *id_B, size_t idlen_B,size_t klen,unsigned char* sharedKey,int A_is_initiator);
  
/**
 SM3上下文结构体的大小
 */
SMLib_EXPORT int SM3CtxSize(void);

/**
 SM3 hash算法，3个接口用法与OpenSSL的MD5算法的接口保持一致。
 digest至少需要分配32字节
 */
SMLib_EXPORT int SM3Init(sm3_ctx_t *ctx);
SMLib_EXPORT int SM3Update(sm3_ctx_t *ctx, const unsigned char* data, size_t data_len);
SMLib_EXPORT int SM3Final(sm3_ctx_t *ctx, unsigned char *digest);

/**
 SM3 hash算法， 内部依次调用了init update和final三个接口
 */
SMLib_EXPORT int SM3(const unsigned char *data, size_t datalen, unsigned char *digest);

/**
 * @brief 基于sm3算法计算HMAC值 ctx init
 * @param key HMAC用的秘钥
 * @param key_len 秘钥长度
 * @return 0 -- OK
 */
SMLib_EXPORT TstHmacSm3Ctx* SM3_HMAC_Init(const unsigned char *key, size_t key_len);
/**
 * @brief 基于sm3算法计算HMAC值 update数据
 * @param ctx hmac上下文结构指针
 * @param data 做HMAC计算的数据
 * @param data_len 数据长度
 * @return 0 -- OK
 */
SMLib_EXPORT int SM3_HMAC_Update(TstHmacSm3Ctx *ctx,const unsigned char *data, size_t data_len);
/**
 * @brief 基于sm3算法计算HMAC值 最终计算HMAC值
 * @param ctx hmac上下文结构指针
 * @param mac 输出的HMAC字节码
 * @return 0 -- OK
 */
SMLib_EXPORT int SM3_HMAC_Final(TstHmacSm3Ctx *ctx, unsigned char mac[SM3_HMAC_SIZE]);
/**
 * @brief 基于sm3算法计算HMAC值
 * @param data 做HMAC计算的数据
 * @param data_len 数据长度
 * @param key HMAC用的秘钥
 * @param key_len 秘钥长度
 * @param mac 输出的HMAC字节码
 * @return 0 -- OK
 */
SMLib_EXPORT int SM3_HMAC(const unsigned char *data, size_t data_len,const unsigned char *key, size_t key_len,unsigned char mac[SM3_HMAC_SIZE]);

/**
 密钥导出函数。sharelen不可大于1024字节。
*/
SMLib_EXPORT int SM3KDF(const unsigned char *share, size_t sharelen, unsigned char *outkey, size_t keylen);

/**
 基于SM3的PBKDF2算法接口，可用于口令哈希。

 @param msg 函数入参 - 待哈希的原始消息
 @param msglen 函数入参 - 待哈希的原始消息长度
 @param salt 函数入参 - 外部添加的盐值，可为空，建议添加动态盐以增强安全性，该函数内部也会添加静态盐
 @param saltlen 函数入参 - 外部添加的盐值长度
 @param roundtimes 函数入参 - 算法的轮次，数值越大，构建彩虹表攻击的难度越大，但数值越大其性能开销越大
 @param outbuf 函数输出- 输出长度为32字节数据
*/

SMLib_EXPORT int SM3BasedPBKDF2(const unsigned char* msg,int msglen,const unsigned char *salt,int saltlen,int roundtimes,unsigned char* outbuf);
  
/**
 生成16字节128bit的SM4 Key，也可调用该接口生成SM4 CBC模式的初始化向量iv，iv长度和key长度一致
 @param outKey  函数出参 - 16字节密钥。
 */
SMLib_EXPORT int generateSM4Key(unsigned char *outKey);
  
/**
 SM4 CBC模式对称加解密。加密，使用PKCS#7填充标准
 @param in  函数入参 - 明文
 @param inlen  函数入参 - 明文长度
 @param out  函数出参 - 密文 - 为out分配的内存长度为明文长度+16字节可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 @param iv  函数入参 - 初始化向量
 */
SMLib_EXPORT int SM4_CBC_Encrypt(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key,const unsigned char *iv);
    
/**
 SM4 CBC模式对称加解密。解密，使用PKCS#7填充标准
 @param in  函数入参 - 密文
 @param inlen  函数入参 - 密文长度
 @param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 @param iv  函数入参 - 初始化向量
 
 */
SMLib_EXPORT int SM4_CBC_Decrypt(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key,const unsigned char *iv);
    
/**
 SM4 CBC模式对称加解密。加密，无填充。请保证明文为16字节整数倍，否则加密会失败。
 @param in 函数入参 - 明文
 @param inlen 函数入参 - 明文长度
 @param out  函数出参 - 密文 - 为out分配的内存长度为明文长度可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
 @param key 函数入参 - 秘钥（128bit）
 @param iv 函数入参 - 初始化向量
 */
SMLib_EXPORT int SM4_CBC_Encrypt_NoPadding(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key, const unsigned char *iv);

/**
 SM4 CBC模式对称解密，无填充。请保证密文为16字节整数倍，否则解密会失败。
 @param in  函数入参 - 密文
 @param inlen  函数入参 - 密文长度
 @param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 @param iv  函数入参 - 初始化向量
 */
SMLib_EXPORT int SM4_CBC_Decrypt_NoPadding(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key,const unsigned char *iv);

/**
 SM4 ECB模式对称加解密。加密
 @param in  函数入参 - 明文
 @param inlen  函数入参 - 明文长度
 @param out  函数出参 - 密文 - 为out分配的内存长度为明文长度+16字节可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 */
SMLib_EXPORT int SM4_ECB_Encrypt(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key);

/**
 SM4 ECB模式对称加解密。解密
 @param in  函数入参 - 密文
 @param inlen  函数入参 - 密文长度
 @param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 */
SMLib_EXPORT int SM4_ECB_Decrypt(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key);
    
/**
 SM4 ECB模式对称加密，无填充。请保证明文为16字节整数倍，否则加密会失败。
 @param in  函数入参 - 明文
 @param inlen  函数入参 - 明文长度
 @param out  函数出参 - 密文 - 为out分配的内存长度为明文长度可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 */
SMLib_EXPORT int SM4_ECB_Encrypt_NoPadding(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key);

/**
 SM4 ECB模式对称解密，无填充。请保证密文为16字节整数倍，否则解密会失败。
 @param in  函数入参 - 密文
 @param inlen  函数入参 - 密文长度
 @param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 */
SMLib_EXPORT int SM4_ECB_Decrypt_NoPadding(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key);

/**
 SM4 CTR模式对称加解密。加密，CTR模式不需要填充。
 @param in 函数入参 - 明文
 @param inlen 函数入参 - 明文长度
 @param out  函数出参 - 密文 - 为out分配的内存长度为明文长度可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
 @param key 函数入参 - 秘钥（128bit）
 @param iv 函数入参 - 初始化向量
 */
SMLib_EXPORT int SM4_CTR_Encrypt_NoPadding(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key, const unsigned char *iv);

/**
 SM4 CTR模式对称加解密。解密，CTR模式不需要填充。
 @param in  函数入参 - 密文
 @param inlen  函数入参 - 密文长度
 @param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
 @param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
 @param key  函数入参 - 秘钥（128bit）
 @param iv  函数入参 - 初始化向量
 */
SMLib_EXPORT int SM4_CTR_Decrypt_NoPadding(const unsigned char *in, size_t inlen,unsigned char *out, size_t *outlen, const unsigned char *key,const unsigned char *iv);

/**
 SM4 GCM模式对称加解密。加密，使用PKCS7填充，实际上GCM模式可不填充，非短明文加密推荐使用SM4_GCM_Encrypt_NoPadding替代。
@param in  函数入参 - 明文
@param inlen  函数入参 - 明文长度
@param out  函数出参 - 密文 - 为out分配的内存长度为明文长度+16字节可保证安全
@param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
@param tag  函数出参 - GMAC值，即消息验证码
@param taglen  既作函数入参也作为函数出参 - GMAC长度，通常取16字节 - 入参请将*taglen置为你希望的tag长度，同时你需要保证tag所指向的内存空间足够，函数返回后*taglen将被置为实际长度，taglen如果为0将会返回失败，taglen大于16返回后将会置为16
@param key  函数入参 - 秘钥（128bit）
@param iv  函数入参 - 初始化向量,GCM模式的向量长度与CBC模式不同，不一定需要使用128bit，该接口内部默认使用了8字节
@param aad  函数入参 - 附加验证消息
@param aadlen  函数入参 - 附加验证消息长度
@return 成功为0，一般加密失败是由参数错误导致
*/
SMLib_EXPORT int SM4_GCM_Encrypt(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen, unsigned char *tag, size_t *taglen, const unsigned char *key, const unsigned char *iv,const unsigned char *aad, size_t aadlen);

/**
 SM4 GCM模式对称加解密。解密，使用PKCS7填充，实际上GCM模式可不填充。
@param in  函数入参 - 密文
@param inlen  函数入参 - 密文长度
@param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
@param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
@param tag  函数入参 - GMAC值，即消息验证码
@param taglen  函数入参 - GMAC长度，通常取16字节
@param key  函数入参 - 秘钥（128bit）
@param iv  函数入参 - 初始化向量,GCM模式的向量长度与CBC模式不同，不一定需要使用128bit，该接口内部默认使用了8字节
@param aad  函数入参 - 附加验证消息
@param aadlen  函数入参 - 附加验证消息长度
@return 成功为0，GCM的解密失败主要是tag校验失败
*/
SMLib_EXPORT int SM4_GCM_Decrypt(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen,const unsigned char *tag, size_t taglen, const unsigned char *key, const unsigned char *iv,const unsigned char *aad, size_t aadlen);

/**
SM4 GCM模式对称加解密。加密，无填充，明文长度无要求。
@param in  函数入参 - 明文
@param inlen  函数入参 - 明文长度
@param out  函数出参 - 密文 - 为out分配的内存长度为明文长度可保证安全
@param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为密文的实际长度
@param tag  函数出参 - GMAC值，即消息验证码
@param taglen  既作函数入参也作为函数出参 - GMAC长度，通常取16字节 - 入参请将*taglen置为你希望的tag长度，同时你需要保证tag所指向的内存空间足够，函数返回后*taglen将被置为实际长度，taglen如果为0将会返回失败，taglen大于16返回后将会置为16
@param key  函数入参 - 秘钥（128bit）
@param iv  函数入参 - 初始化向量,GCM模式的向量长度与CBC模式不同，不一定需要使用128bit，该接口内部默认使用了8字节
@param aad  函数入参 - 附加验证消息
@param aadlen  函数入参 - 附加验证消息长度
@return 成功为0，一般加密失败是由参数错误导致
*/
SMLib_EXPORT int SM4_GCM_Encrypt_NoPadding(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen, unsigned char *tag, size_t *taglen, const unsigned char *key, const unsigned char *iv,const unsigned char *aad, size_t aadlen);

/**
SM4 GCM模式对称加解密。解密，无填充，密文长度无要求。
@param in  函数入参 - 密文
@param inlen  函数入参 - 密文长度
@param out  函数出参 - 明文 - 为out分配的内存长度与inlen一致可保证安全
@param outlen  函数入参和出参 - 这是一个UNIX C风格的函数参数用法，入参请将*outlen置为out指针所指向内存的分配大小，函数返回后*outlen将被置为明文的实际长度
@param tag  函数入参 - GMAC值，即消息验证码
@param taglen  函数入参 - GMAC长度，通常取16字节
@param key  函数入参 - 秘钥（128bit）
@param iv  函数入参 - 初始化向量,GCM模式的向量长度与CBC模式不同，不一定需要使用128bit，该接口内部默认使用了8字节
@param aad  函数入参 - 附加验证消息
@param aadlen  函数入参 - 附加验证消息长度
@return 返回解密是否失败，GCM的解密失败主要是tag校验失败
*/
SMLib_EXPORT int SM4_GCM_Decrypt_NoPadding(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen,const unsigned char *tag, size_t taglen, const unsigned char *key, const unsigned char *iv,const unsigned char *aad, size_t aadlen);


/**
 以下接口为支持任意长度(>0)iv的GCM接口
 按照NIST SP800-38D标准实现GCM部分算法
 RFC5647标准iv推荐使用12字节，96bit
 */
SMLib_EXPORT int SM4_GCM_Encrypt_NIST_SP800_38D(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen, unsigned char *tag, size_t *taglen, const unsigned char *key, const unsigned char *iv,size_t ivlen,const unsigned char *aad, size_t aadlen);
SMLib_EXPORT int SM4_GCM_Decrypt_NIST_SP800_38D(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen,const unsigned char *tag, size_t taglen, const unsigned char *key, const unsigned char *iv,size_t ivlen,const unsigned char *aad, size_t aadlen);
SMLib_EXPORT int SM4_GCM_Encrypt_NoPadding_NIST_SP800_38D(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen, unsigned char *tag, size_t *taglen, const unsigned char *key, const unsigned char *iv,size_t ivlen,const unsigned char *aad, size_t aadlen);
SMLib_EXPORT int SM4_GCM_Decrypt_NoPadding_NIST_SP800_38D(const unsigned char *in, size_t inlen, unsigned char *out, size_t *outlen,const unsigned char *tag, size_t taglen, const unsigned char *key, const unsigned char *iv,size_t ivlen,const unsigned char *aad, size_t aadlen);

/* ---------------------------------------------------------------- 以下为SM2兼容编码相关接口 ---------------------------------------------------------------- */

/*
  SM2非对称加密的结果由C1,C2,C3三部分组成。其中C1是生成随机数的计算出的椭圆曲线点,C2是密文数据,C3是SM3的摘要值。
  C1||C3||C2的ASN1编码格式为目前最新标准规范格式，旧版本标准规范格式为C1||C2||C3
 */
 typedef enum SM2CipherMode
 {
    SM2CipherMode_C1C3C2_ASN1,
    SM2CipherMode_C1C3C2,
    SM2CipherMode_C1C2C3_ASN1,
    SM2CipherMode_C1C2C3,
    SM2CipherMode_04C1C3C2,
    SM2CipherMode_04C1C2C3
} SM2CipherMode;

/**
 SM2非对称加解密算法，加密的兼容接口
 @param ctx 函数入参 - 上下文
 @param in 函数入参 - 待加密消息
 @param inlen 函数入参 - 消息长度(字节单位)
 @param strPubKey 函数入参 - 公钥
 @param pubkeyLen 函数入参 - 公钥长度
 @param out 函数出参 - 密文
 @param outlen 函数出参 - 密文长度
 @param mode 密文输出格式
 @return  0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2EncryptWithMode(sm2_ctx_t *ctx, const unsigned char *in, size_t inlen, const char *strPubKey, size_t pubkeyLen, unsigned char *out, size_t *outlen,SM2CipherMode mode);

/**
 SM2非对称加解密算法，解密的兼容接口
 @param ctx 函数入参 - 上下文
 @param in  函数入参 - 待解密密文
 @param inlen  函数入参 - 密文长度(字节单位)
 @param strPriKey  函数入参 - 私钥
 @param prikeyLen  函数入参 - 私钥长度
 @param out  函数出参 - 明文
 @param outlen  函数出参 - 明文长度
 @param mode  密文格式
 @return  0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2DecryptWithMode(sm2_ctx_t *ctx, const unsigned char *in, size_t inlen, const char *strPriKey, size_t prikeyLen,unsigned char *out, size_t *outlen,SM2CipherMode mode);

/*
  SM2签名结果由R和S分量组成，标准规定需采用ASN1编码，但仍然提供SM2SignMode_RS模式，以便兼容那些没有使用ASN1编码的版本。
 */
 typedef enum SM2SignMode
 {
    SM2SignMode_RS_ASN1,
    SM2SignMode_RS
} SM2SignMode;

/**
 SM2签名验签算法，签名的兼容接口
 @param ctx 函数入参 - 上下文
 @param msg 函数入参 - 待签名消息
 @param msglen 函数入参 - 待签名消息长度
 @param id 函数入参 - 用户ID(作用是加入到签名hash中，对于传入值无特殊要求)
 @param idlen 函数入参 - 用户ID长度
 @param strPubKey 函数入参 - 公钥(作用是加入到签名hash中)
 @param pubkeyLen 函数入参 - 公钥长度
 @param strPriKey 函数入参 - 私钥
 @param prikeyLen 函数入参 - 私钥长度
 @param sig 函数出参 - 签名结果
 @param siglen 函数出参 - 签名结果长度
 @param mode 签名格式
 */
SMLib_EXPORT int SM2SignWithMode(sm2_ctx_t *ctx, const unsigned char *msg, size_t msglen, const char *id, size_t idlen, const char *strPubKey, size_t pubkeyLen, const char *strPriKey, size_t prikeyLen,unsigned char *sig, size_t *siglen,SM2SignMode mode);

/**
 SM2签名验签算法，验签的兼容接口
 @param ctx 函数入参 - 上下文
 @param msg 函数入参 - 待验签内容
 @param msglen 函数入参 - 待验签内容长度
 @param id 函数入参 - 用户ID
 @param idlen 函数入参 - 用户ID长度
 @param sig 函数入参 - 签名结果
 @param siglen 函数入参 - 签名结果长度
 @param strPubKey 函数入参 - 公钥
 @param pubkeyLen 函数入参 - 公钥长度
 @param mode 签名格式
 @return 0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2VerifyWithMode(sm2_ctx_t *ctx, const unsigned char *msg, size_t msglen,const char *id, size_t idlen, const unsigned char *sig, size_t siglen, const char *strPubKey, size_t pubkeyLen,SM2SignMode mode);

/* ---------------------------------------------------------------- 以下为通用随机数接口 ---------------------------------------------------------------- */

/**
 产生随机数。

 @param bytes_buf 函数出参 - 随机数的buf
 @param bytes_len 函数入参 - 需要获取的随机数长度，字节单位
*/
SMLib_EXPORT int generateRandomBytes(unsigned char *bytes_buf,int bytes_len);

/* ---------------------------------------------------------------- 以下为SM2外部随机熵相关接口，以便特殊业务需求调用 ---------------------------------------------------------------- */

/**
  签名时使用的随机数由外部设置
*/
SMLib_EXPORT int SM2SetRandomDataCtx(sm2_ctx_t *ctx,const char *sign_random);

/**
  判断外部设置的签名时随机数是否有效
*/
SMLib_EXPORT int IsSM2CtxRandomDataVaild(sm2_ctx_t *sm2ctx);

/**
 为SM2增加外部熵源，该接口主要用于输入外部随机熵，一般情况下无需调用该接口，模块内部使用的随机熵通常情况下可保证熵值足够。
 
 @param ctx 函数入参 - 上下文
 @param buf 函数入参 - 熵buf
 @param buflen 函数入参 - 熵buf长度
*/
SMLib_EXPORT int SM2ReSeed(sm2_ctx_t *ctx, const unsigned char *buf, size_t buflen);


/* -------------------------------------- 以下接口为本地证书的轻量级管理接口，提供证书请求串的生成、证书的导入导出、解析证书关键信息，以及进行本地证书链的验证功能  ------------------------------------- */

 typedef enum SM2CSRMode
 {
    SM2CSRMode_Single,/*国密单证*/
    SM2CSRMode_Double /*国密双证*/
} SM2CSRMode;

typedef struct {
   const char* dir;
   void* def;
   void* sm2handler;
}sm2_cert_ctx_t;

typedef struct {
   unsigned long not_before;
   unsigned long not_after;
 }sm2_valid_time_t;
//-----------------------------------证书生成项及类型定义----------------------------------
#define TYPE_CERTITEM_NECESSARY_MAX           10      //必要参数type最大值
#define TYPE_CERTITEM_NECESSARY_COUNT         4       //必要参数type数量
#define TYPE_CERTITEM_CSR_DER                  1      //(unsigned char*)   //cert request der  
#define TYPE_CERTITEM_SERIALNUMBER             2      //(char*)            //序列号
#define TYPE_CERTITEM_SIGN_PRIKEY              3      //(char*)            //私钥
#define TYPE_CERTITEM_VALID_DAYS               4      //(int*)            //有效天数
#define TYPE_CERTITEM_USEAGE                  20      //(int*)            //签名证书(1), 加解密证书(2), 两者都有(3)
#define TYPE_CERTITEM_CA_DER                  21      //(unsigned char*)   //ca cert der,自签名时不需要
#define TYPE_CERTITEM_SUBJECT_BEGIN          100      //subject下项目type begin
#define TYPE_CERTITEM_SUBJECT_END            199      //subject下项目type end
#define TYPE_CERTITEM_SUBJECT_DESCRIPTION    101      //(char*)            //subject下description项目
#define TYPE_CERTITEM_SUBJECT_CN             102      //(char*)            //subject下commonname项目

//---------------------------证书读取项目及定义-----------------------------------
#define TYPE_READCERTITEM_PUBKEY                  2002      //(char*)            //公钥
#define TYPE_READCERTITEM_SERIALNUMBER            TYPE_CERTITEM_SERIALNUMBER      //(char*)            //序列号
#define TYPE_READCERTITEM_SUBJECT_DESCRIPTION    TYPE_CERTITEM_SUBJECT_DESCRIPTION      //(char*)            //subject下description项目
#define TYPE_READCERTITEM_SUBJECT_CN             TYPE_CERTITEM_SUBJECT_CN      //(char*)            //subject下commonname项目

typedef struct stTCSMItem {
    int type;
    void* value;
    int valuelen;
} TstTCSMItem;                     


SMLib_EXPORT int SM2CertCtxSize(void);
/**
初始化证书管理上下文sm2_cert_ctx_t
 
@param ctx 函数入参 - 上下文
@param dir 函数入参 - 目录，指定存储证书信息的目录，不导入证书时可以传NULL
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertInitCtx(sm2_cert_ctx_t *ctx,const char* dir);

/**
销毁证书管理上下文sm2_cert_ctx_t
 
@param ctx 函数入参 - 上下文
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertFreeCtx(sm2_cert_ctx_t *ctx);

/**
生成SM2的证书请求串，生成的CSR以DER格式编码，支持普通单证书和国密双证书
 
@param ctx 函数入参 - 上下文
@param country_name 函数入参 - CountryName
@param province 函数入参 - StateOrProvinceName
@param locality_name 函数入参 - LocalityName
@param organization_name 函数入参 - OrganizationName
@param organization_unit_name 函数入参 - OrganizationalUnitName
@param common_name 函数入参 - CommonName
@param email 函数入参 - Email
@param challenge_password 函数入参 - ChallengePassword
@param public_key 函数入参 - 公钥
@param private_key 函数入参 - 私钥，用作CSR签名
@param temp_public_key 函数入参 - 临时公钥，国密双证书需要填写该字段，单证书填NULL，双证书使用该公钥用作交换公私钥，用于解密CA下发的加密密钥对
@param out 函数出参 - CSR输出，请保证out指向的内存空间足够，例如可以分配8K空间，der编码
@param outlen 函数出参 - CSR输出的长度
@param mode 模式 -  选择生成SM2单证书或双证书
@return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2GenerateCSR(sm2_cert_ctx_t *ctx,const char* country_name,const char* province,const char* locality_name,const char* organization_name,
                                const char* organization_unit_name,const char* common_name,const char* email,const char* challenge_password,
                                const char* public_key,const char*  private_key,const char* temp_public_key,unsigned char* out,int *outlen,SM2CSRMode mode);

/**
导入证书
 
@param ctx 函数入参 - 上下文
@param cert_data 函数入参 - 证书数据，为der编码二进制数据
@param cert_data_len 函数入参 - 证书数据长度
@param cert_id 函数出参 - 导入后生成的证书ID，后续可通过证书ID进行访问
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertImport(sm2_cert_ctx_t *ctx,unsigned char* cert_data,int cert_data_len,char* cert_id);

/**
导出证书
 
@param ctx 函数入参 - 上下文
@param cert_data 函数出参 - 证书数据，为der编码二进制数据
@param cert_data_len 函数出参 - 证书数据长度
@param cert_id 函数入参 - 证书ID
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertExport(sm2_cert_ctx_t *ctx,unsigned char* cert_data,int* cert_data_len,const char* cert_id);

/**
 * @brief 检查csr是否sm2 pubkey，并输出pubkey
 * 
 * @param ctx input 上下文
 * @param der input csr der data
 * @param derlen input csr der data len
 * @param pubkey output !NULL时输出公钥buf，size > 130; NULL时，不输出
 * @return 0表示成功，其他值为错误码
 */
SMLib_EXPORT int SM2CSRCheckSm2Pubkey(sm2_cert_ctx_t *ctx, unsigned char*der, int derlen, char* pubkey);
/**
删除证书
 
@param ctx 函数入参 - 上下文
@param cert_id 函数入参 - 证书ID
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertDelete(sm2_cert_ctx_t *ctx,const char* cert_id);

/**
根据证书ID读取公钥
 
@param ctx 函数入参 - 上下文
@param cert_id 函数入参 - 证书ID
@param public_key 函数出参 - 公钥
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertReadPublicKey(sm2_cert_ctx_t *ctx,const char* cert_id,char* public_key);

/**
根据证书ID读取有效期
 
@param ctx 函数入参 - 上下文
@param cert_id 函数入参 - 证书ID
@param valid_time 函数出参 - 有效期
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertReadValidTime(sm2_cert_ctx_t *ctx,const char* cert_id,sm2_valid_time_t* valid_time);

/**
判断证书是否为根证书，即是否为自签发证书
 
@param ctx 函数入参 - 上下文
@param cert_id 函数入参 - 证书ID
@param b_root_cert 函数出参 - 是否根证书，1为根证书，0为不是根证书
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CertIsRoot(sm2_cert_ctx_t *ctx,const char* cert_id,int* b_root_cert);

/**
进行本地证书链的验证，请确保你所验证的证书的所有上级证书均已通过导入接口事先导入，否则会由于证书链不完整验证失败
证书链的验证这里未验证是否过期，原因为获取的本地时间戳可能不准。
如果调用方需要验证是否过期。可通过SM2CertReadValidTime接口读取有效期并与可信的时间戳进行对比。
证书链的验证这里未验证是否吊销。
如果验证方需要验证是否吊销。可自行前往CA请求证书状态。
 
@param ctx 函数入参 - 上下文
@param cert_id 函数入参 - 证书ID
#return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2VerifyCertChain(sm2_cert_ctx_t *ctx,const char* cert_id);

/**
 * @brief 根据csr生成证书
 * @param ctx input 上下文
 * @param der_csr input cert request der
 * @param der_csr_len input cert request der len
 * @param der_ca input ca cert der
 * @param der_ca_len imput  ca len
 * @param serialNumber input 序列号,maxlen=20 bytes
 * @param prikey input ca 私钥
 * @param valid_days input 有效天数
 * @param isSign input 是否签名用证书
 * @param pcertder output cert der
 * @param derlen input&output input buf len and output cert der len
 * @return int --0-- successful
 */
SMLib_EXPORT int SM2CertGenerate(sm2_cert_ctx_t *ctx, const unsigned char* der_csr, int der_csr_len,
                        const unsigned char* der_ca, int der_ca_len, const char* serialNumber,
                        const char* prikey, int valid_days, int isSign,
                        unsigned char* pcertder, int* derlen);

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
SMLib_EXPORT int SM2CertGenerate2(sm2_cert_ctx_t *ctx, TstTCSMItem items[], int itemcount, unsigned char* pcertder, int* derlen);
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
 * @return  int --0-- successful
 */
SMLib_EXPORT int SM2GetCertItem(sm2_cert_ctx_t *ctx, const unsigned char* der, int certlen, int itemid, int* itemcount, char* outstr, int* len);
/**
将DER编码格式的CSR转换为PEM格式
 
@param der 函数入参 - der CSR数据
@param derlen 函数入参 - der CSR数据的长度
@param pem 函数出参 - pem格式的CSR数据
@return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CSRConvertDER2PEM(const unsigned char* der,size_t derlen,char* pem);

/**
将PEM编码格式的CSR转换为DER格式
 
@param pem 函数入参 - der CSR数据
@param der 函数入参 - der CSR数据的长度
@param derlen 函数出参 - pem格式的CSR数据
@return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CSRConvertPEM2DER(const char* pem,unsigned char* der,size_t *derlen);

/**
将DER编码格式的证书转换为PEM格式
 
@param der 函数入参 - der CSR数据
@param derlen 函数入参 - der CSR数据的长度
@param pem 函数出参 - pem格式的CSR数据
@return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CRTConvertDER2PEM(const unsigned char* der,size_t derlen,char* pem);

/**
将PEM编码格式的证书转换为DER格式
 
@param pem 函数入参 - der CSR数据
@param der 函数入参 - der CSR数据的长度
@param derlen 函数出参 - pem格式的CSR数据
@return 0表示成功，其他值为错误码
*/
SMLib_EXPORT int SM2CRTConvertPEM2DER(const char* pem,unsigned char* der,size_t *derlen);

/**
 * @brief der 2 pem for prikey
 * 
 * @param der der格式私钥
 * @param derlen der格式私钥长度
 * @param pem pem格式私钥
 * @return int   --0-- successful
 */
SMLib_EXPORT int SM2PrikeyConvertDER2PEM(const unsigned char* der,size_t derlen,char* pem);
/**
 * @brief pem 2 der for prikeypem
 * 
 * @param pem  pem格式私钥
 * @param der der格式私钥
 * @param derlen  der格式私钥长度
 * @return int  --0-- successful
 */
SMLib_EXPORT int SM2PrikeyDerGetFromPem(const char* pem, unsigned char* der, size_t* derlen);
/**
 * @brief 从key der 中获取prikey 字符串,未压缩格式
 * 
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式私钥
 * @param der_len  der格式私钥长度
 * @param prikeystr  私钥字符串
 * @return int int --0-- successful
 */
SMLib_EXPORT int SM2PrikeyStrGetFromDer(sm2_cert_ctx_t *ctx, const unsigned char* der, int der_len, char* prikeystr);
/**
 * @brief der to str for prikey
 * 
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式私钥
 * @param der_len der格式私钥长度
 * @param prikeystr 私钥字符串
 * @return int  --0-- successful
 */
SMLib_EXPORT int SM2PrikeyStr2Der(sm2_cert_ctx_t *ctx, unsigned char* der, int* der_len, const char* prikeystr);

/**
 * @brief 解析prikey pem格式字符串并输出私钥等信息
 *        目前支持格式EC PRIVATE KEY和PRIVATE KEY
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param keypem  [input] pem格式字符串 
 * @param outPrikey [output] 输出的私钥字串，hex格式
 * @param extdata  [output] 扩展输出，备用
 * @return SMLib_EXPORT 
 */
SMLib_EXPORT int SM2ParsePrikeyPem(sm2_cert_ctx_t *ctx, const char* keypem, char* outPrikey, char** extdata);

/**
 * @brief der 2 pem for pubkey
 * 
 * @param der der格式公钥
 * @param derlen der格式公钥长度
 * @param pem pem格式公钥
 * @return int   --0-- successful
 */
SMLib_EXPORT int SM2PubkeyConvertDER2PEM(const unsigned char* der,size_t derlen,char* pem);
/**
 * @brief pem 2 der for pubkeypem
 * 
 * @param pem  pem格式公钥
 * @param der der格式公钥
 * @param derlen  der格式公钥长度
 * @return int  --0-- successful
 */
SMLib_EXPORT int SM2PubkeyDerGetFromPem(const char* pem, unsigned char* der, size_t* derlen);
/**
 * @brief 从key der 中获取pubkey 字符串,未压缩格式
 * 
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式公钥
 * @param der_len  der格式公钥长度
 * @param pubkeystr  公钥字符串
 * @return int int --0-- successful
 */
SMLib_EXPORT int SM2PubkeyStrGetFromDer(sm2_cert_ctx_t *ctx, const unsigned char* der, int der_len, char* pubkeystr);
/**
 * @brief der to str for pubkey
 * 
 * @param ctx 上下文，该格式转换需进行ASN1编码解析，因而需传入上下文
 * @param der der格式公钥
 * @param der_len der格式公钥长度
 * @param pubkeystr 公钥字符串
 * @return int  --0-- successful
 */
SMLib_EXPORT int SM2PubkeyStr2Der(sm2_cert_ctx_t *ctx, unsigned char* der, int* der_len, const char* pubkeystr);


#ifdef  __cplusplus
}
#endif
#endif
