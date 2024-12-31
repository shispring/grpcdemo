package test

import (
	"context"
	"fmt"
	pb "grpcdemo/helloworld"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func TestGrpcClient(t *testing.T) {
	// TLS连接
	creds, err := credentials.NewClientTLSFromFile("/Users/shijianpeng/work/private/grpcdemo/keys/ca.crt", "www.demo.cn")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	// 连接服务器
	conn, err := grpc.NewClient(":8972", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	// 调用服务端的SayHello
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "CN"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}

	fmt.Printf("Greeting: %s !\n", r.Message)
}
