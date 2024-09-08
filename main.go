package main 

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	// "github.com/gorilla/handlers"
	// "fmt"

	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	middlewares "github.com/RaihanMalay21/middlewares_TB_Berkah_Jaya"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya/controller"
)

func main() {
	r := mux.NewRouter()
	config.DB_Connection()
	
	r.HandleFunc("/berkahjaya/get/hadiah", controller.Hadiah).Methods("GET", "OPTIONS")
	// r.Use(corsMiddlewares)
	
	api := r.PathPrefix("/berkahjaya").Subrouter()
	api.Use(middlewares.JWTMiddleware)
	api.HandleFunc("/gifts/have/change/user", controller.GiftHasExchanged).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/data", controller.DataUser).Methods("GET", "OPTIONS")
	api.HandleFunc("/proses/poin/verify", controller.NotaUserCanceled).Methods("GET", "OPTIONS")
	api.HandleFunc("/scan/poin", controller.InputNota).Methods("POST", "OPTIONS")
	api.HandleFunc("/tukar/poin/hadiah", controller.ExchangePoin).Methods("POST", "OPTIONS") 
	api.HandleFunc("/user/proses/hadiah", controller.HadiahUser).Methods("GET", "OPTIONS")
	api.HandleFunc("/user/remove/nota/not/valid", controller.RemoveSubmissionPoin).Methods("POST", "OPTIONS")
	api.HandleFunc("/change/password", controller.ChangePassword).Methods("POST", "OPTIONS")
	
	log.Fatal(http.ListenAndServe(":8081", r))
}

// func corsMiddlewares(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")
// 		fmt.Println("Origin received:", origin)

// 		allowedOrigins := "http://localhost:3000"

// 		if origin == allowedOrigins {
// 			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
// 			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")
// 			w.Header().Set("Access-Control-Allow-Credentials", "true")
// 		}

// 		if r.Method == http.MethodOptions {
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }