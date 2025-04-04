package grpc

import (
	"context"

	account_protos "github.com/kalom60/bill-aggregator/account-linking-service/internal/grpc/pb/account-protos"
)

func (p *AccountServer) GetLinkedAccounts(ctx context.Context, req *account_protos.AccountRequest) (*account_protos.AccountResponse, error) {
	id := req.GetUserId()

	linkedAccounts, err := p.client.FetchLinkedAccountsByUserID(id)
	if err != nil {
		return nil, err
	}

	var pbAccounts []*account_protos.LinkedAccount
	for _, acc := range linkedAccounts {
		pbAccounts = append(pbAccounts, &account_protos.LinkedAccount{
			Id:                  acc.ID,
			UserId:              acc.UserID,
			ProviderId:          acc.ProviderID,
			AccountIdentifier:   acc.AccountIdentifier,
			EncryptedCredential: acc.EncryptedCredential,
		})
	}

	return &account_protos.AccountResponse{
		Accounts: pbAccounts,
	}, nil
}
