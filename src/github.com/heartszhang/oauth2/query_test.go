package oauth2

import (
	"testing"
)

func TestHttpQueryEncode(t *testing.T) {
	x := struct {
		a int
		b struct {
			c float64
		}
		d string
		e bool
		f *int
		g []string
	}{
		a: 1,
		d: "dstring",
		e: true,
		g: []string{"sg1", "sg2"},
	}
	t.Log(HttpQueryEncode(x))
}
