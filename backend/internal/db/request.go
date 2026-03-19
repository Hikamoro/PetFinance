package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Transaction struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

func Income(apiHash string, amount int, description string, DB *sql.DB) error {
	_, err := DB.Exec("INSERT INTO income (api_hash, amount, description) VALUES ($1, $2, $3)", apiHash, amount, description)
	if err != nil {
		return err
	}
	return nil
}

func Expens(apiHash string, amount int, description string, DB *sql.DB) error {
	_, err := DB.Exec("INSERT INTO expenses (api_hash, amount, description) VALUES ($1, $2, $3)", apiHash, amount, description)
	if err != nil {
		return err
	}
	return nil
}

func GetBalance(apiHash string, DB *sql.DB) (int, error) {
	var balance int
	err := DB.QueryRow("SELECT COALESCE((SELECT SUM(amount) FROM income WHERE api_hash = $1), 0) - COALESCE((SELECT SUM(amount) FROM expenses WHERE api_hash = $1), 0)", apiHash).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func CheckIncome(apiHash string, DB *sql.DB) ([]Transaction, error) {
	rows, err := DB.Query("SELECT amount, description FROM income WHERE api_hash = $1", apiHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.Amount, &t.Description); err != nil {
			return nil, err
		}
		incomes = append(incomes, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return incomes, nil
}

func CheckExpens(apiHash string, DB *sql.DB) ([]Transaction, error) {
	rows, err := DB.Query("SELECT amount, description FROM expenses WHERE api_hash = $1", apiHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.Amount, &t.Description); err != nil {
			return nil, err
		}
		expenses = append(expenses, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}
