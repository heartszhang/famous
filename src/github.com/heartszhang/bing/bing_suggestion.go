package bing

import (
	"github.com/heartszhang/curl"
	"github.com/heartszhang/oauth2"
)

type bing_suggestion_result struct {
	AS struct {
		Query       string
		FullResults int
		Results     []struct {
			Type     string
			Suggests []struct {
				Txt  string
				Type string
				Sk   string
				HCS  float64
			}
		}
	}
}

func export_bing_suggestions(result bing_suggestion_result) []string {
	v := make([]string, 0)
	for _, r := range result.AS.Results {
		for _, s := range r.Suggests {
			v = append(v, s.Txt)
		}
	}
	return v
}

const (
	bing_suggest = "http://api.bing.com/qsonhs.aspx"
)

func BingSuggestion(q string) ([]string, error) {
	p := struct {
		market string `param:"mkt"` //zh-CN
		cp     int    `param:"cp"`  //2
		o      string `param:"o"`   //a+ds+ds+p
		query  string `param:"q"`
	}{market: "zh-CN", cp: 2, query: q}

	qu := oauth2.HttpQueryEncode(p)
	uri := bing_suggest + "?" + qu
	var val bing_suggestion_result
	err := curl.NewCurl("").GetAsJson(uri, &val)
	x := export_bing_suggestions(val)
	return x, err
}
