package models

// AccountResponse 账户操作的通用响应
type AccountResponse struct {
	ReturnCode  string `json:"returnCode"`  // 返回码
	ReturnMsg   string `json:"returnMsg"`   // 返回信息
	DataContent string `json:"dataContent"` // 返回数据（加密）
	SignStr     string `json:"signStr"`     // 签名字符串
}

// AccountRequest 账户操作的通用请求
type AccountRequest struct {
	Header map[string]string      `json:"header"` // 请求头
	Body   map[string]interface{} `json:"body"`   // 请求体
}

// AccountOpenRequest 开户请求参数
type AccountOpenRequest struct {
	NotifyUrl string `json:"notifyUrl"` // 异步通知地址

	// 商户账户信息
	TransSerialNo       string `json:"transSerialNo"`       // 请求流水号
	LoginNo             string `json:"loginNo"`             // 登录号(用户ID)建议长度8位以上 唯一
	Email               string `json:"email"`               // 邮箱
	SelfEmployed        string `json:"selfEmployed"`        // 是否个体户 默认为false
	CustomerName        string `json:"customerName"`        // 商户名称（营业执照上的名称）
	AliasName           string `json:"aliasName"`           // 商户名称别名  (选传)
	CertificateNo       string `json:"certificateNo"`       // 证件号码
	CertificateType     string `json:"certificateType"`     // 证件类型 营业执照:LICENSE
	CorporateName       string `json:"corporateName"`       // 法人姓名
	CorporateCertType   string `json:"corporateCertType"`   // 法人证件类型  身份证:ID
	CorporateCertId     string `json:"corporateCertId"`     // 法人身份证号码
	CorporateMobile     string `json:"corporateMobile"`     // 法人手机号
	IndustryId          string `json:"industryId"`          // 所属行业
	ContactName         string `json:"contactName"`         // 联系人姓名
	ContactMobile       string `json:"contactMobile"`       // 联系人手机号
	CardNo              string `json:"cardNo"`              // 卡号
	BankName            string `json:"bankName"`            // 银行名称
	DepositBankProvince string `json:"depositBankProvince"` // 开户行省份
	DepositBankCity     string `json:"depositBankCity"`     // 开户行城市
	DeositBankName      string `json:"depositBankName"`     // 开户支行名称
	RegisterCapital     string `json:"registerCapital"`     // 注册资本
	CardUserName        string `json:"cardUserName"`        // 持卡人姓名  (个体绑法人对私必传)

}

// OpenAccountQueryRequest 开户查询请求参数
type OpenAccountQueryRequest struct {
	CertificateNo   string `json:"certificateNo"`   // 证件号码（社会信用代码）
	CertificateType string `json:"certificateType"` // 证件类型 只能取”ID”或”LICENSE”
	PlatformNo      string `json:"platformNo"`      // 平台号(主商户号)
	LoginNo         string `json:"loginNo"`         // 登录号(传此参数以上三个参数必填)
	AccType         string `json:"accType"`         // 账户类型 1个人，2商户
}

// AccountModifyRequest 修改账户请求参数
type AccountModifyRequest struct {
	// 请求头
	MemberId   string `json:"memberId"`   // 商户号
	TerminalId string `json:"terminalId"` // 终端号
	ServiceTp  string `json:"serviceTp"`  // 服务类型
	VerifyType string `json:"verifyType"` // 验证类型

	// 请求体
	Version         string `json:"version"`         // 版本号
	ContractNo      string `json:"contractNo"`      // 合同号（客户账户）
	AcctType        string `json:"acctType"`        // 账户类型：1个人，2机构
	AcctName        string `json:"acctName"`        // 账户名称
	LegalPersonName string `json:"legalPersonName"` // 法人姓名
	BizLicenseCode  string `json:"bizLicenseCode"`  // 营业执照号
	Phone           string `json:"phone"`           // 手机号
	Email           string `json:"email"`           // 邮箱
	IdCardType      string `json:"idCardType"`      // 证件类型
	IdCardCode      string `json:"idCardCode"`      // 证件号码
	BankCode        string `json:"bankCode"`        // 银行编码
	CardType        string `json:"cardType"`        // 卡类型
	BankCardNo      string `json:"bankCardNo"`      // 银行卡号
	ExpireDate      string `json:"expireDate"`      // 有效期
	CVV             string `json:"cvv"`             // 安全码
}

