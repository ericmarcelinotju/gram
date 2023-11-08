package utils

import (
	"math/rand"
	"reflect"
	"regexp"
	"time"
)

func IsEmailValid(e string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func ConvertToMap(model interface{}) (ret map[string]interface{}) {

	ret = make(map[string]interface{})

	modelReflect := reflect.ValueOf(model)

	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()

	var fieldData interface{}

	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Ptr:
			fieldData = ConvertToMap(field.Interface())
		default:
			fieldData = field.Interface()
		}

		ret[modelRefType.Field(i).Name] = fieldData
	}

	return ret
}

func RandTime(minDate, maxDate time.Time) time.Time {
	min := minDate.Unix()
	max := maxDate.Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
