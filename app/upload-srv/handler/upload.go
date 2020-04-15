package handler

import (
	"context"
	"os"
	"sync"

	uploadMd "cs/app/upload-srv/model/upload"
	uploadPb "cs/app/upload-srv/proto/upload"
	"cs/plugin/db"

	//"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	once          sync.Once
	uploadService uploadMd.Service
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

// Return file mate of the file
func (e *Upload) FileDetail(ctx context.Context, info *uploadPb.FileMate, info2 *uploadPb.FileMate) error {
	log.Info("[Upload][FileDetail]:Start...")
	err := uploadService.FileDetail(db.DB(), info)
	if err != nil {
		log.Errorf("[Upload][FileDetail]:%s", err.Error())
		return err
	}
	log.Info("[Upload][FileDetail]:End...")
	*info2 = *info
	return nil
}

// Single file write
func (e *Upload) WriteImage(ctx context.Context, stream uploadPb.Upload_WriteImageStream) error {
	log.Info("[Upload][SendBytes]:Start...")
	var (
		file     *os.File
		location string
		fileInfo uploadPb.FileMate
		err      error
	)
	//recv the msg to create file
	if err := stream.RecvMsg(&fileInfo); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err.Error())
		return err
	}
	log.Infof("文件名:%s，文件大小%d", fileInfo.Filename, fileInfo.Size)
	if file, location, err = uploadService.CreateFile(fileInfo.Filename); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err.Error())
		return err
	}
	defer file.Close()
	for {
		// Recv the msg
		recv, err := stream.Recv()
		if err != nil {
			log.Errorf("[Upload][SendBytes]:%s", err.Error())
			return err
		}
		if recv.Content == nil {
			break
		}
		err = uploadService.Write(file, recv.Content)
		if err != nil {
			log.Errorf("[Upload][SendBytes]:%s", err)
			return err
		}
	}
	// Read file to hash
	hash, err := uploadService.Hash(file)
	if err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err)
		return err
	}
	// Assembly file meta
	info := uploadPb.FileMate{
		Filename:   fileInfo.Filename,
		Filesha256: hash,
		Size:       fileInfo.Size,
		Location:   location,
	}
	// Write the file mate to mysql
	if err = uploadService.WriteDB(db.DB(), &info); err != nil {
		log.Error("[Upload][SendBytes]:%s", err)
		return err
	}
	log.Info("[Upload][SendBytes]:End...")
	// Notify the file mate
	return stream.SendMsg(&info)
}
