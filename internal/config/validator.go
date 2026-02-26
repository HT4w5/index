package config

import (
	"time"

	"github.com/docker/go-units"
	"github.com/go-playground/validator/v10"
)

func validateByteSize(fl validator.FieldLevel) bool {
	bs := fl.Field().String()
	_, err := units.FromHumanSize(bs)
	return err == nil
}

func validateDuration(fl validator.FieldLevel) bool {
	du := fl.Field().String()
	_, err := time.ParseDuration(du)
	return err == nil
}
