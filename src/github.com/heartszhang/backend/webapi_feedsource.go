package backend

import (
	"log"
	"net/http"
	"net/url"
)

func init() {

	http.HandleFunc("/api/feed_source/all.json", webapi_feedsource_all)
	http.HandleFunc("/api/feed_source/subscribe.json", webapi_feedsource_subscribe)
	http.HandleFunc("/api/feed_source/unsubscribe.json", webapi_feedsource_unsubscribe)
	http.HandleFunc("/api/feed_source/find.json", webapi_feedsource_find) // ?q=
	http.HandleFunc("/api/feed_source/show.json", webapi_feedsource_find) // ?q=
}

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
	if u, err := url.Parse(q); u.IsAbs() && err == nil {
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

func webapi_feedsource_all(w http.ResponseWriter, r *http.Request) {
	switch fs, err := feedsource_all(); err {
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
