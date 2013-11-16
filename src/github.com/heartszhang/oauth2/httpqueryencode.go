package oauth2

import (
	"net/url"
	"reflect"
	"strconv"
)

func extract_struct_params(val reflect.Value, tag string) url.Values {
	v := make(url.Values)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		ftype := val.Type().Field(i)
		param := ftype.Tag.Get("param")
		v = merge(v, extract_value_params(field, param))
	}
	return v
}

func merge(l url.Values, r url.Values) url.Values {
	rlt := make(map[string][]string)
	for k, v := range l {
		for _, i := range v {
			rlt[k] = append(rlt[k], i)
		}
	}
	for k, v := range r {
		for _, i := range v {
			rlt[k] = append(rlt[k], i)
		}
	}
	return rlt
}

func extract_value_params(val reflect.Value, tag string) url.Values {
	v := make(url.Values)
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
			v.Add(tag, strconv.FormatFloat(val.Float(), 'f', -1, 64))
		}
	case reflect.String:
		s := val.String()
		if tag != "" && s != "" {
			v.Add(tag, s)
		}
	}
	return v
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

func HttpQueryValues(i interface{}) url.Values {
	return extract_value_params(reflect.ValueOf(i), "")
}
