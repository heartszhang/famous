package backend

import "net/http"

func init() {
	http.HandleFunc("/api/update/popup.json", webapi_update_popup)
}

// no parameters
func webapi_update_popup(w http.ResponseWriter, r *http.Request) {
	switch v, err := update_popup(); err {
	case nil:
		webapi_write_as_json(w, v)
	default:
		webapi_write_error(w, err)
	}
}
