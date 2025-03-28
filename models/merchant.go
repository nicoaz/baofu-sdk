package models

// PayResponse 支付响应
type Response struct {
	ReturnCode  string `json:"returnCode"`  // 返回码
	ReturnMsg   string `json:"returnMsg"`   // 返回信息
	DataContent string `json:"dataContent"` // 返回数据（JSON字符串）
	SignStr     string `json:"signStr"`     // 签名字符串
}

type MerchantWXReportReq struct {
	//AgentMerId string `json:"agentMerId"` // 代理商商户号
	//AgentTerId string `json:"agentTerId"` // 代理商终端号
	MerId      string     `json:"merId"`      // 交易商户号
	TerId      string     `json:"terId"`      // 交易终端号
	ReportType string     `json:"reportType"` // 报备类型 WECHAT
	ReportNo   string     `json:"reportNo"`   // 报备编号
	ReportInfo ReportInfo `json:"reportInfo"` // 报备信息
	BctMerId   string     `json:"bctMerId"`   // 宝财通二级商户号
}

type ReportInfo struct {
	MerchantName      string `json:"merchant_name"`      // 商户名称
	MerchantShortname string `json:"merchant_shortname"` // 商户简称
	ServicePhone      string `json:"service_phone"`      // 客服电话
	//Contact             string       `json:"contact"`               // 联系人
	//ContactPhone        string       `json:"contact_phone"`         // 联系电话
	//ContactEmail        string       `json:"contact_email"`         // 联系邮箱
	ChannelId   string `json:"channel_id"`   // 渠道商商户号
	ChannelName string `json:"channel_name"` // 渠道商商户名称
	Business    string `json:"business"`     // 经营类目
	//ContactWechatidType string       `json:"contact_wechatid_type"` // 联系人微信账号类型
	//	ContactWechatid     string       `json:"contact_wechatid"`      // 联系人微信帐号
	ServiceCodes        []string     `json:"service_codes"`         // 申请服务，可传送所有需要开通的服务，详见《微信服务类型》，JSON数组格式 [“JSAPI”,”APPLET”]
	AddressInfo         AddressInfo  `json:"address_info"`          // 地址信息，JSON格式
	BusinessLicense     string       `json:"business_license"`      // 商户证件编号
	BusinessLicenseType string       `json:"business_license_type"` // 商户证件类型
	BankcardInfo        BankcardInfo `json:"bankcard_info"`         // 银行结算卡信息
}

type AddressInfo struct {
	CityCode     string `json:"cityCode"`     // 城市编码
	DistrictCode string `json:"districtCode"` // 区县编码
	ProvinceCode string `json:"provinceCode"` // 省份编码
	Address      string `json:"address"`      // 详细地址
	// Longitude    string `json:"longitude"`    // 经度
	// Latitude     string `json:"latitude"`     // 纬度
	// AddressType  string `json:"addressType"`  // 地址类型
}

type BankcardInfo struct {
	CardNo   string `json:"cardNo"`   // 银行卡号
	CardName string `json:"cardName"` // 银行卡持卡人姓名
	BankName string `json:"bankName"` // 银行开户行名称
}

// MerchantReportQueryRequest 商户报备查询请求参数
type MerchantReportQueryRequest struct {
	//AgentMerId  string `json:"agentMerId"`  // 代理商商户号
	//AgentTerId  string `json:"agentTerId"`  // 代理商终端号
	MerId      string `json:"merId"`      // 交易商户号
	TerId      string `json:"terId"`      // 交易终端号
	ReportType string `json:"reportType"` // 报备类型 WECHAT
	ReportNo   string `json:"reportNo"`   // 报备编号
}

// MerchantBindSubConfigRequest 绑定授权目录请求参数
type MerchantBindSubConfigRequest struct {
	//AgentMerId  string `json:"agentMerId"`  // 代理商商户号
	//AgentTerId  string `json:"agentTerId"`  // 代理商终端号
	MerId       string `json:"merId"`       // 交易商户号
	TerId       string `json:"terId"`       // 交易终端号
	SubMchId    string `json:"subMchId"`    // 商户识别码
	AuthType    string `json:"authType"`    // 授权类型 AUTH JSAPI APPLET
	AuthContent string `json:"authContent"` // 授权内容
	Remark      string `json:"remark"`      // 备注
}
