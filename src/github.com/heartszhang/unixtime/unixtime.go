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
		time.RFC822Z,
		time.ANSIC,
		time.RFC850,
		time.UnixDate,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC1123,
		time.RFC1123Z,
		time.RubyDate,
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
	return v, fmt.Errorf("unrecognized timeformat")
}
