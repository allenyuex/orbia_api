package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/bcrypt"
)

// ValidateWalletAddress 验证钱包地址格式
func ValidateWalletAddress(address string) bool {
	// 以太坊地址格式：0x开头，42个字符，包含0-9和a-f
	matched, _ := regexp.MatchString("^0x[a-fA-F0-9]{40}$", address)
	return matched
}

// VerifySignature 验证钱包签名
func VerifySignature(walletAddress, message, signature string) error {
	// 验证钱包地址格式
	if !ValidateWalletAddress(walletAddress) {
		return errors.New("invalid wallet address format")
	}

	// 移除0x前缀
	if strings.HasPrefix(signature, "0x") {
		signature = signature[2:]
	}

	// 解码签名
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signature format: %v", err)
	}

	if len(sigBytes) != 65 {
		return errors.New("signature must be 65 bytes")
	}

	// 调整recovery ID
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// 构造消息哈希（以太坊签名消息格式）
	messageHash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// 恢复公钥
	pubKey, err := crypto.SigToPub(messageHash.Bytes(), sigBytes)
	if err != nil {
		return fmt.Errorf("failed to recover public key: %v", err)
	}

	// 从公钥生成地址
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	// 比较地址（不区分大小写）
	if !strings.EqualFold(recoveredAddress, walletAddress) {
		return errors.New("signature verification failed")
	}

	return nil
}

// GenerateSignMessage 生成用于签名的消息
func GenerateSignMessage(walletAddress string, timestamp int64) string {
	return fmt.Sprintf("Welcome to Orbia!\n\nWallet: %s\nTimestamp: %d\n\nThis request will not trigger a blockchain transaction or cost any gas fees.", walletAddress, timestamp)
}

// HashPassword 对密码进行哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码哈希
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateKeyPair 生成ECDSA密钥对（用于测试）
func GenerateKeyPair() (*ecdsa.PrivateKey, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, "", err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return privateKey, address, nil
}