package backend

import (
	"testing"
)

func TestFeedCategory(t *testing.T) {
	fco := new_feedcategory_operator()
	id, err := fco.Save("what")
	if id != nil {
		t.Log(id)
	}
	fc, err := fco.All()
	if err != nil {
		t.Error(err)
	}
	t.Log(fc)
	fco.Drop("what")
	if err != nil {
		t.Error(err)
	}
}
