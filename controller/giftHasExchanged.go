package controller

import (
	"log"
	"net/http" 
	config "github.com/RaihanMalay21/config-TB_Berkah_Jaya"
	helper "github.com/RaihanMalay21/helper_TB_Berkah_Jaya"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func GiftHasExchanged(w http.ResponseWriter, r *http.Request) {
	// mengambil id user di session
	session, err := config.Store.Get(r, "berkah-jaya-session")
	if err != nil {
		log.Println("Error cant get session:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idUser := session.Values["id"].(uint)

	// mengambil hadiah yang sudah di tukar oleh users dan tidak dapat di tukar kembali 
	var hadiahHaveChange []models.HadiahUser
	if err := config.DB.Where("user_id = ?", idUser).Find(&hadiahHaveChange).Error; err != nil {
		log.Println("Error cant get data hadiah have change:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.Response(w, hadiahHaveChange, http.StatusOK)
}