package cast

func StrSliceToPtr(s []string) *[]string {
	if s == nil {
		return nil
	}

	return &s
}
