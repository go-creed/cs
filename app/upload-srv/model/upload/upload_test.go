package upload

import (
	"fmt"
	"testing"
)

func TestUploadExt(t *testing.T) {
	Init()
	ext, b := s.imagesExt("sdf.png")
	fmt.Print(ext,b)
}
