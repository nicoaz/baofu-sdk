package consts

// 账户服务
const (
	AccountServiceHostTest = "https://vgw.baofoo.com/union-gw/api/{报文编号}/transReq.do"
	AccountServiceHostProd = "https://public.baofu.com/union-gw/api/{报文编号}/transReq.do"

	// 开户
	MethodOpenAccount = "T-1001-013-01"
	// 开户查询 T-1001-013-03
	MethodOpenAccountQuery = "T-1001-013-03"
	// 余额查询 T-1001-013-06
	MethodBalanceQuery = "T-1001-013-06"
	// 账户提现 T-1001-013-14
	MethodWithdraw = "T-1001-013-14"
	// 提现查询 T-1001-013-15
	MethodWithdrawQuery = "T-1001-013-15"
	// 账户间转账 T-1001-013-13
	MethodTransfer = "T-1001-013-13"
	// 账户间转账查询 T-1001-013-10
	MethodTransferQuery = "T-1001-013-10"
)

// 聚合支付服务
const (
	PaymentServiceHostTest = "https://mch-juhe.baofoo.com/api"
	PaymentServiceHostProd = "https://juhe.baofoo.com/api"

	// 统一下单
	MethodUnifiedOrder = "unified_order"
	// 分账
	MethodShareAfterPayOrder = "share_after_pay"
	// 关闭订单
	MethodOrderClose = "order_close"
	// 退款
	MethodOrderRefund = "order_refund"
	// 支付订单查询
	MethodOrderQuery = "order_query"
	// share_query 分账订单查询
	MethodShareQuery = "share_query"
	// refund_query 退款订单查询
	MethodRefundQuery = "refund_query"
)

// 聚合报备服务
const (
	ReportServiceHostTest = "https://mch-juhe.baofoo.com/mch-service/api"
	ReportServiceHostProd = "https://juhe.baofoo.com/mch-service/api"

	// 报备认证 merchant_report
	MethodMerchantReport = "merchant_report"
	// 报备信息查询 merchant_report_query
	MethodMerchantReportQuery = "merchant_report_query"
	// 绑定授权目录 bind_sub_config
	MethodBindSubConfig = "bind_sub_config"
)
