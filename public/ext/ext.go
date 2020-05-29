package ext

import (
	"path/filepath"
	"strings"
)

type FileExter interface {
	Valid(path string) (name, ext string, valid bool)
}

func GetExt(path string) (name, ext string) {
	base := filepath.Base(path)
	split := strings.Split(base, ".")
	switch len(split) {
	case 1:
		return split[0], ""
	case 2:
		return split[0], split[1]
	}
	return
}
