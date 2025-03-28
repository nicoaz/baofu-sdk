package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/nicoaz/baofu-sdk/config"
	"github.com/nicoaz/baofu-sdk/consts"
	"github.com/nicoaz/baofu-sdk/models"
	"github.com/nicoaz/baofu-sdk/utils"
)

// MerchantService 商户报备服务
type MerchantService struct {
	config     *config.Config
	httpClient *utils.HTTPClient
}

// NewMerchantService 创建商户报备服务
func NewMerchantService(config *config.Config) *MerchantService {
	return &MerchantService{
		config:     config,
		httpClient: utils.NewHTTPClient(),
	}
}

// getHost 获取服务主机地址
func (s *MerchantService) getHost() string {
	if s.config.ReleaseEnv {
		return consts.ReportServiceHostProd
	}
	return consts.ReportServiceHostTest
}

// MerchantWxReport 商户报备微信
func (s *MerchantService) MerchantWxReport(request *models.MerchantWXReportReq) (string, error) {
	fmt.Println("==========================")
	fmt.Println("商户报备微信")
	fmt.Println("==========================")

	request.MerId = s.config.MerchantID
	request.TerId = s.config.TerminalID
	request.ReportType = "WECHAT"
	request.ReportNo = utils.GetTransid("BBWX")

	// 将业务内容转为JSON
	bizContentJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("业务参数JSON编码失败: %v", err)
	}

	// 生成签名
	signStr, err := utils.Sign(string(bizContentJSON), s.config.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodMerchantReport)
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("bizContent", string(bizContentJSON))
	mapParams.Set("charset", "UTF-8")
	mapParams.Set("signStr", signStr)
	mapParams.Set("version", "1.0")
	mapParams.Set("format", "json")
	mapParams.Set("signType", "RSA")
	mapParams.Set("signSn", "1")
	mapParams.Set("ncrptnSn", "1")
	mapParams.Set("timestamp", time.Now().Format("20060102150405"))

	// 发送请求
	response, err := s.httpClient.Post(s.getHost(), mapParams)
	if err != nil {
		return "", fmt.Errorf("发送商户报备请求失败: %v", err)
	}

	// 解析响应
	var merchantResponse models.Response
	err = json.Unmarshal([]byte(response), &merchantResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if merchantResponse.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("商户报备请求失败: %s", merchantResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(merchantResponse.DataContent, merchantResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return "", fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		fmt.Println("警告：响应签名验证失败")
	}

	return merchantResponse.DataContent, nil
}

// MerchantReportQuery 商户报备查询
func (s *MerchantService) MerchantReportQuery(request *models.MerchantReportQueryRequest) (string, error) {
	fmt.Println("==========================")
	fmt.Println("商户报备查询")
	fmt.Println("==========================")

	request.MerId = s.config.MerchantID
	request.TerId = s.config.TerminalID

	// 将业务内容转为JSON
	bizContentJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("业务参数JSON编码失败: %v", err)
	}

	// 生成签名
	signStr, err := utils.Sign(string(bizContentJSON), s.config.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodMerchantReportQuery)
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("bizContent", string(bizContentJSON))
	mapParams.Set("charset", "UTF-8")
	mapParams.Set("signStr", signStr)
	mapParams.Set("version", "1.0")
	mapParams.Set("format", "json")
	mapParams.Set("signType", "RSA")
	mapParams.Set("signSn", "1")
	mapParams.Set("ncrptnSn", "1")
	mapParams.Set("timestamp", time.Now().Format("20060102150405"))

	// 发送请求
	response, err := s.httpClient.Post(s.getHost(), mapParams)
	if err != nil {
		return "", fmt.Errorf("发送商户报备查询请求失败: %v", err)
	}

	// 解析响应
	var merchantResponse models.Response
	err = json.Unmarshal([]byte(response), &merchantResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if merchantResponse.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("商户报备查询请求失败: %s", merchantResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(merchantResponse.DataContent, merchantResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return "", fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		fmt.Println("警告：响应签名验证失败")
	}

	fmt.Println("商户报备查询请求成功")

	return merchantResponse.DataContent, nil
}

// BindSubConfig 绑定授权目录
func (s *MerchantService) BindSubConfig(request *models.MerchantBindSubConfigRequest) (string, error) {
	fmt.Println("==========================")
	fmt.Println("绑定授权目录")
	fmt.Println("==========================")

	request.MerId = s.config.MerchantID
	request.TerId = s.config.TerminalID

	// 将业务内容转为JSON
	bizContentJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("业务参数JSON编码失败: %v", err)
	}

	// 生成签名
	signStr, err := utils.Sign(string(bizContentJSON), s.config.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodBindSubConfig)
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("bizContent", string(bizContentJSON))
	mapParams.Set("charset", "UTF-8")
	mapParams.Set("signStr", signStr)
	mapParams.Set("version", "1.0")
	mapParams.Set("format", "json")
	mapParams.Set("signType", "RSA")
	mapParams.Set("signSn", "1")
	mapParams.Set("ncrptnSn", "1")
	mapParams.Set("timestamp", time.Now().Format("20060102150405"))

	// 发送请求
	response, err := s.httpClient.Post(s.getHost(), mapParams)
	if err != nil {
		return "", fmt.Errorf("发送绑定授权目录请求失败: %v", err)
	}

	// 解析响应
	var merchantResponse models.Response
	err = json.Unmarshal([]byte(response), &merchantResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if merchantResponse.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("绑定授权目录请求失败: %s", merchantResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(merchantResponse.DataContent, merchantResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return "", fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		fmt.Println("警告：响应签名验证失败")
	}

	fmt.Println("绑定授权目录请求成功")

	return merchantResponse.DataContent, nil
}
