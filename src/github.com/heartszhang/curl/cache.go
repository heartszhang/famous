package curl

type Cache struct {
	Uri          string `json:"uri" bson:"uri"`
	Mime         string `json:"mime,omitempty" bson:"mime,omitempty"`
	Charset      string `json:"charset,omitempty" bson:"charset,omitempty"`
	Local        string `json:"local,omitempty" bson:"local,omitempty"`
	Disposition  string `json:"disposition,omitempty" bson:"disposition,omitempty"`
	LocalUtf8    string `json:"local_utf8,omitempty" bson:"local_utf8,omitempty"`
	LastModified string `json:"last_modified,omitempty" bson:"last_modified,omitempty"`
	ETag         string `json:"etag,omitempty" bson:"etag,omitempty"`
	Length       int64  `json:"length" bson:"length"`
	LengthUtf8   int64  `json:"length_utf8" bson:"length_utf8"`
	StatusCode   int    `json:"status_code" bson:"status_code"`
}
