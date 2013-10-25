package main

import (
	"fmt"
	"github.com/heartszhang/backend"
)

func main() {
	x, err := backend.Subscribe("http://fullrss.net/a/http/184.154.128.245/rss.php?fid=2")
	fmt.Println(x, err)
}
