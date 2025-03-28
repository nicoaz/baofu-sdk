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
func (s *PaymentService) CreateUnifiedOrder(amount int, payCode, goodsDesc, clientIP, SubAppID, SubOpenID, notifyURL string) (*models.UnifiedOrderDataContent, error) {

	// 构建业务内容
	bizContent := models.BizContent{
		MerID:        s.config.MerchantID,
		TerID:        s.config.TerminalID,
		OutTradeNo:   utils.GetTransid("pay"),
		TxnAmt:       amount,
		TxnTime:      utils.GetTimeFormat("YmdHis"),
		TotalAmt:     amount,
		TimeExpire:   "120",
		ProdType:     "SHARING", // SHARING:分账产品,ORDINARY:普通产品
		OrderType:    "7",
		PayCode:      payCode,
		SubMchID:     "756405755", // 子商户号
		NotifyURL:    notifyURL,
		ForbidCredit: "0", // 1：禁止,0：不禁止
		RiskInfo: models.RiskInfo{
			LocationPoint: "",
			ClientIP:      clientIP,
		},
	}

	// 设置支付扩展信息
	payExtend := models.PayExtend{}
	if len(payCode) >= 6 && payCode[:6] == "WECHAT" {
		// 微信支付
		payExtend.Body = goodsDesc

		if payCode == "WECHAT_JSAPI" {
			payExtend.SubAppID = SubAppID   // 实际项目中应设置正确的AppID
			payExtend.SubOpenID = SubOpenID // 实际项目中应设置正确的OpenID
		} else if payCode == "WECHAT_MICROPAY" {
			payExtend.TerminalInfo = &models.TerminalInfo{
				PayCode:  payCode,
				DeviceID: "",
			}
		}
	} else if len(payCode) >= 6 && payCode[:6] == "ALIPAY" {
		// 支付宝
		payExtend.Subject = goodsDesc
		payExtend.AreaInfo = ""

		if payCode == "ALIPAY_JSAPI" {
			payExtend.BuyerID = "" // 实际项目中应设置正确的买家ID
		} else if payCode == "ALIPAY_MICROPAY" {
			payExtend.TerminalInfo = &models.TerminalInfo{
				PayCode:  payCode,
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
	mapParams.Set("method", "unified_order")
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

	var unifiedOrderDataContent models.UnifiedOrderDataContent
	err = json.Unmarshal([]byte(payResponse.DataContent), &unifiedOrderDataContent)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &unifiedOrderDataContent, nil
}

// QueryOrder 查询订单
func (s *PaymentService) QueryOrder(outTradeNo string) (string, error) {
	fmt.Println("==========================")
	fmt.Println("订单查询")
	fmt.Println("==========================")

	// 构建请求内容
	content := fmt.Sprintf("{\"merId\":\"%s\",\"terId\":\"%s\",\"outTradeNo\":\"%s\"}",
		s.config.MerchantID, s.config.TerminalID, outTradeNo)

	// 生成签名
	signStr, err := utils.Sign(content, s.config.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", "query_order")
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("outTradeNo", outTradeNo)
	mapParams.Set("signStr", signStr)
	mapParams.Set("version", "1.0")
	mapParams.Set("timestamp", time.Now().Format("20060102150405"))

	// 发送请求
	response, err := s.httpClient.Post(s.getHost(), mapParams)
	if err != nil {
		return "", fmt.Errorf("发送订单查询请求失败: %v", err)
	}

	// 解析响应
	var payResponse models.PayResponse
	err = json.Unmarshal([]byte(response), &payResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if payResponse.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("订单查询失败: %s", payResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(payResponse.DataContent, payResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return "", fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		fmt.Println("警告：响应签名验证失败")
	}

	fmt.Println("订单查询成功")

	return payResponse.DataContent, nil
}

// RefundOrder 退款请求
func (s *PaymentService) RefundOrder(outTradeNo, outRefundNo string, refundAmount int, refundReason string) (string, error) {
	fmt.Println("==========================")
	fmt.Println("退款请求")
	fmt.Println("==========================")

	// 构建请求内容
	content := fmt.Sprintf("{\"merId\":\"%s\",\"terId\":\"%s\",\"outTradeNo\":\"%s\",\"outRefundNo\":\"%s\",\"refundAmount\":%d,\"refundReason\":\"%s\"}",
		s.config.MerchantID, s.config.TerminalID, outTradeNo, outRefundNo, refundAmount, refundReason)

	// 生成签名
	signStr, err := utils.Sign(content, s.config.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	// 构建请求参数
	mapParams := url.Values{}
	mapParams.Set("method", "refund")
	mapParams.Set("merId", s.config.MerchantID)
	mapParams.Set("terId", s.config.TerminalID)
	mapParams.Set("outTradeNo", outTradeNo)
	mapParams.Set("outRefundNo", outRefundNo)
	mapParams.Set("refundAmount", fmt.Sprintf("%d", refundAmount))
	mapParams.Set("refundReason", refundReason)
	mapParams.Set("signStr", signStr)
	mapParams.Set("version", "1.0")
	mapParams.Set("timestamp", time.Now().Format("20060102150405"))

	// 发送请求
	response, err := s.httpClient.Post(s.getHost(), mapParams)
	if err != nil {
		return "", fmt.Errorf("发送退款请求失败: %v", err)
	}

	// 解析响应
	var payResponse models.PayResponse
	err = json.Unmarshal([]byte(response), &payResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查返回码
	if payResponse.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("退款请求失败: %s", payResponse.ReturnMsg)
	}

	// 验证响应签名
	verify, err := utils.VerifySign(payResponse.DataContent, payResponse.SignStr, s.config.BFPublicKey)
	if err != nil {
		return "", fmt.Errorf("验证响应签名失败: %v", err)
	}
	if !verify {
		fmt.Println("警告：响应签名验证失败")
	}

	fmt.Println("退款请求成功")

	return payResponse.DataContent, nil
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
