package op

import (
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	options := Options{}

	_ = options.Init(
		SetEs(Es{Index: "xxx"}),
		SetDisk(Disk{Path: "path"}),
	)
	fmt.Println(options)
}
