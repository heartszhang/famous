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
	xd := `<time_s><time>Tue, 25 Nov 2008 10:47:32 +0800</time></time_s>`
	jd := `{"time":"Mon, 6 May 2013 15:23:32 +0800"}`
	var xv, jv Time_s
	err := xml.Unmarshal([]byte(xd), &xv)
	t.Log(xv, err)
	err = json.Unmarshal([]byte(jd), &jv)
	t.Log(jv, err)
}
