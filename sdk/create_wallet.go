package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

// Identity represents a blockchain user's identity
type Identity struct {
	Cert string
	Key  string
	MSP  string
}

func main() {
	// 文件路径
	certPath := "./sdk/wallet/appUser/msp/signcerts/cert.pem"
	keyPath := "./sdk/wallet/appUser/msp/keystore/private.pem"

	// 加载证书
	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatalf("Failed to read certificate: %v", err)
	}

	// 加载私钥
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	// 验证证书格式
	if err := validateCertificate(cert); err != nil {
		log.Fatalf("Invalid certificate: %v", err)
	}

	// 验证私钥格式
	if err := validatePrivateKey(key); err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	// 创建用户身份
	identity := Identity{
		Cert: string(cert),
		Key:  string(key),
		MSP:  "Org1MSP",
	}

	// 将身份输出到日志或存储
	log.Printf("User identity created successfully: %+v\n", identity)

	// 如果需要存储，可以将 `identity` 写入文件或数据库
	fmt.Println("Identity successfully created!")
}

// validateCertificate validates the format of a PEM-encoded certificate
func validateCertificate(cert []byte) error {
	block, _ := pem.Decode(cert)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("not a valid PEM-encoded certificate")
	}
	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}
	return nil
}

// validatePrivateKey validates the format of a PEM-encoded private key
func validatePrivateKey(key []byte) error {
	block, _ := pem.Decode(key)
	if block == nil || (block.Type != "PRIVATE KEY" && block.Type != "EC PRIVATE KEY") {
		return fmt.Errorf("not a valid PEM-encoded private key")
	}
	return nil
}
