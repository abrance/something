package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpcv1 "git.ouryun.cn/something/grpc/proto/pb/base/grpc"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *grpcv1.HelloRequest) (*grpcv1.HelloReply, error) {
	return &grpcv1.HelloReply{Message: fmt.Sprintf("Hello, %s!", in.Name)}, nil
}

func main() {
	// 加载服务器的证书和私钥
	serverCert, err := tls.LoadX509KeyPair("./server.crt", "./server.key")
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}

	// 加载根 CA 证书（用于验证客户端）
	caCert, err := ioutil.ReadFile("./root-ca.crt")
	if err != nil {
		log.Fatalf("Failed to read root CA cert: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 创建并配置 TLS 凭证
	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		//ClientAuth:   tls.NoClientCert,
		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  caCertPool,
	}

	// 增强 TLS 配置的安全性
	serverTLSConfig.BuildNameToCertificate()

	// 使用 TLS 凭证创建 gRPC 服务器选项
	creds := credentials.NewTLS(serverTLSConfig)

	// 创建 gRPC 服务器监听地址，并应用凭证选项
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))

	// 注册你的 gRPC 服务实现到服务器上...
	// ...
	grpcv1.RegisterGreeterServer(s, &server{})

	// 开始服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
