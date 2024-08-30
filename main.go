package main 

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	config "github.com/RaihanMalay21/config-tb-berkah-jaya"
	middlewares "github.com/RaihanMalay21/middlewares_TB_Berkah_Jaya"
	controller "github.com/RaihanMalay21/server-customer-TB-Berkah-Jaya/controller"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/berkahjaya/get/hadiah", controller.Hadiah).Methods("GET")

	config.DB_Connection()
	api := r.PathPrefix("/berkahjaya").Subrouter()
	api.HandleFunc("/gifts/have/change/user", controller.GiftHasExchanged).Methods("GET")
	api.HandleFunc("/users/data", controller.DataUser).Methods("GET")
	api.HandleFunc("/proses/poin/verify", controller.NotaUserCanceled).Methods("GET")
	api.HandleFunc("/scan/poin", controller.InputNota).Methods("POST")
	api.HandleFunc("/tukar/poin/hadiah", controller.ExchangePoin).Methods("POST") 
	api.HandleFunc("/user/proses/hadiah", controller.HadiahUser).Methods("GET")
	api.HandleFunc("/user/remove/nota/not/valid", controller.RemoveSubmissionPoin).Methods("POST")

	api.HandleFunc("/change/password", controller.ChangePassword).Methods("POST")
	api.Use(middlewares.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", 
		handlers.CORS(
			handlers.AllowedOrigins([]string{"http://localhost:3000", "https://fe-tb-berkah-jaya-igcfjdj5fa-uc.a.run.app"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowCredentials(),
		)(r)))

}