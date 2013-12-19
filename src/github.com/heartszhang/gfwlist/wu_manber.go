package gfwlist

import "strings"

const (
	wu_manber_m          = 5 // The minimum length of a pattern
	wu_manber_b          = 2 // 2 or 3  the size of the block
	wu_manber_shift_size = 65536
)

/*
Starts by comparing the last character
Looking text in blocks instead of one by one char
Hash functions and tables are used

Three tables to build: a SHIFT table, a HASH table, and a PREFIX table

*/
type wu_manber struct {
	shift    [wu_manber_shift_size]byte
	patterns map[int][]wm_pattern // pattern hash table
	size     int
}

type wm_pattern struct {
	pattern string
	prefix  string // prefix of pattern, length is B
}

func (this *wu_manber) initialize() {
	for i := 0; i < wu_manber_shift_size; i++ {
		this.shift[i] = wu_manber_m - wu_manber_b + 1
	}
	this.patterns = make(map[int][]wm_pattern)
}

func (this *wu_manber) wm_hash(str string, index int) int {
	a := int(str[index])
	b := int(str[index+1]) // this means that wu_manber_b = 2
	key := a<<8 + b
	return key
}

func (this *wu_manber) add(pattern string) {
	if len(pattern) < wu_manber_m { // too short pattern is ignored
		//		log.Println("ignore short pattern", pattern)
		return
	}

	this.size++
	for i := 0; i < wu_manber_m-1; i++ {
		// Compute a hash value h based on the current B
		key := this.wm_hash(pattern, i)
		c := this.shift[key]
		// the position that X ends in pattern
		shift := byte(wu_manber_m - wu_manber_b - i)
		if shift < c {
			this.shift[key] = shift
		}
		if shift == 0 {
			this.patterns[key] = append(this.patterns[key], new_wu_pattern(pattern))
		}
	}
}

func (this *wu_manber) match(txt string) bool {
	if this.size == 0 {
		return false
	}
	//	log.Println("do-match", txt)
	txtlen := len(txt) // T[m-b] + L + T[m]
	ix := wu_manber_m - wu_manber_b
	for ix <= txtlen-wu_manber_b {
		key := this.wm_hash(txt, ix)
		if shift := int(this.shift[key]); shift > 0 {
			ix += shift
			//			log.Println("shift", txt, shift, ix)
		} else {
			//			log.Println("match-end", txt[ix:])
			patterns := this.patterns[key]
			start := ix + wu_manber_b - wu_manber_m
			prefix := txt[start : start+wu_manber_b] // prefix-length is wu_manber_b
			for _, pattern := range patterns {
				if pattern.prefix == prefix { // prefix match
					//					log.Println("do-pattern", pattern.pattern, pattern.prefix, txt, txt[ix:ix+2])
					leftp := pattern.pattern[len(prefix):]
					leftt := txt[start+len(prefix):] // do match
					if strings.HasPrefix(leftt, leftp) {
						return true
					}
				}
			}
			ix++
		}
	}
	return false
}

const prime_rk = 16777619

//rabin-k hash algorithm
func make_key(sep string) uint32 {
	l := len(sep)
	h := uint32(0)
	for i := 0; i < l; i++ {
		h = h*prime_rk + uint32(sep[i])
	}
	return h
}

func int_min(l, r int) int {
	if l < r {
		return l
	}
	return r
}
func new_wu_pattern(pattern string) wm_pattern {
	prefix := pattern[:2]
	return wm_pattern{pattern, prefix}
}
