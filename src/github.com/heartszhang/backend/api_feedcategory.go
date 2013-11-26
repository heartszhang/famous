package backend

// select a idle category_id, assigned to category
func feedcategory_create(name string) (string, error) {
	fco := new_feedcategory_operator()
	uid, err := fco.save(name)

	if uid == nil {
		return "", err
	}
	return uid.(string), err
}

// id : isn't root or all, drop the category whoes name is name
// id : other, drop categories
// name : can be empty. if id is root or all, name cann't be empty
func feedcategory_drop(name string) error {
	return nil
}

func feedcategory_all() ([]string, error) {
	fco := new_feedcategory_operator()
	return fco.all()
}
