package mapper

import (
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/dto"
)

// konversi from dto to domain
func ToHadiahDomain(hadiah *dto.Hadiah) domain.Hadiah {
	return domain.Hadiah{
		ID: hadiah.ID,
		CreatedAt: hadiah.CreatedAt,
		UpdatedAt: hadiah.UpdatedAt,
		Nama_Barang: hadiah.Nama_Barang,
		Harga_Hadiah: hadiah.Harga_Hadiah,
		Poin: hadiah.Poin,
		Image: hadiah.Image,
		Deskripsi: hadiah.Deskripsi,
	}
}

// konversi from domain to dto
func ToHadiahDTO(hadiah domain.Hadiah) dto.Hadiah {
	return dto.Hadiah{
		ID: hadiah.ID,
		CreatedAt: hadiah.CreatedAt,
		UpdatedAt: hadiah.UpdatedAt,
		Nama_Barang: hadiah.Nama_Barang,
		Harga_Hadiah: hadiah.Harga_Hadiah,
		Poin: hadiah.Poin,
		Image: hadiah.Image,
		Deskripsi: hadiah.Deskripsi,
	}
}

func ToArrayUserDTO(user []domain.User) []dto.User {
	var datas []dto.User

	for _, data := range user {
		dto := ToUserDTO(data)
		datas = append(datas, dto)
	}

	return datas
}