package backend

import (
	bing "github.com/heartszhang/bingsearchservice"
)

func suggest_bing(q string) ([]string, error) {
	return bing.BingSuggestion(q)
}
