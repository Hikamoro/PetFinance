package db

import (
	"database/sql"
	"petFinance/backend/crypto"

	// "fmt"
	// "log"
	_ "github.com/lib/pq"
)

type Users struct {
	Id           int
	Login        string
	PasswordHash string
	ApiHash      string
}

func AddUser(login, password string, DB *sql.DB) error {
	passwordHash := crypto.XorCrypto(password, "zxcursed")
	apiHash := crypto.XorCrypto(login+passwordHash, "zxcursed")
	_, err := DB.Exec("INSERT INTO users (login, password_hash, api_hash) VALUES ($1, $2, $3)", login, passwordHash, apiHash)
	if err != nil {
		return err
	}
	return nil
}
func GetUserById(id int, DB *sql.DB) (string, error) {
	var login string
	err := DB.QueryRow("SELECT login FROM users WHERE id = $1", id).Scan(&login)
	if err != nil {
		return "", err
	}
	return login, nil
}

func GetAllUsers(DB *sql.DB) ([]Users, error) {
	rows, err := DB.Query("SELECT id, login, password_hash, api_hash FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []Users
	for rows.Next() {
		var u Users
		if err := rows.Scan(&u.Id, &u.Login, &u.PasswordHash, &u.ApiHash); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func UpdateUser(id int, newLogin, newPassword string, DB *sql.DB) error {
	passwordHash := crypto.XorCrypto(newPassword, "zxcursed")
	apiHash := crypto.XorCrypto(newLogin+passwordHash, "zxcursed")
	_, err := DB.Exec("UPDATE users SET login = $1, password_hash = $2, api_hash = $3 WHERE id = $4", newLogin, passwordHash, apiHash, id)
	if err != nil {
		return err
	}
	return nil
}
func DeleteUser(id int, DB *sql.DB) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
func GetNameByApiHash(apiHash string, DB *sql.DB) (string, error) {
	var name string
	err := DB.QueryRow("SELECT login FROM users WHERE api_hash = $1", apiHash).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
