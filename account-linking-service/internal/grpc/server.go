package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/kalom60/bill-aggregator/account-linking-service/internal/database"
	account_protos "github.com/kalom60/bill-aggregator/account-linking-service/internal/grpc/pb/account-protos"
	"google.golang.org/grpc"
)

type AccountServer struct {
	account_protos.UnimplementedAccountServiceServer
	client database.Service
}

func StartGRPCListen(client database.Service) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", "50001"))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()

	account_protos.RegisterAccountServiceServer(grpcServer, &AccountServer{
		client: client,
	})

	log.Printf("gRPC server started on port: %s", "50001")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}
}
