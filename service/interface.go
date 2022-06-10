package service

import (
	"context"
)

type AlgorandService interface {
	InitGenesis(ctx context.Context) error
	GetCurrentRound(ctx context.Context) (uint64, error)
	ProcessTransactions(ctx context.Context, start, current, increment uint64)
	ProcessPayments(ctx context.Context)
}
