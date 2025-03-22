package entity

import (
	"regexp"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
)

func RegisterValidators(validate *validator.Validate) error {
	err := validate.RegisterValidation("real_name", func(fl validator.FieldLevel) bool {
		name := fl.Field().String()

		re := regexp.MustCompile(`^[A-Za-zА-Яа-яЁё\s\-]{1,100}$`)

		return re.MatchString(name)
	})

	return errors.Wrap(err, "register validators")
}
