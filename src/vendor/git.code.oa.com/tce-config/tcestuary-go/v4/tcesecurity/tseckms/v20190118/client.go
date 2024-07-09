// All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20190118

import (
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common"
	tchttp "github.com/tencentyun/tcecloud-sdk-go/tcecloud/common/http"
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common/profile"
)

const APIVersion = "2019-01-18"

type Client struct {
	common.Client
}

// Deprecated
func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
	cpf := profile.NewClientProfile()
	client = &Client{}
	client.Init(region).WithSecretId(secretId, secretKey).WithProfile(cpf)
	return
}

func NewClient(credential *common.Credential, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
	client = &Client{}
	client.Init(region).
		WithCredential(credential).
		WithProfile(clientProfile)
	return
}

func NewArchiveKeyRequest() (request *ArchiveKeyRequest) {
	request = &ArchiveKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "ArchiveKey")
	return
}

func NewArchiveKeyResponse() (response *ArchiveKeyResponse) {
	response = &ArchiveKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 对密钥进行归档，被归档的密钥只能用于解密，不能加密
func (c *Client) ArchiveKey(request *ArchiveKeyRequest) (response *ArchiveKeyResponse, err error) {
	if request == nil {
		request = NewArchiveKeyRequest()
	}
	response = NewArchiveKeyResponse()
	err = c.Send(request, response)
	return
}

func NewAsymmetricRsaDecryptRequest() (request *AsymmetricRsaDecryptRequest) {
	request = &AsymmetricRsaDecryptRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "AsymmetricRsaDecrypt")
	return
}

func NewAsymmetricRsaDecryptResponse() (response *AsymmetricRsaDecryptResponse) {
	response = &AsymmetricRsaDecryptResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 使用指定的RSA非对称密钥的私钥进行数据解密，密文必须是使用对应公钥加密的。处于Enabled 状态的非对称密钥才能进行解密操作。
func (c *Client) AsymmetricRsaDecrypt(request *AsymmetricRsaDecryptRequest) (response *AsymmetricRsaDecryptResponse, err error) {
	if request == nil {
		request = NewAsymmetricRsaDecryptRequest()
	}
	response = NewAsymmetricRsaDecryptResponse()
	err = c.Send(request, response)
	return
}

func NewAsymmetricRsaEncryptRequest() (request *AsymmetricRsaEncryptRequest) {
	request = &AsymmetricRsaEncryptRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "AsymmetricRsaEncrypt")
	return
}

func NewAsymmetricRsaEncryptResponse() (response *AsymmetricRsaEncryptResponse) {
	response = &AsymmetricRsaEncryptResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 使用指定的RSA非对称密钥的公钥进行数据加密。
func (c *Client) AsymmetricRsaEncrypt(request *AsymmetricRsaEncryptRequest) (response *AsymmetricRsaEncryptResponse, err error) {
	if request == nil {
		request = NewAsymmetricRsaEncryptRequest()
	}
	response = NewAsymmetricRsaEncryptResponse()
	err = c.Send(request, response)
	return
}

func NewAsymmetricSm2DecryptRequest() (request *AsymmetricSm2DecryptRequest) {
	request = &AsymmetricSm2DecryptRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "AsymmetricSm2Decrypt")
	return
}

func NewAsymmetricSm2DecryptResponse() (response *AsymmetricSm2DecryptResponse) {
	response = &AsymmetricSm2DecryptResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 使用指定的SM2非对称密钥的私钥进行数据解密，密文必须是使用对应公钥加密的。处于Enabled 状态的非对称密钥才能进行解密操作。传入的密文的长度不能超过256字节。
func (c *Client) AsymmetricSm2Decrypt(request *AsymmetricSm2DecryptRequest) (response *AsymmetricSm2DecryptResponse, err error) {
	if request == nil {
		request = NewAsymmetricSm2DecryptRequest()
	}
	response = NewAsymmetricSm2DecryptResponse()
	err = c.Send(request, response)
	return
}

func NewAsymmetricSm2EncryptRequest() (request *AsymmetricSm2EncryptRequest) {
	request = &AsymmetricSm2EncryptRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "AsymmetricSm2Encrypt")
	return
}

