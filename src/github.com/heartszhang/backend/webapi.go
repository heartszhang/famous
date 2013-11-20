package backend

import (
	"encoding/json"
	"fmt"
	feed "github.com/heartszhang/feedfeed"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func init() {
	http.HandleFunc("/api/tick.json", webapi_tick)
	http.HandleFunc("/api/meta.json", webapi_meta)
	http.HandleFunc("/api/meta/cleanup.json", webapi_meta_cleanup)
	http.HandleFunc("/api/feed_category/all.json", webapi_feedcategory_all)
	http.HandleFunc("/api/feed_category/create.json", webapi_feedcategory_create)
	http.HandleFunc("/api/feed_category/drop.json", webapi_feedcategory_drop)
	http.HandleFunc("/api/feed_tag/all.json", webapi_feedtag_all)
	http.HandleFunc("/api/feed_source/all.json", webapi_feedsource_all)
	http.HandleFunc("/api/feed_source/subscribe.json", webapi_feedsource_subscribe)
	http.HandleFunc("/api/feed_source/unsubscribe.json", webapi_feedsource_unsubscribe)
	//	http.HandleFunc("/api/feed_source/entries_since.json", webapi_feedsource_entries_since)
	//	http.HandleFunc("/api/feed_source/entries_since.json", webapi_feedentries_since)

	http.HandleFunc("/api/feed_entry/unread.json", webapi_feedentry_unread)
	http.HandleFunc("/api/feed_entry/mark.json", webapi_feedentry_mark)
	http.HandleFunc("/api/feed_entry/umark.json", webapi_feedentry_umark)
	http.HandleFunc("/api/feed_entry/full_text.json", webapi_feedentry_fulltext)
	http.HandleFunc("/api/feed_entry/media.json", webapi_feedentry_media)
	http.HandleFunc("/api/feed_entry/drop.json", webapi_feedentry_drop)
	//	http.HandleFunc("/api/meta/categories.json", webapi_meta_categories)

	http.HandleFunc("/api/image/description.json", webapi_image_description)
	http.HandleFunc("/api/image/thumbnail.json", webapi_image_thumbnail)       // ?uri= return image/jpeg
	http.HandleFunc("/api/image/origin.json", webapi_image_origin)             // ?uri=, return image/xxx
	http.HandleFunc("/api/image/dimension.json", webapi_image_dimension)       // ?uri=, return FeedMedia
	http.HandleFunc("/api/image/video.thumbnail", webapi_image_videothumbnail) //?uri=, return image/xxx

	http.HandleFunc("/api/link/origin.json", webapi_link_origin)          // ?uri=
	http.HandleFunc("/api/suggest/bing.json", webapi_suggest_bing)        // ?q=
	http.HandleFunc("/api/feed_source/find.json", webapi_feedsource_find) // ?q=
	http.HandleFunc("/api/feed_source/show.json", webapi_feedsource_find) // ?q=
	http.HandleFunc("/exit.json", webapi_exit)
	http.HandleFunc("/", webapi_home)
}

const uint64_bits int = 64

func webapi_feedsource_show(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	switch fs, err := feedsource_show(q); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}
func webapi_feedsource_find(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if _, err := url.Parse(q); err == nil {
		webapi_feedsource_show(w, r)
		return
	}
	switch fs, err := feedsource_find(q); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}

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

// uri : /api/feed_category/all.json
func webapi_feedcategory_all(w http.ResponseWriter, r *http.Request) {
	switch fc, err := feedcategory_all(); err {
	case nil:
		webapi_write_as_json(w, fc)
	default:
		webapi_write_error(w, err)
	}
	log.Println(r.URL.RequestURI())
}

// uri: /api/feed_entry/unread.json/{uri}/{count}/{page}
func webapi_feedentry_unread(w http.ResponseWriter, r *http.Request) {
	//	category, err := strconv.ParseUint(r.URL.Query().Get("category"), 0, uint64_bits)
	uri := r.URL.Query().Get("uri")
	log.Println("feedentry-unread", uri)
	count, _ := strconv.ParseUint(r.URL.Query().Get("count"), 0, 0)
	page, _ := strconv.ParseUint(r.URL.Query().Get("page"), 0, 0)
	switch fe, err, sc := feedentry_unread(uri, uint(count), uint(page)); err {
	case nil:
		webapi_write_as_json(w, fe)
	default:
		webapi_write_error_code(w, err, sc)
	}
}

// uri: /feed_source/subscribe.json?uri=
func webapi_feedsource_subscribe(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("feedsource-sub", uri)
	source_type := source_type_map(r.URL.Query().Get("source_type"))
	switch fs, err := feedsource_subscribe(uri, source_type); err {
	case nil:
		log.Println(fs, err)
		webapi_write_as_json(w, fs)
	default:
		log.Println(err)
		webapi_write_error(w, err)
	}
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

// uri: /feed_entry/mark.json/{entry_id}/{flags}
func webapi_feedentry_mark(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("entry_uri")
	log.Println("feedentry-mark", uri)
	f, err := strconv.ParseUint(r.URL.Query().Get("flags"), 0, 0)
	flag := uint(f)
	if err == nil {
		flag, err = feedentry_mark(uri, flag)
	}
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, struct{}{})
	}
}

