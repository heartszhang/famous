package backend

import (
	"github.com/heartszhang/bing"
)

func suggest_bing(q string) ([]string, error) {
	return bing.BingSuggestion(q)
}
