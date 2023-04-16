package sql

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Config struct {
	Placeholder         Placeholder
	LogicalErrorHandler func(err error)
}

var cfg = Config{
	Placeholder:         Dollar,
	LogicalErrorHandler: PrintLogicalError,
}

func Configure(f func(Config) Config) {
	cfg = f(cfg)
}

func Dollar(index int) string {
	return "$" + strconv.Itoa(index+1)
}

func Question(index int) string {
	return "?"
}

func Colon(index int) string {
	return ":" + strconv.Itoa(index+1)
}

func AtP(index int) string {
	return "@p" + strconv.Itoa(index+1)
}

func PrintLogicalError(err error) {
	stackTrace := trace(2)
	header := fmt.Sprintf("DIRTY LOGIC: %s", err.Error())
	body := make([]string, 0, 2*len(stackTrace))
	for _, frame := range stackTrace {
		body = append(body, frame.Func())
		body = append(body, fmt.Sprintf("\t%s:%d", frame.File(), frame.Line()))
	}
	fmt.Println(box(header, body))
}

func box(header string, body []string) string {
	size := utf8.RuneCountInString(header) + 4
	for _, str := range body {
		sizeCandidate := utf8.RuneCountInString(str) + 2
		if sizeCandidate > size {
			size = sizeCandidate
		}
	}
	if size%2 == 1 {
		size++
	}
	sb := strings.Builder{}
	count := size - utf8.RuneCountInString(header)%2
	for index := 0; index < count; index++ {
		sb.WriteString("*")
	}
	sb.WriteString(header)
	for index := 0; index < count; index++ {
		sb.WriteString("*")
	}
	for _, str := range body {
		sb.WriteString("*")
		sb.WriteString(str)
		sb.WriteString("*")
	}
	for index := 0; index < size; index++ {
		sb.WriteString("*")
	}
	return sb.String()
}

func PanicOnLogicalError(err error) {
	panic(err)
}
