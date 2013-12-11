package backend

import (
	"github.com/heartszhang/feed"
	"github.com/qiniu/log"
	"net/http"
	"strconv"
)

func init() {
	http.HandleFunc("/api/feed_entry/unread.json", webapi_feedentry_unread)
	http.HandleFunc("/api/feed_entry/mark.json", webapi_feedentry_mark)
	http.HandleFunc("/api/feed_entry/umark.json", webapi_feedentry_umark)
	http.HandleFunc("/api/feed_entry/fulldoc.json", webapi_feedentry_fulldoc)
	http.HandleFunc("/api/feed_entry/media.json", webapi_feedentry_media)
	http.HandleFunc("/api/feed_entry/drop.json", webapi_feedentry_drop)
	// same as feed_entry/unread
	http.HandleFunc("/api/feed_entry/source/unread.json", webapi_feedentry_source_unread)
	http.HandleFunc("/api/feed_entry/source/mark_read.json", webapi_feedentry_source_markread)
	http.HandleFunc("/api/feed_entry/source/unread_count.json", webapi_feedentry_source_unreadcount)
	http.HandleFunc("/api/feed_entry/sources/unread_count.json", webapi_feedentry_sources_unreadcount)
	http.HandleFunc("/api/feed_entry/category/unread.json", webapi_feedentry_category_unreadcount)
	http.HandleFunc("/api/feed_entry/categories/unread.json", webapi_feedentry_categories_unreadcount)
	http.HandleFunc("/api/feed_entry/category/mark_read.json", webapi_feedentry_category_markread)
	http.HandleFunc("/api/feed_entry/categories/unread_count.json", webapi_feedentry_categories_unreadcount)
	http.HandleFunc("/api/feed_entry/category/unread_count.json", webapi_feedentry_category_unreadcount)
}

func webapi_feedentry_category_unreadcount(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	cate := r.URL.Query().Get("category")
	switch v, err := new_feedentry_operator().unread_count_category(cate); err {
	case nil:
		webapi_write_as_json(w, v)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedentry_categories_unreadcount(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	switch v, err := new_feedentry_operator().unread_count_categories(); err {
	case nil:
		webapi_write_as_json(w, v)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedentry_sources_unreadcount(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	switch v, err := new_feedentry_operator().unread_count_sources(); err {
	case nil:
		webapi_write_as_json(w, v)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedentry_source_unreadcount(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")
	switch c, err := new_feedentry_operator().unread_count(uri); err {
	case nil:
		webapi_write_as_json(w, c)
	default:
		webapi_write_error(w, err)
	}
}

//?uri=&count=&page=
func webapi_feedentry_source_unread(w http.ResponseWriter, r *http.Request) {
	webapi_feedentry_unread(w, r)
}

// uri: /api/feed_entry/unread.json?uri=&count=&page=
func webapi_feedentry_unread(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")

	count, _ := strconv.ParseInt(r.URL.Query().Get("count"), 0, 0)
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 0, 0)
	if count == 0 {
		count = backend_context.config.EntryDefaultPageCount
	}
	log.Println("feedentry-unread(uri, count, page)", uri, count, page)
	switch fe, err, sc := feedentry_unread(uri, int(count), int(page)); err {
	case nil:
		webapi_write_as_json(w, fe)
	default:
		webapi_write_error_code(w, err, sc)
	}
}

// uri: /feed_entry/mark.json/{entry_id}/{flags}
func webapi_feedentry_mark(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("entry_uri")
	f, err := strconv.ParseUint(r.URL.Query().Get("flags"), 0, 0)
	flag := uint(f)
	if err == nil {
		flag, err = feedentry_mark(uri, flag)
	}
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, struct{}{})
	}
}

func webapi_feedentry_category_markread(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	cate := r.URL.Query().Get("category")
	read, err := strconv.ParseInt(r.URL.Query().Get("flags"), 0, 0)
	flag := uint(read)
	switch err = feedentry_category_mark(cate, flag); err {
	case nil:
		webapi_write_as_json(w, flag)
	default:
		webapi_write_error(w, err)
	}
}

func webapi_feedentry_source_markread(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	src := r.URL.Query().Get("uri")
	read, err := strconv.ParseInt(r.URL.Query().Get("flags"), 0, 0)
	flag := uint(read)
	switch err = feedentry_source_mark(src, flag); err {
	case nil:
		webapi_write_as_json(w, flag)
	default:
		webapi_write_error(w, err)
	}
}

// uri: /feed_entry/umark.json/{entry_id}/{flags}
func webapi_feedentry_umark(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")
	f, err := strconv.ParseInt(r.URL.Query().Get("flags"), 0, 0)
	flag := uint(f)
	if err == nil {
		flag, err = feedentry_umark(uri, flag)
	}
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, flag)
	}
}

// uri: /feed_entry/full_text.json/{entry_uri}
func webapi_feedentry_fulldoc(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri") // web-url
	ff, err := feedentry_fulldoc(uri)
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, ff)
	}
}

// uri: /feed_entry/media.json/{entry_id}/{url}/{media_type:[0-9]+}
func webapi_feedentry_media(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri := r.URL.Query().Get("uri")
	id := r.URL.Query().Get("entry_id")
	media_type, err := strconv.ParseUint(r.URL.Query().Get("media_type"), 0, 0)
	if err != nil {
		media_type = uint64(feed.Feed_media_type_unknown)
	}
	v, err := feedentry_media(uri, id, uint(media_type))
	if err != nil {
		webapi_write_error(w, err)
	} else {
		webapi_write_as_json(w, v)
	}
}

// uri: /feed_entry/drop.json/{entry_id}
// id is mongo's _id
func webapi_feedentry_drop(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	id := r.URL.Query().Get("entry_id")
	err := feedentry_drop(id)
	webapi_write_error(w, err)
}
