package ext

type excel struct{}

func Excel() FileExter {
	return excel{}
}
func (excel) Valid(path string) (name, ext string, valid bool) {
	name, ext = GetExt(path)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".ico", ".svg":
		return name, ext, true
	default:
		return "", "", false
	}
}
