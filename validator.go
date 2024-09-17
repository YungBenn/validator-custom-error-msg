package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Help valid struct with tag custom
const tagCustom = "error"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func errorTagFunc[T interface{}](obj interface{}, snp string, fieldname, actualTag string) error {
	o := obj.(T)

	if !strings.Contains(snp, fieldname) {
		return nil
	}

	fieldArr := strings.Split(snp, ".")
	rsf := reflect.TypeOf(o)

	for i := 1; i < len(fieldArr); i++ {
		field, found := rsf.FieldByName(fieldArr[i])
		if found {
			if fieldArr[i] == fieldname {
				customMessage := field.Tag.Get(tagCustom)
				if customMessage != "" {
					return fmt.Errorf("%s", customMessage)
				}
				return nil
			} else {
				if field.Type.Kind() == reflect.Ptr {
					// If the field type is a pointer, dereference it
					rsf = field.Type.Elem()
				} else {
					rsf = field.Type
				}
			}
		}
	}
	return nil
}

func ValidateFunc[T interface{}](obj interface{}, validate *validator.Validate) ([]ValidationError, error) {
	o := obj.(T)
	var validationErrors []ValidationError

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Validate:", r)
		}
	}()

	if err := validate.Struct(o); err != nil {
		errorValid := err.(validator.ValidationErrors)
		for _, e := range errorValid {
			// snp  X.Y.Z
			snp := e.StructNamespace()
			errmgs := errorTagFunc[T](obj, snp, e.Field(), e.ActualTag())
			if errmgs != nil {
				validationErrors = append(validationErrors, ValidationError{Field: e.Field(), Message: errmgs.Error()})
			} else {
				validationErrors = append(validationErrors, ValidationError{Field: e.Field(), Message: e.Error()})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors, errors.New("validation errors")
	}

	return nil, nil
}
