package unixtime

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

type Time_s struct {
	Time UnixTime `xml:"time"`
}

func TestUnixTimeMarshal(t *testing.T) {
	xd := `<time_s><time>Thu, 28 Nov 2013 20:06:22 -0500</time></time_s>`
	jd := `{"time":"Thu, 28 Nov 2013 20:06:22 -0500"}`
	var xv, jv Time_s
	err := xml.Unmarshal([]byte(xd), &xv)
	t.Log(xv, err)
	err = json.Unmarshal([]byte(jd), &jv)
	t.Log(jv, err)
}
