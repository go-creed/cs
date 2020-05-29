package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"cs/app/upload-srv/cache"
	uploadMd "cs/app/upload-srv/model/upload"
	uploadPb "cs/app/upload-srv/proto/upload"
	"cs/plugin/db"
	"cs/plugin/rd"
	//"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	once          sync.Once
	uploadService uploadMd.Service
)

const (
	chunkSize = 5 * 1024 * 1024
)

func Init() {
	var err error
	once.Do(func() {
		uploadService, err = uploadMd.GetService()
		if err != nil {
			log.Fatal("[Upload] Handler Init Failure , %s", err)
			return
		}
	})
}

type Upload struct{}

func (e *Upload) FileMerge(ctx context.Context, request *uploadPb.MergeRequest, mate *uploadPb.FileMate) error {
	log.Info("[Upload][FileMerge]:Start...")
	chunk, err := cache.ReadChunk(rd.Cache(), request.UploadId)
	if err != nil {
		log.Errorf("[Upload][FileMerge]: 缓存数据加载失败 %s", err.Error())
		return err
	}
	for _, chunk := range chunk.Chunks {
		if chunk.Size != 0 {
			err = errors.New("some block files are incomplete")
			log.Errorf("[Upload][FileMerge]: 分块验证 %s", err.Error())
			return err
		}
	}

	if err = uploadService.
		MergeFile(
			fmt.Sprintf("/%d/%s/%s",
				request.UserId, request.UploadId, chunk.FileName),
			chunk.Filesha256,
		); err != nil {
		log.Errorf("[Upload][FileMerge]: 文件合并失败 %s", err.Error())
		return err
	}

	//// TODO 合并之后返回 对应的文件路径等文件信息
	//uploadService.MergeFile()
	return nil
}

func (e *Upload) FileChunkVerify(ctx context.Context, request *uploadPb.ChunkRequest, response *uploadPb.ChunkResponse) error {
	log.Info("[Upload][FileChunkVerify]:Start...")
	chunk, err := cache.ReadChunk(rd.Cache(), request.UploadId)
	if err != nil {
		log.Errorf("[Upload][FileChunkVerify]: 缓存数据加载失败 %s", err.Error())
		return err
	}
	if !(chunk.Filesha256 == request.Filesha256 &&
		chunk.FileName == request.FileName &&
		chunk.Count >= request.Index) {
		err = errors.New("file's chunk verification failed！")
		log.Errorf("[Upload][FileChunkVerify]: 缓存数据校验失败 %s", err.Error())
		return err
	}

	response.Size = chunk.Index(request.Index).Size
	log.Info("[Upload][FileChunkVerify]:End...")
	return nil
}

func (e *Upload) uploadId(userId int64) string {
	return fmt.Sprintf("CHUNK_%d_%x", userId, time.Now().UnixNano())
}

func (e *Upload) FileChunk(ctx context.Context, request *uploadPb.ChunkRequest, response *uploadPb.ChunkResponse) (err error) {
	log.Info("[Upload][FileChunk]:Start...")
	response.Size = request.Size
	response.UploadId = e.uploadId(request.UserId)
	response.Filesha256 = request.Filesha256
	response.ChunkSize = chunkSize
	response.ChunkCount = int64(math.Ceil(float64(request.Size) / chunkSize))
	response.FileName = request.FileName

	chunk := cache.Chunk{
		FileName:   response.FileName,
		Size:       response.Size,
		Count:      response.ChunkCount,
		Filesha256: response.Filesha256,
	}

	chunk.Chunks = make([]cache.ChunkSize, response.ChunkCount)
	for i := int64(0); i < response.ChunkCount; i++ {
		chunk.Chunks[i].Index = i + 1
		if i == response.ChunkCount-1 {
			chunk.Chunks[i].Size = response.Size - (i)*chunkSize
		} else {
			chunk.Chunks[i].Size = chunkSize
		}
	}

	if err = chunk.Write(rd.Cache(), response.UploadId); err != nil {
		return err
	}
	log.Info("[Upload][FileDetail]:End...")
	return nil
}

// Return file mate of the file
func (e *Upload) FileDetail(ctx context.Context, info *uploadPb.FileMate, info2 *uploadPb.FileMate) error {
	log.Info("[Upload][FileDetail]:Start...")
	err := uploadService.FileDetail(db.DB(), info)
	if err != nil {
		log.Errorf("[Upload][FileDetail]:%s", err.Error())
		return err
	}
	log.Info("[Upload][FileDetail]:End...")
	fmt.Println(info2, info)
	*info2 = *info
	return nil
}

func (e *Upload) WriteBytes(ctx context.Context, stream uploadPb.Upload_WriteBytesStream) error {
	log.Info("[Upload][SendBytes]:Start...")
	var (
		file      *os.File
		location  string
		chunkInfo uploadPb.ChunkRequest
		err       error
	)
	//recv the msg to create file
	if err := stream.RecvMsg(&chunkInfo); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err.Error())
		return err
	}
	if file, location, err = uploadService.CreateFile(fmt.Sprintf("%d/%s/%s", chunkInfo.UserId, chunkInfo.UploadId, chunkInfo.FileName)); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err.Error())
		return err
	}
	defer file.Close()
	var size int
	for {
		recv, err := stream.Recv()
		if err != nil {
			log.Errorf("[Upload][SendBytes]:%s", err.Error())
			return err
		}
		if recv.Content == nil {
			break
		}
		size += len(recv.Content)
		err = uploadService.Write(file, recv.Content)
		if err != nil {
			log.Errorf("[Upload][SendBytes]:%s", err)
			return err
		}
	}
	fmt.Println(location)

	if err = cache.UpdateIndex(
		rd.Cache(),
		chunkInfo.UploadId,
		chunkInfo.Index,
		chunkInfo.Size-int64(size)); err != nil {
		log.Errorf("[Upload][SendBytes]: 更新当前分块下载进度%s", err)
		return err
	}
	log.Info("[Upload][SendBytes]:End...")

	return stream.SendMsg(&uploadPb.ChunkResponse{
		FileName: chunkInfo.FileName,
		Size:     int64(size),
	})
}
