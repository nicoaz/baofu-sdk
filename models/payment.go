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

type QueryOrderData struct {
	AgentMerID  string   `json:"agentMerId"` // 代理商商户号
	AgentTerID  string   `json:"agentTerId"` // 代理商终端号
	MerID       string   `json:"merId"`      // 商户号
	TerID       string   `json:"terId"`      // 终端号
	TradeNo     string   `json:"tradeNo"`    // 宝付交易号
	OutTradeNo  string   `json:"outTradeNo"` // 商户订单号
	TxnState    TxnState `json:"txnState"`   // 订单状态
	FinishTime  string   `json:"finishTime"` // 完成时间
	SuccAmt     int      `json:"succAmt"`    // 成功金额
	FeeAmt      int      `json:"feeAmt"`     // 支付手续费
	InstFeeAmt  int      `json:"instFeeAmt"` // 分期手续费
	ResultCode  string   `json:"resultCode"` // 业务结果
	ErrCode     string   `json:"errCode"`    // 错误代码
	ErrMsg      string   `json:"errMsg"`     // 错误描述
	ReqChlNo    string   `json:"reqChlNo"`   // 请求渠道订单号
	PayCode     string   `json:"payCode"`    // 支付方式
	ChlRetParam struct {
		// 支付宝返回
		BuyerID   string `json:"buyerId"`   // 买家支付宝用户号
		AccountNo string `json:"accountNo"` // 支付宝会返回
		// 微信返回
		OpenID    string `json:"openId"`    // 用户在服务商公众号appid下的唯一标识
		SubOpenID string `json:"subOpenid"` // 微信平台的sub_openid
	} `json:"chlRetParam"` // 渠道返回参数
	ClearingDate string `json:"clearingDate"` // 清算日期
}

type CloseOrderData struct {
	AgentMerID string `json:"agentMerId"` // 代理商商户号
	AgentTerID string `json:"agentTerId"` // 代理商终端号
	MerID      string `json:"merId"`      // 商户号
	TerID      string `json:"terId"`      // 终端号
	TradeNo    string `json:"tradeNo"`    // 宝付订单号
	OutTradeNo string `json:"outTradeNo"` // 商户订单号
	ResultCode string `json:"resultCode"` // 业务结果 SUCCESS：成功，FAIL：失败
	ErrCode    string `json:"errCode"`    // 错误代码 当业务结果FAIL时，返回错误代码
	ErrMsg     string `json:"errMsg"`     // 错误描述 当业务结果为FAIL时，返回错误描述
}

// RefundRequest 退款请求
type RefundRequest struct {
	// AgentMerId       string `json:"agentMerId,omitempty"`       // 代理商商户号
	// AgentTerId       string `json:"agentTerId,omitempty"`       // 代理商终端号
	MerId         string `json:"merId"`                   // 商户号
	TerId         string `json:"terId"`                   // 终端号
	MerchantName  string `json:"merchantName,omitempty"`  // 商户名称
	OriginTradeNo string `json:"originTradeNo,omitempty"` // 原支付订单宝付交易号
	//OriginOutTradeNo string `json:"originOutTradeNo,omitempty"` // 原支付订单商户订单号
	OutTradeNo string `json:"outTradeNo"`          // 退款订单号
	NotifyUrl  string `json:"notifyUrl,omitempty"` // 服务端通知地址
	RefundAmt  int    `json:"refundAmt"`           // 退款金额 单位：分，退款金额不得大于用户实际付款金额
	TotalAmt   int    `json:"totalAmt"`            // 退款总金额 如包含营销信息，则退款总金额=退款金额+营销退款总金额，反之退款总金额=退款金额
	TxnTime    string `json:"txnTime"`             // 交易时间 订单交易时间

	SharingRefundInfo []SharingRefundInfo `json:"sharingRefundInfo,omitempty"` // 分账退款信息
	// MktRefundInfo     []MktRefundInfo     `json:"mktRefundInfo,omitempty"`     // 营销退款信息
	//AdvanceAmt        int                 `json:"advanceAmt,omitempty"`        // 垫资金额
	RefundReason string `json:"refundReason"` // 退款原因
}

type SharingRefundInfo struct {
	SharingMerId string `json:"sharingMerId"` // 宝付支付分配的商户号
	SharingAmt   int    `json:"sharingAmt"`   // 分账金额，单位：分，如1元则传入100
}

type MktRefundInfo struct {
	MktMerId string `json:"mktMerId"` // 宝付支付分配的商户号
	MktAmt   int    `json:"mktAmt"`   // 分账金额，单位：分，如1元则传入100
}

type RefundResponse struct {
	OriginTradeNo    string      `json:"originTradeNo"`    // 原支付订单宝付交易号
	OriginOutTradeNo string      `json:"originOutTradeNo"` // 原支付订单商户订单号
	OutTradeNo       string      `json:"outTradeNo"`       // 商户退款订单号
	TradeNo          string      `json:"tradeNo"`          // 宝付退款交易号
	RefundAmt        int         `json:"refundAmt"`        // 退款金额
	TotalAmt         int         `json:"totalAmt"`         // 退款总金额
	ResultCode       string      `json:"resultCode"`       // 业务结果 SUCCESS
	RefundState      RefundState `json:"refundState"`      // 订单状态 REFUND
	ErrCode          string      `json:"errCode"`          // 错误代码
	ErrMsg           string      `json:"errMsg"`           // 错误描述
	ReqReserved      string      `json:"reqReserved"`      // 请求方保留域
}

type RefundState string

const (
	RefundStateSuccess     RefundState = "SUCCESS"      // 退款成功
	RefundStateRefund      RefundState = "REFUND"       // 退款受理成功
	RefundStateRefundError RefundState = "REFUND_ERROR" // 退款失败
	RefundStateAbnormal    RefundState = "ABNORMAL"     // 退款异常，返回此状态的退款订单，请稍后发起查询。
)

type RefundQueryData struct {
	TradeNo     string      `json:"tradeNo"`     // 宝付订单号
	OutTradeNo  string      `json:"outTradeNo"`  // 商户订单号
	RefundState RefundState `json:"refundState"` // 订单状态
	FinishTime  string      `json:"finishTime"`  // 完成时间
	SuccAmt     string      `json:"succAmt"`     // 成功金额 单位：分，订单状态为成功时才有值
	ResultCode  string      `json:"resultCode"`  // 业务结果 SUCCESS：成功 FAIL：失败
	ErrCode     string      `json:"errCode"`     // 错误代码 当业务结果FAIL时，返回错误代码
	ErrMsg      string      `json:"errMsg"`      // 错误描述 当业务结果为FAIL时，返回错误描述
}
