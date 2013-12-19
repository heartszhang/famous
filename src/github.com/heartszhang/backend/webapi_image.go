package backend

import (
	"net/http"
	"strconv"

	"github.com/qiniu/log"
)

func init() {
	http.HandleFunc("/api/image/description.json", webapi_image_description)
	http.HandleFunc("/api/image/thumbnail.json", webapi_image_thumbnail)       // ?uri= return image/jpeg
	http.HandleFunc("/api/image/origin.json", webapi_image_origin)             // ?uri=, return image/xxx
	http.HandleFunc("/api/image/dimension.json", webapi_image_dimension)       // ?uri=, return FeedMedia
	http.HandleFunc("/api/image/video.thumbnail", webapi_image_videothumbnail) // ?uri=, return image/xxx
	http.HandleFunc("/api/image/icon", webapi_image_icon)                      // ?uri=, return image/xxx
}

// uri: /image/description.json?uri=
func webapi_image_description(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")

	v, err := image_description(uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /image/dimension.json?uri=
func webapi_image_dimension(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")
	v, err := image_dimension(uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}
func webapi_image_videothumbnail(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")
	switch vt, err := image_videothumbnail(uri); err {
	case nil:
		webapi_image_entity(vt.Image, w, r)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_image_entity(uri string, w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	switch cache, err := image_description(uri); err {
	case nil:
		w.Header().Set("content-type", cache.Mime)
		w.Header().Set("x-image-meta-property-height", strconv.Itoa(cache.Height))
		w.Header().Set("x-image-meta-property-width", strconv.Itoa(cache.Width))
		w.Header().Set("x-image-meta-property-alter", cache.Origin)

		http.ServeFile(w, r, cache.Thumbnail)
	default:
		webapi_write_error(w, err)
	}
}

// /api/image/thumbnail.json?uri=
func webapi_image_thumbnail(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	url := r.URL.Query().Get("uri")
	webapi_image_entity(url, w, r)
}

// /api/image/origin.json?uri=
func webapi_image_origin(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")
	switch cache, err := image_description(uri); err {
	case nil:
		w.Header().Set("content-type", cache.Mime)
		w.Header().Set("x-image-meta-property-height", strconv.Itoa(cache.Height))
		w.Header().Set("x-image-meta-property-width", strconv.Itoa(cache.Width))
		w.Header().Set("x-image-meta-property-alter", cache.Thumbnail)
		http.ServeFile(w, r, cache.Origin)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_image_icon(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")
	switch cache, err := image_icon(uri); err {
	default:
		webapi_write_error_code(w, err, cache.StatusCode)
	case nil:
		w.Header().Set("content-type", cache.Mime)
		http.ServeFile(w, r, cache.Local)
	}
}
