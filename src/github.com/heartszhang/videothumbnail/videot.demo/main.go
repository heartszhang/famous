package main

import (
	"flag"
	"fmt"
	vt "github.com/heartszhang/videothumbnail"
)

var (
	uri = flag.String("uri", "http://player.56.com/v_MTAwNzI4NDMx.swf/1030_justin0842.swf", "video player address for 56")
)

func main() {
	flag.Parse()
	v, e := vt.DescribeVideo(*uri)
	fmt.Println(v, e)
}
