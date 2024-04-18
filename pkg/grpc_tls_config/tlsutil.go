package grpc_tls_config

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func BuildInsecureClientTlsConfig() credentials.TransportCredentials {
	clientTLSConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	return credentials.NewTLS(clientTLSConfig)
}

func BuildSecureClientTlsConfig(ClientCert, ClientKey, CACert string) credentials.TransportCredentials {
	// 客户端加载自己的证书和私钥
	clientCert, err := tls.LoadX509KeyPair(ClientCert, ClientKey)
	_ = clientCert
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}
	// 加载信任的根 CA 证书
	caCert, err := ioutil.ReadFile(CACert)
	if err != nil {
		log.Fatalf("Failed to read root CA cert: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientTLSConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}
	return credentials.NewTLS(clientTLSConfig)
}

func BuildServerTlsConfig(ServerCert, ServerKey, CACert string) credentials.TransportCredentials {
	// 加载服务器的证书和私钥
	serverCert, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}
	// 加载根 CA 证书（用于验证客户端）
	caCert, err := ioutil.ReadFile(CACert)
	if err != nil {
		log.Fatalf("Failed to read root CA cert: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	DefaultTLSCipherSuites := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	}
	// 创建并配置 TLS 凭证
	serverTLSConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		CipherSuites: DefaultTLSCipherSuites,
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		//ClientAuth:   tls.NoClientCert,
		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  caCertPool,
	}

	serverTLSConfig.BuildNameToCertificate()
	return credentials.NewTLS(serverTLSConfig)
}
