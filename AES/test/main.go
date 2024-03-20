package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"os"
)

// 生成一个32字节的AES-256密钥
func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// AES加密
func encryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个AES-GCM实例
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 创建一个随机Nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密并附加Nonce
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// AES解密
func decryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个AES-GCM实例
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// 分离Nonce和实际的密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	text := []byte("DANGEROUS")
	key, err := generateKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cipherText, err := encryptAES(key, text)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(cipherText)
	_text, err := decryptAES(key, cipherText)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(_text))

	fmt.Println(math.MaxInt8 >> 1)
}
