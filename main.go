package main 

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	// "fmt"
	// "os"

	config "github.com/RaihanMalay21/config-tb-berkah-jaya"
	middlewares "github.com/RaihanMalay21/middlewares_TB_Berkah_Jaya"
	controller "github.com/RaihanMalay21/server-customer-TB-Berkah-Jaya/controller"
)

func main() {
	r := mux.NewRouter()
	config.DB_Connection()
	
	r.HandleFunc("/berkahjaya/get/hadiah", controller.Hadiah).Methods("GET")
	r.Use(corsMiddlewares)
	
	api := r.PathPrefix("/berkahjaya").Subrouter()
	api.Use(middlewares.JWTMiddleware)
	api.HandleFunc("/gifts/have/change/user", controller.GiftHasExchanged).Methods("GET")
	api.HandleFunc("/users/data", controller.DataUser).Methods("GET")
	api.HandleFunc("/proses/poin/verify", controller.NotaUserCanceled).Methods("GET")
	api.HandleFunc("/scan/poin", controller.InputNota).Methods("POST")
	api.HandleFunc("/tukar/poin/hadiah", controller.ExchangePoin).Methods("POST") 
	api.HandleFunc("/user/proses/hadiah", controller.HadiahUser).Methods("GET")
	api.HandleFunc("/user/remove/nota/not/valid", controller.RemoveSubmissionPoin).Methods("POST")
	api.HandleFunc("/change/password", controller.ChangePassword).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":8080", r))
}

func corsMiddlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println("Origin received:", origin)

		allowedOrigins := "https://fe-tb-berkah-jaya-750892348569.us-central1.run.app"

		if origin == allowedOrigins {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}