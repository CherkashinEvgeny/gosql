package sqlx

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func params(str string) []string {
	var result []string
	tokens := tokenize(str, brackets)
	index := 0
	for index < len(tokens) {
		var param string
		var ok bool
		param, ok, index = readId(tokens, index)
		if ok {
			result = append(result, param)
		} else {
			_, _, index = readToken(tokens, index)
		}
	}
	return result
}

func format(str string, params map[string]any) string {
	tokens := tokenize(str, brackets)
	sb := strings.Builder{}
	index := 0
	for index < len(tokens) {
		var paramId string
		var ok bool
		var data string
		paramId, ok, index = readId(tokens, index)
		if ok {
			param, found := params[paramId]
			if found {
				data = fmt.Sprintf("%v", param)
			} else {
				data = fmt.Sprintf("{%s}", paramId)
			}
		} else {
			data, _, index = readToken(tokens, index)
		}
		sb.WriteString(data)
	}
	return sb.String()
}

func readId(tokens []string, startIndex int) (string, bool, int) {
	index := startIndex
	if index >= len(tokens) {
		return "", false, startIndex
	}
	if tokens[index] != openBracket {
		return "", false, startIndex
	}
	index++
	for {
		if index >= len(tokens) {
			return "", false, startIndex
		}
		token := tokens[index]
		if token == closeBracket {
			return strings.Join(tokens[startIndex+1:index], ""), true, index + 1
		}
		index++
	}
}

func readToken(tokens []string, index int) (string, bool, int) {
	if index >= len(tokens) {
		return "", false, index
	}
	return tokens[index], false, index + 1
}

const (
	openBracket  = "ob"
	closeBracket = "cb"
)

var brackets = map[rune]string{
	'{': openBracket,
	'}': closeBracket,
}

func tokenize(str string, escapeeMap map[rune]string) []string {
	tokens := make([]string, 0, utf8.RuneCountInString(str))
	index := 0
	escapee := false
	for index < len(str) {
		r, n := utf8.DecodeRuneInString(str[index:])
		switch {
		case !escapee && r == '\\':
			escapee = true
		case escapee:
			tokens = append(tokens, str[index:index+n])
			escapee = false
		default:
			token, found := escapeeMap[r]
			if found {
				tokens = append(tokens, token)
			} else {
				tokens = append(tokens, str[index:index+n])
			}
		}
		index += n
	}
	return tokens
}
