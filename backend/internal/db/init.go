package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
	DB_USER     = "postgres"
	DB_PASSWORD = "zxc"
	DB_NAME     = "finance_app"
)

func InitDB() {
	// подключаемся к postgres чтобы создать БД
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// создаем БД если нет
	_, err = db.Exec("CREATE DATABASE " + DB_NAME)
	if err != nil {
		fmt.Println("База возможно уже существует")
	}

	fmt.Println("Проверка базы данных завершена")

	// подключаемся уже к нужной БД
	psqlInfo = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Подключение к БД успешно")

	createTables()
}

func createTables() {
	// таблица пользователей
	usersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		api_hash TEXT UNIQUE NOT NULL
	);
	`

	// таблица финансов
	// financeTable := `
	// CREATE TABLE IF NOT EXISTS finance(
	// 	id SERIAL PRIMARY KEY,
	// 	api_hash TEXT NOT NULL,
	// 	id_user INT REFERENCES users(id) ON DELETE CASCADE,
	// 	expenses INT DEFAULT 0,
	// 	income INT DEFAULT 0,
	// 	priority TEXT
	// );
	// `
	Expenses := `
	CREATE TABLE IF NOT EXISTS expenses(
		id SERIAL PRIMARY KEY,
		api_hash TEXT NOT NULL,
		amount INT NOT NULL,
		description TEXT
	);
	`
	Income := `
	CREATE TABLE IF NOT EXISTS income(
		id SERIAL PRIMARY KEY,
		api_hash TEXT NOT NULL,
		amount INT NOT NULL,
		description TEXT
	);
	`
	_, err := DB.Exec(Expenses)
	if err != nil {
		log.Fatal("Ошибка создания таблицы Expenses:", err)
	}

	_, err = DB.Exec(Income)
	if err != nil {
		log.Fatal("Ошибка создания таблицы Income:", err)
	}

	_, err = DB.Exec(usersTable)
	if err != nil {
		log.Fatal("Ошибка создания таблицы users:", err)
	}

	// _, err = DB.Exec(financeTable)
	// if err != nil {
	// 	log.Fatal("Ошибка создания таблицы finance:", err)
	// }

	fmt.Println("Таблицы успешно проверены/созданы")
}
