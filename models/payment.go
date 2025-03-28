package models

type UnifiedOrderRequest struct {
	OutTradeNo   string // 商户订单号
	Amount       int    // 交易金额（分）
	PayCode      string // 支付方式
	GoodsDesc    string // 商品描述
	ClientIP     string // 客户端IP
	SubAppID     string // 收单二级商户ID
	SubOpenID    string // 二级商户下用户OpenID
	Attach       string // 附加数据
	ForbidCredit string // 是否禁止信用卡支付 1禁止 0或者""不禁用
	NotifyURL    string // 异步通知地址
}

// BizContent 业务参数
type BizContent struct {
	AgentMerID   string    `json:"agentMerId,omitempty"` // 代理商户号
	AgentTerID   string    `json:"agentTerId,omitempty"` // 代理终端号
	MerID        string    `json:"merId"`                // 商户号
	TerID        string    `json:"terId"`                // 终端号
	OutTradeNo   string    `json:"outTradeNo"`           // 商户订单号
	TxnAmt       int       `json:"txnAmt"`               // 交易金额（分）
	TxnTime      string    `json:"txnTime"`              // 交易时间
	TotalAmt     int       `json:"totalAmt"`             // 订单总金额（分）
	TimeExpire   string    `json:"timeExpire"`           // 订单有效期（分钟）
	ProdType     string    `json:"prodType"`             // 产品类型
	OrderType    string    `json:"orderType"`            // 订单类型 7
	PayCode      string    `json:"payCode"`              // 支付方式
	PayExtend    PayExtend `json:"payExtend"`            // 支付扩展信息
	SubMchID     string    `json:"subMchId,omitempty"`   // 子商户号
	NotifyURL    string    `json:"notifyUrl"`            // 异步通知地址
	ForbidCredit string    `json:"forbidCredit"`         // 是否禁止信用卡支付
	Attach       string    `json:"attach,omitempty"`     // 附加数据
	RiskInfo     RiskInfo  `json:"riskInfo,omitempty"`   // 风控信息
}

// PayExtend 支付扩展信息
type PayExtend struct {
	// 微信支付参数
	Body         string        `json:"body,omitempty"`          // 商品描述
	SubAppID     string        `json:"sub_appid,omitempty"`     // 子应用ID
	SubOpenID    string        `json:"sub_openid,omitempty"`    // 子用户OpenID
	TerminalInfo *TerminalInfo `json:"terminal_info,omitempty"` // 终端信息

	// 支付宝参数
	Subject  string `json:"subject,omitempty"`   // 商品名称
	AreaInfo string `json:"area_info,omitempty"` // 区域信息
	BuyerID  string `json:"buyer_id,omitempty"`  // 买家ID
}

// TerminalInfo 终端信息
type TerminalInfo struct {
	PayCode  string `json:"pay_code"`  // 支付方式
	DeviceID string `json:"device_id"` // 设备ID
}

// RiskInfo 风控信息
type RiskInfo struct {
	LocationPoint string `json:"locationPoint"` // 地理位置
	ClientIP      string `json:"clientIp"`      // 客户端IP
}

// PayResponse 支付响应
type PayResponse struct {
	ReturnCode  string `json:"returnCode"`  // 返回码
	ReturnMsg   string `json:"returnMsg"`   // 返回信息
	DataContent string `json:"dataContent"` // 返回数据（JSON字符串）
	SignStr     string `json:"signStr"`     // 签名字符串
}

type UnifiedOrderDataContent struct {
	AgentMerID  string      `json:"agentMerId"`  // 代理商商户号
	AgentTerID  string      `json:"agentTerId"`  // 代理商终端号
	MerID       string      `json:"merId"`       // 商户号
	TerID       string      `json:"terId"`       // 终端号
	OutTradeNo  string      `json:"outTradeNo"`  // 商户订单号
	TxnState    TxnState    `json:"txnState"`    // 订单状态 SUCCESS WAIT_PAYING
	TradeNo     string      `json:"tradeNo"`     // 宝付交易号 宝付交易号
	ReqChlNo    string      `json:"reqChlNo"`    // 请求渠道订单号
	PayCode     string      `json:"payCode"`     // 支付方式
	ChlRetParam ChlRetParam `json:"chlRetParam"` // 渠道返回参数
	ResultCode  string      `json:"resultCode"`  // 业务结果
	ErrCode     string      `json:"errCode"`     // 错误代码
	ErrMsg      string      `json:"errMsg"`      // 错误描述
}

type ChlRetParam struct {
	WcPayData string `json:"wc_pay_data"` // 微信支付参数
	OrderID   int    `json:"order_id"`    // 订单ID
	PrepayID  string `json:"prepay_id"`   // 预支付ID
}

// QueryOrderRequest 订单查询请求
type QueryOrderRequest struct {
	Method     string `json:"method"`     // 接口方法名
	MerID      string `json:"merId"`      // 商户号
	TerID      string `json:"terId"`      // 终端号
	OutTradeNo string `json:"outTradeNo"` // 商户订单号
	SignStr    string `json:"signStr"`    // 签名字符串
	Version    string `json:"version"`    // 版本号
	Timestamp  string `json:"timestamp"`  // 时间戳
}

// RefundRequest 退款请求
type RefundRequest struct {
	Method       string `json:"method"`       // 接口方法名
	MerID        string `json:"merId"`        // 商户号
	TerID        string `json:"terId"`        // 终端号
	OutTradeNo   string `json:"outTradeNo"`   // 原商户订单号
	OutRefundNo  string `json:"outRefundNo"`  // 商户退款单号
	RefundAmount int    `json:"refundAmount"` // 退款金额（分）
	RefundReason string `json:"refundReason"` // 退款原因
	SignStr      string `json:"signStr"`      // 签名字符串
	Version      string `json:"version"`      // 版本号
	Timestamp    string `json:"timestamp"`    // 时间戳
}
