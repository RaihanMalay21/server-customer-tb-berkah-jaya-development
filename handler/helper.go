package handler

import (
	"net/http"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"crypto/sha256"
	"mime/multipart"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/dto"
)

func DecodeJsonToMap(r *http.Request, response map[string]interface{}) (map[string]string, error) {
	var data map[string]string
	
	json := json.NewDecoder(r.Body)
	if err := json.Decode(&data); err != nil {
		response["message"] = "Cant Decode json to map: " + err.Error()
		return nil, err
	}

	return data, nil
}

func DecodeJsonToStructHadiah(r *http.Request, response map[string]interface{}) (dto.Hadiah, error) {
	var data dto.Hadiah
	
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&data); err != nil {
		return dto.Hadiah{}, err
	}
	
	return data, nil
}

// mengambil file from request form data dan mengembalikan file, fileHeader, ext, dan nameonly dalam bentuk hashing string dan status code
func GetFile(r *http.Request, response map[string]interface{}) (multipart.File, *multipart.FileHeader, string, string, int) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			response["message"] = "Error File tidak ditemukan" + err.Error()
			return nil, nil, "", "", http.StatusBadRequest
		} else {
			response["message"] = err.Error()
			return nil, nil, "", "", http.StatusInternalServerError
		}
	}

	ext := filepath.Ext(fileHeader.Filename)

	nameOnly := filepath.Base(fileHeader.Filename[:len(fileHeader.Filename)-len(ext)])
	hasher := sha256.Sum256([]byte(nameOnly))
	hashnameOnlyString := hex.EncodeToString(hasher[:])

	return file, fileHeader, ext, hashnameOnlyString, http.StatusOK
} 

