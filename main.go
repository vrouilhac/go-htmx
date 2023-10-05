package main

import (
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vrouilhac/webserver/databases"
)

type PageData struct {
	Title string
}

type TransactionPageData struct {
	Title        string
	Transactions []databases.Transaction
	HasAccounts  bool
}

func HandleGetHome(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	transactions := databases.GetTransactions()
	accounts := databases.GetAccounts()

	err := templ.ExecuteTemplate(w, "base", TransactionPageData{
		Title:        "Welcome to Budget App",
		Transactions: transactions,
		HasAccounts:  len(accounts) != 0,
	})

	if err != nil {
		log.Println(err)
	}
}

func HandleGetTransactions(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	transactions := databases.GetTransactions()
	accounts := databases.GetAccounts()

	err := templ.ExecuteTemplate(w, "transactions-list", TransactionPageData{
		Title:        "Welcome to Budget App",
		Transactions: transactions,
		HasAccounts:  len(accounts) != 0,
	})

	if err != nil {
		log.Println(err)
	}
}

func HandleHome(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetHome(w, r, templates)
		}
	}
}

type AddModalData struct {
	Accounts []databases.Account
}

func HandleGetAddModal(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	accounts := databases.GetAccounts()

	log.Println(accounts)

	err := templ.ExecuteTemplate(w, "add-modal", AddModalData{
		Accounts: accounts,
	})

	if err != nil {
		log.Println(err)
	}
}

func HandleAddModal(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetAddModal(w, r, templates)
		}
	}
}

func HandleGetTransactionsDetails(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	id := strings.TrimPrefix(r.URL.Path, "/components/transaction-row/details/")
	err, transaction := databases.GetTransactionByID(id)

	if err != nil || transaction == nil {
		return
	}

	// First Get Transaction
	t_err := templ.ExecuteTemplate(w, "transaction-row-details", transaction)

	if t_err != nil {
		log.Println(err)
	}
}

func HandleTransactionDetails(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetTransactionsDetails(w, r, templates)
		}
	}
}

func HandleGetTransactionsRow(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	id := strings.TrimPrefix(r.URL.Path, "/components/transaction-row/")
	err, transaction := databases.GetTransactionByID(id)

	if err != nil || transaction == nil {
		return
	}

	t_err := templ.ExecuteTemplate(w, "transaction-row", transaction)

	if t_err != nil {
		log.Println(t_err)
	}
}

func HandleTransactionRow(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetTransactionsRow(w, r, templates)
		}
	}
}

func HandleStatic() {
	fileServer := http.FileServer(http.Dir("./dist"))
	fileHandler := http.StripPrefix("/static/", fileServer)
	http.Handle("/static/", fileHandler)
}

func HandlePostTransaction(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	err := r.ParseForm()

	if err != nil {
		log.Println("Error: Can't read FormData")
		return
	}

	amount, parseErr := strconv.ParseFloat(r.PostForm.Get("amount"), 32)
	description := r.PostForm.Get("description")

	if parseErr != nil {
		log.Println("Error: Can't read FormData")
		return
	}

	id := strconv.Itoa(rand.Intn(1000))

	var op databases.Operation

	if amount < 0 {
		op = databases.Sub
		amount = math.Abs(amount)
	} else {
		op = databases.Add
	}

	databases.AddTransaction(databases.Transaction{
		ID:          id,
		Amount:      float32(amount),
		Operation:   op,
		Description: description,
		Date:        time.Now().Unix(),
	})
	w.Header().Set("HX-Trigger", "update-transactions")
}

func HandleDeleteTransaction(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	id := strings.TrimPrefix(r.URL.Path, "/transactions/")

	err := databases.DeleteTransactionById(id)

	if err != nil {
		log.Println("Error: Deleting item")
		err, transaction := databases.GetTransactionByID(id)

		if err != nil || transaction == nil {
			return
		}

		t_err := templ.ExecuteTemplate(w, "transaction-row", transaction)

		if t_err != nil {
			log.Println(err)
		}
	}

	transactions := databases.GetTransactions()

	if len(transactions) == 0 {
		templ.ExecuteTemplate(w, "empty-rows", nil)
	}
}

func HandleTransactions(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetTransactions(w, r, templates)
		case "POST":
			HandlePostTransaction(w, r, templates)
		case "DELETE":
			HandleDeleteTransaction(w, r, templates)
		}
	}
}

type AccountData struct {
	Accounts []databases.Account
}

func HandleGetAccountsComponent(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	accounts := databases.GetAccounts()

	log.Println(accounts)

	err := templ.ExecuteTemplate(w, "accounts-list", AccountData{
		Accounts: accounts,
	})

	if err != nil {
		log.Println(err)
	}
}

func HandleAccountsComponent(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetAccountsComponent(w, r, templates)
		}
	}
}

func HandlePostAccount(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	newAccount := databases.Account{
		ID:   strconv.Itoa(rand.Intn(1000)),
		Name: "Unnamed",
	}

	databases.AddAccount(newAccount)

	w.Header().Set("HX-Trigger", "update-accounts")
}

func HandleAccounts(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			HandlePostAccount(w, r, templates)
		case "PUT":
			HandlePutAccountEdit(w, r, templates)
		}
	}
}

func HandleGetAccountEdit(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	accountId := strings.TrimPrefix(r.URL.Path, "/components/accounts/edit/")
	a_err, account := databases.GetAccountById(accountId)

	if a_err != nil {
		log.Println("Error: Account")
		return
	}

	err := templ.ExecuteTemplate(w, "account-edit", account)

	if err != nil {
		log.Println("Error: Account")
		return
	}
}

func HandlePutAccountEdit(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	r.ParseForm()

	accountId := strings.TrimPrefix(r.URL.Path, "/accounts/")
	name := r.Form.Get("account-name")

	log.Println(name)

	err := databases.UpdateAccountById(accountId, name)

	if err != nil {
		log.Println("Error")
		return
	}

	a_err, account := databases.GetAccountById(accountId)

	if a_err != nil {
		log.Println("Error")
		return
	}

	t_err := templ.ExecuteTemplate(w, "account-row", account)

	if t_err != nil {
		log.Println(t_err)
	}
}

func HandleAccountEditComponent(templates *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetAccountEdit(w, r, templates)
		}
	}
}

func CreateServer(templates *template.Template) {
	HandleStatic()
	http.HandleFunc("/", HandleHome(templates))
	http.HandleFunc("/components/add-modal", HandleAddModal(templates))
	http.HandleFunc("/components/transaction-row/details/", HandleTransactionDetails(templates))
	http.HandleFunc("/components/transaction-row/", HandleTransactionRow(templates))
	http.HandleFunc("/components/accounts", HandleAccountsComponent(templates))
	http.HandleFunc("/components/accounts/edit/", HandleAccountEditComponent(templates))
	http.HandleFunc("/transactions", HandleTransactions(templates))
	http.HandleFunc("/transactions/", HandleTransactions(templates))
	http.HandleFunc("/accounts", HandleAccounts(templates))
	http.HandleFunc("/accounts/", HandleAccounts(templates))
}

func main() {
	var templates, temp_err = template.ParseGlob("tmpls/*/*.html") // Somehow **/* does not work for filepath.Glob
	templates, temp_err = templates.ParseGlob("tmpls/*.html")

	if temp_err != nil {
		log.Fatal("Error: cannot load templates")
	}

	CreateServer(templates)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Error")
		return
	}
}
