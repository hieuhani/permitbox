package strutil

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unsafe"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

func SnakeToCamel(snakeCase string) string {
	words := strings.Split(snakeCase, "_")
	tc := cases.Title(language.English)
	if len(words) > 1 {
		for i := 1; i < len(words); i++ {
			words[i] = tc.String(words[i])
		}
	}
	return strings.Join(words, "")
}

func Slugify(s string) string {
	var buf bytes.Buffer

	for _, r := range s {
		switch {
		case r > unicode.MaxASCII:
			continue
		case unicode.IsLetter(r):
			buf.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r), r == '_', r == '-':
			buf.WriteRune(r)
		case unicode.IsSpace(r):
			buf.WriteRune('-')
		}
	}

	return buf.String()
}

func FormatInt(i any) (string, error) {
	n, err := toInt64(i)
	if err != nil {
		return "", err
	}

	return printer.Sprintf("%d", n), nil
}

func FormatFloat(f float64, dp int) string {
	format := "%." + strconv.Itoa(dp) + "f"
	return printer.Sprintf(format, f)
}

func toInt64(i any) (int64, error) {
	switch v := i.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	// Note: uint64 not supported due to risk of truncation.
	case string:
		return strconv.ParseInt(v, 10, 64)
	}

	return 0, fmt.Errorf("unable to convert type %T to int", i)
}

// StringToBytes converts a string to a byte slice without a memory allocation, the returned slice must not be modified.
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts a byte slice to a string without a memory allocation, the returned string must not be modified.
func BytesToString(s []byte) string {
	if len(s) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(s), len(s))
}
