package lexer

import "regexp"

var (
	integerRegexp = regexp.MustCompile(`^0|[1-9]\d*$`)
	floatRegexp   = regexp.MustCompile(`^0|[1-9]\d*\.\d+$`)
	stringRegexp  = regexp.MustCompile(`"([^"])*"`)
)

func isInteger(raw string) bool {
	return integerRegexp.MatchString(raw)
}

func isFloat(raw string) bool {
	return floatRegexp.MatchString(raw)
}

func isNumber(raw string) bool {
	return isInteger(raw) || isFloat(raw)
}

func isString(raw string) bool {
	return stringRegexp.MatchString(raw)
}
