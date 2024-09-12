package controller

import (
	"log"
	"net/http"

	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func DataUser(w http.ResponseWriter, r *http.Request) {
	// mengambil id user dari session
	// session, err := config.Store.Get(r, "berkah-jaya-session")
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// id := session.Values["id"]

	id, err := helper.GetIDFromToken(r)
	if err != nil {
		log.Println(err)
		message := map[string]interface{}{"message": err.Error()}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}

	var getDataUser models.User
	if err := config.DB.First(&getDataUser, id).Error; err != nil{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := models.User{
		UserName: getDataUser.UserName,
		Email: getDataUser.Email,
		NoWhatshapp: getDataUser.NoWhatshapp,
		Poin: getDataUser.Poin,
	}

	helper.Response(w, data, http.StatusOK)
	return
}