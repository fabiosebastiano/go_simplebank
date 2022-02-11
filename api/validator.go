package api

import (
	"github.com/fabiosebastiano/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validateCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {

	if currency, isOk := fieldLevel.Field().Interface().(string); isOk {
		//check if currency is supported
		return util.IsSupportedCurrency(currency)
	} else {
		return false
	}
}
