package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	grpcv1 "git.ouryun.cn/something/grpc/proto/pb/base/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 客户端加载自己的证书和私钥
	clientCert, err := tls.LoadX509KeyPair("./client1.crt", "./client1.key")
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	// 加载信任的根 CA 证书
	caCert, err := ioutil.ReadFile("./root-ca.crt")
	if err != nil {
		log.Fatalf("Failed to read root CA cert: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 创建并配置 TLS 凭证
	clientTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool, // 信任服务端证书的根 CA
	}

	// 创建带有 TLS 凭证的 DialOption
	creds := credentials.NewTLS(clientTLSConfig)

	// 连接到 gRPC 服务器
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建客户端实例并调用 gRPC 方法...
	// ...
	c := grpcv1.NewGreeterClient(conn)
	name := "world"
	resp, err := c.SayHello(context.Background(), &grpcv1.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Failed to greet: %v", err)
	}
	log.Println("greeting: %s", resp.Message)
}
