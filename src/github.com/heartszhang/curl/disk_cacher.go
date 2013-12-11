package curl

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
)

type disk_cacher struct {
	data_folder string
}

func (this *disk_cacher) name(uri string) string {
	return filepath.Join(resolve_dir(this.data_folder), "cache."+url.QueryEscape(uri))
}

func (this *disk_cacher) load_index(uri string) (*Cache, error) {
	name := this.name(uri)
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cache := &Cache{StatusCode: 410} // no longer avail
	err = json.NewDecoder(f).Decode(cache)
	return cache, err
}
func (this *disk_cacher) save_index(uri string, cache Cache) error {
	name := this.name(uri)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(&cache)
	return err
}

func (this *disk_cacher) remove(uri string) error {
	name := this.name(uri)
	err := os.Remove(name)
	return err
}
