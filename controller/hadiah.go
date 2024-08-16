package controller

import (
	"log"
	"net/http"
	
	config "github.com/RaihanMalay21/config-tb-berkah-jaya"
	helper "github.com/RaihanMalay21/helper_TB_Berkah_Jaya"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func Hadiah(w http.ResponseWriter, r *http.Request) {
	// inialisasi field hadiah
	var gethadiah []models.Hadiah
	if err := config.DB.Find(&gethadiah).Error; err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// response dan kirim data ke client dalam bentuk json
	helper.Response(w, gethadiah, http.StatusOK)
}