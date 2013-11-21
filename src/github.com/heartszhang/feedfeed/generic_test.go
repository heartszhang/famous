package feedfeed

import (
	"testing"
)

func TestMakeFeedSource(t *testing.T) {
	fp := `g:\a.xml`
	fs, err := MakeFeedSource(fp)
	t.Log(fs, err)
	if err != nil {
		t.Error(err)
	}
}
