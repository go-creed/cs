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

func (e *Upload) FileMerge(ctx context.Context, request *uploadPb.ChunkResponse, response *uploadPb.FileMate) error {

	chunk, err := cache.ReadChunk(rd.Cache(), request.UploadId)
	if err != nil {
		return err
	} else if chunk == false {
		return nil
	}
	// TODO 合并之后返回 对应的文件路径等文件信息

	return nil
}

func (e *Upload) FileChunkLegitimate(ctx context.Context, request *uploadPb.ChunkResponse, response *uploadPb.ChunkResponse) error {
	log.Info("[Upload][FileChunkLegitimate]:Start...")
	mapUpload, err := cache.ReadMapUpload(rd.Cache(), request.UploadId)
	if err != nil {
		log.Errorf("[Upload][FileChunkLegitimate]: 缓存数据加载失败 %s", err.Error())
		return err
	}
	if !(mapUpload.GetChunkCount() >= request.ChunkCount &&
		mapUpload.GetFilesha256() == request.Filesha256 &&
		mapUpload.GetUploadId() == request.GetUploadId()) {
		err = errors.New("file's chunk verification failed！")
		log.Errorf("[Upload][FileChunkLegitimate]: 缓存数据校验失败 %s", err.Error())
		return err
	}
	*response = *request
	response.Size = mapUpload.GetChunkSize()
	log.Info("[Upload][FileChunkLegitimate]:End...")
	return nil
}

func (e *Upload) uploadId(userId int64) string {
	return fmt.Sprintf("MP_%d_%x", userId, time.Now().UnixNano())
}

func (e *Upload) FileChunk(ctx context.Context, request *uploadPb.ChunkRequest, response *uploadPb.ChunkResponse) (err error) {
	log.Info("[Upload][FileChunk]:Start...")
	response.Size = request.Size
	response.UploadId = e.uploadId(request.UserId)
	response.Filesha256 = request.Filesha256
	response.ChunkSize = chunkSize
	response.ChunkCount = int64(math.Ceil(float64(request.Size) / chunkSize))
	if err := cache.NewMapUpload().
		SetSize(response.Size).
		SetUploadId(response.UploadId).
		SetFilesha256(response.Filesha256).
		SetChunkSize(response.ChunkSize).
		SetChunkCount(response.ChunkCount).Write(rd.Cache(), response.UploadId); err != nil {
		log.Errorf("[Upload][FileChunk]:%s", err.Error())
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
		file     *os.File
		location string
		fileInfo uploadPb.FileRequest
		err      error
	)
	//recv the msg to create file
	if err := stream.RecvMsg(&fileInfo); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err.Error())
		return err
	}
	if file, location, err = uploadService.CreateFile(fileInfo.Filename); err != nil {
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
	log.Info("[Upload][SendBytes]:End...")

	upload, err := cache.ReadMapUpload(rd.Cache(), fileInfo.UploadId)
	if err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err)
		return err
	}
	if err = upload.SetChunk(rd.Cache(), fileInfo.UploadId, fileInfo.Index); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err)
		return err
	}
	return stream.SendMsg(&uploadPb.ChunkResponse{
		FileName: fileInfo.Filename,
		Size:     int64(size),
	})
}
