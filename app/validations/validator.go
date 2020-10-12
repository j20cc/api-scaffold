package validations

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
)

type customValidator struct {
	tag          string
	fn           func(validator.FieldLevel) bool
	translations map[string]string
}

type validators struct {
	m []customValidator
}

func (val *validators) add(v ...customValidator) {
	val.m = append(val.m, v...)
}

func (val *validators) getTranslations(tag string) map[string]string {
	var translations map[string]string
	for _, v := range val.m {
		//translations = append(translations, v.translations[locale])
		if v.tag == tag {
			translations = v.translations
		}
	}
	return translations
}

var customValidators = new(validators)

func init() {
	customValidators.add(hasUser)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		entranslator := en.New()
		cntranslator := zh.New()
		uni = ut.New(entranslator, cntranslator)

		//修改tag
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		//注册自定验证和翻译
		for _, vali := range customValidators.m {
			_ = v.RegisterValidation(vali.tag, vali.fn)
			for lang, translation := range vali.translations {
				trans := GetTranslator(lang)
				switch lang {
				case "zh":
					_ = zhtranslations.RegisterDefaultTranslations(validate, trans)
					break
				case "en":
					_ = entranslations.RegisterDefaultTranslations(validate, trans)
					break
				}
				//自定义翻译
				_ = v.RegisterTranslation(vali.tag, trans, func(ut ut.Translator) error {
					return ut.Add(vali.tag, translation, true)
				}, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T(vali.tag, fe.Field())
					return t
				})
			}
		}
	}
}

// GetTranslator get validator translator
func GetTranslator(locale string) ut.Translator {
	if locale != "en" && locale != "zh" {
		locale = "en"
	}
	trans, _ := uni.GetTranslator(locale)
	return trans
}
