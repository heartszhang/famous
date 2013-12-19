package oauth2

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Config is the configuration of an OAuth consumer.
type OAuthConfig struct {
	ClientId     string // ClientId is the OAuth client identifier used when communicating with oauth-provider
	ClientSecret string // ClientSecret is the OAuth client secret used when communicating with oauth provider
	Scope        string // Scope identifies the level of access being requested. may be space-delimited
	AuthUrl      string // AuthURL is the URL the user will be directed to in order to grant access
	TokenUrl     string // TokenURL is the URL used to retrieve OAuth tokens.
	RedirectUrl  string // RedirectURL is the URL to which the user will be returned to
}

type WebAuthorizationBroker interface {
	Authorize(force bool, extras *map[string]string) (*Token, error)
}

func NewWebAuthorizationBroker(config OAuthConfig, store TokenStorage) WebAuthorizationBroker {
	return &local_webauth_broker{config, store, local_codereceiver{}}
}

type local_webauth_broker struct {
	config        OAuthConfig
	token_storage TokenStorage // TokenCache allows tokens to be cached for subsequent requests.
	code_receiver local_codereceiver
}

func (this *local_webauth_broker) Authorize(force bool, extras *map[string]string) (*Token, error) {
	var token *Token
	var err error
	if !force {
		token, err = this.load_token()
	}
	if token != nil && (token.RefreshToken != "" || !token.Expired()) {
		return token, nil
	}
	code_req_url := this.authcode_url(this.code_receiver.redirect_url(), extras)
	code, err := this.code_receiver.receive_code(code_req_url)
	if err != nil {
		return token, err
	}
	token, err = this.exchange_code(code)
	if token != nil && err != nil {
		this.store_token(*token)
	}
	return token, err
}

func (this *local_webauth_broker) exchange_code(code string) (*Token, error) {
	vals := url.Values{"grant_type": {"authorization_code"},
		"redirect_uri":  {this.config.RedirectUrl},
		"code":          {code},
		"client_id":     {this.config.ClientId},
		"client_secret": {this.config.ClientSecret},
	}
	if this.config.Scope != "" {
		vals.Add("scope", this.config.Scope)
	}
	resp, err := http.DefaultClient.PostForm(this.config.TokenUrl, vals)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, &oauth_error{resp.Status, "", ""}
	}
	token := Token{Issued: time.Now()}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&token)
	return &token, err
}

func (this *local_webauth_broker) load_token() (*Token, error) {
	if this.token_storage == nil {
		return nil, &oauth_error{"broker.load-token", "no storage", ""}
	}
	return this.token_storage.Token()
}
func (this *local_webauth_broker) store_token(t Token) {
	if this.token_storage == nil {
		return
	}
	this.token_storage.PutToken(t)
}

type local_codereceiver struct {
	listener net.Listener
	code     string
	err      error
}

func (this *local_codereceiver) redirect_url() string {
	var err error
	this.listener, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		return ""
	}
	addr := this.listener.Addr().(*net.TCPAddr)
	return fmt.Sprintf("http://localhost:%v/authorize", addr.Port)
}

const (
	close_page_response = `<html>
  <head><title>OAuth 2.0 Authentication Token Received</title></head>
  <body>
    Received verification code. Closing...
    <script type='text/javascript'>
      window.setTimeout(function() {
          window.open('', '_self', '');
          window.close();
        }, 1000);
      if (window.opener) { window.opener.checkToken(); }
    </script>
  </body>
</html>`
)

func (this *local_codereceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.code = r.FormValue("code")
	if this.code == "" {
		this.err = &oauth_error{r.FormValue("error"), r.FormValue("error_description"), r.FormValue("error_uri")}
	}
	w.Header().Set("content-type", "text/html;charset=utf-8")
	io.WriteString(w, close_page_response)
	f, cf := w.(http.Flusher)
	if cf {
		f.Flush()
	}
	c, _, err := w.(http.Hijacker).Hijack()
	if err == nil {
		c.Close()
	}
	this.listener.Close()
}

func (this *local_codereceiver) receive_code(authcodeurl string) (code string, err error) {
	defer this.listener.Close()
	err = exec.Command("cmd", "/c", "start", strings.Replace(authcodeurl, "&", "^&", -1)).Start()
	if err != nil {
		return "", err
	}
	err = http.Serve(this.listener, this)
	switch {
	case this.err != nil:
		return "", this.err
	case this.code == "":
		return "", err
	default:
		return this.code, nil
	}
}

// AuthCodeURL returns a URL that the end-user should be redirected to,
// so that they may obtain an authorization code.
func (this *local_webauth_broker) authcode_url(state string, additionals *map[string]string) string {
	url_, err := url.Parse(this.config.AuthUrl)
	if err != nil {
		return ""
	}
	vals := url.Values{
		"response_type": {"code"},
		"client_id":     {this.config.ClientId},
		"redirect_uri":  {this.config.RedirectUrl},
		"state":         {state},
	}
	if this.config.Scope != "" {
		vals.Set("scope", this.config.Scope)
	}
	if additionals != nil {
		for k, v := range *additionals {
			vals.Set(k, v)
		}
	}
	q := vals.Encode()
	if url_.RawQuery == "" {
		url_.RawQuery = q
	} else {
		url_.RawQuery += "&" + q
	}
	return url_.String()
}

type oauth_error struct {
	prefix      string
	description string
	uri         string
}

func (oe oauth_error) Error() string {
	return oe.prefix + ": " + oe.description + oe.uri
}

// Cache specifies the methods that implement a Token cache.
type TokenStorage interface {
	Token() (*Token, error)
	PutToken(Token) error
}

// FileTokenStorage implements Cache. Its value is the name of the file in which
// the Token is stored in JSON format.
type FileTokenStorage string

func (f FileTokenStorage) Token() (*Token, error) {
	file, err := os.Open(string(f))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	tok := &Token{}
	if err := json.NewDecoder(file).Decode(tok); err != nil {
		return nil, err
	}
	return tok, nil
}

func (f FileTokenStorage) PutToken(tok *Token) error {
	file, err := os.OpenFile(string(f), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(file).Encode(tok); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

// Token contains an end-user's tokens.
// This is the data you must store to persist authentication.
type Token struct {
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Id           string    `json:"id_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	ExpiresIn    int64     `json:"expires_in"`
	Issued       time.Time `json:"issued"`
}

func (t *Token) Expired() bool {
	if t.ExpiresIn == 0 {
		return false
	}
	return t.Issued.Add(time.Duration(t.ExpiresIn) * time.Second).Before(time.Now())
}
