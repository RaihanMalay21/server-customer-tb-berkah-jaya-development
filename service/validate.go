package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/dto"
	"github.com/RaihanMalay21/server-customer-tb-berkah-jaya-development/helper"
)

func ValidateStructHadiah(data *dto.Hadiah) []map[string]string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	trans := helper.TranslatorIDN()

	var errors []map[string]string

	if err := validate.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			NameField := err.StructField()
			errTranlate := err.Translate(trans)
			errorMap := map[string]string{
				NameField: errTranlate,
			}
			errors = append(errors, errorMap)
		}
	}

	return errors
}