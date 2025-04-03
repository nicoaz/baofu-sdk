package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nicoaz/baofu-sdk/config"
	"github.com/nicoaz/baofu-sdk/consts"
	"github.com/nicoaz/baofu-sdk/models"
	"github.com/nicoaz/baofu-sdk/utils"
)

// AccountService 账户服务
type AccountService struct {
	config     *config.Config
	httpClient *utils.HTTPClient
}

// NewAccountService 创建账户服务
func NewAccountService(config *config.Config) *AccountService {
	return &AccountService{
		config:     config,
		httpClient: utils.NewHTTPClient(),
	}
}

func (s *AccountService) getHost(method string) string {
	if s.config.ReleaseEnv {
		return strings.Replace(consts.AccountServiceHostProd, "{报文编号}", method, 1)
	}
	return strings.Replace(consts.AccountServiceHostTest, "{报文编号}", method, 1)
}

// OpenAccount 开户接口
func (s *AccountService) OpenAccount(req *models.AccountOpenRequest) (string, error) {
	fmt.Println("==========================")
	fmt.Println("宝付账簿个人/机构开户接口")
	fmt.Println("==========================")

	// 构建Header参数
	headerPost := make(map[string]string)
	headerPost["memberId"] = s.config.MerchantID
	headerPost["terminalId"] = s.config.TerminalID
	headerPost["serviceTp"] = consts.MethodOpenAccount
	headerPost["verifyType"] = "1" // 加密方式目前只有1种，请填：1

	// 构建请求数据
	contentData := make(map[string]interface{})
	contentData["header"] = headerPost

	// 构建Body数据
	bodyData := make(map[string]interface{})
	bodyData["version"] = "4.1.0"         // 版本号
	bodyData["accType"] = "2"             // 账户类型:1个人,2商户
	bodyData["noticeUrl"] = req.NotifyUrl // 异步通知地址
	bodyData["businessType"] = "BCT2.0"   // 宝财通2.0

	// 构建账户信息
	accInfo := make(map[string]string)

	// 商户账户信息
	accInfo["transSerialNo"] = utils.GetTransid("TSN") // 请求流水号
	accInfo["loginNo"] = req.LoginNo                   // 登录号(用户ID)建议长度8位以上 唯一
	accInfo["email"] = req.Email                       // 邮箱
	accInfo["selfEmployed"] = req.SelfEmployed         // 是否个体户 默认为false
	accInfo["customerName"] = req.CustomerName         // 商户名称（营业执照上的名称）
	//accInfo["aliasName"] = ""                          // 商户名称别名  (选传)
	accInfo["certificateNo"] = req.CertificateNo         // 证件号码
	accInfo["certificateType"] = req.CertificateType     // 证件类型 营业执照:LICENSE
	accInfo["corporateName"] = req.CorporateName         // 法人姓名
	accInfo["corporateCertType"] = req.CorporateCertType // 法人证件类型  身份证:ID
	accInfo["corporateCertId"] = req.CorporateCertId     // 法人身份证号码
	accInfo["corporateMobile"] = req.CorporateMobile     // 法人手机号
	accInfo["industryId"] = req.IndustryId               // 所属行业
	//accInfo["contactName"] = ""                        // 联系人姓名
	//accInfo["contactMobile"] = ""                      // 联系人手机号
	accInfo["cardNo"] = req.CardNo                           // 卡号
	accInfo["bankName"] = req.BankName                       // 银行名称
	accInfo["depositBankProvince"] = req.DepositBankProvince // 开户行省份
	accInfo["depositBankCity"] = req.DepositBankCity         // 开户行城市
	accInfo["depositBankName"] = req.DeositBankName          // 开户支行名称
	accInfo["registerCapital"] = req.RegisterCapital         // 注册资本
	//accInfo["cardUserName"] = ""                       // 持卡人姓名  (个体绑法人对私必传)
	// 设置平台相关信息
	accInfo["platformNo"] = ""                 // 平台号(主商户号) (代理模式必传)
	accInfo["platformTerminalId"] = ""         // 终端号(代理模式必传)
	accInfo["qualificationTransSerialNo"] = "" // 资质文件流水,businessType为宝财通2.0非必填

	// 将账户信息添加到请求体中
	bodyData["accInfo"] = accInfo
	contentData["body"] = bodyData

	// 将请求数据转换为JSON
	jsonObject, err := json.Marshal(contentData)
	if err != nil {
		return "", err
	}
	fmt.Println("JSON：", string(jsonObject))

	// 加密请求数据
	dataContent, err := utils.EncryptByPFXFile(string(jsonObject), s.config.PrivateKey)
	if err != nil {
		return "", err
	}
	headerPost["content"] = dataContent

	// 发送请求
	response, err := utils.Post(headerPost, s.getHost(consts.MethodOpenAccount), "json")
	if err != nil {
		return "", err
	}
	fmt.Println("返回：", string(response))

	if len(response) == 0 {
		return "", fmt.Errorf("返回异常！")
	}

	// 解密返回数据
	rPostString, err := utils.DecryptByCERFile(string(response), s.config.BFPublicKey, s.config.BFPublicKeyPem)
	if err != nil {
		return "", err
	}
	fmt.Println("解密明文：", rPostString)

	return rPostString, nil
}

