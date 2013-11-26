package main

import (
	"fmt"
	"github.com/heartszhang/oauth2"
)

const (
	baidu_id        = "1685584"
	client_id       = "DHMKfwY1xbBMS2WvMz44CLyc"
	client_secret   = "Hn6SqT20RtBPMIcufOHG0IDE0xs6p45g"
	baidu_oauth_url = "https://openapi.baidu.com/oauth/2.0/authorize"
	baidu_token_url = "https://openapi.baidu.com/oauth/2.0/token"
	redirect_url    = "http://iweizhi2.duapp.com/authorize"
)

func main() {
	oc := oauth2.OAuthConfig{
		ClientId:     client_id,
		ClientSecret: client_secret,
		AuthUrl:      baidu_oauth_url,
		TokenUrl:     baidu_token_url,
		RedirectUrl:  redirect_url,
	}
	auth := oauth2.NewWebAuthorizationBroker(oc, nil)
	token, err := auth.Authorize(true, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token.AccessToken)
	fmt.Println(token.RefreshToken)
}
