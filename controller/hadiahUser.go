package controller

import (
	"log"
	"net/http"

	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)


func HadiahUser(w http.ResponseWriter, r *http.Request) {
	var hadiahUser []models.HadiahUser

	// session, err := config.Store.Get(r, "berkah-jaya-session")
	// if err != nil {
	// 	log.Println("cannot sign in to session")
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

	if err := config.DB.Preload("Hadiah").Where("user_id = ? AND (gifts_arrive = ? OR status = ?)", idUser, "NO", "unfinished").Find(&hadiahUser).Error; err != nil {
		log.Println("Error Can't Get data hadiah user from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}  

	helper.Response(w, hadiahUser, http.StatusOK)
	return
}