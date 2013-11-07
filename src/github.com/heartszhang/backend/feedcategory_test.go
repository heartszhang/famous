package backend

import (
	"testing"
)

func TestFeedCategory(t *testing.T) {
	t.Skip("skipping category test")
	fco := new_feedcategory_operator()
	id, err := fco.save("what")
	if id != nil {
		t.Log(id)
	}
	fc, err := fco.all()
	if err != nil {
		t.Error(err)
	}
	t.Log(fc)
	fco.drop("what")
	if err != nil {
		t.Error(err)
	}
}