// 开户查询接口
func (s *AccountService) OpenAccountQuery(req *models.OpenAccountQueryRequest) (string, error) {
	fmt.Println("==========================")
	fmt.Println("宝付账簿开户查询接口")
	fmt.Println("==========================")

	// 构建Header参数
	headerPost := make(map[string]string)
	headerPost["memberId"] = s.config.MerchantID
	headerPost["terminalId"] = s.config.TerminalID
	headerPost["serviceTp"] = consts.MethodOpenAccountQuery
	headerPost["verifyType"] = "1" // 加密方式目前只有1种，请填：1

	// 构建请求数据
	contentData := make(map[string]interface{})
	contentData["header"] = headerPost

	// 构建Body数据
	bodyData := make(map[string]interface{})
	bodyData["version"] = "4.0.0"
	bodyData["certificateNo"] = req.CertificateNo
	bodyData["certificateType"] = req.CertificateType
	bodyData["platformNo"] = req.PlatformNo
	bodyData["loginNo"] = req.LoginNo
	bodyData["accType"] = req.AccType

	contentData["body"] = bodyData

	// 将请求数据转换为JSON
	jsonObject, err := json.Marshal(contentData)
	if err != nil {
		return "", err
	}
	fmt.Println("JSON：", string(jsonObject))

	// 加密请求数据
	dataContent, err := utils.EncryptByPFXFile(string(jsonObject), s.config.PrivateKey)
	if err != nil {
		return "", err
	}
	headerPost["content"] = dataContent

	// 发送请求
	response, err := utils.Post(headerPost, s.getHost(consts.MethodOpenAccountQuery), "json")
	if err != nil {
		return "", err
	}
	fmt.Println("返回：", string(response))

	if len(response) == 0 {
		return "", fmt.Errorf("返回异常！")
	}

	// 解密返回数据
	rPostString, err := utils.DecryptByCERFile(string(response), s.config.BFPublicKey, s.config.BFPublicKeyPem)
	if err != nil {
		return "", err
	}
	fmt.Println("解密明文：", rPostString)

	return rPostString, nil
}

// BalanceQuery 余额查询接口
func (s *AccountService) BalanceQuery(req *models.BalanceQueryRequest) (*models.BalanceQueryResponse, error) {

	// 构建Header参数
	headerPost := make(map[string]string)
	headerPost["memberId"] = s.config.MerchantID
	headerPost["terminalId"] = s.config.TerminalID
	headerPost["serviceTp"] = consts.MethodBalanceQuery
	headerPost["verifyType"] = "1" // 加密方式目前只有1种，请填：1

	// 构建请求数据
	contentData := make(map[string]interface{})
	contentData["header"] = headerPost

	// 构建Body数据
	bodyData := make(map[string]interface{})
	bodyData["version"] = "4.0.0"
	bodyData["contractNo"] = req.ContractNo
	bodyData["accType"] = req.AcctType

	contentData["body"] = bodyData

	// 将请求数据转换为JSON
	jsonObject, err := json.Marshal(contentData)
	if err != nil {
		return nil, err
	}
	fmt.Println("JSON: ", string(jsonObject))

	// 加密请求数据
	dataContent, err := utils.EncryptByPFXFile(string(jsonObject), s.config.PrivateKey)
	if err != nil {
		return nil, err
	}
	headerPost["content"] = dataContent

	// 发送请求
	response, err := utils.Post(headerPost, s.getHost(consts.MethodBalanceQuery), "json")
	if err != nil {
		return nil, err
	}
	fmt.Println("返回：", string(response))

	if len(response) == 0 {
		return nil, fmt.Errorf("返回异常！")
	}

	// 解密返回数据
	rPostString, err := utils.DecryptByCERFile(string(response), s.config.BFPublicKey, s.config.BFPublicKeyPem)
	if err != nil {
		return nil, err
	}
	fmt.Println("解密明文：", rPostString)

	var balanceQueryResponse models.BalanceQueryResponse
	err = json.Unmarshal([]byte(rPostString), &balanceQueryResponse)
	if err != nil {
		return nil, err
	}
	return &balanceQueryResponse, nil
}

