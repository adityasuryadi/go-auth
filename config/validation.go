package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// contract
type Validation interface {
	ValidateRequest(request interface{}) interface{}
}

func NewValidation(db *gorm.DB) Validation {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		// fmt.Println(fl.StructFieldName())
		// // get parameter dari tag struct validate
		table := fl.Param()
		field := strings.ToLower(fl.StructFieldName())
		var total int64
		err := db.Table(table).Where(""+field+" = ?", fl.Field()).Count(&total).Error
		if err != nil {
			fmt.Println(err)
		}
		// // Return true if the count is zero (i.e., the value is unique)
		return total == 0
	})

	// validasi gte than now
	validate.RegisterValidation("gtenow", func(fl validator.FieldLevel) bool {
		bookingDate, _ := time.Parse("2006-01-02", fl.Field().String())
		now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		diff := bookingDate.Sub(now)
		if diff >= 0 {
			return true
		}
		return false
	})

	return &ValidationImpl{
		Validate: validate,
	}
}

func (validateImpl *ValidationImpl) ValidateRequest(request interface{}) interface{} {
	err := validateImpl.Validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		out := make([]ErrorMessage, len(validationErrors))
		for i, fieldError := range validationErrors {
			out[i] = ErrorMessage{
				Field:   fieldError.Field(),
				Message: GetErrorMsg(fieldError),
			}
		}
		return out
	}
	return nil
}

type ValidationImpl struct {
	Validate *validator.Validate
}
