package backend

import (
	"encoding/json"
	"github.com/gorilla/pat"
	"net/http"
	"strconv"
)

func init() {
	p := pat.New()
	p.Get("/feeds/entries_since.json/{since_unixtime:[0-9]+}/{category:[0-9]+}/{count:[0-9]+}/{page:[0-9]+}", webapi_feeds_entries_since)

	p.Get("/api/feed_source/subscribe.json/{uri}/{source_type}/{category}", webapi_feed_source_subscribe)
	p.Get("/api/meta.json", webapi_meta)
	p.Get("/api/feed_entry/mark.json/{entry_id}/{flags}", webapi_feed_entry_mark)
	p.Get("/api/feed_entry/umark.json/{entry_id}/{flags}", webapi_feed_entry_umark)
	p.Get("/api/feed_entry/full_text.json/{entry_id}/{uri}", webapi_feed_entry_fulltext)
	p.Get("/api/feed_entry/image.json/{entry_id}/{url}", webapi_feed_entry_image)
	p.Get("/api/feed_entry/media.json/{entry_id}/{url}/{media_type:[0-9]+}", webapi_feed_entry_media)
	p.Get("/api/feed_entry/drop.json/{entry_id}", webapi_feed_entry_drop)
	p.Get("/api/feed_category/create.json/{name}", webapi_feed_category_create)
	p.Get("/api/feed_catetory/drop.json/{name}", webapi_feed_category_drop)
	p.Get("/api/tick.json", webapi_tick)
	p.Get("/api/feed_source/unsubscribe.json/{uri}/{source_type}/{category}", webapi_feed_source_unsubscribe)
	p.Get("/api/meta/categories.json", webapi_meta_categories)
	p.Get("/api/meta/cleanup.json", webapi_meta_cleanup)

	p.Get("/", webapi_home)
	http.Handle("/", p)
}

// uri: /feeds/entries_since.json/{since_unixtime:[0-9]+}/{category:[0-9]+}/{count:[0-9]+}/{page:[0-9]+}
func webapi_feeds_entries_since(w http.ResponseWriter, r *http.Request) {
	since, err := strconv.ParseInt(r.URL.Query().Get(":since_unixtime"), 0, uint64_bits)
	if err != nil {
		since = unixtime_now()
	}
	category, err := strconv.ParseUint(r.URL.Query().Get(":category"), 0, uint64_bits)
	if err != nil {
		category = feed_category_root
	}
	count, err := strconv.ParseUint(r.URL.Query().Get(":count"), 0, 0)
	if err != nil {
		count = 20
	}
	page, err := strconv.ParseUint(r.URL.Query().Get(":page"), 0, 0)
	if err != nil {
		page = 0
	}
	fe, err := feeds_entries_since(since, category, uint(count), uint(page))

	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, fe)
	}
}

// uri: /feed_source/subscribe.json/{uri}/{source_type}/{category}
func webapi_feed_source_subscribe(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get(":uri")
	source_type := source_type_map(r.URL.Query().Get(":source_type"))
	category, err := strconv.ParseUint(r.URL.Query().Get(":category"), 0, uint64_bits)
	if err != nil {
		category = feed_category_root
	}

	fs, err := feed_source_subscribe(url, source_type, category)

	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, fs)
	}
}

// uri: /meta.json
func webapi_meta(w http.ResponseWriter, r *http.Request) {
	m, err := meta()

	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, m)
	}
}

// uri: /feed_entry/mark.json/{entry_id}/{flags}
func webapi_feed_entry_mark(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":entry_id")
	f, err := strconv.ParseInt(r.URL.Query().Get(":flags"), 0, 0)
	flag := int(f)
	if err == nil {
		flag, err = feed_entry_mark(id, flag)
	}
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, struct{}{})
	}
}

// uri: /feed_entry/umark.json/{entry_id}/{flags}
func webapi_feed_entry_umark(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	f, err := strconv.ParseInt(r.URL.Query().Get(":flags"), 0, 0)
	flag := int(f)
	if err == nil {
		flag, err = feed_entry_umark(id, flag)
	}
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, struct{}{})
	}
}

// uri: /feed_entry/full_text.json/{entry_id}/{uri}
func webapi_feed_entry_fulltext(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get(":uri")
	id := r.URL.Query().Get(":entry_id")
	ff, err := feed_entry_fulltext(url, id)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, ff)
	}
}

// uri: /feed_entry/image.json/{entry_id}/{url}
func webapi_feed_entry_image(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get(":uri")
	id := r.URL.Query().Get(":entry_id")
	v, err := feed_entry_image(url, id)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /feed_entry/media.json/{entry_id}/{url}/{media_type:[0-9]+}

func webapi_feed_entry_media(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get(":uri")
	id := r.URL.Query().Get(":entry_id")
	media_type, err := strconv.ParseUint(r.URL.Query().Get(":media_type"), 0, 0)
	if err != nil {
		media_type = uint64(feed_media_type_unknown)
	}
	v, err := feed_entry_media(url, id, uint(media_type))
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /feed_entry/drop.json/{entry_id}

// id is mongo's _id
func webapi_feed_entry_drop(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":entry_id")
	err := feed_entry_drop(id)
	webapi_write_error(w, err)
}

// uri: /feed_category/create.json/{name}
func webapi_feed_category_create(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(":name")
	v, err := feed_category_create(name)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /feed_catetory/drop.json/{name}
func webapi_feed_category_drop(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(":name")
	err := feed_category_drop(name)
	webapi_write_error(w, err)
}

// uri: /tick.json
func webapi_tick(w http.ResponseWriter, r *http.Request) {
	v, err := tick()
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /feed_source/unsubscribe.json/{uri}/{source_type}/{category}
func webapi_feed_source_unsubscribe(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get(":uri")
	source_type := source_type_map(r.URL.Query().Get(":source_type"))
	category, err := strconv.ParseUint(r.URL.Query().Get(":category"), 0, uint64_bits)
	if err != nil {
		category = feed_category_none
	}
	err = feed_source_unsubscribe(url, source_type, category)
	webapi_write_error(w, err)
}

// uri: /meta/categories.json
func webapi_meta_categories(w http.ResponseWriter, r *http.Request) {
	v, err := meta_categories()
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /meta/cleanup.json
func webapi_meta_cleanup(w http.ResponseWriter, r *http.Request) {
	err := meta_cleanup()
	webapi_write_error(w, err)
}

func webapi_write_as_json(w http.ResponseWriter, body interface{}) {
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(body)
}

func webapi_write_error(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		webapi_write_as_json(w, BackendError{Reason: err.Error(), Code: -1})
	} else {
		webapi_write_as_json(w, BackendError{})
	}
}

// uri: /
func webapi_home(w http.ResponseWriter, r *http.Request) {
	webapi_write_as_json(w, BackendStatus())
}