func NewAsymmetricSm2EncryptResponse() (response *AsymmetricSm2EncryptResponse) {
	response = &AsymmetricSm2EncryptResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 使用指定的SM2非对称密钥的公钥进行数据加密。传入的密文的长度不能超过160字节。
func (c *Client) AsymmetricSm2Encrypt(request *AsymmetricSm2EncryptRequest) (response *AsymmetricSm2EncryptResponse, err error) {
	if request == nil {
		request = NewAsymmetricSm2EncryptRequest()
	}
	response = NewAsymmetricSm2EncryptResponse()
	err = c.Send(request, response)
	return
}

func NewCancelKeyArchiveRequest() (request *CancelKeyArchiveRequest) {
	request = &CancelKeyArchiveRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "CancelKeyArchive")
	return
}

func NewCancelKeyArchiveResponse() (response *CancelKeyArchiveResponse) {
	response = &CancelKeyArchiveResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 取消密钥归档，取消后密钥的状态变为Enabled。
func (c *Client) CancelKeyArchive(request *CancelKeyArchiveRequest) (response *CancelKeyArchiveResponse, err error) {
	if request == nil {
		request = NewCancelKeyArchiveRequest()
	}
	response = NewCancelKeyArchiveResponse()
	err = c.Send(request, response)
	return
}

func NewCancelKeyDeletionRequest() (request *CancelKeyDeletionRequest) {
	request = &CancelKeyDeletionRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "CancelKeyDeletion")
	return
}

func NewCancelKeyDeletionResponse() (response *CancelKeyDeletionResponse) {
	response = &CancelKeyDeletionResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 取消CMK的计划删除操作
func (c *Client) CancelKeyDeletion(request *CancelKeyDeletionRequest) (response *CancelKeyDeletionResponse, err error) {
	if request == nil {
		request = NewCancelKeyDeletionRequest()
	}
	response = NewCancelKeyDeletionResponse()
	err = c.Send(request, response)
	return
}

func NewCreateKeyRequest() (request *CreateKeyRequest) {
	request = &CreateKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "CreateKey")
	return
}

func NewCreateKeyResponse() (response *CreateKeyResponse) {
	response = &CreateKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 创建用户管理数据密钥的主密钥CMK（Custom Master Key）。
func (c *Client) CreateKey(request *CreateKeyRequest) (response *CreateKeyResponse, err error) {
	if request == nil {
		request = NewCreateKeyRequest()
	}
	response = NewCreateKeyResponse()
	err = c.Send(request, response)
	return
}

func NewDecryptRequest() (request *DecryptRequest) {
	request = &DecryptRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "Decrypt")
	return
}

