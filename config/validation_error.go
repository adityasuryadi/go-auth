package config

import (
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Message string
}

func (validationError ValidationError) Error() string {
	return validationError.Message
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field tidak boleh kosong"
	case "lte":
		return "harus lebih kecil dari " + fe.Param()
	case "gtenow":
		return "harus lebih besar dari tanggal hari ini"
	case "gte":
		return "harus lebih besar dari " + fe.Param()
	case "email":
		return "format email salah"
	case "unique":
		return fe.Field() + " exist"
	case "min":
		return "minimal " + fe.Param() + " karakter"
	case "max":
		return "max " + fe.Param() + " Kb"
	case "image_validation":
		return "Harus Image"
	case "number":
		return "harus numeric"
	case "eqfield":
		return "field harus sama dengan " + fe.Param()
	}
	return "Unknown error"
}
