package rest

import (
	"errors"
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	_ "github.com/lib/pq"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

var translatorNotFoundErr = errors.New("translator not found")

type reqValidator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

type fieldValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type validationError struct {
	Details []fieldValidationError `json:"details"`
}

func newReqValidator() (reqValidator, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return reqValidator{}, translatorNotFoundErr
	}

	val := validator.New()

	return reqValidator{Validator: val, Translator: trans}, nil
}

func validateRequest(r interface{}, rv *reqValidator) (bool, validationError) {
	if err := rv.Validator.Struct(r); err != nil {
		fieldErrors := []fieldValidationError{}
		for _, e := range err.(validator.ValidationErrors) {
			fieldError := fieldValidationError{
				Field: e.Namespace(),
				Error: e.Translate(rv.Translator),
			}
			fieldErrors = append(fieldErrors, fieldError)
		}
		return false, validationError{Details: fieldErrors}
	}
	return true, validationError{}
}

func registerCustomTranslations(v *validator.Validate, trans ut.Translator) {
	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("passwd", "{0} is not strong enough", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})
}

func registerCustomValidations(v *validator.Validate) {
	passwdMinLength := 6
	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > passwdMinLength
	})
}
