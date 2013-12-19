package backend

import "net/http"

func init() {
	http.HandleFunc("/api/feed_tag/all.json", webapi_feedtag_all)
	http.HandleFunc("/api/feed_tag/drop.json", webapi_feedtag_all)
	http.HandleFunc("/api/feed_tag/rename.json", webapi_feedtag_all)
}

func webapi_feedtag_all(w http.ResponseWriter, r *http.Request) {
	switch ft, err := feedtag_all(); err {
	case nil:
		webapi_write_as_json(w, ft)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedtag_rename(w http.ResponseWriter, r *http.Request) {
	switch ft, err := feedtag_all(); err {
	case nil:
		webapi_write_as_json(w, ft)
	default:
		webapi_write_error(w, err)
	}
}
