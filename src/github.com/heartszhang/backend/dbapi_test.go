package backend

import (
	"fmt"
	"github.com/heartszhang/feedfeed"
	"testing"
)

func TestFeedSoruceImport(t *testing.T) {
	t.Skip()
	fs, err := feedfeed.CreateFeedsCategoryOpml("feedly.opml")
	if err != nil {
		t.Error(err)
		return
	}
	dbo := new_feedsource_operator()

	fs, err = dbo.save(fs)
	if err != nil {
		t.Error(err)
	}
	t.Log(fs)
}

func TestFeedSourceOps(t *testing.T) {
	t.Skip()
	dbo := new_feedsource_operator()
	uris := []string{"http://is.gd/e3zMW", "http://www.voachinese.com/rss/?count=20&zoneid=1915"}
	fs, err := dbo.findbatch(uris)
	if err != nil {
		t.Error(err)
	}
	t.Logf("feedsources: %v", fs)
}

/*
func TestFeedEntryMark(t *testing.T) {
	//	dbo := NewFeedEntryOperator()

}
*/
func TestFeedContentTouch(t *testing.T) {
	t.Skip("skipped")
	dbo := new_feedcontent_operator()
	x, err := dbo.touch(111101)
	if err != nil {
		t.Fatal(err)
	}
	x, err = dbo.touch(111101)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current touch", x)
}

func ExampleFeedEntries_unread() {
	dbo := new_feedentry_operator()
	fss, err := dbo.unread_count_sources()
	fmt.Println(fss, err)
	//Output:
}
func ExampleFeedEntry_unread() {
	dbo := new_feedentry_operator()
	fss, err := dbo.unread_count("http://feeds.feedburner.com/chinadigitaltimes/ThSg")
	fmt.Println(fss, err)
	//Output:
}
