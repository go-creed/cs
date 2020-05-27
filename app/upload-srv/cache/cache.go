package cache

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

const (
	size        = "size"
	uploadId    = "uploadId"
	filesha256  = "filesha256"
	chunkSize   = "chunk_size"
	chunkCount  = "chunk_count"
	chunkPrefix = "chunk_"
)

type cache struct {
	sync.RWMutex
	data map[string]interface{}
}

func (um *cache) chunkField(index int64) string {
	return fmt.Sprintf("%s%d", chunkPrefix, index)
}

func (um *cache) SetChunk(rd *redis.Client, key string, index int64) error {
	um.set(um.chunkField(index), true)
	return rd.HSetNX(key, um.chunkField(index), true).Err()
}

func ReadChunk(rd *redis.Client, key string) (bool, error) {
	result, err := rd.HGetAll(key).Result()
	if err != nil {
		return false, err
	}
	var count int
	var countLimit int
	for k := range result {
		if strings.HasPrefix(k, chunkPrefix) {
			count++
		} else if k == chunkCount {
			countLimitStr, _ := result[k]
			countLimit, _ = strconv.Atoi(countLimitStr)
		}
	}
	return countLimit == count, nil
}

func ReadMapUpload(rd *redis.Client, key string) (*cache, error) {
	upload := NewMapUpload()
	result, err := rd.HMGet(key, size, uploadId, filesha256, chunkSize, chunkCount).Result()
	if err != nil {
		return upload, err
	}
	for _, v := range result {
		if v == nil {
			return upload, redis.Nil
		}
	}
	s, _ := strconv.Atoi(result[0].(string))
	cs, _ := strconv.Atoi(result[3].(string))
	cc, _ := strconv.Atoi(result[4].(string))
	upload.
		SetSize(int64(s)).
		SetUploadId(result[1].(string)).
		SetFilesha256(result[2].(string)).
		SetChunkSize(int64(cs)).
		SetChunkCount(int64(cc))
	return upload, nil
}

func (um *cache) Chunk() {

}

func (um *cache) Write(rd *redis.Client, key string) error {
	if err := rd.HMSet(key, um.GetMap()).Err(); err != nil {
		return err
	} else {
		if err := rd.Expire(key, time.Minute*60).Err(); err != nil {
			return err
		}
	}
	return nil
}

func NewMapUpload() *cache {
	upload := &cache{}
	upload.data = make(map[string]interface{})
	return upload
}

func (um *cache) GetMap() map[string]interface{} {
	return um.data
}

func (um *cache) set(k string, v interface{}) {
	um.Lock()
	defer um.Unlock()
	um.data[k] = v
}
func (um *cache) get(k string) interface{} {
	um.RLock()
	defer um.RUnlock()
	return um.data[k]
}

func (um *cache) SetSize(v int64) *cache {
	um.set(size, v)
	return um
}
func (um *cache) GetSize() int64 {
	s, _ := um.get(size).(int64)
	return s
}
func (um *cache) SetUploadId(v string) *cache {
	um.set(uploadId, v)
	return um
}
func (um *cache) GetUploadId() string {
	uid, _ := um.get(uploadId).(string)
	return uid
}
func (um *cache) SetFilesha256(v string) *cache {
	um.set(filesha256, v)
	return um
}
func (um *cache) GetFilesha256() string {
	f, _ := um.get(filesha256).(string)
	return f
}
func (um *cache) SetChunkSize(v int64) *cache {
	um.set(chunkSize, v)
	return um
}
func (um *cache) GetChunkSize() int64 {
	c, _ := um.get(chunkSize).(int64)
	return c
}
func (um *cache) SetChunkCount(v int64) *cache {
	um.set(chunkCount, v)
	return um
}
func (um *cache) GetChunkCount() int64 {
	c, _ := um.get(chunkCount).(int64)
	return c
}
