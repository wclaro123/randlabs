package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/algorand/go-algorand-sdk/client/v2/indexer"

	"github.com/wclaro123/randlabs/domain"
	"github.com/wclaro123/randlabs/repository"
)

const paymentTxType = "pay"

type algorandService struct {
	client   *indexer.Client
	repo     repository.Repository
	payments chan []domain.Payment
	errors   chan error
}

func NewAlgorandService(client *indexer.Client, repo repository.Repository, errors chan error) AlgorandService {
	payments := make(chan []domain.Payment)
	return &algorandService{client: client, repo: repo, payments: payments, errors: errors}
}

func (a algorandService) InitGenesis(ctx context.Context) error {
	err := a.repo.InitDb(ctx)
	if err != nil {
		return fmt.Errorf("failed db init with error: %w", err)
	}

	r, err := http.Get("https://node.testnet.algoexplorerapi.io/genesis")
	if err != nil {
		return fmt.Errorf("failed getting genesis file info with error: %w", err)
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed reading body from genesis file with error: %w", err)
	}

	var content domain.Genesis
	_ = json.Unmarshal(body, &content)

	for _, alloc := range content.Alloc {
		err = a.repo.UpsertWallet(ctx, domain.Wallet{
			Account: alloc.Addr,
			Amount:  alloc.State.Algo,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a algorandService) GetCurrentRound(ctx context.Context) (uint64, error) {
	response, err := a.client.HealthCheck().Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting current round: %w", err)
	}

	return response.Round, nil
}

func (a algorandService) ProcessTransactions(ctx context.Context, start, current, increment uint64) {
	for i := start; i < current; i += increment {
		log.Printf("processing transaction round %d", i/increment)
		txs, err := a.getTransactions(ctx, i, i+increment)
		if err != nil {
			a.errors <- err
		}

		a.payments <- txs.GetPayments()
	}
}

func (a algorandService) ProcessPayments(ctx context.Context) {
	round := 0
	for payment := range a.payments {
		log.Printf("processing payment round %d", round)
		err := a.savePayments(ctx, payment)
		if err != nil {
			a.errors <- err
		}
		round++
	}
}

func (a algorandService) getTransactions(ctx context.Context, min, max uint64) (domain.Transactions, error) {
	response, err := a.client.
		SearchForTransactions().
		TxType(paymentTxType).
		MinRound(min).
		MaxRound(max).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error searching for transactions: %w", err)
	}

	txs := make([]domain.Transaction, len(response.Transactions))
	for i, t := range response.Transactions {
		txs[i] = domain.Transaction(t)
	}

	return txs, nil
}

func (a algorandService) savePayments(ctx context.Context, payments []domain.Payment) error {
	for _, payment := range payments {
		err := a.repo.SavePayment(ctx, payment)
		if err != nil {
			return err
		}

		err = a.updateWallet(ctx, payment)
	}

	return nil
}

func (a algorandService) updateWallet(ctx context.Context, payment domain.Payment) error {
	err := a.repo.UpsertWallet(ctx, domain.Wallet{
		Account: payment.Receiver,
		Amount:  int64(payment.Amount),
	})
	if err != nil {
		return fmt.Errorf("failed upserting receiver wallet")
	}

	amount := int64(payment.Amount+payment.Fee) * -1
	err = a.repo.UpsertWallet(ctx, domain.Wallet{
		Account: payment.Sender,
		Amount:  amount,
	})
	if err != nil {
		return fmt.Errorf("failed upserting sender wallet")
	}

	return nil
}
