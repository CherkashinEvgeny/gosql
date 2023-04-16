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
}

type Result = sql.Result

type Rows = sql.Rows

var _ Executor = (*DB)(nil)

func (d *DB) Executor(ctx context.Context) (executor Executor) {
	tx, found := d.extractTxFromContext(ctx)
	if found {
		executor = tx
	} else {
		executor = d
	}
	return executor
}

func (d *DB) Exec(query string, arg ...any) (result Result, err error) {
	return (*sql.DB)(d).Exec(query, arg...)
}

func (d *DB) NamedExec(query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.Exec(query, posArgs...)
}

func (d *DB) ExecContext(ctx context.Context, query string, arg ...any) (result Result, err error) {
	return (*sql.DB)(d).ExecContext(ctx, query, arg...)
}

func (d *DB) NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.ExecContext(ctx, query, posArgs...)
}

func (d *DB) Query(query string, args ...any) (rows *Rows, err error) {
	return (*sql.DB)(d).Query(query, args...)
}

func (d *DB) NamedQuery(query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.Query(query, posArgs...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (rows *Rows, err error) {
	return (*sql.DB)(d).QueryContext(ctx, query, args...)
}

func (d *DB) NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.QueryContext(ctx, query, posArgs...)
}

var _ Executor = (*Tx)(nil)

func (t *Tx) Exec(query string, arg ...any) (result Result, err error) {
	return (*sql.Tx)(t).Exec(query, arg...)
}

func (t *Tx) NamedExec(query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.Exec(query, posArgs...)
}

func (t *Tx) ExecContext(ctx context.Context, query string, arg ...any) (result Result, err error) {
	return (*sql.Tx)(t).ExecContext(ctx, query, arg...)
}

func (t *Tx) NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.ExecContext(ctx, query, posArgs...)
}

func (t *Tx) Query(query string, args ...any) (rows *Rows, err error) {
	return (*sql.Tx)(t).Query(query, args...)
}

func (t *Tx) NamedQuery(query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.Query(query, posArgs...)
}

func (t *Tx) QueryContext(ctx context.Context, query string, args ...any) (rows *Rows, err error) {
	return (*sql.Tx)(t).QueryContext(ctx, query, args...)
}

func (t *Tx) NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.QueryContext(ctx, query, posArgs...)
}

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
