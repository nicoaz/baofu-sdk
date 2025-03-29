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

// PaymentService 支付服务
type PaymentService struct {
	config     *config.Config
	httpClient *utils.HTTPClient
}

// NewPaymentService 创建支付服务
func NewPaymentService(config *config.Config) *PaymentService {
	return &PaymentService{
		config:     config,
		httpClient: utils.NewHTTPClient(),
	}
}

func (s *PaymentService) getHost() string {
	if s.config.ReleaseEnv {
		return consts.PaymentServiceHostProd
	}
	return consts.PaymentServiceHostTest
}

// CreateUnifiedOrder 创建统一支付订单
func (s *PaymentService) CreateUnifiedOrder(req *models.UnifiedOrderRequest) (*models.UnifiedOrderDataContent, error) {

	// 构建业务内容
	bizContent := models.BizContent{
		MerID:        s.config.MerchantID,
		TerID:        s.config.TerminalID,
		OutTradeNo:   req.OutTradeNo,
		TxnAmt:       req.Amount,
		TxnTime:      utils.GetTimeFormat("YmdHis"),
		TotalAmt:     req.Amount,
		TimeExpire:   "120",
		ProdType:     "SHARING", // SHARING:分账产品,ORDINARY:普通产品
		OrderType:    "7",
		PayCode:      req.PayCode,
		SubMchID:     "756405755", // 子商户号
		NotifyURL:    req.NotifyURL,
		ForbidCredit: req.ForbidCredit, // 1：禁止,0：不禁止
		Attach:       req.Attach,
		RiskInfo: models.RiskInfo{
			LocationPoint: "",
			ClientIP:      req.ClientIP,
		},
	}

	// 设置支付扩展信息
	payExtend := models.PayExtend{}
	if len(req.PayCode) >= 6 && req.PayCode[:6] == "WECHAT" {
		// 微信支付
		payExtend.Body = req.GoodsDesc

		if req.PayCode == "WECHAT_JSAPI" {
			payExtend.SubAppID = req.SubAppID   // 实际项目中应设置正确的AppID
			payExtend.SubOpenID = req.SubOpenID // 实际项目中应设置正确的OpenID
		} else if req.PayCode == "WECHAT_MICROPAY" {
			payExtend.TerminalInfo = &models.TerminalInfo{
				PayCode:  req.PayCode,
				DeviceID: "",
			}
		}
	} else if len(req.PayCode) >= 6 && req.PayCode[:6] == "ALIPAY" {
		// 支付宝
		payExtend.Subject = req.GoodsDesc
		payExtend.AreaInfo = ""

		if req.PayCode == "ALIPAY_JSAPI" {
			payExtend.BuyerID = "" // 实际项目中应设置正确的买家ID
		} else if req.PayCode == "ALIPAY_MICROPAY" {
			payExtend.TerminalInfo = &models.TerminalInfo{
				PayCode:  req.PayCode,
				DeviceID: "",
			}
		}
	}
	bizContent.PayExtend = payExtend

	// 将业务内容转为JSON
	bizContentJSON, err := json.Marshal(bizContent)
	if err != nil {
		return nil, fmt.Errorf("业务参数JSON编码失败: %v", err)
	}

	// 生成签名
	signStr, err := utils.Sign(string(bizContentJSON), s.config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodUnifiedOrder)
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
		return nil, fmt.Errorf("发送支付请求失败: %v", err)
	}

	// 解析响应
	var payResponse models.PayResponse
	err = json.Unmarshal([]byte(response), &payResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if payResponse.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("支付请求失败: %s", payResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(payResponse.DataContent, payResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return nil, fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		return nil, fmt.Errorf("签名验不通过")
	}

	var data models.UnifiedOrderDataContent
	err = json.Unmarshal([]byte(payResponse.DataContent), &data)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &data, nil
}

// QueryOrder 查询订单
// tradeNo 宝付交易号
func (s *PaymentService) QueryOrder(tradeNo string) (*models.QueryOrderData, error) {

	// 构建请求内容
	content := fmt.Sprintf("{\"merId\":\"%s\",\"terId\":\"%s\",\"tradeNo\":\"%s\"}",
		s.config.MerchantID, s.config.TerminalID, tradeNo)

	// 生成签名
	signStr, err := utils.Sign(content, s.config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodOrderQuery)
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("bizContent", content)
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
		return nil, fmt.Errorf("发送订单查询请求失败: %v", err)
	}

	// 解析响应
	var payResponse models.PayResponse
	err = json.Unmarshal([]byte(response), &payResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if payResponse.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("订单查询失败: %s", payResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(payResponse.DataContent, payResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return nil, fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		return nil, fmt.Errorf("签名验不通过")
	}

	var data models.QueryOrderData
	err = json.Unmarshal([]byte(payResponse.DataContent), &data)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &data, nil
}

// CloseOrder 订单关闭
func (s *PaymentService) CloseOrder(outTradeNo string) (*models.CloseOrderData, error) {

	// 构建请求内容
	content := fmt.Sprintf("{\"merId\":\"%s\",\"terId\":\"%s\",\"tradeNo\":\"%s\"}",
		s.config.MerchantID, s.config.TerminalID, outTradeNo)

	// 生成签名
	signStr, err := utils.Sign(content, s.config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodOrderQuery)
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("bizContent", content)
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
		return nil, fmt.Errorf("发送退款请求失败: %v", err)
	}

	// 解析响应
	var payResponse models.PayResponse
	err = json.Unmarshal([]byte(response), &payResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if payResponse.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("退款请求失败: %s", payResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(payResponse.DataContent, payResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return nil, fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		return nil, fmt.Errorf("签名验不通过")
	}

	var data models.CloseOrderData
	err = json.Unmarshal([]byte(payResponse.DataContent), &data)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &data, nil
}

// RefundOrder 退款请求
func (s *PaymentService) RefundOrder(req *models.RefundRequest) (*models.RefundResponse, error) {

	// 构建请求内容
	req.MerId = s.config.MerchantID
	req.TerId = s.config.TerminalID
	b, _ := json.Marshal(req)
	content := string(b)

	// 生成签名
	signStr, err := utils.Sign(content, s.config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodOrderRefund)
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("bizContent", content)
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
		return nil, fmt.Errorf("发送退款请求失败: %v", err)
	}

	// 解析响应
	var payResponse models.PayResponse
	err = json.Unmarshal([]byte(response), &payResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if payResponse.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("退款请求失败: %s", payResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(payResponse.DataContent, payResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return nil, fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		return nil, fmt.Errorf("签名验不通过")
	}

	var data models.RefundResponse
	err = json.Unmarshal([]byte(payResponse.DataContent), &data)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &data, nil
}

// QueryRefundOrder 查询退款订单
// tradeNo 商户系统内部退款订单号
func (s *PaymentService) QueryRefundOrder(outTradeNo string) (*models.RefundQueryData, error) {

	// 构建请求内容
	content := fmt.Sprintf("{\"merId\":\"%s\",\"terId\":\"%s\",\"outTradeNo\":\"%s\"}",
		s.config.MerchantID, s.config.TerminalID, outTradeNo)

	// 生成签名
	signStr, err := utils.Sign(content, s.config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", consts.MethodRefundQuery)
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("bizContent", content)
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
		return nil, fmt.Errorf("发送退款订单查询请求失败: %v", err)
	}

	// 解析响应
	var payResponse models.PayResponse
	err = json.Unmarshal([]byte(response), &payResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if payResponse.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("退款订单查询失败: %s", payResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(payResponse.DataContent, payResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return nil, fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		return nil, fmt.Errorf("签名验不通过")
	}

	var data models.RefundQueryData
	err = json.Unmarshal([]byte(payResponse.DataContent), &data)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &data, nil
}

// VerifyNotify 验证异步通知
func (s *PaymentService) VerifyNotify(notifyData, signature string) (bool, error) {
	// 验证签名
	verify, err := utils.VerifySign(notifyData, signature, s.config.BFPublicKey)
	if err != nil {
		return false, fmt.Errorf("通知签名验证失败: %v", err)
	}
	if !verify {
		fmt.Println("警告：通知签名验证失败")
	}

	return true, nil
}
