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
/*
func TestFeedSourceOps(t *testing.T) {
	dbo := new_feedsource_operator()
	fs, err := dbo.all()
	if err != nil {
		t.Error(err)
	}
	t.Logf("feedsources count : %v", len(fs))
}

func TestFeedEntryMark(t *testing.T) {
	//	dbo := NewFeedEntryOperator()

}
*/
func TestFeedContentTouch(t *testing.T) {
	dbo := new_feedcontent_operator()
	x, err := dbo.touch(0)
	if err != nil {
		t.Fatal(err)
	}
	x, err = dbo.touch(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(x)
}
