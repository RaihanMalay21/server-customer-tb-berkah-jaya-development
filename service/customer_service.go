package service

import (
	"errors"
	"net/http"
	"strconv"
	"gorm.io/gorm"
	"mime/multipart"
	"golang.org/x/crypto/bcrypt"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/repository"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/dto"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/mapper"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/domain"
)

type ServiceCustomer interface {
	GetHadiah(response map[string]interface{}) ([]dto.Hadiah, int)
	GetGiftHasExchanged(userID uint, response map[string]interface{}) ([]dto.HadiahUser, int)
	GetDataUser(userID uint, response map[string]interface{}) (dto.User, int)
	GetPembeliansNotaCanceled(userID uint, response map[string]interface{}) ([]dto.Pembelian, int)
	GetProsesHadiahUser(userID uint, response map[string]interface{}) (*[]dto.HadiahUser, int)
	InputNota(userID uint, file multipart.File, ext string, onlyFileName string, response map[string]interface{}) int
	RemoveSubmissionPoin(pembelianID uint, response map[string]interface{}) int
	ExchangePoin(data *dto.Hadiah, userID uint, response map[string]interface{}) int
	ChangePassword(userID uint, data map[string]string, response map[string]interface{}) (int, error)
}

type serviceCustomer struct {
	repo repository.RepositoryCustomer
}

func NewServiceCustomer(repo repository.RepositoryCustomer) ServiceCustomer {
	return &serviceCustomer{repo: repo}
}

func (sc *serviceCustomer) GetHadiah(response map[string]interface{}) ([]dto.Hadiah, int) {
	hadiah, err := sc.repo.GetHadiah()
	if err != nil {
		response["message"] = err.Error()
		return nil, http.StatusInternalServerError
	}

	var hadiahDTO []dto.Hadiah
	for _, data := range hadiah {
		dto := mapper.ToHadiahDTO(data)
		hadiahDTO = append(hadiahDTO, dto)
	}

	return hadiahDTO, http.StatusOK
}

func (sc *serviceCustomer) GetGiftHasExchanged(userID uint, response map[string]interface{}) ([]dto.HadiahUser, int) {
	hadiahUser, err := sc.repo.GetHadiahUser(userID)
	if err != nil {
		response["message"] = err.Error()
		return nil, http.StatusInternalServerError
	}

	var hadiahUserDTO []dto.HadiahUser
	for _, datas := range hadiahUser {
		data := mapper.ToPengajuanHadiahDTO(&datas)
		hadiahUserDTO = append(hadiahUserDTO, data)
	} 

	return hadiahUserDTO, http.StatusOK
}

func (sc *serviceCustomer) GetDataUser(userID uint, response map[string]interface{}) (dto.User, int) {
	user, err := sc.repo.GetDataUser(userID)
	if err != nil {
		response["message"] = err.Error()
		return dto.User{}, http.StatusInternalServerError
	}

	userDTO := mapper.ToUserDTO(user)

	return userDTO, http.StatusOK
}

func (sc *serviceCustomer) GetPembeliansNotaCanceled(userID uint, response map[string]interface{}) ([]dto.Pembelian, int) {
	datas, err := sc.repo.GetPembeliansNotaCanceled(userID)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		response["message"] = "Tidak ada pembelian atau nota yang di cancel"
		return nil, http.StatusBadRequest
	case err != nil:
		response["message"] = err.Error()
		return nil, http.StatusInternalServerError
	}

	var datasDTO []dto.Pembelian
	for _, data := range datas {
		dto := mapper.ToPembelianDTO(&data)
		datasDTO = append(datasDTO, dto)
	}

	return datasDTO, http.StatusOK
}

func (sc *serviceCustomer) GetProsesHadiahUser(userID uint, response map[string]interface{}) (*[]dto.HadiahUser, int) {
	datas, err := sc.repo.GetProsesHadiahUser(userID)
	if err != nil {
		response["message"] = err.Error()
		return nil, http.StatusInternalServerError
	}

	var datasDTO []dto.HadiahUser
	for _, data := range datas {
		dto := mapper.ToPengajuanHadiahDTO(&data)
		datasDTO = append(datasDTO, dto)
	}

	return &datasDTO, http.StatusOK
}

