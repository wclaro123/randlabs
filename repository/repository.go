package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/wclaro123/randlabs/domain"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) Repository {
	return &repository{db: db}
}

func (r repository) SavePayment(ctx context.Context, payment domain.Payment) error {
	query := "INSERT INTO balance.payment (receiver, sender, amount, fee) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(ctx, query,
		payment.Receiver,
		payment.Sender,
		payment.Amount,
		payment.Fee,
	)

	if err != nil {
		return fmt.Errorf("failed to insert payment with error: %w", err)
	}

	return nil
}

func (r repository) UpsertWallet(ctx context.Context, wallet domain.Wallet) error {
	queryReceiver := `INSERT INTO balance.wallet (account, amount) VALUES ($1, $2)
				ON CONFLICT (account) DO
				UPDATE SET amount = wallet.amount + EXCLUDED.amount;`
	_, err := r.db.Exec(ctx, queryReceiver,
		wallet.Account,
		wallet.Amount,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert wallet with error: %w", err)
	}

	return nil
}

func (r repository) InitDb(ctx context.Context) error {
	query := `drop schema if exists balance cascade;
			
			create schema balance;
			
			create table balance.payment
			(
				id       serial
					constraint payment_pk
						primary key,
				receiver varchar not null,
				sender   varchar not null,
				amount   bigint,
				fee      int
			);
			
			create unique index payment_id_uindex
				on balance.payment (id);
			
			create table balance.wallet
			(
				id      serial
					constraint wallet_pk
						primary key,
				account varchar not null,
				amount  bigint
			);
			
			create unique index wallet_account_uindex
				on balance.wallet (account);
			
			create unique index wallet_id_uindex
				on balance.wallet (id);
			`

	_, err := r.db.Exec(ctx, query)
	return err
}
