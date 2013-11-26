package bing

import (
	"github.com/heartszhang/curl"
	"github.com/heartszhang/oauth2"
	"log"
	"net/http"
	"net/url"
)

type BingSearchService interface {
	Web(params BingSearchWebParameters) (curl.Cache, error)
	Image(params BingSearchImageParameters) (curl.Cache, error)
	Video(p BingSearchVideoParameters) (curl.Cache, error)
	News(p BingSearchNewsParameters) (curl.Cache, error)
	SpellingSuggestion(p BingSearchParameters) (curl.Cache, error)
	RelatedSearch(p BingSearchParameters) (curl.Cache, error)
	Composite(p BingSearchCompositeParameters) (curl.Cache, error)
}
type bing_search struct {
	temp_folder string
	account_key string
}

// bing seach service need http basic-authorization use account_key as username and password
func (this bing_search) RoundTrip(r *http.Request) {
	r.URL.User = url.UserPassword(this.account_key, this.account_key)
}
func NewBingSearchService(temp, acckey string) BingSearchService {
	return &bing_search{temp, acckey}
}

func (this bing_search) Web(params BingSearchWebParameters) (curl.Cache, error) {
	q := oauth2.HttpQueryEncode(params)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyNoProxy, 0, this)
	return c.Get(BingSearchServiceRoot + ServiceOperationWeb + "?" + q)
}

func (this bing_search) News(params BingSearchNewsParameters) (curl.Cache, error) {
	q := oauth2.HttpQueryEncode(params)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyNoProxy, 0, this)
	return c.Get(BingSearchServiceRoot + ServiceOperationNews + "?" + q)
}

func (this bing_search) Image(params BingSearchImageParameters) (curl.Cache, error) {
	q := oauth2.HttpQueryEncode(params)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyNoProxy, 0, this)
	return c.Get(BingSearchServiceRoot + ServiceOperationImage + "?" + q)
}

func (this bing_search) Video(params BingSearchVideoParameters) (curl.Cache, error) {
	q := oauth2.HttpQueryEncode(params)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyNoProxy, 0, this)
	return c.Get(BingSearchServiceRoot + ServiceOperationVideo + "?" + q)
}
func new_string(v string) *string {
	return &v
}
func (this bing_search) SpellingSuggestion(params BingSearchParameters) (curl.Cache, error) {
	params.Format = new_string(SearchResultFormatJson)
	q := oauth2.HttpQueryEncode(params)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyNoProxy, 0, this)
	x := BingSearchServiceRoot + ServiceOperationSpellingSuggestion + "?" + q
	log.Println(x)
	return c.Get(x)
}

func (this bing_search) RelatedSearch(params BingSearchParameters) (curl.Cache, error) {
	params.Format = new_string(SearchResultFormatJson)
	q := oauth2.HttpQueryEncode(params)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyNoProxy, 0, this)
	return c.Get(BingSearchServiceRoot + ServiceOperationRelatedSearch + "?" + q)
}

func (this bing_search) Composite(params BingSearchCompositeParameters) (curl.Cache, error) {
	q := oauth2.HttpQueryEncode(params)
	c := curl.NewCurlerDetail(this.temp_folder, curl.CurlProxyPolicyNoProxy, 0, this)
	return c.Get(BingSearchServiceRoot + ServiceOperationComposite + "?" + q)
}

type BingSearchParameters struct {
	Top       *int    `param:"$top"`
	Skip      *int    `param:"$skip"`
	Format    *string `param:"$format"`
	Query     string  `param:"Query"`
	Market    string  `param:"Market"`
	Latitude  *int    `param:"Latitude"`
	Longitude *int    `param:"Longitude"`
	Adult     *int    `param:"Adult"`
	Options   string  `param:"Options"`
}

type BingSearchWebParameters struct {
	BingSearchParameters
	WebFileType string `param:"WebFileType"`
	WebOptions  string `param:"WebOptions"`
}

type BingSearchCompositeParameters struct {
	BingSearchParameters
	Sources string `param:"Sources"`
}

