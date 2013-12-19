package backend

import (
	"net/http"
	"strconv"
)

func init() {
	http.HandleFunc("/api/image/description.json", webapi_image_description)        // ?uri=, return application/json
	http.HandleFunc("/api/image/thumbnail.jpeg", webapi_image_thumbnail)            // ?uri= return image/jpeg
	http.HandleFunc("/api/image/origin.jpeg", webapi_image_origin)                  // ?uri=, return image/xxx
	http.HandleFunc("/api/image/dimension.json", webapi_image_dimension)            // ?uri=, return FeedMedia
	http.HandleFunc("/api/image/video.thumbnail.jpeg", webapi_image_videothumbnail) // ?uri=, return image/jpeg
	http.HandleFunc("/api/image/icon", webapi_image_icon)                           // ?uri=, return image/xxx
}

// uri: string, required
func webapi_image_description(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")

	v, err := image_description(uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: string, requried
func webapi_image_dimension(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	v, err := image_dimension(uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: string, required
func webapi_image_videothumbnail(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	switch vt, err := image_videothumbnail(uri); err {
	case nil:
		webapi_image_entity(vt.Image, w, r)
	default:
		webapi_write_error(w, err)
	}
}

// uri: string, required
func webapi_image_entity(uri string, w http.ResponseWriter, r *http.Request) {
	switch cache, err := image_description(uri); err {
	case nil:
		image_writeheader(w, r, cache.Height, cache.Width, cache.Mime, cache.Thumbnail, cache.Origin)
	default:
		webapi_write_error(w, err)
	}
}
func image_writeheader(w http.ResponseWriter, r *http.Request, height, width int, mime, local, alter string) {
	w.Header().Set("content-type", mime)
	w.Header().Set("x-image-meta-property-height", strconv.Itoa(height))
	w.Header().Set("x-image-meta-property-width", strconv.Itoa(width))
	w.Header().Set("x-image-meta-property-alter", alter)

	http.ServeFile(w, r, local)
}

// uri: string, required
func webapi_image_thumbnail(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("uri")
	webapi_image_entity(url, w, r)
}

// uri: string, required
func webapi_image_origin(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	switch cache, err := image_description(uri); err {
	case nil:
		image_writeheader(w, r, cache.Height, cache.Width, cache.Mime, cache.Origin, cache.Thumbnail)
	default:
		webapi_write_error(w, err)
	}
}

// uri: string, required
func webapi_image_icon(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	switch img, err := image_icon(uri); err {
	default:
		webapi_write_error(w, err)
	case nil:
		w.Header().Set("content-type", img.Mime)
		http.ServeFile(w, r, img.Origin)
	}
}
