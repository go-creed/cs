package ext

type image struct{}

func Image() FileExter {
	return image{}
}
func (image) Valid(path string) (name, ext string, valid bool) {
	name, ext = GetExt(path)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".ico", ".svg":
		return name, ext, true
	default:
		return "", "", false
	}
}
