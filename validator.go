package validator

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"reflect"
	"strings"
)

type ValidationErrorMutator func(*ValidationErrors)

type ErrorFactory func(e validator.FieldError) error

type TagErrorMap map[string]ErrorFactory

type FieldValidation func(level validator.FieldLevel) error

type Validator struct {
	*validator.Validate
	tagErrorMap             map[string]error
	tagErrorMessageOverride map[string]ErrorFactory
	defaultErrorMessage     ErrorFactory
}

// Override the error message factory for a particular validation.
func (v *Validator) OverrideErrorMessage(tag string, factory ErrorFactory) {
	v.tagErrorMessageOverride[tag] = factory
}

// Set a default error message template for all validator errors.
func (v *Validator) DefaultErrorMessage(factory ErrorFactory) {
	v.defaultErrorMessage = factory
}

// Create a custom field tag validator that returns an error message.
func (v *Validator) SetFieldTagValidator(tag string, validation FieldValidation) {
	_ = v.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		err := validation(fl)
		if err != nil {
			v.tagErrorMap[tag] = err
			return false
		}
		return true
	})
}

func (v *Validator) Struct(s interface{}) error {
	v.tagErrorMap = map[string]error{}
	err := v.Struct(s)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		result := make(ValidationErrors, len(vErrs))
		for i, e := range vErrs {
			if m, ok := v.tagErrorMessageOverride[e.Tag()]; ok {
				result[i] = m(e)
			} else if m, ok := v.tagErrorMap[e.Tag()]; ok {
				result[i] = m
			} else {
				result[i] = v.defaultErrorMessage(e)
			}
		}
		return result
	}
	return nil
}

func NewValidator() *Validator {
	v := &Validator{
		Validate:    validator.New(),
		tagErrorMap: map[string]error{},
		tagErrorMessageOverride: map[string]ErrorFactory{
			"oneof": func(e validator.FieldError) error {
				return errors.New(fmt.Sprintf("%s is invalid. Valid values are %s.", e.Value(), e.Param()))
			},
			"tcp4_addr": func(e validator.FieldError) error {
				return errors.New(fmt.Sprintf("%s is not a valid tcp address", e.Value()))
			},
			"hostname": func(e validator.FieldError) error {
				return errors.New(fmt.Sprintf("%s is not a valid hostname", e.Value()))
			},
			"required_with": func(e validator.FieldError) error {
				return errors.New(fmt.Sprintf("%s is required with %s", strings.ToLower(e.Field()), strings.ToLower(e.Param())))
			},
			"required_without": func(e validator.FieldError) error {
				return errors.New(fmt.Sprintf("%s is required when %s is not set", strings.ToLower(e.Field()), strings.ToLower(e.Param())))
			},
			"required": func(e validator.FieldError) error {
				return errors.New(fmt.Sprintf("%s is required.", strings.ToLower(e.Field())))
			},
		},
		defaultErrorMessage: func(e validator.FieldError) error {
			return errors.New(fmt.Sprintf("%s is an invalid %s", e.Value(), strings.ToLower(e.Field())))
		},
	}

	v.SetFieldTagValidator("file", func(fl validator.FieldLevel) error {
		_, err := os.Stat(fl.Field().String())
		if err != nil {
			return errors.New(fmt.Sprintf("Could not find file %s.", fl.Field().String()))
		}
		return nil
	})

	v.SetFieldTagValidator("required_when", func(level validator.FieldLevel) error {
		params := strings.Split(level.Param(), " ")
		fieldName := params[0]
		expectedVal := params[1]
		field := level.Parent().FieldByName(fieldName)
		if reflect.Zero(level.Field().Type()).Interface() != level.Field().Interface() && field.Interface() != expectedVal {
			return errors.New(fmt.Sprintf("%s is required when %s is set to %s",
				strings.ToLower(level.StructFieldName()), strings.ToLower(fieldName), expectedVal))
		}
		return nil
	})

	return v
}
