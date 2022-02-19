package main

import (
	"calculator/sumpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func main() {
	fmt.Println("== Set up Calculator Server")

	list, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalln("Error settimg up server")

	}
	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})
	if err := s.Serve(list); err != nil {
		log.Fatalln("Failed to serve sumService")
	}
}

func (*server) Sum(ctx context.Context, req *sumpb.Add) (*sumpb.SumResponse, error) {
	fmt.Println("Sum service was invoked with %v", req)
	num1 := req.GetSumRequest().GetNum_1()
	num2 := req.GetSumRequest().GetNum_2()
	num3 := num1 + num2

	res := &sumpb.SumResponse{
		NumSum: num3,
	}

	return res, nil
}
