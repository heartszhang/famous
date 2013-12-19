package backend

import "net/http"

func init() {
	http.HandleFunc("/api/feed_source/all.json", webapi_feedsource_all)
	http.HandleFunc("/api/feed_source/subscribe.json", webapi_feedsource_subscribe)
	http.HandleFunc("/api/feed_source/unsubscribe.json", webapi_feedsource_unsubscribe)
	http.HandleFunc("/api/feed_source/by_category.json", webapi_feedsource_by_category)
	http.handleFunc("/api/feed_source/addto_category.json", webapi_feedsource_addto_category)
	http.handleFunc("/api/feed_source/removefrom_category.json", webapi_feedsource_removefrom_category)
	http.HandleFunc("/api/feed_source/by_tag.json", webapi_feedsource_by_tag)
	http.handleFunc("/api/feed_source/add_tag.json", webapi_feedsource_add_tag)
	http.handleFunc("/api/feed_source/remove_tag.json", webapi_feedsource_remove_tag)
	http.handleFunc("/api/feed_source/count.json", webapi_feedsource_add_tag)
	http.handleFunc("/api/feed_source/all/count.json", webapi_feedsource_add_tag)
}

// uri: /feed_source/subscribe.json?uri=
func webapi_feedsource_subscribe(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	source_type := source_type_map(r.URL.Query().Get("source_type"))
	switch fs, err := feedsource_subscribe(uri, source_type); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
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
}

// uri: /feed_source/unsubscribe.json/{uri}
func webapi_feedsource_unsubscribe(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	err := feedsource_unsubscribe(uri)
	webapi_write_error(w, err)
}

func webapi_feedsource_by_category(w http.ResponseWriter, r *http.Request) {
	switch fs, err := feedsource_all(); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedsource_addto_category(w http.ResponseWriter, r *http.Request) {
	switch fs, err := feedsource_all(); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedsource_removefrom_category(w http.ResponseWriter, r *http.Request) {
	switch fs, err := feedsource_all(); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedsource_by_tag(w http.ResponseWriter, r *http.Request) {
	switch fs, err := feedsource_all(); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedsource_add_tag(w http.ResponseWriter, r *http.Request) {
	switch fs, err := feedsource_all(); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedsource_remove_tag(w http.ResponseWriter, r *http.Request) {
	switch fs, err := feedsource_all(); err {
	case nil:
		webapi_write_as_json(w, fs)
	default:
		webapi_write_error(w, err)
	}
}
