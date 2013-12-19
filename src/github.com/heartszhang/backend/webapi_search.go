package backend

import "net/http"

func init() {
	http.HandleFunc("/api/search/feed_source/find.json", webapi_search_feedsource_find) // ?q=
	http.HandleFunc("/api/search/feed_source/show.json", webapi_search_feedsource_show) // ?q=
}

// q: string, required, rss or atom uri
func webapi_search_feedsource_show(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	switch fs, err := feedsource_show(q); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}

// q: string, required
func webapi_search_feedsource_find(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	switch fs, err := feedsource_find(q); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}
