package controller

import (
	"net/http"
	"log"
	"crypto/sha256"
	"path/filepath"
	"errors"
	"os"
	"io"
	"encoding/hex"
	"gorm.io/gorm"
	"strconv"

	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func InputNota(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file")
	if err != nil {
		log.Println("error on line 10 function InputNota")
		message := map[string]interface{}{"message": "File Tidak Ditemukan"}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// mengambil ext dari file 
	ext := filepath.Ext(handle.Filename)
	if ext == "" || (ext != ".png" && ext != ".jpg" && ext != ".gif" && ext != ".jpeg" && ext != ".jpe") {
		log.Println("error on line 21 function input nota")
		http.Error(w, "Error Tidak Dapat Upload gambar, pastikan jenis gambar adalah png, jpg, jpeg, jpe gif", http.StatusInternalServerError)
		return
	}

	// mengambil nama filenya saja
	nameOnly := filepath.Base(handle.Filename[:len(handle.Filename)-len(ext)])
	hasher := sha256.Sum256([]byte(nameOnly))
	hashnameOnlyString := hex.EncodeToString(hasher[:])

	// mengambil id dari session
	// session, err := config.Store.Get(r, "berkah-jaya-session")
	// if err != nil {
	// 	log.Println("error on line 27 function input nota")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // mendapatkan id dari session
	// idUser := session.Values["id"]

	idUser, err := helper.GetIDFromToken(r)
	if err != nil {
		log.Println(err)
		message := map[string]interface{}{"message": err.Error()}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}	
	
	// konversi id user dari interface ke int
	// idUserUint, ok := idUser.(uint)
	// if !ok {
	// 	log.Println("error on line 55 function input nota")
	// 	http.Error(w, "cannot Get ID user", http.StatusInternalServerError)
	// 	return
	// }
	// id int 
	idInt := int(idUser)
	// konver idUserInt ke strring
	idUserString := strconv.Itoa(idInt)
	images := hashnameOnlyString + idUserString + ext

	notaUser := models.Pembelian {
		UserID: idUser,
		Image: images,
	}

	// mengecek name file agar tidak ada kesamaan 
	var fileExist []models.Pembelian
	err = config.DB.Where("user_id = ?", notaUser.UserID).Find(&fileExist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notaUser.Image = images
	} else {
		lenFileString := strconv.Itoa(len(fileExist))
		images := hashnameOnlyString + lenFileString + idUserString + ext
		notaUser.Image = images
	}

	tx := config.DB.Begin()

	// insert ke database
	if err := tx.Omit("keterangan_nota_cancel_id").Create(&notaUser).Error; err != nil {
		log.Println("error on line 90 function input nota")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// path folder dimana file akan dibuat nantinya
	pathFile := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", notaUser.Image)

	// creta img
	outeFile, err := os.Create(pathFile)
	if err != nil {
		log.Println("error on line 67 function input nota")
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outeFile.Close()

	// copy image 
	if _, err := io.Copy(outeFile, file); err != nil {
		// menghapus img 
		if err := os.Remove(pathFile); err != nil {
			log.Println("error on line 75 function input nota")
		}
		tx.Rollback()
		log.Println("error on line 75 function input nota")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Error cant commit transaction:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := map[string]string{"message" : "succesfuly Upload Nota Pembayaran, poin sedang dalam proses Kalkulasi"}
	helper.Response(w, message, http.StatusOK)
	return
}