// uri: /feed_entry/umark.json/{entry_id}/{flags}
func webapi_feedentry_umark(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("entry_uri")
	log.Println("feedentry-umark", uri)
	f, err := strconv.ParseInt(r.URL.Query().Get("flags"), 0, 0)
	flag := uint(f)
	if err == nil {
		flag, err = feedentry_umark(uri, flag)
	}
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, struct{}{})
	}
}

// uri: /feed_entry/full_text.json/{entry_uri}
func webapi_feedentry_fulltext(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	entry_uri := r.URL.Query().Get("entry_uri")
	log.Println("feed-entry-fullt", entry_uri)
	ff, err := feedentry_fulltext(uri, entry_uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, ff)
	}
}

// uri: /feed_entry/media.json/{entry_id}/{url}/{media_type:[0-9]+}
func webapi_feedentry_media(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	id := r.URL.Query().Get("entry_id")
	log.Println("feedentry-media", uri, id)
	media_type, err := strconv.ParseUint(r.URL.Query().Get("media_type"), 0, 0)
	if err != nil {
		media_type = uint64(feed.Feed_media_type_unknown)
	}
	v, err := feedentry_media(uri, id, uint(media_type))
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /feed_entry/drop.json/{entry_id}
// id is mongo's _id
func webapi_feedentry_drop(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("entry_id")
	log.Println("feedentry-drop", id)
	err := feedentry_drop(id)
	webapi_write_error(w, err)
}

// uri: /feed_category/create.json/{name}
func webapi_feedcategory_create(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	log.Println("feed-cat", name)
	v, err := feedcategory_create(name)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /feed_catetory/drop.json/{name}
func webapi_feedcategory_drop(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	log.Println("feed-cat", name)
	err := feedcategory_drop(name)
	webapi_write_error(w, err)
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

func webapi_feedsource_all(w http.ResponseWriter, r *http.Request) {
	fso := new_feedsource_operator()
	switch fs, err := fso.all(); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
	log.Println(r.URL.RequestURI())
}

// uri: /feed_source/unsubscribe.json/{uri}
func webapi_feedsource_unsubscribe(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("feed-unsub", uri)
	err := feedsource_unsubscribe(uri)
	webapi_write_error(w, err)
}

// uri: /meta/cleanup.json
func webapi_meta_cleanup(w http.ResponseWriter, r *http.Request) {
	err := meta_cleanup()
	webapi_write_error(w, err)
	log.Println(r.URL.RequestURI())
}

// uri: /image/description.json?uri=
func webapi_image_description(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("img-desc", uri)
	v, err := image_description(uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /image/dimension.json?uri=
func webapi_image_dimension(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("img-dim", uri)
	v, err := image_dimension(uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}
func webapi_image_videothumbnail(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("img-vthumb", uri)
	switch vt, err := image_videothumbnail(uri); err {
	case nil:
		webapi_image_entity(vt.Image, w, r)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_image_entity(uri string, w http.ResponseWriter, r *http.Request) {
	log.Println("img-entry", uri)
	switch cache, err := image_description(uri); err {
	case nil:
		w.Header().Set("content-type", cache.Mime)
		w.Header().Set("x-image-meta-property-height", strconv.Itoa(cache.Height))
		w.Header().Set("x-image-meta-property-width", strconv.Itoa(cache.Width))
		w.Header().Set("x-image-meta-property-alter", cache.OriginLocal)

		http.ServeFile(w, r, cache.ThumbnailLocal)
	default:
		webapi_write_error_code(w, err, cache.Code)
	}
}

// /api/image/thumbnail.json?uri=
func webapi_image_thumbnail(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("uri")
	webapi_image_entity(url, w, r)
}

// /api/image/origin.json?uri=
func webapi_image_origin(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("img-origin", uri)
	switch cache, err := image_description(uri); err {
	case nil:
		w.Header().Set("content-type", cache.Mime)
		w.Header().Set("x-image-meta-property-height", strconv.Itoa(cache.Height))
		w.Header().Set("x-image-meta-property-width", strconv.Itoa(cache.Width))
		w.Header().Set("x-image-meta-property-alter", cache.ThumbnailLocal)
		http.ServeFile(w, r, cache.OriginLocal)
	default:
		webapi_write_error_code(w, err, cache.Code)
	}
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
		webapi_write_as_json(w, BackendError{})
	default:
		cw := &counted_writer{w, 0}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(statuscode)
		enc := json.NewEncoder(cw)
		enc.Encode(BackendError{Reason: err.Error(), Code: statuscode})
		w.Header().Set("content-length", strconv.Itoa(cw.content_length))
	}
}

func webapi_write_error(w http.ResponseWriter, err error) {
	webapi_write_error_code(w, err, http.StatusBadGateway)
}

// uri: /exit.json
func webapi_exit(w http.ResponseWriter, r *http.Request) {
	webapi_write_as_json(w, r.URL)

	f, canflush := w.(http.Flusher)
	if canflush {
		f.Flush()
	}

	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Fatalf("error while shutting down: %v", err)
	}

	conn.Close()

	log.Println("Shutting down")
	os.Exit(0)
}

// uri: /
func webapi_home(w http.ResponseWriter, r *http.Request) {
	webapi_write_as_json(w, r.URL)
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
