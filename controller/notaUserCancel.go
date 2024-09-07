package controller

import (
	"log"
	"net/http"

	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	helper "github.com/RaihanMalay21/helper_TB_Berkah_Jaya"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func NotaUserCanceled(w http.ResponseWriter, r *http.Request) {
	// retreaving id user from session
	session, err := config.Store.Get(r, "berkah-jaya-session")
	if err != nil {
		log.Println("Error cant retreaving id fron session function ProsesPoinUser")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idUser := session.Values["id"].(uint)

	// retreaving data pembelians from database
	var pembelians []models.Pembelian
	if err := config.DB.Where("user_id = ? and status = ?", idUser, "cancel").Preload("KeteranganNotaCancel").Find(&pembelians).Error; err != nil {
		log.Println("Error cant retreaving data pembelian from database function ProsesPoinUser")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.Response(w, pembelians, http.StatusOK)
}