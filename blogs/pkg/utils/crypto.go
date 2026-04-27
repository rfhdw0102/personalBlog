package utils

import (
	"blogs/pkg/logger"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"go.uber.org/zap"
)

// DecryptPassword 前端RSA加密密码 → 后端解密
func DecryptPassword(reqPassword string, privateKeyPEM string) (string, error) {
	// 解析PEM私钥
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logger.Info("无效的RSA私钥")
		return "", errors.New("无效的RSA私钥")
	}
	// 解析PKCS1私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Info("私钥解析失败", zap.Error(err))
		return "", errors.New("私钥解析失败")
	}
	// Base64解码前端传过来的加密密码
	decodedPwd, err := base64.StdEncoding.DecodeString(reqPassword)
	if err != nil {
		logger.Info("密码base64解码失败", zap.Error(err))
		return "", errors.New("密码base64解码失败")
	}
	// RSA解密
	plainPwd, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedPwd)
	if err != nil {
		logger.Info("密码解密失败", zap.Error(err))
		return "", errors.New("密码解密失败")
	}
	return string(plainPwd), nil
}
