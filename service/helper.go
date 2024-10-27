package service

import (
	"io"
	"os"
	"net/http"
	"mime/multipart"
	"golang.org/x/crypto/bcrypt"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
)

func CompareHashPassword(data1 string, data2 string, response map[string]interface{}) (int, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(data1), []byte(data2)); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
				response["message"] = "Password Salah Silahkan Coba Kembali"
				response["field"] = "passwordBefore"
				return http.StatusBadRequest, err
		default:
			response["message"] = err.Error()
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func CreateImage(file multipart.File, image string, response map[string]interface{}) error {
	pathFile  := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", image)

	outeFile, err := os.Create(pathFile)
	if err != nil {
		response["message"] = "Error cant Create file" + err.Error()
		return err
	}

	if _, err := io.Copy(outeFile, file); err != nil {
		response["message"] = "Error Cant copy file " + err.Error()
		os.Remove(pathFile)
		return err
	}

	return nil
} 

func RemoveFile(image string, response map[string]interface{}) error {
	pathFile := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", image)

	if err := os.Remove(pathFile); err != nil {
		response["message"] = err.Error()
		return err
	}

	return nil
}
