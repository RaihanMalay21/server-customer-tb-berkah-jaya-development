package controller

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/go-playground/validator/v10"

	config "github.com/RaihanMalay21/config-TB_Berkah_Jaya"
	helper "github.com/RaihanMalay21/helper_TB_Berkah_Jaya"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func ExchangePoin(w http.ResponseWriter, r *http.Request) {
	// mengambil data barang yang akan di change with point
	var Hadiah models.Hadiah
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&Hadiah); err != nil {
		log.Println("error Cannot decode data barang:", err)
		message := map[string]interface{}{"message": err}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}

	// validate struct hadiah
	validate := validator.New(validator.WithRequiredStructEnabled())
	trans := helper.TranslatorIDN()

	hadiahs := &models.Hadiah{
		Nama_Barang: Hadiah.Nama_Barang,
		Poin: Hadiah.Poin, 
		Image: Hadiah.Image,
		Deskripsi: Hadiah.Deskripsi,
	}

	if err := validate.Struct(hadiahs); err != nil {
		errors := err.(validator.ValidationErrors)
		errMessage := errors.Translate(trans)
		log.Println("Error validate struct", errMessage)
		message := map[string]interface{}{"message": errMessage}
		helper.Response(w, message, http.StatusBadRequest)
		return
	}

	// mengambil poin from tabel barang untuk authentikasi poin barang
	var hadiah models.Hadiah
	if err := config.DB.Model(models.Hadiah{}).Select("poin").Where("ID = ?", Hadiah.ID).Take(&hadiah).Error; err != nil {
		log.Println("cannot retreaving poin from table barang")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// authentika poin form client with in database
	if Hadiah.Poin != hadiah.Poin {
		log.Println("Error Poin Not same with in database")
		message := map[string]string{"message": "Error Poin different with in Database"}
		helper.Response(w, message, http.StatusBadRequest)
		return
	}

	// take id from session
	session, err := config.Store.Get(r, "berkah-jaya-session")
	if err != nil {
		log.Println("cannot sign in to session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idUser := session.Values["id"].(uint)
    
	// mengambil data poin yang dimiliki oleh user
	var user models.User
	if err := config.DB.Model(models.User{}).Select("poin").Where("ID = ?", idUser).Take(&user).Error; err != nil {
		log.Println("cannot retreaving poin from table user in database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// authentikasi barang yang ingin di tukar oleh user
	if user.Poin < hadiah.Poin {
		log.Println("poin user not enough")
		message := map[string]string{"message": "poin not enough"}
		helper.Response(w, message, http.StatusBadRequest)
		return
	} 

	poinUser := user.Poin - hadiah.Poin

	// inialisasi transaksi database gorm
	tx := config.DB.Begin()
	
	// update poin in table user
	if err := tx.Model(models.User{}).Where("ID = ?", idUser).Update("poin", poinUser).Error; err != nil {
		log.Println("Error don't update poin")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// masukkan idUser ke dalam table pengajuan_hadiah

	// insert on table relation bettwen table hadiah and user
	// apus primaryKey di table struct hadiahUser
	HadiahUser := models.HadiahUser{
		UserID: idUser, 
		HadiahID: Hadiah.ID,
		GiftsArrive: "NO",
		Status: "unfinished",
	}

	if err := tx.Create(&HadiahUser).Error; err != nil {
		log.Println("Cannot insert to table Hadiah User")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()

	message := map[string]string{"Message": "berhasil menukar poin, tungggu pemberitahuan selanjutnya dalam 7 hari kedepan di Email Anda"}
	helper.Response(w, message, http.StatusOK)
}