package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func main() {

	plainText := []byte("DANGEROUS_DATA")
	key := []byte("KEYS____KEYS____KEYS____KEYS____")
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)

	fmt.Print(cipherText)

	text, err := decryptAES(key, cipherText)
	fmt.Print(string(text))
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
