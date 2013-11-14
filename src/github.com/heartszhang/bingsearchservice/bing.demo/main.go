package main

import (
	"fmt"
	//	b "github.com/heartszhang/bingsearchservice"
	//	"github.com/heartszhang/curl"
	"github.com/whee/ddg"
	//	"net/http"
)

const (
	AccountKey = "CUJq/zmtCVYt9qiM7XnU+eOlTjohFJ4x5jwRoO3gJSU"
)

func main() {
	//	x := NewBingSearchService("", AccountKey)
	//	c, err := x.Web(BingSearchWebParameters{BingSearchParameters: BingSearchParameters{Query: `'caoliu'`}})
	//	fmt.Println(c, err)
	r, err := ddg.ZeroClick("caoliu")
	fmt.Println(r, err)
}
