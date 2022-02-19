package main

import (
	"calculator/sumpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("== Init Client Calculator ==")

	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Error dialling service")

	}

	defer cc.Close()
	c := sumpb.NewSumServiceClient(cc)

	req := &sumpb.Add{
		SumRequest: &sumpb.SumRequest{
			Num_1: 68,
			Num_2: 33,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalln("Error calling sum service")
	}
	fmt.Println("Sum of number :: %v", res.NumSum)

}
