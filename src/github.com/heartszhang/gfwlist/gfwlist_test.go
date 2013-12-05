package gfwlist

import (
	"os"
	"testing"
)

const (
	rule_file = `E:\sourcesafe\famous\src\github.com\heartszhang\gfwlist\gfw.txt`
)

func TestWhiteList(t *testing.T) {
	f, err := os.Open(rule_file)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	ruler, _ := NewGfwRuler(f)
	//	sites := []string{"http://www.sina.com.cn", "http://www.cafepress.com", "http://www.hulu.com", "https://youtube.com", "http://twitter.com"}
	sites := []string{"http://chinadigitaltimes.net/chinese"}
	for _, site := range sites {
		v := ruler.IsBlocked(site)
		t.Log(site, v)
	}
}
