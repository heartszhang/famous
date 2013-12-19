package backend

import (
	"net/http"

	"github.com/heartszhang/feed"
)

func init() {
	http.HandleFunc("/api/opml/upload.json", webapi_opml_upload)
}

/* multipartform-encoded file
 * single file only
 * utf-8 encoded
 */
func webapi_opml_upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(_1M)
	if err != nil {
		webapi_write_error(w, err)
	}
	for _, headers := range r.MultipartForm.File {
		for _, header := range headers {
			f, err := header.Open()
			defer f.Close()
			fs, err := feed.OpmlExportFeedSource(f)
			switch err {
			case nil:
				rs := feedsource_mark_subscribed(new_readsources(fs))
				webapi_write_as_json(w, rs)
			default:
				webapi_write_error(w, err)
			}
			return
		}
	}
	webapi_write_error(w, nil)
}
