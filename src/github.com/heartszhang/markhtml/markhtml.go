package markhtml

import (
	"fmt"
	"strings"
)

const (
	_ = iota
	token_type_text
	token_type_elementstart
	token_type_elementend
	token_type_elementclose
	token_type_br
	token_type_eos
)

type token struct {
	token_type int
	data       string
}

//[img][size][/xxx
func TransferText(txt string) string {
	var out string
	var transfered int
	chars := []rune(strings.ToLower(txt))
	for token, i := next_token(chars, 0); token.token_type != token_type_eos; token, i = next_token(chars, i) {
		switch token.token_type {
		case token_type_text:
			out += token.data
		case token_type_br:
			out += "<br/>"
		case token_type_elementstart:
			sub, xi := expect_element_text(chars, i)
			var nn string
			switch token.data {
			case "img":
				nn = `<img src=%q/>`
			case "url":
				nn = `<a href=%q/>`
			case "wmv":
				nn = `<embed src=%q/>`
			}
			if nn != "" {
				i = xi
				transfered++
				out += fmt.Sprintf(nn, sub.data)
			}
		}
	}
	if transfered > 0 {
		return out
	}
	return txt
}

func next_token(chars []rune, start int) (token, int) {
	switch {
	case start >= len(chars):
		return token{token_type: token_type_eos}, start
	case chars[start] == '[':
		return expect_element(chars, start+1)
		//	case chars[start] == 0xd:
		//		fallthrough
		//	case chars[start] == 0xa:
		//		return expect_carrage(chars, start)
	default:
		return expect_element_text(chars, start)
	}
}

func expect_carrage(chars []rune, start int) (token, int) {
	for start < len(chars) && (chars[start] == 0xd || chars[start] == 0xa) {
		start++
	}
	return token{token_type: token_type_br}, start
}
func expect_element_text(chars []rune, start int) (token, int) {
	var s string
	for start < len(chars) {
		if chars[start] == '[' {
			break
		}
		s += string(chars[start])
		start++
	}
	return token{token_type: token_type_text, data: s}, start
}

func expect_element(chars []rune, start int) (token, int) {
	if start >= len(chars) {
		return token{data: "[", token_type: token_type_text}, start
	}
	s := start
	t := token_type_elementstart
	if chars[start] == '/' {
		s += 1
		t = token_type_elementend
	}
	if d, s, ok := expect_right_bracket(chars, s); ok {
		return token{token_type: t, data: d}, s
	}
	return token{data: "[", token_type: token_type_text}, start
}
func expect_right_bracket(chars []rune, start int) (string, int, bool) {
	if start >= len(chars) {
		return "", start, false
	}
	var s string
	for start < len(chars) {
		if chars[start] == ']' {
			start++
			break
		}
		s += string(chars[start])
		start++
	}
	return s, start, start < len(chars)
}
