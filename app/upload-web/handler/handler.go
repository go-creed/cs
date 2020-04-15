package handler

import (
	"errors"
	"net/http"

	uploadSrv "cs/app/upload-srv/proto/upload"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	uploadClient uploadSrv.UploadService
)

func Init() {
	uploadClient = uploadSrv.NewUploadService("go.micro.cs.service.upload", client.DefaultClient)
}

type JSONP struct {
	Error error       `json:"error,omitempty"`
	Msg   interface{} `json:"msg,omitempty"`
}

func FileDetail(ctx *gin.Context) {
	filesha256, b := ctx.GetQuery("filesha256")
	if filesha256 == "" || !b {
		ctx.JSONP(http.StatusBadRequest, JSONP{Error: errors.New("field filesha256 mustn't empty")})
		return
	}
	detail, err := uploadClient.FileDetail(ctx, &uploadSrv.FileMate{
		Filesha256: filesha256,
	})
	if err != nil {
		ctx.JSONP(http.StatusInternalServerError, JSONP{Error: err})
		return
	}
	ctx.JSONP(200, JSONP{Msg: detail})
}
func FileUpload(ctx *gin.Context) {
	log.Info("[Upload][Image]:Start")
	//从 前端直接接受文件的hash值，跟其它一些东西合并看是否能够直接返回
	var (
		filesha256, _ = ctx.GetPostForm("filesha256")
		fileMate      *uploadSrv.FileMate
		err           error
	)
	// You need to first determine whether the file already exists
	if filesha256 != "" {
		fileMate, err = uploadClient.FileDetail(ctx, &uploadSrv.FileMate{
			Filesha256: filesha256,
		})
		if err != nil {
			ctx.JSONP(http.StatusInternalServerError, JSONP{Error: err})
			return
		}
		if fileMate.Id != 0 {
			log.Infof("[Upload][File]:文件 %s ，filesha256 %s 已经存在，停写入", fileMate.Filename, fileMate.Filesha256)
			ctx.JSONP(http.StatusOK, JSONP{Msg: "The file already exists！"})
			return
		}
	}
	// Read file from stream
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Errorf("[Upload][Image]:读取文件失败 %s", err)
		ctx.JSONP(http.StatusInternalServerError, JSONP{Error: err})
		return
	}
	log.Infof("[Upload][Image]:size:%d,fileName:%s", header.Size, header.Filename)
	//创建一个发送字节数据包的连接通道
	sendBytes, err := uploadClient.WriteImage(ctx)
	if err != nil {
		log.Errorf("[Upload][Image]: 创建通道失败 %s", err)
		return
	}
	defer sendBytes.Close()
	fileInfo := uploadSrv.FileMate{
		Filename: header.Filename,
		Size:     header.Size,
	}
	if err := sendBytes.SendMsg(&fileInfo); err != nil {
		log.Errorf("[Upload][Image]: 发送文件信息失败 %s", err)
		return
	}
	//设置通道的容量10块
	b := make(chan []byte, 10)
	go func() {
		var (
			n int
		)
		for {
			bt := make([]byte, 1024*1024) //开辟一个大小1m的缓冲空间
			switch n, err = file.Read(bt); true {
			case n < 0:
				log.Errorf("[Upload][Image]:数据读取失败 %s", err)
				fallthrough
			case n == 0:
				close(b)
				log.Errorf("[Upload][Image]:数据读取结束")
				return
			case n > 0:
				b <- bt
			}
		}
	}()

	isClose := make(chan struct{}) //无缓冲通道，等待输入某个东西来关闭
	go func() {
		for x := range b {
			err := sendBytes.Send(&uploadSrv.Bytes{Content: x, Size: header.Size})
			if err != nil {
				log.Errorf("[Upload][Image]:数据发送失败 %s", err)
				return
			}
		}
		isClose <- struct{}{}
		if err = sendBytes.Send(&uploadSrv.Bytes{
			Content: nil,
		}); err != nil {
			log.Errorf("[Upload][Image]:数据发送失败 %s", err)
			return
		}
	}()
	<-isClose
	if err := sendBytes.RecvMsg(&fileInfo); err != nil {
		log.Errorf("[Upload][Image]:数据发送失败 %s", err)
		ctx.JSONP(http.StatusInternalServerError, JSONP{Error: err})
		return
	}
	ctx.JSONP(200, fileInfo)
}
