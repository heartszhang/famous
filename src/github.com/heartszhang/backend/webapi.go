package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func init() {
	http.HandleFunc("/api/tick.json", webapi_tick)
	http.HandleFunc("/api/update/popup.json", webapi_meta) // no param
	http.HandleFunc("/api/meta.json", webapi_meta)
	http.HandleFunc("/api/meta/cleanup.json", webapi_meta_cleanup)

	http.HandleFunc("/api/feed_tag/all.json", webapi_feedtag_all)

	http.HandleFunc("/api/link/origin.json", webapi_link_origin)          // ?uri=
	http.HandleFunc("/api/suggest/bing.json", webapi_suggest_bing)        // ?q=
}

const uint64_bits int = 64

func webapi_feedtag_all(w http.ResponseWriter, r *http.Request) {
	switch ft, err := feedtag_all(); err {
	case nil:
		webapi_write_as_json(w, ft)
	default:
		webapi_write_error(w, err)
	}
	log.Println(r.URL.RequestURI())
}

func webapi_suggest_bing(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	switch sg, err := suggest_bing(q); err {
	case nil:
		webapi_write_as_json(w, sg)
	default:
		webapi_write_error(w, err)
	}
}


func webapi_assert(w http.ResponseWriter, r *http.Request) {
	webapi_write_error(w, backend_error{"not implemented yet", -1})
}
// uri: /meta.json
func webapi_meta(w http.ResponseWriter, r *http.Request) {
	m, err := meta()

	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, m)
	}
	log.Println(r.URL.RequestURI())
}

// uri: /tick.json
func webapi_tick(w http.ResponseWriter, r *http.Request) {
	v, err := tick()
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
	log.Println(r.URL.RequestURI())
}

// uri: /meta/cleanup.json
func webapi_meta_cleanup(w http.ResponseWriter, r *http.Request) {
	err := meta_cleanup()
	webapi_write_error(w, err)
	log.Println(r.URL.RequestURI())
}

func webapi_link_origin(w http.ResponseWriter, r *http.Request) {
	webapi_write_error(w, nil)
	log.Println(r.URL.RequestURI())
}

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

//http://address/api/image/thumbnail.json?uri=
func redirect_thumbnail(uri string) string {
	return fmt.Sprintf("http://%v/api/image/thumbnail.json?uri=%v", config.Address(), url.QueryEscape(uri))
}

func imageurl_from_video(uri string) string {
	return fmt.Sprintf("http://%v/api/image/video.thumbnail?uri=%v", config.Address(), url.QueryEscape(uri))
}

//server/api/link/origin.json?uri=
func redirect_link(uri string) string {
	return fmt.Sprintf("http://%v/api/link/origin.json?uri=%v", config.Address(), url.QueryEscape(uri))
}
