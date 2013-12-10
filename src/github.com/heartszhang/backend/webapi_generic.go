package backend

import (
	"log"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/exit.json", webapi_exit)
	http.HandleFunc("/", webapi_home)
}

// uri: /exit.json
func webapi_exit(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
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
