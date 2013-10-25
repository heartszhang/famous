package backend

import (
	"net/http"
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

//rfc850 or time.ANSI
// Sunday, 06-Nov-94 08:49:37 GMT ; RFC 850
// Sun Nov  6 08:49:37 1994; ANSI C
func unixtime_nano_rfc850(text string) int64 {
	t, _ := http.ParseTime(text)
	return t.UnixNano()
}

const time_format_rfc822 = "Mon, 2 Jan 2006 15:04:05 -0700"

// Sun, 06 Nov 1994 08:49:37 GMT  ; RFC 822, updated by RFC 1123
// 1998-05-12T14:15Z/16:00Z // iso8601
func unixtime_nano_rfc822(t string, formats ...string) int64 {
	x, err := time.Parse(time_format_rfc822, t)
	if err == nil {
		return x.UnixNano()
	}
	for _, f := range formats {
		x, err := time.Parse(f, t)
		if err == nil {
			return x.UnixNano()
		}
	}
	return time.Now().UnixNano()
}

func unixtime_now() int64 {
	return time.Now().UnixNano()
}

func unixtime_nano(nano int64) time.Time {
	sec := time.Duration(nano) / time.Second
	n := time.Duration(nano) % time.Second
	return time.Unix(int64(sec), int64(n))
}