// Transfer 账户间转账接口
func (s *AccountService) Transfer(req *models.TransferRequest) (*models.TransferResponse, error) {
	fmt.Println("==========================")
	fmt.Println("宝付账户间转账接口")
	fmt.Println("==========================")

	// 构建Header参数
	headerPost := make(map[string]string)
	headerPost["memberId"] = s.config.MerchantID
	headerPost["terminalId"] = s.config.TerminalID
	headerPost["serviceTp"] = consts.MethodTransfer
	headerPost["verifyType"] = "1" // 加密方式目前只有1种，请填：1

	// 构建请求数据
	contentData := make(map[string]interface{})
	contentData["header"] = headerPost

	// 构建Body数据
	bodyData := make(map[string]interface{})
	bodyData["version"] = "4.0.0"
	bodyData["payerNo"] = req.PayerNo
	bodyData["payeeNo"] = req.PayeeNo
	bodyData["transSerialNo"] = req.TransSerialNo
	bodyData["dealAmount"] = req.DealAmount

	contentData["body"] = bodyData

	// 将请求数据转换为JSON
	jsonObject, err := json.Marshal(contentData)
	if err != nil {
		return nil, err
	}
	fmt.Println("JSON：", string(jsonObject))

	// 加密请求数据
	dataContent, err := utils.EncryptByPFXFile(string(jsonObject), s.config.PrivateKey)
	if err != nil {
		return nil, err
	}
	headerPost["content"] = dataContent

	// 发送请求
	response, err := utils.Post(headerPost, s.getHost(consts.MethodTransfer), "json")
	if err != nil {
		return nil, err
	}
	fmt.Println("返回：", string(response))

	if len(response) == 0 {
		return nil, fmt.Errorf("返回异常！")
	}

	// 解密返回数据
	rPostString, err := utils.DecryptByCERFile(string(response), s.config.BFPublicKey, s.config.BFPublicKeyPem)
	if err != nil {
		return nil, err
	}

	fmt.Println("解密明文：", rPostString)

	var transferResponse models.TransferResponse
	err = json.Unmarshal([]byte(rPostString), &transferResponse)
	if err != nil {
		return nil, err
	}
	return &transferResponse, nil
}

// Withdraw 提现接口
func (s *AccountService) Withdraw(req *models.WithdrawRequest) (string, error) {
	fmt.Println("==========================")
	fmt.Println("宝付账簿提现接口")
	fmt.Println("==========================")

	// 构建Header参数
	headerPost := make(map[string]string)
	headerPost["memberId"] = s.config.MerchantID
	headerPost["terminalId"] = s.config.TerminalID
	headerPost["serviceTp"] = consts.MethodWithdraw
	headerPost["verifyType"] = "1" // 加密方式目前只有1种，请填：1

	// 构建请求数据
	contentData := make(map[string]interface{})
	contentData["header"] = headerPost

	// 构建Body数据
	bodyData := make(map[string]interface{})
	bodyData["version"] = "4.0.0"
	bodyData["srcAcctNo"] = req.SrcAcctNo
	bodyData["transAmt"] = req.TransAmt
	bodyData["transId"] = req.TransId
	bodyData["transDate"] = req.TransDate
	bodyData["transTime"] = req.TransTime
	bodyData["curType"] = req.CurType
	bodyData["transSummary"] = req.TransSummary
	bodyData["reservedExpand"] = req.ReservedExpand
	bodyData["cardNo"] = req.CardNo
	bodyData["cardName"] = req.CardName
	bodyData["cardBankCode"] = req.CardBankCode
	bodyData["directFlag"] = req.DirectFlag

	bodyData["returnUrl"] = req.NotifyUrl

	contentData["body"] = bodyData

	// 将请求数据转换为JSON
	jsonObject, err := json.Marshal(contentData)
	if err != nil {
		return "", err
	}
	fmt.Println("JSON：", string(jsonObject))

	// 加密请求数据
	dataContent, err := utils.EncryptByPFXFile(string(jsonObject), s.config.PrivateKey)
	if err != nil {
		return "", err
	}
	headerPost["content"] = dataContent

	// 发送请求
	response, err := utils.Post(headerPost, s.getHost(consts.MethodWithdraw), "json")
	if err != nil {
		return "", err
	}
	fmt.Println("返回：", string(response))

	if len(response) == 0 {
		return "", fmt.Errorf("返回异常！")
	}

	// 解密返回数据
	rPostString, err := utils.DecryptByCERFile(string(response), s.config.BFPublicKey, s.config.BFPublicKeyPem)
	if err != nil {
		return "", err
	}
	fmt.Println("解密明文：", rPostString)

	return rPostString, nil
}
