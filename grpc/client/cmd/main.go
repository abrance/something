package main

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"log"

	grpcv1 "git.ouryun.cn/something/grpc/proto/pb/base/grpc"
	"google.golang.org/grpc"
)

func main() {
	//// 客户端加载自己的证书和私钥
	//clientCert, err := tls.LoadX509KeyPair("./client1.crt", "./client1.key")
	//_ = clientCert
	//if err != nil {
	//	log.Fatalf("Failed to load client certificate: %v", err)
	//}
	//
	//// 加载信任的根 CA 证书
	//caCert, err := ioutil.ReadFile("./root-ca.crt")
	//if err != nil {
	//	log.Fatalf("Failed to read root CA cert: %v", err)
	//}
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)

	clientTLSConfig := &tls.Config{
		//Certificates:       []tls.Certificate{clientCert},
		//RootCAs:            caCertPool, // 信任服务端证书的根 CA
		InsecureSkipVerify: true,
	}
	creds := credentials.NewTLS(clientTLSConfig)
	//_ = creds

	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
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
