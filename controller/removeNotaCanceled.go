package controller

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func RemoveSubmissionPoin (w http.ResponseWriter, r *http.Request) {
	var pembelian models.Pembelian
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&pembelian); err != nil {
		log.Println("Error Cant Decode json:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// menghapus secara soft delete di table pembelian
	if err := config.DB.Delete(&pembelian).Error; err != nil {
		log.Println("Error cant delete data pembelian:", err.Error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := map[string]string{"message": "Berhasil Menghapus Data"}
	helper.Response(w, message, http.StatusOK)
}