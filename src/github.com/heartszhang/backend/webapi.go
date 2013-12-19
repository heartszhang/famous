package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func init() {
	http.HandleFunc("/api/tick.json", webapi_tick)
	http.HandleFunc("/api/update/popup.json", webapi_update_popup)
}

func webapi_assert(w http.ResponseWriter, r *http.Request) {
	webapi_write_error(w, backend_error{"not implemented yet", -1})
}

func webapi_tick(w http.ResponseWriter, r *http.Request) {
	switch v, err := tick(); err {
	case nil:
		webapi_write_error(w, err)
	default:
		webapi_write_as_json(w, v)
	}
}

const (
	_1M = 1024 * 1024
)

type counted_writer struct {
	io.Writer
	content_length int
}

func (this *counted_writer) Write(p []byte) (n int, err error) {
	n, e := this.Writer.Write(p)
	this.content_length += n
	return n, e
}

func webapi_write_as_json(w http.ResponseWriter, body interface{}) {
	cw := &counted_writer{w, 0}
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(cw)
	enc.Encode(body)
	w.Header().Set("content-length", strconv.Itoa(cw.content_length))
}

// sc may be ok, if err == nil, sc would be ignored
func webapi_write_error_code(w http.ResponseWriter, err error, statuscode int) {
	if statuscode == 0 || statuscode == http.StatusOK {
		statuscode = http.StatusBadGateway
	}
	switch err {
	case nil: // ignore statuscode
		webapi_write_as_json(w, backend_error{})
	default:
		cw := &counted_writer{w, 0}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(statuscode)
		enc := json.NewEncoder(cw)
		enc.Encode(backend_error{Reason: err.Error(), Code: statuscode})
		w.Header().Set("content-length", strconv.Itoa(cw.content_length))
	}
}

func webapi_write_error(w http.ResponseWriter, err error) {
	webapi_write_error_code(w, err, http.StatusBadGateway)
}

//http://address/api/image/thumbnail.jpeg?uri=
func redirect_thumbnail(uri string) string {
	return fmt.Sprintf("http://%v/api/image/thumbnail.jpeg?uri=%v", backend_context.config.Address(), url.QueryEscape(uri))
}

func redirect_origin(uri string) string {
	return fmt.Sprintf("http://%v/api/image/origin.jpeg?uri=%v", backend_context.config.Address(), url.QueryEscape(uri))
}

func imageurl_from_video(uri string) string {
	return fmt.Sprintf("http://%v/api/image/video.thumbnail.jpeg?uri=%v", backend_context.config.Address(), url.QueryEscape(uri))
}

//server/api/link/origin.json?uri=
func redirect_link(uri string) string {
	return fmt.Sprintf("http://%v/api/link/origin.json?uri=%v", backend_context.config.Address(), url.QueryEscape(uri))
}
