package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	pb "aaronroh.com/m/proto/v1"
	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 10000, "The server port")
	jsonFile = flag.String("json_file", "", "Json file containing list of content")
)

type RouteServer struct {
	pb.UnimplementedRouteServer
	savedContents []*pb.Content
}

func (s *RouteServer) GetInfo(ctx context.Context, req *pb.Content) (*pb.Content, error) {
	log.Printf("GetInfo - %v", req)
	return &pb.Content{Message: "Hi!"}, nil
}

func (s *RouteServer) ListInfo(req *pb.Content, stream pb.Route_ListInfoServer) error {
	log.Printf("ListInfo - %v", req)

	for _, content := range s.savedContents {
		if err := stream.Send(content); err != nil {
			return err
		}
	}
	return nil
}

func (s *RouteServer) loadContents(filePath string) {
	if filePath == "" {
		log.Fatalf("Must set jsonFile option")
	}

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatalf("Failed to load Contents: %v", err)
	}

	if err := json.Unmarshal(data, &s.savedContents); err != nil {
		log.Fatalf("Failed to load: %v", err)
	}
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := &RouteServer{}
	s.loadContents(*jsonFile)

	grpcServer := grpc.NewServer()
	pb.RegisterRouteServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
