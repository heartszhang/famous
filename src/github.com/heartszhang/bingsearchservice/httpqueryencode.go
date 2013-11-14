package bingsearchservice

import (
	"net/url"
	"reflect"
	"strconv"
)

func extract_struct_params(val reflect.Value, tag string) (v url.Values) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		ftype := val.Type().Field(i)
		param := ftype.Tag.Get("param")
		v = merge(v, extract_value_params(field, param))
	}
	return
}

func merge(l url.Values, r url.Values) (rlt url.Values) {
	for k, v := range l {
		for _, i := range v {
			rlt.Add(k, i)
		}
	}
	for k, v := range r {
		for _, i := range v {
			rlt.Add(k, i)
		}
	}
	return
}

func extract_value_params(val reflect.Value, tag string) (v url.Values) {
	switch val.Kind() {
	case reflect.Struct:
		v = merge(v, extract_struct_params(val, tag))
	case reflect.Interface, reflect.Ptr:
		if !val.IsNil() {
			v = merge(v, extract_value_params(val.Elem(), tag))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if tag != "" {
			v.Add(tag, strconv.FormatInt(val.Int(), 10))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if tag != "" {
			v.Add(tag, strconv.FormatUint(val.Uint(), 10))
		}
	case reflect.Float32, reflect.Float64:
		if tag != "" {
			v.Add(tag, strconv.FormatFloat(val.Float(), 0, 0, 0))
		}
	case reflect.String:
		s := val.String()
		if tag != "" && s != "" {
			v.Add(tag, s)
		}
	}
	return
}

/*
//interface, ptr
// struct
// simple
*/
func HttpQueryEncode(i interface{}) string {
	val := reflect.ValueOf(i)
	v := extract_value_params(val, "")
	return v.Encode()
}
