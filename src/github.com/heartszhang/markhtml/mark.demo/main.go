package main

import (
	"fmt"
	"github.com/heartszhang/markhtml"
)

const (
	x = `[align=center][img]http://107.imagebam.com/download/2XaHuuG5h7N3JUlJu2Q4oA/28490/284899367/%3F%3F8%3F%3F%3F%3F-%3F%3F[/img][/align] [align=center][color=red][b][size=6]世界杯附[/size][/b][/color][/align] [table=100%][tr][td][size=4][color=purple][b][align=center]比赛时间 [/align][/b][/color][/siz ..`
)

func main() {
	fmt.Println(markhtml.TransferText(x))
}
