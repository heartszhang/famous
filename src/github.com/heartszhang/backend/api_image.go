package backend

import (
	"net/url"
	"os"
	"strings"

	"code.google.com/p/go.net/html"
	"github.com/heartszhang/curl"
	"github.com/heartszhang/feed"
	vt "github.com/heartszhang/videothumbnail"
)

// /api/image/video.thumbnail?uri=
func image_videothumbnail(uri string) (vt.VideoDescription, error) {
	return vt.DescribeVideo(uri)
}

// /feed/entry/image.json/{url}/{entry_id}
func image_description(uri string) (feed.FeedImage, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)

	if err == nil {
		return v, nil
	}
	c := curl.NewCurlerDetail(backend_context.config.ImageFolder, 0, 0, nil, backend_context.ruler)
	cache, err := c.Get(uri)

	v.Mime = cache.Mime
	v.Origin = cache.Local

	if err != nil {
		return v, err
	}
	v.Thumbnail, v.Mime, v.Width, v.Height, err = curl.NewThumbnail(cache.Local, backend_context.config.ThumbnailFolder, backend_context.config.ThumbnailWidth, 0)
	go imgo.save(uri, v)
	return v, err
}

func image_icon(uri string) (curl.Cache, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return curl.Cache{}, err
	}
	c := curl.NewCurlerDetail(backend_context.config.ImageFolder, 0, 0, nil, backend_context.ruler)
	candis := []string{u.RequestURI(), "/", "/favicon.ico"}
	for _, candi := range candis {
		x, _ := u.Parse(candi)
		cache, _ := c.Get(x.String())
		if cache.Mime == "text/html" {
			cache, _ = icon_from_link_rel(cache.Local)
		}
		if strings.Split(cache.Mime, "/")[0] == "image" {
			return cache, nil
		}
	}
	return curl.Cache{}, new_backenderror(-1, "icon cannot resolve")
}
func icon_from_link_rel(local string) (curl.Cache, error) {
	f, err := os.Open(local)
	if err != nil {
		return curl.Cache{}, err
	}
	defer f.Close()
	doc, err := html.Parse(f)
	if err != nil {
		return curl.Cache{}, err
	}
	de := node_query_select(doc, "html")
	head := node_query_select(de, "head")
	links := node_query_selects(head, "link")
	var href string
	for _, link := range links {
		rel := node_get_attribute(link, "rel")
		if rel == "icon" || rel == "shortcut icon" || rel == "apple-touch-icon" {
			href = node_get_attribute(link, "href")
			break
		}
	}
	if href != "" {
		c := curl.NewCurlerDetail(backend_config().ImageFolder, 0, 0, nil, backend_context.ruler)
		return c.Get(href)
	}
	return curl.Cache{}, new_backenderror(-1, "icon cannot resolved in html")
}
func image_description_cached(uri string) (feed.FeedImage, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)
	return v, err
}

func image_dimension(uri string) (feed.FeedImage, error) {
	imgo := new_imagecache_operator()
	v, err := imgo.find(uri)
	if err == nil {
		return v, err
	}
	v.Mime, v.Width, v.Height, _, err = curl.DescribeImage(uri)
	return v, err
}

func node_query_select(node *html.Node, name string) *html.Node {
	if node == nil {
		return nil
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == name {
			return child
		}
	}
	return nil
}
func node_query_selects(node *html.Node, name string) []*html.Node {
	if node == nil {
		return nil
	}
	var v []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == name {
			v = append(v, child)
		}
	}
	return v
}

func node_get_attribute(node *html.Node, name string) string {
	if node == nil {
		return ""
	}
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}
