package videothumbnail

import (
	"fmt"
	"net/url"
)

func DescribeVideo(uri string) (v VideoDescription, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	domain := domain_from_host(u.Host)
	processor := processor_from_domain(domain)
	v, err = processor.process(u)
	return
}

type thumbnail_processor interface {
	process(uri *url.URL) (VideoDescription, error)
}
type processor struct{
	origin_uri string
}

type processor_dummy processor
type processor_56 processor
type processor_qq processor
type processor_youku processor

func (this processor_dummy) process(uri *url.URL) (VideoDescription, error) {
	return VideoDescription{}, vd_error{"unknown-processor", uri.String()}
}

func processor_from_domain(domain string) thumbnail_processor {
	switch domain {
	case "56.com":
		return processor_56{domain}
	case "youku.com":
		return processor_youku{domain}
	case "qq.com":
		return processor_qq{domain}
	case "tudou.com":
		fallthrough
	case "ku6.com":
		fallthrough
	case "letv.com":
		fallthrough
	case "sina.com.cn":
		fallthrough
	case "iqiyi.com":
		fallthrough
	case "sohu.com":
		fallthrough
	case "xunlei.com":
		fallthrough
	default:
		return processor_dummy{}
	}
}

type vd_error struct {
	reason string
	location   string
}

func (this vd_error) Error() string {
	return fmt.Sprintf("%v: %v", this.location, this.reason)
}
