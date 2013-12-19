package backend

import "net/http"

func init() {
	http.HandleFunc("/api/link/origin.json", webapi_link_origin)   // ?uri=
	http.HandleFunc("/api/link/content.json", webapi_link_content) // ?uri=
}

// uri: string, required
func webapi_link_origin(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	http.Redirect(w, r, uri, http.StatusOK)
}

// uri: string, required
func webapi_link_content(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	http.Redirect(w, r, uri, http.StatusOK)
}
