package grpc

import (
	"context"
	"time"

	provider_protos "github.com/kalom60/bill-aggregator/account-linking-service/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProviderClient struct {
	conn   *grpc.ClientConn
	client provider_protos.ProviderServiceClient
}

func NewProviderClient(grpcAddr string) (*ProviderClient, error) {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := provider_protos.NewProviderServiceClient(conn)
	return &ProviderClient{
		client: client,
	}, nil
}

func (c *ProviderClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *ProviderClient) IsProviderExist(providerID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &provider_protos.ProviderRequest{ProviderId: providerID}
	resp, err := c.client.IsProviderExist(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.Exist, nil
}
