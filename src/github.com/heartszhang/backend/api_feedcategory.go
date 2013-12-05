package backend

// create a category
// if category has existed, nothing will be done
func feedcategory_create(name string) error {
	fco := new_feedcategory_operator()
	_, err := fco.save(name)
	return err
}

// TODO: to be implemented
func feedcategory_drop(name string) error {
	panic("not implemented")
	return nil
}

func feedcategory_all() ([]string, error) {
	fco := new_feedcategory_operator()
	return fco.all()
}
