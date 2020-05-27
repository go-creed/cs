package handler

import "fmt"

func filepath(uploadId string, filename string, index int64) string {
	return fmt.Sprintf("%d_%s_%s", index, filename, uploadId)
}
