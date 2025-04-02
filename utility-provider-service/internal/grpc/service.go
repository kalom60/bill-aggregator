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
