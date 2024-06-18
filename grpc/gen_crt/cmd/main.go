package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// 创建根CA证书
	rootCert, rootKey := createRootCACert()

	// 使用根CA为服务器签发证书
	serverCert, serverKey := createSignedCert(rootCert, rootKey, "myServe")

	// 使用根CA为第一个客户端签发证书
	client1Cert, client1Key := createSignedCert(rootCert, rootKey, "myClient1")

	// 使用根CA为第二个客户端签发证书
	client2Cert, client2Key := createSignedCert(rootCert, rootKey, "myClient2")

	// 将生成的证书和私钥写入文件（请确保替换路径为你希望保存的位置）
	writeCertificate(rootCert, "root-ca.crt")
	writePrivateKey(rootKey, "root-ca.key")
	writeCertificate(serverCert, "server.crt")
	writePrivateKey(serverKey, "server.key")
	writeCertificate(client1Cert, "client1.crt")
	writePrivateKey(client1Key, "client1.key")
	writeCertificate(client2Cert, "client2.crt")
	writePrivateKey(client2Key, "client2.key")
}

// 创建根CA证书
func createRootCACert() (*x509.Certificate, *rsa.PrivateKey) {
	caKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "myServer",
			Organization: []string{"Ouryun"},
		},
		NotBefore:             time.Now().AddDate(-10, 0, 0),
		NotAfter:              time.Now().AddDate(3000, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		panic(err)
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		panic(err)
	}

	return caCert, caKey
}

// 使用CA证书和密钥创建并返回已签名的证书与私钥
func createSignedCert(caCert *x509.Certificate, caKey *rsa.PrivateKey, hostname string) (*x509.Certificate, *rsa.PrivateKey) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:   hostname,
			Organization: []string{"Ouryun"},
		},
		DNSNames:    []string{hostname, "localhost"},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now().AddDate(-10, 0, 0),
		NotAfter:    time.Now().AddDate(3000, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &key.PublicKey, caKey)
	if err != nil {
		panic(err)
	}
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		panic(err)
	}

	return cert, key
}

// 将证书或私钥写入到指定文件
func writeCertificate(cert *x509.Certificate, path string) {
	certOut, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	fmt.Println("Wrote", path)
}

func writePrivateKey(key *rsa.PrivateKey, path string) {
	keyOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	fmt.Println("Wrote", path)
}