type BingSearchWebEntry struct {
	ID          string `json:"ID"`
	Title       string `json:"Title,omitempty"`
	Description string `json:"Description,omitempty"`
	DisplayUrl  string `json:"DisplayUrl,omitempty"`
	Url         string `json:"Url,omitempty"`
}

type BingSearchImageParameters struct {
	BingSearchParameters
	ImageFilters string `param:"ImageFilters"`
}

type BingSearchImageEntry struct {
	ID          string `json:"ID,omitempty"`
	Title       string `json:"Title,omitempty"`
	MediaUrl    string `json:"MediaUrl,omitempty"`
	SourceUrl   string `json:"SourceUrl,omitempty"`
	DisplayUrl  string `json:"DisplayUrl,omitempty"`
	Width       int32  `json:"Width"`
	Height      int32  `json:"Height"`
	FileSize    int64  `json:"FileSize"`
	ContentType string `json:"ContentType,omitempty"`
	//	Thumbnail   *BingSearchThumbnail `json:"Thumbnail,omitempty"`
}

type BingSearchVideoParameters struct {
	BingSearchParameters
	VideoFilters string `param:"VideoFilters"`
	VideoSortBy  string `param:"VideoSortBy"`
}

type BingSearchVideoEntry struct {
	ID         string `json:",omitempty"`
	Title      string `json:",omitempty"`
	MediaUrl   string `json:",omitempty"`
	DisplayUrl string `json:",omitempty"`
	RunTime    int32  `json:",omitempty"`
	//	Thumbnail  *BingSearchThumbnail `json:",omitempty"`
}

type BingSearchNewsParameters struct {
	BingSearchParameters
	NewsCategory         string `param:"NewsCategory"`
	NewsLocationOverride string `param:"NewsLocationOverride"`
	NewsSortBy           string `param:"NewsSortBy"`
}

type BingSearchNewsEntry struct {
	ID          string `json:",omitempty"`
	Title       string `json:",omitempty"`
	Url         string `json:",omitempty"`
	Source      string `json:",omitempty"`
	Description string `json:",omitempty"`
	Date        string `json:",omitempty"`
}

type BingSearchRelatedEntry struct {
	ID      string `json:",omitempty"`
	Title   string `json:",omitempty"`
	BingUrl string `json:",omitempty"`
}

const (
	BingSearchServiceRoot        = "https://api.datamarket.azure.com/Bing/Search/v1/"
	BingSearchWebOnlyServiceRoot = "https://api.datamarket.azure.com/Bing/SearchWeb/v1/"

	ServiceOperationWeb                = "Web"
	ServiceOperationImage              = "Image"
	ServiceOperationVideo              = "Video"
	ServiceOperationNews               = "News"
	ServiceOperationSpellingSuggestion = "SpellingSuggestions"
	ServiceOperationRelatedSearch      = "RelatedSearch"
	ServiceOperationComposite          = "Composite"

	SearchResultFormatJson                  = "JSON"
	SearchResultFormatAtom                  = "Atom"
	OptionsDisableLocationDetection         = "DisableLocationDetection"
	OptionsEnableHighlighting               = "EnableHighlighting"
	WebSearchOptionsDisableHostCollapsing   = "DisableHostCollapsing"
	WebSearchOptionsDisableQueryAlterations = "DisableQueryAlterations"
	AdultOff                                = "Off"
	AdultModerate                           = "Moderate"
	AdultStrict                             = "Strict"
	MarketZh                                = "zh-CN"
	NewsCategoryBusiness                    = "rt_Business"
	NewsCategoryEntertainment               = "rt_Entertainment"
	NewsCategoryHealth                      = "rt_Health"
	NewsCategoryPolitics                    = "rt_Politics"
	NewsCategorySports                      = "rt_Sports"
	NewsCategoryUS                          = "rt_US"
	NewsCategoryWorld                       = "rt_World"
	NewsCategoryScience                     = "rt_ScienceAndTechnology"
	SortByDate                              = "Date"
	SortByRelevance                         = "Relevance"
)
