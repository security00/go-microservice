package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc_simple/proto"
	"log"
	"os"
	"time"
)

const (
	Address     = ":50051"
	DefaultName = "kevin"
)

func main() {
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	log.Printf("conn is %v\n", conn)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)
	name := DefaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{
		Name: name,
	})
	if err != nil {
		log.Fatalf("could not hello: %v", err)
	}
	log.Printf("Received Server message: %v", r.GetMessage())
}
