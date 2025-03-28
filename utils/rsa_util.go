package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/farmerx/gorsa"
	"golang.org/x/crypto/pkcs12"
)

// LoadPFX 加载pfx文件
func LoadPFX(pfxPath string, password string) (*rsa.PrivateKey, error) {
	certBytes, err := os.ReadFile(pfxPath)
	if err != nil {
		return nil, err
	}
	pkey, _, _ := pkcs12.Decode(certBytes, password)
	privateKey, ok := pkey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parse private key fail")
	}

	return privateKey, nil
}

// LoadPrivateKey 加载 pem 格式私钥
func LoadPrivateKey(pem string) (*rsa.PrivateKey, error) {
	return DecodePrivateKey([]byte(pem))
}

// LoadPublicCert 加载 cer 格式公钥
func LoadPublicCert(cer string) (*rsa.PublicKey, error) {
	publicKey, err := DecodePublicKey([]byte(cer))
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// PublicCert2Pem 将cer格式公钥转换为pem格式
func PublicCert2Pem(publicKey *rsa.PublicKey) ([]byte, error) {

	// 编码为 PEM 格式
	pubKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	pubKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyDER,
	}

	return pem.EncodeToMemory(pubKeyPEM), nil
}

// EncryptByPFXFile 使用PFX文件加密
func EncryptByPFXFile(content string, privateKey *rsa.PrivateKey) (string, error) {

	if privateKey == nil {
		return "", errors.New("private key is nil")
	}

	// 先base64编码
	content = base64.StdEncoding.EncodeToString([]byte(content))

	// 使用私钥签名数据
	// 分段加密
	blockSize := privateKey.Size() - 11
	signatures := ""
	for i := 0; i < len(content); i += blockSize {
		end := i + blockSize
		if end > len(content) {
			end = len(content)
		}
		block := content[i:end]
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, []byte(block))
		if err != nil {
			return "", err
		}
		// 转为16进制编码
		str := hex.EncodeToString(signature)
		signatures += str
	}

	return signatures, nil
}

// DecryptByCERFile 使用CER文件解密
func DecryptByCERFile(content string, publicKey *rsa.PublicKey, publicPEM []byte) (string, error) {

	gorsa.RSA.SetPublicKey(string(publicPEM))

	// 分块解密
	blockSize := publicKey.Size() * 8
	if blockSize <= 0 {
		return "", errors.New("密钥长度无效，无法确定块大小")
	}
	totalLen := len(content)
	decryptResult := ""
	encryptSubStarLen := 0

	for encryptSubStarLen < totalLen {
		// 确保不会超出数据范围
		endPos := encryptSubStarLen + blockSize
		if endPos > totalLen {
			endPos = totalLen
		}

		// 从十六进制转换为二进制
		hexData := content[encryptSubStarLen:endPos]

		binData, err := hex.DecodeString(hexData)
		if err != nil {
			return "", fmt.Errorf("十六进制解码失败: %w", err)
		}

		pubdecrypt, err := gorsa.RSA.PubKeyDECRYPT(binData)
		if err != nil {
			return "", err
		}

		decryptResult += string(pubdecrypt)
		encryptSubStarLen += blockSize
	}

	// Base64解码
	resultBytes, err := base64.StdEncoding.DecodeString(decryptResult)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}

	return string(resultBytes), nil
}

// Sign 使用私钥对数据进行签名
// 注意：这里是简化实现，实际情况中需要处理pfx证书解析等复杂逻辑
func Sign(data string, privateKey *rsa.PrivateKey) (string, error) {

	if privateKey == nil {
		return "", errors.New("private key is nil")
	}

	// 计算数据的SHA256哈希值
	hashed := sha256.Sum256([]byte(data))

	// 使用私钥对哈希值进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	hexSignature := hex.EncodeToString(signature)

	return hexSignature, nil
}

// VerifySign 验证签名
func VerifySign(data, signStr string, publicKey *rsa.PublicKey) (bool, error) {

	if publicKey == nil {
		return false, errors.New("public key is nil")
	}

	// 计算数据的SHA256哈希值
	hashed := sha256.Sum256([]byte(data))

	h, _ := hex.DecodeString(signStr)
	// 验证签名
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], h)
	return err == nil, nil

}

func DecodePrivateKey(pemContent []byte) (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("pem.Decode(%s)：pemContent decode error", pemContent)
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		pk8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("私钥解析出错 [%s]", pemContent)
		}
		var ok bool
		privateKey, ok = pk8.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("私钥解析出错 [%s]", pemContent)
		}
	}
	return privateKey, nil
}

func DecodePublicKey(pemContent []byte) (publicKey *rsa.PublicKey, err error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("pem.Decode(%s)：pemContent decode error", pemContent)
	}
	switch block.Type {
	case "CERTIFICATE":
		pubKeyCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParseCertificate(%s)：%w", pemContent, err)
		}
		pubKey, ok := pubKeyCert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("公钥证书提取公钥出错 [%s]", pemContent)
		}
		publicKey = pubKey
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParsePKIXPublicKey(%s),err:%w", pemContent, err)
		}
		pubKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("公钥解析出错 [%s]", pemContent)
		}
		publicKey = pubKey
	case "RSA PUBLIC KEY":
		pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParsePKCS1PublicKey(%s)：%w", pemContent, err)
		}
		publicKey = pubKey
	}
	return publicKey, nil
}
