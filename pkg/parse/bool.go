package parse

import "strconv"

func ParseBoolOrFalse(s string) bool {
	if b, err := strconv.ParseBool(s); err != nil {
		return false
	} else {
		return b
	}
}

func ParseBoolOrTrue(s string) bool {
	if b, err := strconv.ParseBool(s); err != nil {
		return true
	} else {
		return b
	}
}
