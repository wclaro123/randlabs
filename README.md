# Specification: process payments
## Summary
Sequentially scan blocks in Algorand’s testnet, iterating payment transactions. Use this data to
keep a “balance sheet” for every address that participated in a payment.
Requirements
Implement a program in Go that:
1. Iterates blocks sequentially in the testnet network.
2. Recognizes payment transactions within each block.
3. Maintains a table that acts like a “balance sheet”.
   a. Should use PostgreSQL.
   b. For every address that participated in a transaction, its current balance should be
   stored.
   c. The initial state of the database can be populated with the information from the
   genesis file

# Solution

### Prerequisites

#### Setup a postgres container running
```
docker compose up
```

#### Start program running

```
go run main.go
```

### Possible improvements

- Add unit tests
- Add environment variables for constants  
- Handle transactions on the repository to ensure only full rounds block are processed and stored.
- Add parameters to handle the range of rounds to run. 