// BalanceQueryRequest 余额查询请求参数
type BalanceQueryRequest struct {
	AcctType   string `json:"acctType"`   // 账户类型：1个人，2商户
	ContractNo string `json:"contractNo"` // 合同号（客户账户）
}

type BalanceQueryResponse struct {
	Body struct {
		RetCode      int     `json:"retCode"`      // 返回码 1 成功 0 失败
		ErrorCode    string  `json:"errorCode"`    // 错误码
		ErrorMsg     string  `json:"errorMsg"`     // 错误原因
		AvailableBal float64 `json:"availableBal"` // 账簿可用余额,单位：元;可用于提现
		PendingBal   float64 `json:"pendingBal"`   // 在途资金余额,单位：元
		CurrBal      float64 `json:"currBal"`      // 账簿余额,单位：元;账簿余额=可用余额(availableBal)+在途余额(pendingBal)+冻结金额
	} `json:"body"`
	Header struct {
		MemberId    string `json:"memberId"`    // 商户号
		TerminalId  string `json:"terminalId"`  // 终端号
		ServiceTp   string `json:"serviceTp"`   // 服务类型
		SysRespCode string `json:"sysRespCode"` // 返回码
		SysRespDesc string `json:"sysRespDesc"` // 返回信息
	} `json:"header"`
}

// TransferRequest 转账请求参数
type TransferRequest struct {
	PayerNo       string  `json:"payerNo"`       // 付款方账号
	PayeeNo       string  `json:"payeeNo"`       // 收款方账号
	TransSerialNo string  `json:"transSerialNo"` // 交易流水号
	DealAmount    float64 `json:"dealAmount"`    // 交易金额 BigDecimal 单位元
}

type TransferResponse struct {
	Body struct {
		RetCode       int     `json:"retCode"`       // 返回码 1 成功 0 失败
		ErrorCode     string  `json:"errorCode"`     // 错误码
		ErrorMsg      string  `json:"errorMsg"`      // 错误原因
		TransSerialNo string  `json:"transSerialNo"` // 请求流水号
		BusinessNo    string  `json:"businessNo"`    // 业务流水号
		PayerNo       string  `json:"payerNo"`       // 付款方(二级子商户号)
		PayeeNo       string  `json:"payeeNo"`       // 收款方(二级子商户号)
		DealAmount    float64 `json:"dealAmount"`    // 转账金额,单位：元
		FeeAmount     float64 `json:"feeAmount"`     // 手续费金额,单位：元
		State         int     `json:"state"`         // 订单状态 1成功 2失败
		TransRemark   string  `json:"transRemark"`   // 失败原因
	} `json:"body"`
	Header struct {
		MemberId    string `json:"memberId"`    // 商户号
		TerminalId  string `json:"terminalId"`  // 终端号
		ServiceTp   string `json:"serviceTp"`   // 服务类型
		SysRespCode string `json:"sysRespCode"` // 返回码
		SysRespDesc string `json:"sysRespDesc"` // 返回信息
	} `json:"header"`
}

// WithdrawRequest 提现请求参数
type WithdrawRequest struct {
	// 请求头
	MemberId   string `json:"memberId"`   // 商户号
	TerminalId string `json:"terminalId"` // 终端号
	ServiceTp  string `json:"serviceTp"`  // 服务类型
	VerifyType string `json:"verifyType"` // 验证类型

	// 请求体
	Version        string `json:"version"`        // 版本号
	SrcAcctNo      string `json:"srcAcctNo"`      // 源账户号
	TransAmt       string `json:"transAmt"`       // 转账金额
	TransId        string `json:"transId"`        // 交易ID
	TransDate      string `json:"transDate"`      // 交易日期
	TransTime      string `json:"transTime"`      // 交易时间
	CurType        string `json:"curType"`        // 币种
	TransSummary   string `json:"transSummary"`   // 交易摘要
	ReservedExpand string `json:"reservedExpand"` // 扩展字段
	CardNo         string `json:"cardNo"`         // 银行卡号
	CardName       string `json:"cardName"`       // 银行卡姓名
	CardBankCode   string `json:"cardBankCode"`   // 银行编码
	DirectFlag     string `json:"directFlag"`     // 直连标志

	NotifyUrl string `json:"notifyUrl"` // 异步通知地址
}
