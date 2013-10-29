package backend

import (
	//	"github.com/heartszhang/feedfeed"
	"testing"
)

/*
func TestFeedSoruceImport(t *testing.T) {
	t.Skip("import skippted")
	fs, err := feedfeed.CreateFeedsCategoryOpml("feedly.opml")
	if err != nil {
		t.Error(err)
		return
	}
	dbo := NewFeedSourceOperator()
	fs, err = dbo.Save(fs)
	if err != nil {
		t.Error(err)
	}
	t.Log(fs)
}
*/

func TestFeedSourceOps(t *testing.T) {
	dbo := NewFeedSourceOperator()
	fs, err := dbo.AllSources()
	if err != nil {
		t.Error(err)
	}
	t.Logf("feedsources count : %v", len(fs))
}

func TestFeedEntryMark(t *testing.T) {
	//	dbo := NewFeedEntryOperator()

}
