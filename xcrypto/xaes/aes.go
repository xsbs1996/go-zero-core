package xaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidKey = errors.New("xaes: invalid key size") // ErrInvalidKey 表示 AES 密钥长度不合法。
	ErrInvalidIV  = errors.New("xaes: invalid iv size")  // ErrInvalidIV 表示 CBC 初始向量长度不合法。
)

// EncryptGCM 使用 AES-GCM 加密，返回 base64 编码密文。
func EncryptGCM(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrInvalidKey
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptGCM 使用 AES-GCM 解密 base64 编码密文。
func DecryptGCM(cipherText string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKey
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("xaes: ciphertext too short")
	}

	nonce := data[:gcm.NonceSize()]
	data = data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, data, nil)
}

// EncryptCBC 使用 AES-CBC 和 PKCS7 填充加密，返回 base64 编码密文。
func EncryptCBC(plainText []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrInvalidKey
	}
	if len(iv) != block.BlockSize() {
		return "", ErrInvalidIV
	}

	plainText = pkcs7Padding(plainText, block.BlockSize())
	cipherText := make([]byte, len(plainText))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptCBC 使用 AES-CBC 和 PKCS7 填充解密 base64 编码密文。
func DecryptCBC(cipherText string, key []byte, iv []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKey
	}
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIV
	}
	if len(data) == 0 || len(data)%block.BlockSize() != 0 {
		return nil, errors.New("xaes: invalid ciphertext size")
	}

	plainText := make([]byte, len(data))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plainText, data)
	return pkcs7UnPadding(plainText, block.BlockSize())
}

// pkcs7Padding 对明文进行 PKCS7 填充。
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// pkcs7UnPadding 移除 PKCS7 填充。
func pkcs7UnPadding(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("xaes: invalid padding size")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize || padding > len(data) {
		return nil, errors.New("xaes: invalid padding")
	}
	for _, v := range data[len(data)-padding:] {
		if int(v) != padding {
			return nil, errors.New("xaes: invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
