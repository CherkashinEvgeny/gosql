package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Executor interface {
	Exec(query string, args ...any) (result Result, err error)
	NamedExec(query string, args map[string]any) (result Result, err error)

	ExecContext(ctx context.Context, query string, args ...any) (result Result, err error)
	NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error)

	Query(query string, args ...any) (rows *Rows, err error)
	NamedQuery(query string, args map[string]any) (rows *Rows, err error)

	QueryContext(ctx context.Context, query string, args ...any) (rows *Rows, err error)
	NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error)

	QueryRow(query string, args ...any) (row *Row)
	NamedQueryRow(query string, args map[string]any) (row *Row)

	QueryRowContext(ctx context.Context, query string, args ...any) (row *Row)
	NamedQueryRowContext(ctx context.Context, query string, args map[string]any) (row *Row)
}

type Result = sql.Result

type Rows = sql.Rows

type Row = sql.Row

func convertToPositionalParameters(str string, args map[string]any) (string, []any) {
	tokens := tokenize(str, tokensMap)
	sb := strings.Builder{}
	positionArgs := make([]any, 0, len(args))
	index := 0
	for index < len(tokens) {
		var argId string
		var ok bool
		var data string
		argId, ok, index = readId(tokens, index)
		if ok {
			arg, found := args[argId]
			if found {
				data = cfg.Placeholder(len(positionArgs))
				positionArgs = append(positionArgs, arg)
			} else {
				data = fmt.Sprintf(":%s", argId)
			}
		} else {
			data, _, index = readToken(tokens, index)
		}
		sb.WriteString(data)
	}
	return sb.String(), positionArgs
}

func readId(tokens []string, startIndex int) (string, bool, int) {
	index := startIndex
	if index >= len(tokens) {
		return "", false, startIndex
	}
	if tokens[index] != idPrefix {
		return "", false, startIndex
	}
	index++
	for {
		if index >= len(tokens) {
			return strings.Join(tokens[startIndex+1:index], ""), true, index
		}
		token := tokens[index]
		if !isIdToken(token) {
			return strings.Join(tokens[startIndex+1:index], ""), true, index
		}
		index++
	}
}

type Placeholder = func(index int) string

func readToken(tokens []string, index int) (string, bool, int) {
	if index >= len(tokens) {
		return "", false, index
	}
	return tokens[index], false, index + 1
}

func isIdToken(token string) bool {
	r, size := utf8.DecodeRuneInString(token)
	if size != len(token) {
		return false
	}
	return !unicode.IsSpace(r)
}

const idPrefix = "idPrefix"

var tokensMap = map[rune]string{
	':': idPrefix,
}

func tokenize(str string, tokensMap map[rune]string) []string {
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
			token, found := tokensMap[r]
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
