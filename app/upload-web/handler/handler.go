package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"

	authPb "cs/app/auth-srv/proto/auth"
	uploadPb "cs/app/upload-srv/proto/upload"
	_const "cs/public/const"
	"cs/public/ecode"
	"cs/public/gin-middleware"
)

var (
	uploadClient uploadPb.UploadService
	authClient   authPb.AuthService
)

func Init() {
	uploadClient = uploadPb.NewUploadService(_const.UploadSrv, client.DefaultClient)
	authClient = authPb.NewAuthService(_const.AuthSrv, client.DefaultClient)
}

func FileDetail(ctx *middleware.MicroContext) {
	filesha256, b := ctx.GetQuery("filesha256")
	if filesha256 == "" || !b {
		middleware.ServerError(ctx, middleware.Response{
			Error: errors.New("field filesha256 mustn't empty"),
		})
		return
	}
	detail, err := uploadClient.FileDetail(ctx, &uploadPb.FileMate{
		Filesha256: filesha256,
	})
	if err != nil {
		middleware.ServerError(ctx, middleware.Response{
			Error: err,
		})
		return
	}
	middleware.Success(ctx, middleware.Response{
		Data: detail,
	})
}

func FileChunk(ctx *middleware.MicroContext) {
	log.Info("[Chunk][File]:Start")
	//接收一个文件大小就行了
	var (
		params struct {
			FileSize   int64  `json:"file_size" form:"file_size" validate:"min=1,max=10"`
			Filesha256 string `json:"filesha_256" form:"filesha_256"`
		}
		err error
	)
	if err = ctx.ShouldBind(&params); err != nil {
		log.Errorf("[Chunk][File]:解析参数失败 %s", err)
		middleware.RequestError(ctx, middleware.Response{})
		return
	}
	if params.FileSize <= 0 {
		err = errors.New("session verification failed")
		middleware.RequestError(ctx, middleware.Response{
			Error: err,
		})
		return
	}
	if chunk, err := uploadClient.FileChunk(ctx, &uploadPb.ChunkRequest{
		Filesha256: params.Filesha256,
		UserId:     ctx.UserId,
		Size:       params.FileSize,
	}); err != nil {
		middleware.ServerError(ctx, middleware.Response{
			Error: err,
		})
	} else {
		middleware.Success(ctx, middleware.Response{Data: chunk})
	}
	log.Info("[Chunk][File]:End")
	return
}

func FileUpload(ctx *middleware.MicroContext) {
	log.Info("[Upload][File]:Start")
	//从 前端直接接受文件的hash值，跟其它一些东西合并看是否能够直接返回
	var (
		params struct {
			UploadId   string `json:"upload_id" form:"upload_id" validate:"ne=''"`
			Filesha256 string `json:"filesha256" form:"filesha256" validate:"ne=''"`
			Offset     int64  `json:"offset" form:"offset" validate:"gt=0"`
		}
		fileMate *uploadPb.FileMate
		err      error
	)
	ctx.Bind(&params)
	if err = validator.New().Struct(params); err != nil {
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrInternalServer,err)})
		return
	}
	// You need to first determine whether the file already exists
	if params.Filesha256 != "" {
		fileMate, err = uploadClient.FileDetail(ctx, &uploadPb.FileMate{
			Filesha256: params.Filesha256,
		})
		if err != nil {
			middleware.ServerError(ctx, middleware.Response{Error: err})
			return
		}
		if fileMate.Id != 0 {
			log.Infof("[Upload][File]:文件 %s ，filesha256 %s 已经存在，停写入", fileMate.Filename, fileMate.Filesha256)
			middleware.ServerError(ctx, middleware.Response{Error: err})
			return
		}
	}
	// Read file from stream
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Errorf("[Upload][Image]:读取文件失败 %s", err)
		middleware.ServerError(ctx, middleware.Response{Error: err})
		return
	}
	uploadClient.FileChunkLegitimate(ctx, &uploadPb.ChunkResponse{
		Filesha256: params.Filesha256,
		UploadId:   params.UploadId,
		ChunkCount: params.Offset,
	})
	log.Infof("[Upload][Image]:size:%d,fileName:%s", header.Size, header.Filename)
	//创建一个发送字节数据包的连接通道
	sendBytes, err := uploadClient.WriteImage(ctx)
	if err != nil {
		log.Errorf("[Upload][Image]: 创建通道失败 %s", err)
		return
	}
	defer sendBytes.Close()
	fileInfo := uploadPb.FileMate{
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
			err := sendBytes.Send(&uploadPb.Bytes{Content: x, Size: header.Size})
			if err != nil {
				log.Errorf("[Upload][Image]:数据发送失败 %s", err)
				return
			}
		}
		isClose <- struct{}{}
		if err = sendBytes.Send(&uploadPb.Bytes{
			Content: nil,
		}); err != nil {
			log.Errorf("[Upload][Image]:数据发送失败 %s", err)
			return
		}
	}()
	<-isClose
	if err := sendBytes.RecvMsg(&fileInfo); err != nil {
		log.Errorf("[Upload][Image]:数据发送失败 %s", err)
		middleware.ServerError(ctx, middleware.Response{Error: err})
		return
	}
	ctx.JSONP(200, fileInfo)
}
