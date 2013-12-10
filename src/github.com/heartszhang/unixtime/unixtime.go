package unixtime

import (
	"fmt"
	"strconv"
	"time"
)

const (
	_1minute = 1
	_1hour   = _1minute * 60
	_2hours  = _1hour * 2
	_1day    = _1hour * 24
	_1week   = _1day * 7
	_1month  = _1day * 30
	_1year   = _1day * 365
)

type UnixTime int64

func UnixTimeNow() UnixTime {
	return UnixTime(time.Now().Unix())
}
func NewUnixTime(sec int64) UnixTime {
	return UnixTime(sec)
}
func (this *UnixTime) UnmarshalJSON(data []byte) error {
	v, err := unixtime_unmarshal(data, true)
	*this = UnixTime(v)
	return err
}

func (this *UnixTime) UnmarshalText(body []byte) error {
	v, err := unixtime_unmarshal(body, false)
	*this = UnixTime(v)
	return err
}

func unixtime_unmarshal(data []byte, format_quoted bool) (int64, error) {
	str := string(data)
	v, err := strconv.ParseInt(str, 0, 0)
	if err == nil {
		return v, nil
	}
	formats := []string{
		time.RFC822Z,                      // 02 Jan 06 15:04 -0700
		time.ANSIC,                        // Mon Jan _2 15:04:05 2006
		time.RFC850,                       // Monday, 02-Jan-06 15:04:05 MST
		time.UnixDate,                     // Mon Jan _2 15:04:05 MST 2006
		time.RFC3339,                      // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano,                  // 2006-01-02T15:04:05Z07:00
		time.RFC822,                       // 02 Jan 06 15:04 MST
		time.RFC1123,                      // Mon, 02 Jan 2006 15:04:05 MST
		time.RFC1123Z,                     // Mon, 02 Jan 2006 15:04:05 -0700
		time.RubyDate,                     // Mon Jan 02 15:04:05 -0700 2006
		`Mon, _2 Jan 2006 15:04:05 -0700`, // like rfc1123z
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
