package databases

import (
	"time"
)

type NotFoundError string

func (e NotFoundError) Error() string {
	return "Error: Not Found"
}

type Operation string

const (
	Add Operation = "ADD"
	Sub Operation = "SUB"
)

type Transaction struct {
	ID          string
	Description string
	Amount      float32
	Operation   Operation
	Date        int64
}

var transactionsDB []Transaction

func GetTransactions() []Transaction {
	return transactionsDB
}

func AddTransaction(transaction Transaction) {
	transactionsDB = append(transactionsDB, transaction)
}

func DeleteTransactionById(ID string) error {
	index := -1

	for i, transaction := range transactionsDB {
		if transaction.ID == ID {
			index = i
			break
		}
	}

	if index == -1 {
		return NotFoundError("Transaction")
	}

	transactionsDB[index] = transactionsDB[len(transactionsDB)-1]
	transactionsDB = transactionsDB[:len(transactionsDB)-1]

	return nil
}

func GetSumAmount(transactions []Transaction) float32 {
	var total float32 = 0.0

	for _, transaction := range transactions {
		switch transaction.Operation {
		case Add:
			total = total + transaction.Amount
		case Sub:
			total = total - transaction.Amount
		}
	}

	return total
}

func GetTransactionByID(id string) (error, *Transaction) {
	var transaction Transaction

	for _, t := range transactionsDB {
		if t.ID == id {
			transaction = t
			break
		}
	}

	if &transaction != nil {
		return nil, &transaction
	}

	return NotFoundError("Transaction"), nil
}

func (transaction *Transaction) FormatDate() string {
	t := time.Unix(transaction.Date, 0)
	return t.Format("2006/01/02")
}
