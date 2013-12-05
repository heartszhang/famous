package curl

import (
	"net/url"
)

type Curler interface {
	Get(uri string) (Cache, error)
	GetUtf8(uri string) (Cache, error)
	GetAsJson(uri string, val interface{}) error
	GetAsString(uri string) (string, error)
	PostForm(uri string, form url.Values) (int, error)
	PostFormAsString(uri string, form url.Values) (string, error)
	PostFormAsJson(uri string, form url.Values, val interface{}) error
}
