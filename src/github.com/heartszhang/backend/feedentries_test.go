package backend

import (
	"testing"
)

func TestFeedEntries(t *testing.T) {
	t.Skip("skipping in normal state")
	_, err := feedentry_unread("http://yyyyiiii.blogspot.com/feeds/posts/default", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
}
