package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"petFinance/backend/crypto"
	"petFinance/backend/internal/db"

	_ "github.com/lib/pq"
)

type Server struct {
	db *sql.DB
}

// func (s *Server) test() {
// 	_, err := s.db.Exec("INSERT INTO users (login, password_hash, api_hash) VALUES ($1, $2, $3)", "testuser", "hashedpassword", "uniqueapihash")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func (s *Server) IssuesHandler(w http.ResponseWriter, r *http.Request) {
	// Простая CORS поддержка, чтобы фронтенд мог обращаться из другого origin (или file://)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch r.URL.Path {
	case "/register":
		if r.Method != http.MethodPost {
			Tmpl, _ := template.ParseFiles("./frontend/templates/auth/register.html")
			Tmpl.Execute(w, nil)
		} else {
			var req struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			if req.Login == "" || req.Password == "" {
				http.Error(w, "login and password required", http.StatusBadRequest)
				return
			}

			if err := db.AddUser(req.Login, req.Password, s.db); err != nil {
				http.Error(w, "unable to create user: "+err.Error(), http.StatusBadRequest)
				return
			}
			passwordHash := crypto.XorCrypto(req.Password, "zxcursed")
			apiHash := crypto.XorCrypto(req.Login+passwordHash, "zxcursed")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"api_hash": apiHash})
		}
	case "/auth":
		if r.Method != http.MethodPost {
			Tmpl, _ := template.ParseFiles("./frontend/templates/auth/auth.html")
			Tmpl.Execute(w, nil)
		} else {
			var req struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			users, err := db.GetAllUsers(s.db)
			if err != nil {
				http.Error(w, "unable to retrieve users: "+err.Error(), http.StatusInternalServerError)
				return
			}
			for _, user := range users {
				if user.Login == req.Login && user.PasswordHash == crypto.XorCrypto(req.Password, "zxcursed") {
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(map[string]string{"api_hash": user.ApiHash})
					return
				}
			}
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
		}
	case "/":
		Tmpl, _ := template.ParseFiles("./frontend/templates/main/index.html")
		Tmpl.Execute(w, nil)
	case "":
		// var Req struct {
		// 	ApiKey string `json:"apiKey"`
		// }
		// if err := json.NewDecoder(r.Body).Decode(&Req); err != nil {
		// 	http.Error(w, "invalid json", http.StatusBadRequest)
		// 	return
		// }
		// name, _ := db.GetNameByApiHash(Req.ApiKey, s.db)
		// type user struct {
		// 	Username string
		// }
		// r := user{Username: name}
		Tmpl, _ := template.ParseFiles("./frontend/templates/main/index.html")
		Tmpl.Execute(w, nil)
	// case "/auth":
	// 	Tmpl, _ := template.ParseFiles("../../frontend/templates/auth/auth.html")
	// 	Tmpl.Execute(w, nil)
	case "/addIncome":
		if r.Method != http.MethodPost {
			return
		}
		var req struct {
			ApiHash     string `json:"api_hash"`
			Amount      int    `json:"amount"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := db.Income(req.ApiHash, req.Amount, req.Description, s.db); err != nil {
			http.Error(w, "unable to add income: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case "/addExpens":
		if r.Method != http.MethodPost {
			return
		}
		var req struct {
			ApiHash     string `json:"api_hash"`
			Amount      int    `json:"amount"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := db.Expens(req.ApiHash, req.Amount, req.Description, s.db); err != nil {
			http.Error(w, "unable to add expense: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case "/getBalance":
		if r.Method != http.MethodPost {
			return
		}
		var req struct {
			ApiHash string `json:"api_hash"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		balance, err := db.GetBalance(req.ApiHash, s.db)
		if err != nil {
			http.Error(w, "unable to get balance: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int{"balance": balance})
	case "/checkIncome":
		if r.Method != http.MethodPost {
			return
		}
		var req struct {
			ApiHash string `json:"api_hash"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		incomes, err := db.CheckIncome(req.ApiHash, s.db)
		if err != nil {
			http.Error(w, "unable to check income: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string][]db.Transaction{"incomes": incomes})
	case "/checkExpens":
		if r.Method != http.MethodPost {
			return
		}
		var req struct {
			ApiHash string `json:"api_hash"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		expenses, err := db.CheckExpens(req.ApiHash, s.db)
		if err != nil {
			http.Error(w, "unable to check expense: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string][]db.Transaction{"expenses": expenses})
	}
}
func main() {
	// if err := db.InitDB(); err != nil {
	// 	log.Fatal("Ошибка инициализации БД:", err)
	// }
	db.InitDB()
	defer db.DB.Close()
	server := &Server{db: db.DB}
	//server.test()
	// err := db.AddUser("leha", "ebanidolboeb", db.DB)
	// if err != nil {
	// 	log.Fatal("Ошибка добавления пользователя:", err)
	// }
	mux := http.NewServeMux()
	// ("/register", server.IssuesHandler)
	// //mux.HandleFunc("/login", server.IssuesHandler)
	// mux.HandleFunc("/addIncome", server.IssuesHandler)
	// mux.HandleFunc("/addExpens", server.IssuesHandler)
	// mux.HandleFunc("/getBalance", server.IssuesHandler)
	// mux.HandleFunc("/checkIncome", server.IssuesHandler)
	// mux.HandleFunc("/checkExpens", server.IssuesHandler)
	// mux.HandleFunc("/auth", server.IssuesHandler)
	mux.HandleFunc("/", server.IssuesHandler)
	//mux.HandleFunc("", server.IssuesHandler)

	// Статика фронтенда (отвечает на /, /index.html и т.п.)
	//mux.Handle("/", http.FileServer(http.Dir("./frontend")))
	log.Println("Server started on :8080")
	//log.Println(db.GetUserById(0, db.DB))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/pkg/static/"))))
	// http.Handle("/static/", http.FileServer(http.Dir("../../frontend/pkg/")))
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServe(":8080", mux) // поднятие сервера
	if err != nil {
		log.Fatal(err)
	}
}

// type DATABASE struct {
// 	ApiHash  string
// 	IdUser   int
// 	Expenses int    //расходы
// 	Income   int    //доходы
// 	Priority string //важность траты
// }
