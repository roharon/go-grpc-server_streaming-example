package main

import (
	"context"
	"flag"
	"io"
	"log"

	pb "aaronroh.com/m/proto/v1"

	"google.golang.org/grpc"
)

var serverAddr = flag.String("server_addr", "localhost:10000", "The server address with port")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewRouteClient(conn)

	content, err := client.GetInfo(context.Background(), &pb.Content{Message: "Hi GetInfo Unary RPC"})
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%s", content)

	stream, err := client.ListInfo(context.Background(), &pb.Content{Message: "Hi ListInfo Server Stream RPC"})

	if err != nil {
		log.Fatalf("ListInfo - %v", err)
	}

	for {
		content, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListInfo stream - %v", err)
		}

		log.Printf("Content: Message: %s", content.GetMessage())
	}
}
