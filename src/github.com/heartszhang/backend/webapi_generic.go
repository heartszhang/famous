package backend

import (
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/exit.json", webapi_exit)
	http.HandleFunc("/", webapi_home)
}

func webapi_exit(w http.ResponseWriter, r *http.Request) {
	webapi_write_as_json(w, r.URL)

	f, canflush := w.(http.Flusher)
	if canflush {
		f.Flush()
	}

	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
	}

	conn.Close()

	os.Exit(0)
}

func webapi_home(w http.ResponseWriter, r *http.Request) {
	webapi_write_as_json(w, r.URL)
}
