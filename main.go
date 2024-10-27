package main 

import (
	"log"
	"net/http"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/router"
)

func main() {
	r := router.InitRouter()
	
	log.Fatal(http.ListenAndServe(":8081", r))
}