func (sc *serviceCustomer) InputNota(userID uint, file multipart.File, ext string, onlyFileName string, response map[string]interface{}) int {
	userIDString := strconv.Itoa(int(userID))
	notaUser := dto.Pembelian {
		UserID: userID,
	}

	notaUserDomain := mapper.ToPembelianDomain(&notaUser)

	fileExist, err := sc.repo.GetFileNotaUser(&notaUserDomain)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			response["message"] = err.Error()
			return http.StatusInternalServerError
		}
	}
	
	lenExistFile := strconv.Itoa(len(fileExist) + 1)
	notaUserDomain.Image = onlyFileName + lenExistFile + userIDString + ext

	var errs error
	tx := sc.repo.BeginNewTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if errs != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if errs = sc.repo.InputPembelian(&notaUserDomain, tx); errs != nil {
		response["message"] = errs.Error()
		return http.StatusInternalServerError
	}

	 if errs = CreateImage(file, notaUserDomain.Image, response); errs != nil {
		return http.StatusInternalServerError
	}

	response["message"] = "succesfuly Upload Nota Pembayaran, poin sedang dalam proses Kalkulasi"
	return http.StatusOK
}

func (sc *serviceCustomer) ExchangePoin(data *dto.Hadiah, userID uint, response map[string]interface{}) int {
	if err := ValidateStructHadiah(data); len(err) > 0 {
		response["messageField"] = err
		return http.StatusBadRequest
	}

	hadiahDomain := mapper.ToHadiahDomain(data)

	// mengambil poin from tabel barang untuk authentikasi poin barang
	if hadiahSaved, err := sc.repo.GetPoinHadiah(hadiahDomain.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response["message"] = "Error Hadiah Tidak Di Temukan"
			return http.StatusBadRequest
		} else if hadiahDomain.Poin != hadiahSaved.Poin {
			response["message"] = "Poin Hadiah Tidak Sesuai dengan yang tersimpan"
			return http.StatusBadRequest
		} else {
			response["message"] = err.Error()
			return http.StatusInternalServerError
		}
	} 

	userSaved, err := sc.repo.GetPoinUser(userID)
	if err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError
	} 

	if userSaved.Poin < hadiahDomain.Poin {
		response["message"] = "Poin Anda Tidak Cukup"
		return http.StatusBadRequest
	}

	poinUser := userSaved.Poin - hadiahDomain.Poin

	hadiahUser := domain.HadiahUser{
		UserID: userID,
		HadiahID: hadiahDomain.ID,
		GiftsArrive: "NO",
		Status: "unfinished",
	}

	var errs error
	tx := sc.repo.BeginNewTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if errs != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if errs = sc.repo.UpdatePoinUser(userID, poinUser, tx); errs != nil {
		response["message"] = errs.Error()
		return http.StatusInternalServerError
	}

	if errs =  sc.repo.CreateHadiahUser(&hadiahUser, tx); errs != nil {
		response["message"] = errs.Error()
		return http.StatusInternalServerError
	} 

	response["message"] = "berhasil menukar poin, tungggu pemberitahuan selanjutnya dalam 7 hari kedepan di Email Anda"
	return http.StatusOK
}

func (sc *serviceCustomer) RemoveSubmissionPoin(pembelianID uint, response map[string]interface{}) int {
	pembelian, err := sc.repo.GetImagePembelian(pembelianID) 
	if err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError
	}

	var errs error
	tx := sc.repo.BeginNewTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if errs != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if errs = RemoveFile(pembelian.Image, response); errs != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError
	}

	if errs := sc.repo.DeletePembelian(pembelianID, tx); errs != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError
	}

	response["message"] = "Berhasil Menghapus Data"
	return http.StatusOK
}

func (sc *serviceCustomer) ChangePassword(userID uint, data map[string]string, response map[string]interface{}) (int, error) {
	if len(data["passwordNew"]) < 6 {
		response["message"] = "Minimal Panjang Password Baru 6 Karakter"
		response["field"] = "passwordNew"
		return http.StatusBadRequest, nil
	}

	user, err := sc.repo.RetreavingPassword(data["email"])
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response["message"] = "username atau email tidak ditemukan"
			return http.StatusBadRequest, err
		} else {
			response["message"] = err.Error()
			return http.StatusInternalServerError, err
		}
	}

	if statusCode, err := CompareHashPassword(user.Password, data["passwordBefore"], response); statusCode != 200 && err != nil {
		return statusCode, err
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(data["passwordNew"]), bcrypt.DefaultCost)
	newPasswordString := string(hashPassword)

	if err := sc.repo.UpdatePassword(userID, newPasswordString); err != nil {
		response["message"] = "Error Tidak Dapat Mengganti Password"
		return http.StatusInternalServerError, err
	}

	response["message"] = "Berhasil Merubah password"
	return http.StatusOK, nil
}



