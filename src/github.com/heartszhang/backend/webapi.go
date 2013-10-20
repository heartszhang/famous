package backend

import (
	"encoding/json"
	"github.com/gorilla/pat"
	"net/http"
)

func init() {
	p := pat.New()
	p.Get("/", webapi_home)
	http.Handle("/", p)
}

func webapi_home(w http.ResponseWriter, r *http.Request) {
	webapi_write_as_json(w, BackendStatus())
}

func webapi_write_as_json(w http.ResponseWriter, body interface{}) {
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(body)
}
