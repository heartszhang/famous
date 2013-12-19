package gfwlist

import (
	"bufio"
	"io"
	"net/url"
	"strings"
)

type GfwRuler interface {
	IsBlocked(uri string) bool
}
type uri_item struct {
	scheme, hostpath, hostpath_query, full string
	scheme_index                           int
}

type rule interface {
	match(item uri_item) bool
}

// gfwlist_file should be encoded with utf8 or ascii
func NewGfwRuler(f io.Reader) (GfwRuler, error) {
	gc := &gfw_ruler{}
	gc.initialize()
	reader := bufio.NewReader(f)
	var (
		comments, pathes, prefixes, exceptions, regexes, pathqueries int
		err                                                          error
	)
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		line = strings.TrimSpace(line)
		switch {
		case len(line) == 0: // ignore empty line
		case strings.HasPrefix(line, `!`): // comment
			comments++
		case strings.HasPrefix(line, `[`): // ?, ignored
		case strings.HasPrefix(line, `||`): // path match
			gc.add_path(line[2:])
			pathes++
		case strings.HasPrefix(line, `|`): // prefix match
			gc.add_prefix(line[1:])
			prefixes++
		case strings.HasPrefix(line, `@@||`): // exception path match
			gc.add_exception(line[4:])
			exceptions++
		case strings.HasPrefix(line, `/`) && strings.HasSuffix(line, `/`): // regex
			gc.add_regex(strings.Trim(line, "/"))
			regexes++
		default:
			gc.add_default(line)
			pathqueries++
		}
	}

	return gc, nil
}

type gfw_ruler struct {
	blocked_rules    gfw_group
	exceptions_rules gfw_group
}

type gfw_group struct {
	path_rules      [2]wu_manber // http and https
	querypath_rules [2]wu_manber // http and https
	prefix_rules    prefix_group
	regex_rules     regex_group // not used
}

func (this *gfw_ruler) initialize() {
	this.blocked_rules.path_rules[0].initialize()
	this.blocked_rules.path_rules[1].initialize()
	this.exceptions_rules.path_rules[0].initialize()
	this.exceptions_rules.path_rules[1].initialize()
	this.blocked_rules.querypath_rules[0].initialize()
	this.blocked_rules.querypath_rules[1].initialize()
	this.exceptions_rules.querypath_rules[0].initialize()
	this.exceptions_rules.querypath_rules[1].initialize()
	this.blocked_rules.prefix_rules.initialize()
	this.exceptions_rules.prefix_rules.initialize()
}

type prefix_group struct {
	prefix_rules [max_prefix_length + 1]fixed_prefix_group
}

func (this *prefix_group) add(pattern string) {
	len := int_min(len(pattern), max_prefix_length)
	if len <= 0 {
		return
	}
	prefix := pattern[:len]
	key := make_key(prefix)
	this.prefix_rules[len][key] = append(this.prefix_rules[len][key], pattern)
}
func (this *prefix_group) initialize() {
	for i := 0; i < len(this.prefix_rules); i++ {
		this.prefix_rules[i] = make(fixed_prefix_group)
	}
}

type fixed_prefix_group map[uint32]prefix_list
type prefix_list []string

type regex_group struct {
}

func (this *regex_group) add(pattern string) {
}

const (
	http_index = iota
	https_index

	max_prefix_length = 10
)

func (this *gfw_ruler) add_default(line string) {
	this.blocked_rules.querypath_rules[http_index].add(line)
}
func (this *gfw_ruler) add_path(line string) {
	this.blocked_rules.path_rules[http_index].add(line)
	this.blocked_rules.path_rules[https_index].add(line)
}
func (this *gfw_ruler) add_prefix(line string) {
	this.blocked_rules.prefix_rules.add(line)
}
func (this *gfw_ruler) add_exception(line string) {
	this.exceptions_rules.path_rules[0].add(line)
	this.exceptions_rules.path_rules[1].add(line)
}
func (this *gfw_ruler) add_regex(line string) {
	this.blocked_rules.regex_rules.add(line)
}

func (this *gfw_ruler) match(ui uri_item) bool {
	if this.exceptions_rules.match(ui) {
		return false
	}
	if this.blocked_rules.match(ui) {
		return true
	}
	return false
}

func (this *gfw_group) match(ui uri_item) bool {
	if this.prefix_rules.match(ui) {
		return true
	}
	if this.path_rules[ui.scheme_index].match(ui.hostpath) {
		return true
	}
	if this.querypath_rules[ui.scheme_index].match(ui.hostpath_query) {
		return true
	}
	return false
}

func (this *prefix_group) match(ui uri_item) bool {
	hashes := this.make_hashes(ui.full) //[]uint32
	for i := 0; i < len(hashes); i++ {
		if rules, ok := this.prefix_rules[i][hashes[i]]; ok {
			if rules.match(ui) {
				return true
			}
		}
	}
	return false
}
func (this *prefix_group) make_hashes(uri string) []uint32 {
	len := int_min(max_prefix_length, len(uri))
	v := make([]uint32, len+1)
	for i := 1; i <= len; i++ {
		v[i] = make_key(uri[:i])
	}
	return v
}
func (this prefix_list) match(ui uri_item) bool {
	for _, pattern := range []string(this) {
		if strings.HasPrefix(ui.full, pattern) {
			return true
		}
	}
	return false
}

func (this *gfw_ruler) IsBlocked(uri string) bool {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}
	si := -1
	if u.Scheme == "http" {
		si = http_index
	} else if u.Scheme == "https" {
		si = https_index
	} else {
		return false
	}
	ui := uri_item{
		scheme:         u.Scheme,
		hostpath:       u.Host + u.Path,
		hostpath_query: u.Host + u.Path + u.RawQuery,
		full:           uri,
		scheme_index:   si,
	}
	return this.match(ui)
}
