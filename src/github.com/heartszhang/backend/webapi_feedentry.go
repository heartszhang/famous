package backend

import (
	feed "github.com/heartszhang/feedfeed"
	"net/http"
	"log"
	"strconv"
)

func init(){


	http.HandleFunc("/api/feed_entry/unread.json", webapi_feedentry_unread)
	http.HandleFunc("/api/feed_entry/mark.json", webapi_feedentry_mark)
	http.HandleFunc("/api/feed_entry/umark.json", webapi_feedentry_umark)
	http.HandleFunc("/api/feed_entry/full_text.json", webapi_feedentry_fulltext)
	http.HandleFunc("/api/feed_entry/media.json", webapi_feedentry_media)
	http.HandleFunc("/api/feed_entry/drop.json", webapi_feedentry_drop)
}

// uri: /api/feed_entry/unread.json?uri=&count=&page=
func webapi_feedentry_unread(w http.ResponseWriter, r *http.Request) {
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
