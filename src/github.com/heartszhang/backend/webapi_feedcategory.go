package backend

import (
	"net/http"
	"log"
)

func init() {

	http.HandleFunc("/api/feed_category/all.json", webapi_feedcategory_all)
	http.HandleFunc("/api/feed_category/create.json", webapi_feedcategory_create)
	http.HandleFunc("/api/feed_category/drop.json", webapi_feedcategory_drop)
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