func NewDecryptResponse() (response *DecryptResponse) {
	response = &DecryptResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 本接口用于解密密文，得到明文数据。
func (c *Client) Decrypt(request *DecryptRequest) (response *DecryptResponse, err error) {
	if request == nil {
		request = NewDecryptRequest()
	}
	response = NewDecryptResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteImportedKeyMaterialRequest() (request *DeleteImportedKeyMaterialRequest) {
	request = &DeleteImportedKeyMaterialRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "DeleteImportedKeyMaterial")
	return
}

func NewDeleteImportedKeyMaterialResponse() (response *DeleteImportedKeyMaterialResponse) {
	response = &DeleteImportedKeyMaterialResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 用于删除导入的密钥材料，仅对EXTERNAL类型的CMK有效，该接口将CMK设置为PendingImport 状态，并不会删除CMK，在重新进行密钥导入后可继续使用。彻底删除CMK请使用 ScheduleKeyDeletion 接口。
func (c *Client) DeleteImportedKeyMaterial(request *DeleteImportedKeyMaterialRequest) (response *DeleteImportedKeyMaterialResponse, err error) {
	if request == nil {
		request = NewDeleteImportedKeyMaterialRequest()
	}
	response = NewDeleteImportedKeyMaterialResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeKeyRequest() (request *DescribeKeyRequest) {
	request = &DescribeKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "DescribeKey")
	return
}

func NewDescribeKeyResponse() (response *DescribeKeyResponse) {
	response = &DescribeKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 用于获取指定KeyId的主密钥属性详情信息。
func (c *Client) DescribeKey(request *DescribeKeyRequest) (response *DescribeKeyResponse, err error) {
	if request == nil {
		request = NewDescribeKeyRequest()
	}
	response = NewDescribeKeyResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeKeysRequest() (request *DescribeKeysRequest) {
	request = &DescribeKeysRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "DescribeKeys")
	return
}

func NewDescribeKeysResponse() (response *DescribeKeysResponse) {
	response = &DescribeKeysResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 该接口用于批量获取主密钥属性信息。
func (c *Client) DescribeKeys(request *DescribeKeysRequest) (response *DescribeKeysResponse, err error) {
	if request == nil {
		request = NewDescribeKeysRequest()
	}
	response = NewDescribeKeysResponse()
	err = c.Send(request, response)
	return
}

func NewDisableKeyRequest() (request *DisableKeyRequest) {
	request = &DisableKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "DisableKey")
	return
}

func NewDisableKeyResponse() (response *DisableKeyResponse) {
	response = &DisableKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 本接口用于禁用一个主密钥，处于禁用状态的Key无法用于加密、解密操作。
func (c *Client) DisableKey(request *DisableKeyRequest) (response *DisableKeyResponse, err error) {
	if request == nil {
		request = NewDisableKeyRequest()
	}
	response = NewDisableKeyResponse()
	err = c.Send(request, response)
	return
}

func NewDisableKeyRotationRequest() (request *DisableKeyRotationRequest) {
	request = &DisableKeyRotationRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "DisableKeyRotation")
	return
}

func NewDisableKeyRotationResponse() (response *DisableKeyRotationResponse) {
	response = &DisableKeyRotationResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 对指定的CMK禁止密钥轮换功能。
func (c *Client) DisableKeyRotation(request *DisableKeyRotationRequest) (response *DisableKeyRotationResponse, err error) {
	if request == nil {
		request = NewDisableKeyRotationRequest()
	}
	response = NewDisableKeyRotationResponse()
	err = c.Send(request, response)
	return
}

func NewDisableKeysRequest() (request *DisableKeysRequest) {
	request = &DisableKeysRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "DisableKeys")
	return
}

func NewDisableKeysResponse() (response *DisableKeysResponse) {
	response = &DisableKeysResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 该接口用于批量禁止CMK的使用。
func (c *Client) DisableKeys(request *DisableKeysRequest) (response *DisableKeysResponse, err error) {
	if request == nil {
		request = NewDisableKeysRequest()
	}
	response = NewDisableKeysResponse()
	err = c.Send(request, response)
	return
}

func NewEnableKeyRequest() (request *EnableKeyRequest) {
	request = &EnableKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "EnableKey")
	return
}

func NewEnableKeyResponse() (response *EnableKeyResponse) {
	response = &EnableKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 用于启用一个指定的CMK。
func (c *Client) EnableKey(request *EnableKeyRequest) (response *EnableKeyResponse, err error) {
	if request == nil {
		request = NewEnableKeyRequest()
	}
	response = NewEnableKeyResponse()
	err = c.Send(request, response)
	return
}

func NewEnableKeyRotationRequest() (request *EnableKeyRotationRequest) {
	request = &EnableKeyRotationRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "EnableKeyRotation")
	return
}

func NewEnableKeyRotationResponse() (response *EnableKeyRotationResponse) {
	response = &EnableKeyRotationResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 对指定的CMK开启密钥轮换功能。
func (c *Client) EnableKeyRotation(request *EnableKeyRotationRequest) (response *EnableKeyRotationResponse, err error) {
	if request == nil {
		request = NewEnableKeyRotationRequest()
	}
	response = NewEnableKeyRotationResponse()
	err = c.Send(request, response)
	return
}

func NewEnableKeysRequest() (request *EnableKeysRequest) {
	request = &EnableKeysRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "EnableKeys")
	return
}

func NewEnableKeysResponse() (response *EnableKeysResponse) {
	response = &EnableKeysResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 该接口用于批量启用CMK。
func (c *Client) EnableKeys(request *EnableKeysRequest) (response *EnableKeysResponse, err error) {
	if request == nil {
		request = NewEnableKeysRequest()
	}
	response = NewEnableKeysResponse()
	err = c.Send(request, response)
	return
}

func NewEncryptRequest() (request *EncryptRequest) {
	request = &EncryptRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "Encrypt")
	return
}

func NewEncryptResponse() (response *EncryptResponse) {
	response = &EncryptResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 本接口用于加密最多为4KB任意数据，可用于加密数据库密码，RSA Key，或其它较小的敏感信息。对于应用的数据加密，使用GenerateDataKey生成的DataKey进行本地数据的加解密操作
func (c *Client) Encrypt(request *EncryptRequest) (response *EncryptResponse, err error) {
	if request == nil {
		request = NewEncryptRequest()
	}
	response = NewEncryptResponse()
	err = c.Send(request, response)
	return
}

func NewGenerateDataKeyRequest() (request *GenerateDataKeyRequest) {
	request = &GenerateDataKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "GenerateDataKey")
	return
}

func NewGenerateDataKeyResponse() (response *GenerateDataKeyResponse) {
	response = &GenerateDataKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 本接口生成一个数据密钥，您可以用这个密钥进行本地数据的加密。
func (c *Client) GenerateDataKey(request *GenerateDataKeyRequest) (response *GenerateDataKeyResponse, err error) {
	if request == nil {
		request = NewGenerateDataKeyRequest()
	}
	response = NewGenerateDataKeyResponse()
	err = c.Send(request, response)
	return
}

func NewGenerateRandomRequest() (request *GenerateRandomRequest) {
	request = &GenerateRandomRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "GenerateRandom")
	return
}

func NewGenerateRandomResponse() (response *GenerateRandomResponse) {
	response = &GenerateRandomResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 随机数生成接口。
func (c *Client) GenerateRandom(request *GenerateRandomRequest) (response *GenerateRandomResponse, err error) {
	if request == nil {
		request = NewGenerateRandomRequest()
	}
	response = NewGenerateRandomResponse()
	err = c.Send(request, response)
	return
}

func NewGetKeyRotationStatusRequest() (request *GetKeyRotationStatusRequest) {
	request = &GetKeyRotationStatusRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "GetKeyRotationStatus")
	return
}

func NewGetKeyRotationStatusResponse() (response *GetKeyRotationStatusResponse) {
	response = &GetKeyRotationStatusResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 查询指定的CMK是否开启了密钥轮换功能。
func (c *Client) GetKeyRotationStatus(request *GetKeyRotationStatusRequest) (response *GetKeyRotationStatusResponse, err error) {
	if request == nil {
		request = NewGetKeyRotationStatusRequest()
	}
	response = NewGetKeyRotationStatusResponse()
	err = c.Send(request, response)
	return
}

func NewGetParametersForImportRequest() (request *GetParametersForImportRequest) {
	request = &GetParametersForImportRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "GetParametersForImport")
	return
}

func NewGetParametersForImportResponse() (response *GetParametersForImportResponse) {
	response = &GetParametersForImportResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 获取导入主密钥（CMK）材料的参数，返回的Token作为执行ImportKeyMaterial的参数之一，返回的PublicKey用于对自主导入密钥材料进行加密。返回的Token和PublicKey 24小时后失效，失效后如需重新导入，需要再次调用该接口获取新的Token和PublicKey。
func (c *Client) GetParametersForImport(request *GetParametersForImportRequest) (response *GetParametersForImportResponse, err error) {
	if request == nil {
		request = NewGetParametersForImportRequest()
	}
	response = NewGetParametersForImportResponse()
	err = c.Send(request, response)
	return
}

func NewGetPublicKeyRequest() (request *GetPublicKeyRequest) {
	request = &GetPublicKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "GetPublicKey")
	return
}

func NewGetPublicKeyResponse() (response *GetPublicKeyResponse) {
	response = &GetPublicKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 该接口用户获取 KeyUsage为ASYMMETRIC_DECRYPT_RSA_2048 和 ASYMMETRIC_DECRYPT_SM2 的非对称密钥的公钥信息，使用该公钥用户可在本地进行数据加密，使用该公钥加密的数据只能通过KMS使用对应的私钥进行解密。只有处于Enabled状态的非对称密钥才可能获取公钥。
func (c *Client) GetPublicKey(request *GetPublicKeyRequest) (response *GetPublicKeyResponse, err error) {
	if request == nil {
		request = NewGetPublicKeyRequest()
	}
	response = NewGetPublicKeyResponse()
	err = c.Send(request, response)
	return
}

func NewImportKeyMaterialRequest() (request *ImportKeyMaterialRequest) {
	request = &ImportKeyMaterialRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "ImportKeyMaterial")
	return
}

func NewImportKeyMaterialResponse() (response *ImportKeyMaterialResponse) {
	response = &ImportKeyMaterialResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 用于导入密钥材料。只有类型为EXTERNAL 的CMK 才可以导入，导入的密钥材料使用 GetParametersForImport 获取的密钥进行加密。可以为指定的 CMK 重新导入密钥材料，并重新指定过期时间，但必须导入相同的密钥材料。CMK 密钥材料导入后不可以更换密钥材料。导入的密钥材料过期或者被删除后，指定的CMK将无法使用，需要再次导入相同的密钥材料才能正常使用。CMK是独立的，同样的密钥材料可导入不同的 CMK 中，但使用其中一个 CMK 加密的数据无法使用另一个 CMK解密。
// 只有Enabled 和 PendingImport状态的CMK可以导入密钥材料。
func (c *Client) ImportKeyMaterial(request *ImportKeyMaterialRequest) (response *ImportKeyMaterialResponse, err error) {
	if request == nil {
		request = NewImportKeyMaterialRequest()
	}
	response = NewImportKeyMaterialResponse()
	err = c.Send(request, response)
	return
}

func NewListAlgorithmsRequest() (request *ListAlgorithmsRequest) {
	request = &ListAlgorithmsRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "ListAlgorithms")
	return
}

func NewListAlgorithmsResponse() (response *ListAlgorithmsResponse) {
	response = &ListAlgorithmsResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 列出当前Region支持的加密方式
func (c *Client) ListAlgorithms(request *ListAlgorithmsRequest) (response *ListAlgorithmsResponse, err error) {
	if request == nil {
		request = NewListAlgorithmsRequest()
	}
	response = NewListAlgorithmsResponse()
	err = c.Send(request, response)
	return
}

func NewListKeyDetailRequest() (request *ListKeyDetailRequest) {
	request = &ListKeyDetailRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "ListKeyDetail")
	return
}

func NewListKeyDetailResponse() (response *ListKeyDetailResponse) {
	response = &ListKeyDetailResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 根据指定Offset和Limit获取主密钥列表详情。
func (c *Client) ListKeyDetail(request *ListKeyDetailRequest) (response *ListKeyDetailResponse, err error) {
	if request == nil {
		request = NewListKeyDetailRequest()
	}
	response = NewListKeyDetailResponse()
	err = c.Send(request, response)
	return
}

func NewReEncryptRequest() (request *ReEncryptRequest) {
	request = &ReEncryptRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "ReEncrypt")
	return
}

func NewReEncryptResponse() (response *ReEncryptResponse) {
	response = &ReEncryptResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 使用指定CMK对密文重新加密。
func (c *Client) ReEncrypt(request *ReEncryptRequest) (response *ReEncryptResponse, err error) {
	if request == nil {
		request = NewReEncryptRequest()
	}
	response = NewReEncryptResponse()
	err = c.Send(request, response)
	return
}

func NewScheduleKeyDeletionRequest() (request *ScheduleKeyDeletionRequest) {
	request = &ScheduleKeyDeletionRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "ScheduleKeyDeletion")
	return
}

func NewScheduleKeyDeletionResponse() (response *ScheduleKeyDeletionResponse) {
	response = &ScheduleKeyDeletionResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// CMK计划删除接口，用于指定CMK删除的时间，可选时间区间为[7,30]天
func (c *Client) ScheduleKeyDeletion(request *ScheduleKeyDeletionRequest) (response *ScheduleKeyDeletionResponse, err error) {
	if request == nil {
		request = NewScheduleKeyDeletionRequest()
	}
	response = NewScheduleKeyDeletionResponse()
	err = c.Send(request, response)
	return
}

func NewSignByAsymmetricKeyRequest() (request *SignByAsymmetricKeyRequest) {
	request = &SignByAsymmetricKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "SignByAsymmetricKey")
	return
}

func NewSignByAsymmetricKeyResponse() (response *SignByAsymmetricKeyResponse) {
	response = &SignByAsymmetricKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 非对称密钥签名。
// 注意：只有成功创建了KeyUsage= ASYMMETRIC_SIGN_VERIFY_SM2 的密钥才可以使用签名功能
func (c *Client) SignByAsymmetricKey(request *SignByAsymmetricKeyRequest) (response *SignByAsymmetricKeyResponse, err error) {
	if request == nil {
		request = NewSignByAsymmetricKeyRequest()
	}
	response = NewSignByAsymmetricKeyResponse()
	err = c.Send(request, response)
	return
}

func NewUpdateAliasRequest() (request *UpdateAliasRequest) {
	request = &UpdateAliasRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "UpdateAlias")
	return
}

func NewUpdateAliasResponse() (response *UpdateAliasResponse) {
	response = &UpdateAliasResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 用于修改CMK的别名。对于处于PendingDelete状态的CMK禁止修改。
func (c *Client) UpdateAlias(request *UpdateAliasRequest) (response *UpdateAliasResponse, err error) {
	if request == nil {
		request = NewUpdateAliasRequest()
	}
	response = NewUpdateAliasResponse()
	err = c.Send(request, response)
	return
}

func NewUpdateKeyDescriptionRequest() (request *UpdateKeyDescriptionRequest) {
	request = &UpdateKeyDescriptionRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "UpdateKeyDescription")
	return
}

func NewUpdateKeyDescriptionResponse() (response *UpdateKeyDescriptionResponse) {
	response = &UpdateKeyDescriptionResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 该接口用于对指定的cmk修改描述信息。对于处于PendingDelete状态的CMK禁止修改。
func (c *Client) UpdateKeyDescription(request *UpdateKeyDescriptionRequest) (response *UpdateKeyDescriptionResponse, err error) {
	if request == nil {
		request = NewUpdateKeyDescriptionRequest()
	}
	response = NewUpdateKeyDescriptionResponse()
	err = c.Send(request, response)
	return
}

func NewVerifyByAsymmetricKeyRequest() (request *VerifyByAsymmetricKeyRequest) {
	request = &VerifyByAsymmetricKeyRequest{
		BaseRequest: &tchttp.BaseRequest{},
	}
	request.Init().WithApiInfo("tseckms", APIVersion, "VerifyByAsymmetricKey")
	return
}

func NewVerifyByAsymmetricKeyResponse() (response *VerifyByAsymmetricKeyResponse) {
	response = &VerifyByAsymmetricKeyResponse{
		BaseResponse: &tchttp.BaseResponse{},
	}
	return
}

// 使用非对称密钥验签
func (c *Client) VerifyByAsymmetricKey(request *VerifyByAsymmetricKeyRequest) (response *VerifyByAsymmetricKeyResponse, err error) {
	if request == nil {
		request = NewVerifyByAsymmetricKeyRequest()
	}
	response = NewVerifyByAsymmetricKeyResponse()
	err = c.Send(request, response)
	return
}
