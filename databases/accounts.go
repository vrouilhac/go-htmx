package databases

type Account struct {
	ID   string
	Name string
}

var accountsDB []Account

func GetAccounts() []Account {
	return accountsDB
}

func AddAccount(account Account) {
	accountsDB = append(accountsDB, account)
}

func GetAccountById(id string) (error, *Account) {
	var account Account

	for _, a := range accountsDB {
		if a.ID == id {
			account = a
			break
		}
	}

	if &account != nil {
		return nil, &account
	}

	return NotFoundError("Account"), nil
}

func UpdateAccountById(id string, name string) error {
	for index, a := range accountsDB {
		if a.ID == id {
			accountsDB[index].Name = name
			return nil
		}
	}

	return NotFoundError("Account")
}
