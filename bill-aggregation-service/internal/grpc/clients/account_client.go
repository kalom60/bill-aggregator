package clients

import (
	account_protos "bill-aggregation-service/internal/grpc/pb/account-protos"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	conn   *grpc.ClientConn
	client account_protos.AccountServiceClient
}

func NewAccountClient(grpcAddr string) (*AccountClient, error) {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := account_protos.NewAccountServiceClient(conn)
	return &AccountClient{
		client: client,
	}, nil
}

type Account struct {
	Id                  string `json:"id"`
	UserId              string `json:"user_id"`
	ProviderId          string `json:"provider_id"`
	AccountIdentifier   string `json:"account_identifier"`
	EncryptedCredential string `json:"encrypted_credential"`
}

func (a *AccountClient) GetLinkedAccounts(ctx context.Context, userID string) ([]Account, error) {
	req := &account_protos.AccountRequest{UserId: userID}
	resp, err := a.client.GetLinkedAccounts(ctx, req)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	for _, acc := range resp.Accounts {
		accounts = append(accounts, Account{
			Id:                  acc.Id,
			UserId:              acc.UserId,
			ProviderId:          acc.ProviderId,
			AccountIdentifier:   acc.AccountIdentifier,
			EncryptedCredential: acc.EncryptedCredential,
		})
	}

	return accounts, nil
}

func (a *AccountClient) Close() {
	if a.conn != nil {
		a.conn.Close()
	}
}
