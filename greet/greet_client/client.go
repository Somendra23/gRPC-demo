package main

import (
	"context"
	"fmt"
	"greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello i am client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Client conneciton is failed..")
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDirStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Hello i am going to call RPC Unary")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Simple",
			LastName:  "Sa",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalln("Errror making a remote call to GreetService")
	}
	fmt.Printf("Successfully received result::  %v ", res.Result)
	//fmt.Printf("Created Client: %v", c)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Server Streaming in Progress==")
	req := &greetpb.GreetManyTimeRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Hola",
			LastName:  "Sam",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while doing server streaming..")
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error receiving stream data")
		}
		fmt.Println("Response from server:: %v ", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Start doing client side streaming..")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jk",
				LastName:  "Jim",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "LK",
				LastName:  "Lim",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ja",
				LastName:  "Jai",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ri",
				LastName:  "Rin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Fi",
				LastName:  "Fin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vi",
				LastName:  "Vin",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("Error streaming cleint")
	}
	for _, req := range requests {
		fmt.Println("Sending Request %v", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalln("Error receiving data from server")
	}

	fmt.Println("LomgResponse received %v ", res)

}

func doBiDirStreaming(c greetpb.GreetServiceClient) {

	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalln("Error calling GreetEveryone")
		return
	}

	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ri",
				LastName:  "Rin",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ji",
				LastName:  "Jin",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ui",
				LastName:  "Uin",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ti",
				LastName:  "Tin",
			},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
