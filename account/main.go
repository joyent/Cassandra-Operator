package account

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/niravpatel27/cassandra-operator-workshop/cassandra"
)

// Account struct
type Account struct {
	ID    gocql.UUID `json:"id,omitempty"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
	Age   int        `json:"age"`
	City  string     `json:"city"`
}

func (acc *Account) GetAccounts() ([]Account, error) {
	var accountList []Account
	m := map[string]interface{}{}

	query := "SELECT id,age,name,city,email FROM accounts"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		accountList = append(accountList, Account{
			ID:    m["id"].(gocql.UUID),
			Age:   m["age"].(int),
			Name:  m["name"].(string),
			Email: m["email"].(string),
			City:  m["city"].(string),
		})
		m = map[string]interface{}{}
	}
	return accountList, nil
}

func (acc *Account) CreateAccount() error {

	fmt.Println("creating a new account")

	if err := cassandra.Session.Query(`
		INSERT INTO accounts (id, name, email, city, age) VALUES (?, ?, ?, ?, ?)`,
		acc.ID, acc.Name, acc.Email, acc.City, acc.Age).Exec(); err != nil {
		return err
	}
	fmt.Println("Account created successfully:", acc.ID, acc.Name)
	return nil
}
