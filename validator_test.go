package validator

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type testStruct struct {
	Proto string `validate:"required,oneof=tcp http"`
	Host  string `validate:"required_without=Addr,required_when=Proto http"`
	Addr  string `validate:"required_without=Host,required_when=Proto tcp"`
}

type ValidatorTest struct {
	suite.Suite
	validator *Validator
}

func (v *ValidatorTest) SetupTest() {
	v.validator = NewValidator().(*Validator)
}

func (v *ValidatorTest) TestRequiredWhen() {
	// -- Given
	//
	t := testStruct{
		Proto: "tcp",
		Host:  "asd",
	}

	// -- When
	//
	err := v.validator.Struct(t)

	// -- Then
	//
	v.EqualError(err, "host is required when proto is set to http")
}

func (v *ValidatorTest) TestRequiredWhenValid() {
	// -- Given
	//
	t := testStruct{
		Proto: "http",
		Host:  "asd",
	}

	// -- When
	//
	err := v.validator.Struct(t)

	// -- Then
	//
	v.NoError(err)
}

func TestValidatorTest(t *testing.T) {
	suite.Run(t, new(ValidatorTest))
}
