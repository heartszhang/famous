package oauth2

import (
	"net/url"
	"reflect"
	"strconv"
)

/*
struct{
	a int
	b string `param:"x.b"`
	c []float64
	d struct{
		e bool   // will be encoded to 'e=false' not 'd.e=false'
	}
	f *int
}
nil pointer will be omitted
all sub-struct are be treated as inline
*/
func HttpQueryEncode(i interface{}) string {
	return HttpQueryValues(i).Encode()
}

func HttpQueryValues(i interface{}) url.Values {
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		panic("only struct value is accepted")
	}
	values := plain_extract_struct(val)
	return values
}

func plain_extract_value(name string, field reflect.Value) url.Values {
	values := make(url.Values)
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		values.Add(name, strconv.FormatInt(field.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		values.Add(name, strconv.FormatUint(field.Uint(), 10))
	case reflect.Bool:
		values.Add(name, strconv.FormatBool(field.Bool()))
	case reflect.Struct:
		values = merge(values, plain_extract_struct(field))
	case reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			sv := field.Index(i)
			values = merge(values, plain_extract_value(name, sv))
		}
	case reflect.String:
		values.Add(name, field.String())
	case reflect.Float32, reflect.Float64:
		values.Add(name, strconv.FormatFloat(field.Float(), 'f', -1, 64))
	case reflect.Ptr, reflect.Interface:
		values = merge(values, plain_extract_pointer(field))
	}
	return values
}

func plain_extract_struct(val reflect.Value) url.Values {
	values := make(url.Values)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		ftype := val.Type().Field(i)
		name := val.Type().Field(i).Name
		p := ftype.Tag.Get("param")
		if p != "" {
			name = p
		}
		values = merge(values, plain_extract_value(name, field))
	}
	return values
}

func plain_extract_pointer(ptr reflect.Value) url.Values {
	values := make(url.Values)
	if ptr.IsNil() {
		return values
	}
	val := ptr.Elem()
	if val.Kind() == reflect.Struct {
		values = merge(values, plain_extract_struct(val))
	}
	return values
}

func merge(l url.Values, r url.Values) url.Values {
	rlt := make(url.Values)
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
