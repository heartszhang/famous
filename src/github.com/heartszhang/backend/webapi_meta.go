package backend

import (
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/api/meta.json", webapi_meta)
	http.HandleFunc("/api/meta/cleanup.json", webapi_meta_cleanup)
}

func webapi_meta(w http.ResponseWriter, r *http.Request) {
	m, err := meta()

	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, m)
	}
}

func webapi_meta_cleanup(w http.ResponseWriter, r *http.Request) {
	err := meta_cleanup()
	webapi_write_error(w, err)
	log.Println(r.URL.RequestURI())
}
