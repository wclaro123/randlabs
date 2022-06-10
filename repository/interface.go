package repository

import (
	"context"

	"github.com/wclaro123/randlabs/domain"
)

type Repository interface {
	InitDb(ctx context.Context) error
	SavePayment(ctx context.Context, payment domain.Payment) error
	UpsertWallet(ctx context.Context, wallet domain.Wallet) error
}
