package cache

import (
	"fmt"
	"testing"

	"cs/plugin/rd"
)

func TestNewChunk(t *testing.T) {

}
func TestReadChunk(t *testing.T) {
	rd.Init()
	cache := rd.Cache()
	chunk, err := ReadChunk(cache, "1")
	if err!=nil{
		t.Fatal(err)
		return
	}
	fmt.Printf("%+v",chunk)
}
func TestChunk_Write(t *testing.T) {
	//测试文件分块信息缓存写入
	rd.Init()
	cache := rd.Cache()
	chunk := Chunk{
		FileName:   "文件名.txt",
		Size:       200,
		Count: 10,
		Filesha256: "sdfrgs23429",
		Chunks: []ChunkSize{
			{
				Index: 1,
				Size:  100,
			},
			{
				Index: 2,
				Size:  100,
			},
		},
	}
	if err := chunk.Write(cache, "1"); err != nil {
		t.Fatal(err)
		return
	}
}
