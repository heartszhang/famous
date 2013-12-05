package curl

import (
	"github.com/qiniu/iconv"
	"io"
	"io/ioutil"
)

type utf8_readcloser struct {
	iconv.Iconv
	*iconv.Reader
}

func new_utf8_reader(charset string, in io.Reader) (io.ReadCloser, error) {
	if charset == "utf-8" || charset == "" {
		return ioutil.NopCloser(in), nil
	}
	converter, err := iconv.Open("utf-8", charset)
	if err != nil {
		return nil, err
	}
	ireader := iconv.NewReader(converter, in, 0)
	return &utf8_readcloser{converter, ireader}, nil
}

func (this *utf8_readcloser) Read(p []byte) (n int, err error) {
	return this.Reader.Read(p)
}

func (this *utf8_readcloser) Close() error {
	return this.Iconv.Close()
}
