package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/database"
	provider_protos "github.com/kalom60/bill-aggregator/utility-provider-service/internal/grpc/pb"
	"google.golang.org/grpc"
)

type ProviderServer struct {
	provider_protos.UnimplementedProviderServiceServer
	client database.Service
}

func StartGRPCListen(client database.Service) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", "50001"))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()

	provider_protos.RegisterProviderServiceServer(grpcServer, &ProviderServer{
		client: client,
	})

	log.Printf("gRPC server started on port: %s", "50001")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}
}
