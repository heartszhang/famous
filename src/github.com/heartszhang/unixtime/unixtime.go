package unixtime

import (
	"fmt"
	"strconv"
	"time"
)

type Time int64

func TimeNow() Time {
	return Time(time.Now().Unix())
}
func NewTime(sec int64) Time {
	return Time(sec)
}
func (this *Time) UnmarshalJSON(data []byte) error {
	v, err := unixtime_unmarshal(data, true)
	*this = Time(v)
	return err
}

func (this *Time) UnmarshalText(body []byte) error {
	v, err := unixtime_unmarshal(body, false)
	*this = Time(v)
	return err
}

const (
	RFC1123Za = `Mon, _2 Jan 2006 15:04:05 -0700`
)

func unixtime_unmarshal(data []byte, format_quoted bool) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	str := string(data)
	if str == "null" { // some times
		return 0, nil
	}
	v, err := strconv.ParseInt(str, 0, 0)
	if err == nil {
		return v, nil
	}
	formats := []string{
		time.RFC822Z,     // 02 Jan 06 15:04 -0700
		time.ANSIC,       // Mon Jan _2 15:04:05 2006
		time.RFC850,      // Monday, 02-Jan-06 15:04:05 MST
		time.UnixDate,    // Mon Jan _2 15:04:05 MST 2006
		time.RFC3339,     // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano, // 2006-01-02T15:04:05Z07:00
		time.RFC822,      // 02 Jan 06 15:04 MST
		time.RFC1123,     // Mon, 02 Jan 2006 15:04:05 MST
		time.RFC1123Z,    // Mon, 02 Jan 2006 15:04:05 -0700
		time.RubyDate,    // Mon Jan 02 15:04:05 -0700 2006
		RFC1123Za,        // like rfc1123z
	}
	for _, format := range formats {
		if format_quoted {
			format = `"` + format + `"`
		}
		t, err := time.Parse(format, str)
		if err == nil {
			return t.Unix(), nil
		}
	}
	return v, fmt.Errorf(str + ": unrecognized timeformat")
}
