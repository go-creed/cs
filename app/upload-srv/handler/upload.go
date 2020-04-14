package handler

import (
	"context"
	"os"

	uploadMd "cs/app/upload-srv/model/upload"
	uploadPb "cs/app/upload-srv/proto/upload"
	//"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	uploadService uploadMd.Service
)

func Init() {
	var err error
	uploadService, err = uploadMd.GetService()
	if err != nil {
		log.Fatal("[Upload] Handler Init Failure , %s", err)
		return
	}
}

type Upload struct{}

//如果是单文件上传，那么只需要将整个文件写入
func (e *Upload) WriteImage(ctx context.Context, stream uploadPb.Upload_WriteImageStream) error {
	log.Info("[Upload][SendBytes]:Start...")
	var (
		file     *os.File
		fileInfo uploadPb.FileInfo
		err      error
	)
	if err := stream.RecvMsg(&fileInfo); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err.Error())
		return err
	}
	log.Infof("文件名:%s，文件大小%d", fileInfo.FileName, fileInfo.Size)
	if file, err = uploadService.CreateFile(fileInfo.FileName); err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err.Error())
		return err
	}
	defer file.Close()
	for {
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
	hash, err := uploadService.Hash(file)
	if err != nil {
		log.Errorf("[Upload][SendBytes]:%s", err)
		return err
	}

	log.Info("[Upload][SendBytes]:End...")
	info := uploadPb.FileInfo{
		FileName:   fileInfo.FileName,
		Filesha256: hash,
		Size:       fileInfo.Size,
	}
	return stream.SendMsg(info)
}
