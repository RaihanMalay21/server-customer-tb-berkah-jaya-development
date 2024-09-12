package controller

import (
	"log"
	"net/http" 
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func GiftHasExchanged(w http.ResponseWriter, r *http.Request) {
	// mengambil id user di session
	// session, err := config.Store.Get(r, "berkah-jaya-session")
	// if err != nil {
	// 	log.Println("Error cant get session:", err.Error())
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

	// mengambil hadiah yang sudah di tukar oleh users dan tidak dapat di tukar kembali 
	var hadiahHaveChange []models.HadiahUser
	if err := config.DB.Where("user_id = ?", idUser).Find(&hadiahHaveChange).Error; err != nil {
		log.Println("Error cant get data hadiah have change:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.Response(w, hadiahHaveChange, http.StatusOK)
}