package cast

func StrToStringPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
