package repository

import (
	"gorm.io/gorm"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/domain"
)

type RepositoryCustomer interface {
	BeginNewTransaction() *gorm.DB
	GetHadiah() ([]domain.Hadiah, error)
	GetHadiahUser(userID uint) ([]domain.HadiahUser, error)
	GetDataUser(userID uint) (domain.User, error)
	GetPembeliansNotaCanceled(userID uint) ([]domain.Pembelian, error)
	GetProsesHadiahUser(userID uint) ([]domain.HadiahUser, error)
	GetFileNotaUser(data *domain.Pembelian) ([]domain.Pembelian, error) 
	InputPembelian(data *domain.Pembelian, tx *gorm.DB) error
	GetPoinHadiah(id uint) (domain.Hadiah, error) 
	GetPoinUser(id uint) (domain.User, error)
	UpdatePoinUser(userID uint, poin float64, tx *gorm.DB) error
	CreateHadiahUser(data *domain.HadiahUser, tx *gorm.DB) error
	GetImagePembelian(id uint) (domain.Pembelian, error)
	DeletePembelian(id uint, tx *gorm.DB) error
	RetreavingPassword(email string) (domain.User, error)
	UpdatePassword(userID uint, password string) error
}

type repositoryCustomer struct {
	db *gorm.DB
}

func NewRepositoryCustomer(db *gorm.DB) RepositoryCustomer {
	return &repositoryCustomer{db: db}
}

func (rc *repositoryCustomer) BeginNewTransaction() *gorm.DB {
	return rc.db.Begin()
}

func (rc *repositoryCustomer) GetHadiah() ([]domain.Hadiah, error) {
	var hadiah []domain.Hadiah
	if err := rc.db.Find(&hadiah).Error; err != nil {
		return nil, err
	}

	return hadiah, nil
}

func (rc *repositoryCustomer) GetHadiahUser(userID uint) ([]domain.HadiahUser, error) {
	var HadiahUser []domain.HadiahUser
	if err := rc.db.Where("user_id = ?", userID).Find(&HadiahUser).Error; err != nil {
		return nil, err
	}

	return HadiahUser, nil
}

func (rc *repositoryCustomer) GetDataUser(userID uint) (domain.User, error) {
	var user domain.User
	if err := rc.db.Select("user_name", "email", "no_whatshapp", "poin").First(&user, userID).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (rc *repositoryCustomer) GetFileNotaUser(data *domain.Pembelian) ([]domain.Pembelian, error) {
	var fileExist []domain.Pembelian
	if err := rc.db.Select("image").Where("user_id = ?", data.UserID).Find(&fileExist).Error; err != nil{
		return nil, err
	}

	return fileExist, nil
}

func (rc *repositoryCustomer) GetPembeliansNotaCanceled(userID uint) ([]domain.Pembelian, error) {
	var data []domain.Pembelian
	if err := rc.db.Where("user_id = ? and status = ?", userID, "cancel").Preload("KeteranganNotaCancel").Find(&data).Error; err != nil {
		return nil, err
	}

	return  data, nil
}

func (rc *repositoryCustomer) GetProsesHadiahUser(userID uint) ([]domain.HadiahUser, error) {
	var data []domain.HadiahUser
	if err := rc.db.Preload("Hadiah").Where("user_id = ? and (status = ? or gifts_arrive = ?)", userID, "unfinished", "NO").Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (rc *repositoryCustomer) InputPembelian(data *domain.Pembelian, tx *gorm.DB) error {
	if err := tx.Omit("keterangan_nota_cancel_id").Create(data).Error; err != nil {
		return err
	}

	return nil
} 

func (rc *repositoryCustomer) GetPoinHadiah(id uint) (domain.Hadiah, error) {
	var data domain.Hadiah
	if err := rc.db.Model(domain.Hadiah{}).Select("poin").Where("id = ?", id).Take(&data).Error; err != nil {
		return domain.Hadiah{}, err
	}

	return data, nil
}

func (rc *repositoryCustomer) GetPoinUser(id uint) (domain.User, error) {
	var data domain.User
	if err := rc.db.Model(domain.User{}).Select("poin").Where("id = ?", id).Take(&data).Error; err != nil {
		return domain.User{}, err
	}

	return data, nil
}

func (rc *repositoryCustomer) UpdatePoinUser(userID uint, poin float64, tx *gorm.DB) error {
	if err := tx.Model(domain.User{}).Where("ID = ?", userID).Update("poin", poin).Error; err != nil {
		return err
	}

	return nil
}

func (rc *repositoryCustomer) CreateHadiahUser(data *domain.HadiahUser, tx *gorm.DB) error {
	if err := tx.Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (rc *repositoryCustomer) GetImagePembelian(id uint) (domain.Pembelian, error) {
	var data domain.Pembelian
	if err := rc.db.Model(domain.Pembelian{}).Where("id = ?", id).Select("image").Take(&data).Error; err != nil {
		return domain.Pembelian{}, err
	}

	return data, nil
}

func (rc *repositoryCustomer) DeletePembelian(id uint, tx *gorm.DB) error {
	var pembelian domain.Pembelian
	if err := tx.Where("id = ?", id).Delete(&pembelian).Error; err != nil {
		return err
	}

	return nil
}

func (rc *repositoryCustomer) RetreavingPassword(email string) (domain.User, error) {
	var dataUser domain.User
	if err := rc.db.Select("password").Find(&dataUser, "email = ?", email).Error; err != nil {
		return domain.User{}, err
	}

	return dataUser, nil
}

func (rc *repositoryCustomer) UpdatePassword(userID uint, password string) error {
	if err := rc.db.Model(&domain.User{}).Where("id = ?", userID).Update("password", password).Error; err != nil {
		return err
	}

	return nil
}

