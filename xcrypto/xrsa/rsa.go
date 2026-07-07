package xrsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

var (
	ErrInvalidPrivateKey = errors.New("xrsa: invalid private key") // ErrInvalidPrivateKey 表示私钥不合法。
	ErrInvalidPublicKey  = errors.New("xrsa: invalid public key")  // ErrInvalidPublicKey 表示公钥不合法。
)

// GenerateKey 生成 RSA 私钥。
//
// 参数：
//   - bits: RSA 密钥位数，例如 2048 或 4096。
//
// 返回值：
//   - *rsa.PrivateKey: 生成的 RSA 私钥。
//   - error: 随机源读取失败或 bits 非法时返回错误。
func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// PrivateKeyToPEM 将 RSA 私钥转换为 PEM 字符串。
//
// 参数：
//   - key: RSA 私钥。
//
// 返回值：
//   - string: PKCS#1 PEM 格式私钥字符串。
func PrivateKeyToPEM(key *rsa.PrivateKey) string {
	data := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: data}
	return string(pem.EncodeToMemory(block))
}

// PublicKeyToPEM 将 RSA 公钥转换为 PEM 字符串。
//
// 参数：
//   - key: RSA 公钥。
//
// 返回值：
//   - string: PKIX PEM 格式公钥字符串。
//   - error: 公钥序列化失败时返回错误。
func PublicKeyToPEM(key *rsa.PublicKey) (string, error) {
	data, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: data}
	return string(pem.EncodeToMemory(block)), nil
}

// ParsePrivateKey 解析 PEM 格式 RSA 私钥。
//
// 参数：
//   - privateKey: PEM 格式私钥字符串，支持 PKCS#1 和 PKCS#8。
//
// 返回值：
//   - *rsa.PrivateKey: 解析出的 RSA 私钥。
//   - error: PEM 格式非法或不是 RSA 私钥时返回错误。
func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, ErrInvalidPrivateKey
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, ErrInvalidPrivateKey
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidPrivateKey
	}
	return key, nil
}

// ParsePublicKey 解析 PEM 格式 RSA 公钥。
//
// 参数：
//   - publicKey: PEM 格式公钥字符串。
//
// 返回值：
//   - *rsa.PublicKey: 解析出的 RSA 公钥。
//   - error: PEM 格式非法或不是 RSA 公钥时返回错误。
func ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, ErrInvalidPublicKey
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, ErrInvalidPublicKey
	}
	key, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidPublicKey
	}
	return key, nil
}

// EncryptOAEP 使用 RSA-OAEP 加密，返回 base64 编码密文。
//
// 参数：
//   - plainText: 待加密明文。
//   - publicKey: RSA 公钥。
//
// 返回值：
//   - string: base64 编码密文。
//   - error: 加密失败时返回错误。
func EncryptOAEP(plainText []byte, publicKey *rsa.PublicKey) (string, error) {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plainText, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptOAEP 使用 RSA-OAEP 解密 base64 编码密文。
//
// 参数：
//   - cipherText: EncryptOAEP 生成的 base64 编码密文。
//   - privateKey: RSA 私钥。
//
// 返回值：
//   - []byte: 解密后的明文。
//   - error: base64 解码或解密失败时返回错误。
func DecryptOAEP(cipherText string, privateKey *rsa.PrivateKey) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
}

// SignPSS 使用 RSA-PSS 签名，返回 base64 编码签名。
//
// 参数：
//   - data: 待签名数据。
//   - privateKey: RSA 私钥。
//
// 返回值：
//   - string: base64 编码签名。
//   - error: 签名失败时返回错误。
func SignPSS(data []byte, privateKey *rsa.PrivateKey) (string, error) {
	sum := sha256.Sum256(data)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, sum[:], nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyPSS 使用 RSA-PSS 验签。
//
// 参数：
//   - data: 原始数据。
//   - signature: SignPSS 生成的 base64 编码签名。
//   - publicKey: RSA 公钥。
//
// 返回值：
//   - error: 验签成功返回 nil；签名格式非法或签名不匹配时返回错误。
func VerifyPSS(data []byte, signature string, publicKey *rsa.PublicKey) error {
	signData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	return rsa.VerifyPSS(publicKey, crypto.SHA256, sum[:], signData, nil)
}
