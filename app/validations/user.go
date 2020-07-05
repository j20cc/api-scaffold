package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var hasUser = customValidator{
	tag: "hasUser",
	fn: func(fl validator.FieldLevel) bool {
		fmt.Println(fl.Param(), fl.Field().Interface().(string))
		return false
	},
	translations: map[string]string{
		"zh": "用户不存在",
		"en": "user not exists",
	},
}
