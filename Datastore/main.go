package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "andreabarchietto.it/oss/go/telemetry_datastore/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMetricServiceServer
}

func (s *server) StreamMetrics(stream pb.MetricService_StreamMetricsServer) error {
	var count int
	for {
		payload, err := stream.Recv()
		if err == io.EOF {
			// Finished receiving stream, return response
			return stream.SendAndClose(&pb.StreamMetricsResponse{
				Success: true,
				Message: fmt.Sprintf("Successfully processed %d metric payloads", count),
			})
		}
		if err != nil {
			return err
		}

		count++
		log.Printf("CPU: %.2f%% | RAM: %d/%d bytes",
			payload.GetCpu().GetUsagePercentage(),
			payload.GetRam().GetUsedBytes(),
			payload.GetRam().GetTotalBytes(),
		)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMetricServiceServer(grpcServer, &server{})

	log.Println("gRPC Server running on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
