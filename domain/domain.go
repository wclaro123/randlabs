package domain

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type (
	Payment struct {
		Receiver        string
		Sender          string
		Amount          uint64
		CloseAccount    bool
		CloseReminderTo string
		Fee             uint64
	}
	Genesis struct {
		Alloc []Alloc `json:"alloc"`
		Fees  string  `json:"fees"`
	}
	Alloc struct {
		Addr  string
		State State `json:"state"`
	}
	State struct {
		Algo int64 `json:"algo"`
	}
	Wallet struct {
		Account string
		Amount  int64
	}
	Transaction  models.Transaction
	Transactions []Transaction
)

func (tx Transaction) GetPayment() Payment {
	return Payment{
		Receiver:        tx.PaymentTransaction.Receiver,
		Sender:          tx.Sender,
		Amount:          tx.PaymentTransaction.Amount,
		CloseAccount:    tx.PaymentTransaction.CloseRemainderTo != "",
		CloseReminderTo: tx.PaymentTransaction.CloseRemainderTo,
		Fee:             tx.Fee,
	}
}

func (txs Transactions) GetPayments() (payments []Payment) {
	for _, tx := range txs {
		payments = append(payments, tx.GetPayment())
	}

	return payments
}
