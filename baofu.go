package baofu

import (
	"github.com/nicoaz/baofu-sdk/config"
	"github.com/nicoaz/baofu-sdk/services"
	"github.com/nicoaz/baofu-sdk/utils"
)

// BaofuClient 宝付支付客户端
type BaofuClient struct {
	Config          *config.Config            // 配置信息
	AccountService  *services.AccountService  // 账户服务
	PaymentService  *services.PaymentService  // 支付服务
	MerchantService *services.MerchantService // 商户报备服务
}

// NewClient 创建宝付支付客户端
// merchantID 商户号
// terminalID 终端号
// privateKey 商户私钥 pem格式，用宝付的 pfx 转 [openssl pkcs12 -in detu1.pfx -nocerts -out detu_private1.pem -nodes]
// publicCert 商户公钥 cer 格式
// bfPubCert 宝付公钥 cer 格式
func NewClient(merchantID, terminalID, privateKey, publicCert, bfPubCert string, opts ...Option) (*BaofuClient, error) {
	// 创建默认配置
	cfg := &config.Config{
		MerchantID: merchantID,
		TerminalID: terminalID,
	}

	// 加载私钥证书
	pri, err := utils.LoadPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	// 加载公钥证书
	pub, err := utils.LoadPublicCert(publicCert)
	if err != nil {
		return nil, err
	}
	// 加载宝付公钥证书
	bfPub, err := utils.LoadPublicCert(bfPubCert)
	if err != nil {
		return nil, err
	}
	// 将公钥证书转换为pem格式
	bfPubPem, err := utils.PublicCert2Pem(bfPub)
	if err != nil {
		return nil, err
	}

	cfg.PrivateKey = pri
	cfg.PublicKey = pub
	cfg.BFPublicKey = bfPub
	cfg.BFPublicKeyPem = bfPubPem

	// 创建客户端实例
	c := &BaofuClient{
		Config:          cfg,
		AccountService:  services.NewAccountService(cfg),
		PaymentService:  services.NewPaymentService(cfg),
		MerchantService: services.NewMerchantService(cfg),
	}

	// 应用选项
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

type Option func(*BaofuClient)

// Release 设置是否为生产环境
func Release(releaseEnv bool) Option {
	return func(c *BaofuClient) {
		c.Config.ReleaseEnv = releaseEnv
	}
}

// Debug 设置是否为调试模式
func Debug(debug bool) Option {
	return func(c *BaofuClient) {
		c.Config.Debug = debug
	}
}
