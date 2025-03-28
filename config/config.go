package config

import "crypto/rsa"

// Config 宝付支付SDK配置
type Config struct {

	// 环境
	ReleaseEnv bool // 是否为生产环境
	Debug      bool // 是否为调试模式
	// 商户信息
	MerchantID string // 商户号
	TerminalID string // 终端号

	// 证书信息

	PrivateKey   *rsa.PrivateKey // 私钥
	PublicKey    *rsa.PublicKey  // 公钥
	PublicKeyPem []byte          // 公钥证书内容

	BFPublicKey    *rsa.PublicKey // 宝付公钥
	BFPublicKeyPem []byte         // 宝付公钥证书内容
	// CertPath     string          // 公钥证书路径
	// PfxPath      string          // 私钥证书路径
	// KeyPassword  string          // 证书密码
}
