package main

import (
	"fmt"
	"github.com/heartszhang/baidu"
	"github.com/heartszhang/pubsub"
)

//http://bcms.api.duapp.com/rest/2.0/bcms/{queue_name}?{query_string}
const (
	client_id     = "DHMKfwY1xbBMS2WvMz44CLyc"
	client_secret = "Hn6SqT20RtBPMIcufOHG0IDE0xs6p45g"
	bcmsq         = "3c6848d073e106f2f2531059251190b8"
	bcms_host     = "bcms.api.duapp.com"
	access_token  = "3.a0d9c307f27262cba7e4c542170d7494.2592000.1387802390.1174464975-1685584"
	refresh_token = "4.84d83482e856eaa6dc29b0935c99660c.315360000.1700570390.1174464975-1685584"
	bcms_urlbase  = "http://bcms.api.duapp.com/rest/2.0/bcms/"
)

//application/x-www-form-urlencoded

func main() {
	//	bq := baidu.NewBcms(bcmsq, access_token, client_id, client_secret)
	bq := baidu.NewBcmsProxy("http://iweizhi2.duapp.com/pop.json")
	var v pubsub.PubsubMessage
	err := bq.FetchOneAsJson(&v)
	fmt.Println(v, err)
}
