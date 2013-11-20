package feedfeed

import (
	"testing"
)

func TestMakeFeedSource(t *testing.T) {
	fp := `E:\SOURCESAFE\famous\src\github.com\heartszhang\backend\run\data\sources\xml.054210088`
	fs, err := MakeFeedSource(fp)
	t.Log(fs, err)
	if err != nil {
		t.Error(err)
	}
}
