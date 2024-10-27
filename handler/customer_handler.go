package handler

import (
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/service"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
)

type HandlerCustomer interface {
	GetHadiah(w http.ResponseWriter, r *http.Request)
	GetGiftHasExchanged(w http.ResponseWriter, r *http.Request)
	GetDataUser(w http.ResponseWriter, r *http.Request)
	GetPembeliansNotaCanceled(w http.ResponseWriter, r *http.Request)
	GetProsesHadiahUser(w http.ResponseWriter, r *http.Request)
	InputNota(w http.ResponseWriter, r *http.Request)
	ExchangePoin(w http.ResponseWriter, r *http.Request)
	RemoveSubmissionPoin(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
}

type handleCustomer struct {
	service service.ServiceCustomer
}

func NewHandlerCustomer(service service.ServiceCustomer) HandlerCustomer {
	return &handleCustomer{service: service}
}


func (hc *handleCustomer) GetHadiah(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	hadiah, statusCode := hc.service.GetHadiah(response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, hadiah, statusCode)
}

func (hc *handleCustomer) GetGiftHasExchanged(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	userID, err := helper.GetIDFromToken(r)
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	hadiahUser, statusCode := hc.service.GetGiftHasExchanged(userID, response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, hadiahUser, statusCode)
}

func (hc *handleCustomer) GetDataUser(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	userID, err := helper.GetIDFromToken(r)
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	data, statusCode := hc.service.GetDataUser(userID, response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, data, statusCode)
}

func (hc *handleCustomer) GetPembeliansNotaCanceled(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	userID, err := helper.GetIDFromToken(r) 
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return 
	}

	data, statusCode := hc.service.GetPembeliansNotaCanceled(userID, response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return 
	}

	helper.Response(w, data, statusCode)
}

func (hc *handleCustomer) GetProsesHadiahUser(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	userID, err := helper.GetIDFromToken(r)
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	data, statusCode := hc.service.GetProsesHadiahUser(userID, response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, data, statusCode)
} 

func (hc *handleCustomer) InputNota(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	file, _, ext, onlyNameFile, statusCode := GetFile(r, response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	userID, err := helper.GetIDFromToken(r)
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	StatusCode := hc.service.InputNota(userID, file, ext, onlyNameFile, response)
	
	helper.Response(w, response, StatusCode)
}

func (hc *handleCustomer) ExchangePoin(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	hadiah, err := DecodeJsonToStructHadiah(r, response)
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	userID, err := helper.GetIDFromToken(r)
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	statusCode := hc.service.ExchangePoin(&hadiah, userID, response)
	
	helper.Response(w, response, statusCode)
}

func (hc *handleCustomer) RemoveSubmissionPoin(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	vars := mux.Vars(r)
	id := vars["id"]
	pembelianID, _:= strconv.ParseUint(id, 10, 0)

	statusCode := hc.service.RemoveSubmissionPoin(uint(pembelianID), response)

	helper.Response(w, response, statusCode)
}

func (hc *handleCustomer) ChangePassword(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	data, err := DecodeJsonToMap(r, response)
	if err != nil {
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	userID, err := helper.GetIDFromToken(r)
	if err != nil {
		response["message"] = err.Error()
		helper.Response(w, response, http.StatusInternalServerError)
		return 
	}

	statusCode, err := hc.service.ChangePassword(userID, data, response)
	if err != nil {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, response, statusCode)
}

