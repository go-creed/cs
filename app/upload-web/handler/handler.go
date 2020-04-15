package handler

import (
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
	Error error       `json:"error"`
	Msg   interface{} `json:"msg"`
}

func UploadImage(ctx *gin.Context) {
	log.Info("[Upload][Image]:Start")
	//从 字段 file 中读取文件
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
	fileInfo := uploadSrv.FileInfo{
		FileName: header.Filename,
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
