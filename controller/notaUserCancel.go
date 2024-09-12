package controller

import (
	"log"
	"net/http"

	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func NotaUserCanceled(w http.ResponseWriter, r *http.Request) {
	// retreaving id user from session
	// session, err := config.Store.Get(r, "berkah-jaya-session")
	// if err != nil {
	// 	log.Println("Error cant retreaving id fron session function ProsesPoinUser")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// idUser := session.Values["id"].(uint)

	idUser, err := helper.GetIDFromToken(r)
	if err != nil {
		log.Println(err)
		message := map[string]interface{}{"message": err.Error()}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}

	// retreaving data pembelians from database
	var pembelians []models.Pembelian
	if err := config.DB.Where("user_id = ? and status = ?", idUser, "cancel").Preload("KeteranganNotaCancel").Find(&pembelians).Error; err != nil {
		log.Println("Error cant retreaving data pembelian from database function ProsesPoinUser")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.Response(w, pembelians, http.StatusOK)
}