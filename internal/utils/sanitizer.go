package utils

import (
	"reflect"

	"github.com/microcosm-cc/bluemonday"
)

func EscapeHTMLString(input string) string {
	p := bluemonday.NewPolicy()
	p.AllowElements()
	return p.Sanitize(input)
}

func escapeField(v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(EscapeHTMLString(v.String()))
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			escapeField(v.Field(i))
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			escapeField(v.Index(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			if value.IsValid() {
				escapeField(value)
				if key.Kind() == reflect.String {
					keyStr := EscapeHTMLString(key.String())
					v.SetMapIndex(reflect.ValueOf(keyStr), value)
					if keyStr != key.String() {
						v.SetMapIndex(key, reflect.Value{})
					}
				}
			}
	}
	case reflect.Ptr:
		if !v.IsNil() {
			escapeField(v.Elem())
		}
	}
}

func EscapeHTMLStruct(s interface{}) {
	value := reflect.ValueOf(s)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	escapeField(value)
}
