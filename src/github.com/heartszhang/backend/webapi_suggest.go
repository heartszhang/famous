package backend

import "net/http"

func init() {
	http.HandleFunc("/api/suggest/bing.json", webapi_suggest_bing) // ?q=
}

// q: string, required
func webapi_suggest_bing(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	switch sg, err := suggest_bing(q); err {
	case nil:
		webapi_write_as_json(w, sg)
	default:
		webapi_write_error(w, err)
	}
}
