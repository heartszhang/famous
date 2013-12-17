package videothumbnail

import (
	"strings"
	"strconv"
)

// the last two node of any host
// sina.com => sina.com
// www.sina.com => sina.com
// www.sina.com.cn => sina.com.cn
// iweizhi2.duapp.com => duapp.com
func domain_from_host(host string) string {
	secs := strings.Split(host, ".")
	secc := len(secs)
	if secc < 3 {
		return host
	}
	return strings.Join(secs[1:], ".")
}

func atoi(x string) int {
	v, _ := strconv.Atoi(x)
	return v
}

func jsjson_foreach_field(s string, f func(n, v string)) {
	fields := strings.Split(s, "\n")
	for _, field := range fields {
		nvs := strings.FieldsFunc(field, func(r rune) bool {
			return strings.ContainsRune(`:," `, r)
		})
		if len(nvs) == 2 {  // `, name = value` or `name = value,` or `name = value`
			name := nvs[0]
			val := nvs[1]
			f(name, val)
		}
	}
}