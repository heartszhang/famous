package feed

import "io"

// file has been converted to utf-8, so we just ignore internal encoding-declaration
func charset_reader_passthrough(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}
