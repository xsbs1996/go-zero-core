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
func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// PrivateKeyToPEM 将 RSA 私钥转换为 PEM 字符串。
func PrivateKeyToPEM(key *rsa.PrivateKey) string {
	data := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: data}
	return string(pem.EncodeToMemory(block))
}

// PublicKeyToPEM 将 RSA 公钥转换为 PEM 字符串。
func PublicKeyToPEM(key *rsa.PublicKey) (string, error) {
	data, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: data}
	return string(pem.EncodeToMemory(block)), nil
}

// ParsePrivateKey 解析 PEM 格式 RSA 私钥。
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
func EncryptOAEP(plainText []byte, publicKey *rsa.PublicKey) (string, error) {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plainText, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptOAEP 使用 RSA-OAEP 解密 base64 编码密文。
func DecryptOAEP(cipherText string, privateKey *rsa.PrivateKey) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
}

// SignPSS 使用 RSA-PSS 签名，返回 base64 编码签名。
func SignPSS(data []byte, privateKey *rsa.PrivateKey) (string, error) {
	sum := sha256.Sum256(data)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, sum[:], nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyPSS 使用 RSA-PSS 验签。
func VerifyPSS(data []byte, signature string, publicKey *rsa.PublicKey) error {
	signData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	return rsa.VerifyPSS(publicKey, crypto.SHA256, sum[:], signData, nil)
}
