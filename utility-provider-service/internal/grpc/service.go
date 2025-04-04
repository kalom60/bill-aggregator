package grpc

import (
	"context"
	"errors"

	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/database"
	provider_protos "github.com/kalom60/bill-aggregator/utility-provider-service/internal/grpc/pb"
)

func (p *ProviderServer) IsProviderExist(ctx context.Context, req *provider_protos.ProviderRequest) (*provider_protos.ProviderResponse, error) {
	id := req.GetProviderId()

	_, err := p.client.FetchProviderByID(id)
	switch {
	case errors.Is(err, database.ErrProviderNotFound):
		return &provider_protos.ProviderResponse{
			Exist: false,
		}, nil
	case err != nil:
		return nil, err
	default:
		return &provider_protos.ProviderResponse{
			Exist: true,
		}, nil
	}
}

func (p *ProviderServer) GetProvider(ctx context.Context, req *provider_protos.GetProviderRequest) (*provider_protos.GetProviderResponse, error) {
	id := req.GetProviderId()

	provider, err := p.client.FetchProviderByID(id)
	if err != nil {
		return nil, err
	}

	return &provider_protos.GetProviderResponse{
		Id:                 provider.ID,
		Name:               provider.Name,
		ApiUrl:             provider.API_URL,
		AuthenticationType: provider.Authentication_Type,
		ApiKey:             provider.API_key,
	}, nil
}
