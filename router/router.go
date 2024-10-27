package router

import (
	"github.com/gorilla/mux"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/middlewares"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/repository"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/service"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/handler"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	config.DB_Connection()
	db := config.DB

	repositoryCustomer := repository.NewRepositoryCustomer(db)
	serviceCustomer := service.NewServiceCustomer(repositoryCustomer)
	handlerCustomer := handler.NewHandlerCustomer(serviceCustomer)
	
	r.HandleFunc("/berkahjaya/get/hadiah", handlerCustomer.GetHadiah).Methods("GET", "OPTIONS")
	
	api := r.PathPrefix("/berkahjaya").Subrouter()
	api.Use(middlewares.JWTMiddleware)
	api.HandleFunc("/gifts/have/change/user", handlerCustomer.GetGiftHasExchanged).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/data", handlerCustomer.GetDataUser).Methods("GET", "OPTIONS")
	api.HandleFunc("/proses/poin/verify", handlerCustomer.GetPembeliansNotaCanceled).Methods("GET", "OPTIONS")
	api.HandleFunc("/scan/poin", handlerCustomer.InputNota).Methods("POST", "OPTIONS")
	api.HandleFunc("/tukar/poin/hadiah", handlerCustomer.ExchangePoin).Methods("POST", "OPTIONS") 
	api.HandleFunc("/user/proses/hadiah", handlerCustomer.GetProsesHadiahUser).Methods("GET", "OPTIONS")
	api.HandleFunc("/user/remove/nota/not/valid/{id}", handlerCustomer.RemoveSubmissionPoin).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/change/password", handlerCustomer.ChangePassword).Methods("POST", "OPTIONS")

	return r
}