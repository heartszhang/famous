package main

import (
	"fmt"
	//	bing "github.com/heartszhang/bingsearchservice"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/oauth2"
	//"github.com/whee/ddg"
	//	"net/http"
)

const (
	AccountKey = "CUJq/zmtCVYt9qiM7XnU+eOlTjohFJ4x5jwRoO3gJSU"
)

const (
	bing_suggest = "http://api.bing.com/qsonhs.aspx"
	//	"http://suggestion.baidu.com/su?wd=cjk&p=3&cb=window.bdsug.sug&sid=&t=1384484295939"
	baidu_suggest = "http://suggestion.baidu.com/su" // use gbk and rfc1521 incompitable content-type declaration
	//	"http://sug.so.360.cn/suggest/word?callback=suggest_so&encodein=utf-8&encodeout=utf-8&word=caoliu#"
	so360_suggest = "http://sug.so.360.cn/suggest/word"
)

func So360Suggestion(q string) ([]string, error) {
	p := struct {
		callback  string `param:"callback"`
		encodein  string `param:"encodein"`
		encodeout string `param:"encodeout"`
		word      string `param:"word"`
	}{word: q, encodein: "utf-8", encodeout: "utf-8"}
	qu := oauth2.HttpQueryEncode(p)
	uri := so360_suggest + "?" + qu
	fmt.Println(uri)
	/*	val := struct {
			Q string   `json:"q"`
			P bool     `json:"p"`
			S []string `json:"s"`
		}{}
	*/
	_, err := curl.NewCurl("", curl.CurlProxyPolicyNoProxy, 0).GetAsString(uri)
	return nil, err
}

func BaiduSuggestion(q string) ([]string, error) {
	p := struct {
		query    string `param:"wd"`
		callback string `param:"cb"`
		sid      string `param:"sid"`
		t        string `param:"t"`
		p        *int   `param:"p"`
	}{query: q}
	qu := oauth2.HttpQueryEncode(p)
	uri := baidu_suggest + "?" + qu
	fmt.Println(uri)
	/*	val := struct {
			Q string   `json:"q"`
			P bool     `json:"p"`
			S []string `json:"s"`
		}{}
	*/
	_, err := curl.NewCurl("", curl.CurlProxyPolicyNoProxy, 0).GetAsString(uri)
	x := make([]string, 0)
	return x, err
}

func main() {
	/*
		sina := oauth2.OAuthConfig{
			ClientId:     "812692320",
			ClientSecret: "9b125acc087a3e372a7028be5bac053a",
			Scope:        "",
			AuthUrl:      "https://api.weibo.com/oauth2/authorize",
			TokenUrl:     "https://api.weibo.com/oauth2/access_token",
			RedirectUrl:  "http://iweizhi2.duapp.com/authorize",
		}
	*/
	/*
		azure := oauth2.OAuthConfig{
			ClientId:     "",
			ClientSecret: "",
			Scope:        "https://api.datamarket.azure.com/",
			AuthUrl:      "https://datamarket.azure.com/embedded/consent",
			TokenUrl:     "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13",
			RedirectUrl:  "http://iweizhi2.duapp.com/authorize",
		}
		broker := oauth2.NewWebAuthorizationBroker(sina, nil)
		token, err := broker.Authorize(true, nil)
		fmt.Println(token, err)
	*/
	t, err := So360Suggestion("中国")
	fmt.Println(t, err)
	/*
		x := bing.NewBingSearchService("", AccountKey)
		c, err := x.SpellingSuggestion(bing.BingSearchParameters{Query: `'china'`})
		fmt.Println(c, err)
		v, err := x.RelatedSearch(bing.BingSearchParameters{Query: `'english'`})
		fmt.Println(v, err)
	*/
	//r, err := ddg.ZeroClick("caoliu")
	//	fmt.Println(r, err)
}
