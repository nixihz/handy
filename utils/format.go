package utils

import (
	"crypto/md5"
	"fmt"
	"reflect"
)

func Md5(str string) string {
	passwordMd5 := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", passwordMd5)
}

func Struct2Map(obj interface{}) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}
