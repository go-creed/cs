package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

const (
	size       = "size"
	uploadId   = "uploadId"
	filesha256 = "filesha256"
	chunkSize  = "chunk_size"
	chunkCount = "chunk_count"
)

type cache struct {
	sync.RWMutex
	data map[string]interface{}
}

func ReadMapUpload(rd *redis.Client, key string) *cache {
	upload := NewMapUpload()
	val := rd.HMGet(key).Val()
	fmt.Println(val)
	return upload
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
func (um *cache) SetChunkSize(v interface{}) *cache {
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
