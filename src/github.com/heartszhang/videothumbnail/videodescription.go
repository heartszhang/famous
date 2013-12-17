package videothumbnail

import ()

type VideoDescription struct {
	Image     string   // what's the difference between image and thumbnail?
	Thumbnail string
	Title     string
	Tags      []string  // or tags?
	Seconds  int  // seconds
}
