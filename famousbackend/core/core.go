package core

import (
	//	"appengine"
	//	"appengine/taskqueue"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/hub/callback", webapi_hubcallback)
	http.HandleFunc("/", webapi_root)
}

func webapi_root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello go google engine")
}

func webapi_hubcallback(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t := r.FormValue("hub.challenge")
		//		w.Write([]byte(t))
		fmt.Fprint(w, t)
	case "POST":
		//		w.Write([]byte("OK"))
		fmt.Fprint(w, "OK")
	}
}
