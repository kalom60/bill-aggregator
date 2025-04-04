package clients

import (
	provider_protos "bill-aggregation-service/internal/grpc/pb/provider-protos"
	"context"

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

type Provider struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	ApiURL             string `json:"api_url"`
	AuthenticationType string `json:"authentication_type"`
	ApiKey             string `json:"api_key"`
}

func (p *ProviderClient) GetProvider(ctx context.Context, provider_id string) (*Provider, error) {
	req := &provider_protos.GetProviderRequest{ProviderId: provider_id}
	resp, err := p.client.GetProvider(ctx, req)
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		Id:                 resp.Id,
		Name:               resp.Name,
		ApiURL:             resp.ApiUrl,
		AuthenticationType: resp.AuthenticationType,
		ApiKey:             resp.ApiKey,
	}

	return provider, nil
}

func (p *ProviderClient) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
