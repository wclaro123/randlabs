package main

import (
	"context"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/jackc/pgx/v4"

	"github.com/wclaro123/randlabs/repository"
	"github.com/wclaro123/randlabs/service"
)

func main() {
	fmt.Println("Connecting to Algorand Test Network")

	const (
		indexerAddress = "https://algoindexer.testnet.algoexplorerapi.io"
		indexerToken   = ""
		roundIncrement = 1000
		connString     = "postgres://postgres:password@localhost:5432/balance"
	)

	conn, err := pgx.Connect(context.Background(), connString)
	handleError(err)

	indexerClient, err := indexer.MakeClient(indexerAddress, indexerToken)
	handleError(err)

	r := repository.NewRepository(conn)
	errors := make(chan error)
	s := service.NewAlgorandService(indexerClient, r, errors)
	ctx := context.Background()

	handleError(s.InitGenesis(ctx))

	currentRound, err := s.GetCurrentRound(ctx)
	handleError(err)

	go s.ProcessTransactions(ctx, 0, currentRound, roundIncrement)

	go s.ProcessPayments(ctx)

	handleError(<-errors)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
