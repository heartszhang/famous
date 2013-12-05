package main

import (
	"flag"
	"fmt"
	"github.com/heartszhang/gfwlist"
	"os"
)

const (
	input = `E:\sourcesafe\famous\src\github.com\heartszhang\gfwlist\gfw.txt`
)

var uri = flag.String("uri", "", "url with scheme")

func main() {
	flag.Parse()
	f, err := os.Open(input)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	ruler, _ := gfwlist.NewGfwRuler(f)
	if len(*uri) != 0 {
		fmt.Println(ruler.IsBlocked(*uri))
	}
}
