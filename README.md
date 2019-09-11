# Validator
A sane wrapper around the standard [validator](https://github.com/go-playground/validator) library. Supports errors for 
custom validators, and a few more useful validators.

## Installation
```
go get github.com/eddieowens/validator
```

## Usage
```go
package main

import (
    "errors"
    "fmt"
    "github.com/eddieowens/validator"
    ogvalid "gopkg.in/go-playground/validator.v9"
)

func main() {
    type myStruct struct {
        Required                    string `validate:"required"`
        AnswerToTheUltimateQuestion int    `validate:"answer_to_the_ultimate_question"`
    }
    
    v := validator.NewValidator()
    
    v.SetFieldTagValidator("answer_to_the_ultimate_question", func(level ogvalid.FieldLevel) error {
        i := level.Field().Int()
        if i == 42 {
            return nil
        }
        return fmt.Errorf("%d is not the ultimate answer to life, the universe, and everything!", i)
    })
    
    s := myStruct{
        Required:                    "required",
        AnswerToTheUltimateQuestion: 42,
    }
    
    if err := v.Struct(s); err != nil {
        panic(err)
    }
    
    // Success!
}
```

A few methods were added to the original validator package to allow you to better manage errors.