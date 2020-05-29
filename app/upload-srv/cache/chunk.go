package cache

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	_fileName   = "file_name"
	_size       = "size"
	_count      = "count"
	_filesha256 = "filesha256"
	_chunk      = "chunk_"
	_count_left = _count + "_left"
)

type ChunkSize struct {
	Index int64 `json:"index"`
	Size  int64 `json:"size"`
}
type Chunk struct {
	FileName   string      `json:"file_name"`
	Size       int64       `json:"size"`
	Count      int64       `json:"count"`
	Filesha256 string      `json:"filesha_256"`
	Chunks     []ChunkSize `json:"chunk_size"`
}

func keyChunkSize(index int64) string {
	return fmt.Sprintf("%s%d", _chunk, index)
}
func UpdateIndex(rd *redis.Client, uploadId string, index int64, remainingSize int64) error {
	return rd.HSet(uploadId, keyChunkSize(index), remainingSize).Err()
}

func (c Chunk) Index(index int64) ChunkSize {
	for _, chunk := range c.Chunks {
		if chunk.Index == index {
			return chunk
		}
	}
	return ChunkSize{}
}

func ReadChunk(rd *redis.Client, uploadId string) (c Chunk, err error) {
	result, err := rd.HGetAll(uploadId).Result()
	if err != nil {
		return
	}
	c.FileName, _ = result[_fileName]
	c.Size, _ = strconv.ParseInt(result[_size], 10, 64)
	c.Count, _ = strconv.ParseInt(result[_count], 10, 64)
	c.Filesha256, _ = result[_filesha256]
	for k, v := range result {
		if strings.HasPrefix(k, _chunk) {
			index, _ := strconv.ParseInt(k[len(_chunk):], 10, 64)
			size, _ := strconv.ParseInt(v, 10, 64)
			c.Chunks = append(c.Chunks, ChunkSize{
				Index: index,
				Size:  size,
			})
		}
	}
	return
}

func (c Chunk) Write(rd *redis.Client, uploadId string) error {
	m := map[string]interface{}{
		_size:       c.Size,
		_fileName:   c.FileName,
		_count_left: c.Count,
		_count:      c.Count,
		_filesha256: c.Filesha256,
	}
	for _, chunk := range c.Chunks {
		m[fmt.Sprintf("%s%d", _chunk, chunk.Index)] = chunk.Size
	}
	pipeline := rd.Pipeline()
	{
		pipeline.HMSet(uploadId, m)
		pipeline.Expire(uploadId, time.Minute*60*24)
	}
	if _, err := pipeline.Exec(); err != nil {
		return err
	}
	return nil
}
