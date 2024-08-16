package controller

import (
	"log"
	"net/http"
	"encoding/json"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

	config "github.com/RaihanMalay21/config-TB_Berkah_Jaya"
	helper "github.com/RaihanMalay21/helper_TB_Berkah_Jaya"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var FieldData map[string]string
	json := json.NewDecoder(r.Body)
	if err := json.Decode(&FieldData); err != nil {
		log.Println("Error function ChangePassword cant decode json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// retreave id from session
	session, err := config.Store.Get(r, "berkah-jaya-session")
	if err != nil {
		log.Println("Error function ChangePassword cant get session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	IDUser := session.Values["id"]

	
	// retreaving password from database
	var dataUser models.User
	if err := config.DB.Select("password").Find(&dataUser, "email = ?", FieldData["email"]).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
				log.Println("Error Funtion ChangePassword:", err)
				message := map[string]string{"message": "username atau email tidak ditemukan"}
				helper.Response(w, message, http.StatusBadRequest)
				return
		default:
			log.Println("Error Function ChangePassword:", err)
			helper.Response(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// camparing pass from client and database
	if err := bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(FieldData["passwordBefore"])); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
				log.Println("Error Function ChangePassword comparation password fail:", err)
				message := map[string]string{"message": "Password Salah Silahkan Coba Kembali", "field": "passwordBefore"}
				helper.Response(w, message, http.StatusBadRequest)
				return			
		default:
			log.Println("Error Function ChangePassword:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// memeriksa apakah new password len min 6
	if len(FieldData["passwordNew"]) < 6 {
		log.Println("Error New Password less than min len caracter")
		message := map[string]string{"message": "Minimal Panjang Password Baru 6 Karakter", "field": "passwordNew"}
		helper.Response(w, message, http.StatusBadRequest)
		return
	}

	// hash new password user
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(FieldData["passwordNew"]), bcrypt.DefaultCost)
	newPasswordString := string(hashPassword)

	// change password that exist in database
	if err := config.DB.Model(&models.User{}).Where("id = ?", IDUser).Update("password", newPasswordString).Error; err != nil {
		log.Println("Error Function ChangePassword Cant update password:", err)
		message := map[string]string{"message": "Tidak Dapat Merubah Password"}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}

	message := map[string]string{"message": "Berhasil Merubah Password"}
	helper.Response(w, message, http.StatusOK)
}