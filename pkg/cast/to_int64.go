package cast

import "strconv"

func StrToInt64(s string) int64 {
	if s == "" {
		return 0
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func StrToInt64Ptr(s string) *int64 {
	if s == "" {
		return nil
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}

	return &i
}
