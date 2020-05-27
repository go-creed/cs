package handler

import (
	"errors"
	"io"

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
			Index      int64  `json:"index" form:"index" validate:"gt=0"`
			FileName   string `json:"file_name" form:"file_name" validate:"ne=''"`
		}
		fileMate *uploadPb.FileMate
		err      error
	)
	_ = ctx.Bind(&params)
	if err = validator.New().Struct(params); err != nil {
		log.Errorf("[Upload][File]: 参数解析失败 %s", err)
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrRequestServer, err)})
		return
	}
	// You need to first determine whether the file already exists
	fileMate, err = uploadClient.FileDetail(ctx, &uploadPb.FileMate{
		Filesha256: params.Filesha256,
	})
	if err != nil {
		log.Errorf("[Upload][File]: 文件详情获取失败 %s", err)
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrInternalServer, err)})
		return
	}
	if fileMate.Id != 0 {
		log.Errorf("[Upload][File]:文件 %s ，filesha256 %s 已经存在，停写入", fileMate.Filename, fileMate.Filesha256)
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrInternalServer, err)})
		return
	}

	legitimate, err := uploadClient.FileChunkLegitimate(ctx, &uploadPb.ChunkResponse{
		Filesha256: params.Filesha256,
		UploadId:   params.UploadId,
		ChunkCount: params.Index,
	})
	if err != nil {
		log.Errorf("[Upload][File]:基础信息验证失败 %s", err.Error())
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrInternalServer, err)})
		return
	}
	//---
	const cap = 1024 * 1024
	var (
		capNow    int
		buf       = make([]byte, cap)
		sendBytes uploadPb.Upload_WriteBytesService
	)

	if sendBytes, err = uploadClient.WriteBytes(ctx); err != nil {
		log.Error("[Upload][File]:创建远程调用失败 %s", err.Error())
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrGrpcServer, err)})
		return
	}
	if err = sendBytes.SendMsg(&uploadPb.FileRequest{
		Filename: filepath(params.UploadId, params.FileName, params.Index),
		UploadId: params.UploadId,
		Index:    params.Index,
	}); err != nil {
		log.Error("[Upload][File]:远程创建文件失败 %s", err.Error())
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrGrpcServer, err)})
		return
	}

	for {
		n, err := ctx.Request.Body.Read(buf)
		if (n + capNow) > int(legitimate.Size) {
			log.Errorf("[Upload][File]:上传的容量超限")
			middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrRequestServer, err)})
			return
		} else if err != nil && err != io.EOF {
			log.Errorf("[Upload][File]:读取文件流失败")
			middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrInternalServer, err)})
			return
		}
		if err2 := sendBytes.Send(&uploadPb.Bytes{Size: int64(n), Content: buf[:n]}); err2 != nil {
			log.Errorf("[Upload][File]:远程写入失败 %s", err2.Error())
			middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrGrpcOp, err2)})
			return
		}
		if err == io.EOF {
			break
		}
	}
	_ = sendBytes.Send(&uploadPb.Bytes{Size: 0, Content: []byte{}})

	var recv uploadPb.ChunkResponse
	if err = sendBytes.RecvMsg(&recv); err != nil {
		log.Errorf("[Upload][File]:远程写入失败 %s", err.Error())
		middleware.ServerError(ctx, middleware.Response{Error: ecode.New(ecode.ErrGrpcOp, err)})
		return
	}
	//if err = sendBytes.Send(&uploadPb.Bytes{
	//	Content: nil,
	//}); err != nil {
	//	log.Errorf("[Upload][File]:发送结束标记 %s", err)
	//	return
	//}
	//设置通道的容量10块
	//b := make(chan []byte, 10)
	//go func() {
	//	var (
	//		n int
	//	)
	//	for {
	//		bt := make([]byte, 1024*1024) //开辟一个大小1m的缓冲空间
	//		switch n, err = file.Read(bt); true {
	//		case n < 0:
	//			log.Errorf("[Upload][Image]:数据读取失败 %s", err)
	//			fallthrough
	//		case n == 0:
	//			close(b)
	//			log.Errorf("[Upload][Image]:数据读取结束")
	//			return
	//		case n > 0:
	//			b <- bt
	//		}
	//	}
	//}()
	//
	//isClose := make(chan struct{}) //无缓冲通道，等待输入某个东西来关闭
	//go func() {
	//	for x := range b {
	//		err := sendBytes.Send(&uploadPb.Bytes{Content: x, Size: header.Size})
	//		if err != nil {
	//			log.Errorf("[Upload][Image]:数据发送失败 %s", err)
	//			return
	//		}
	//	}
	//	isClose <- struct{}{}
	//	if err = sendBytes.Send(&uploadPb.Bytes{
	//		Content: nil,
	//	}); err != nil {
	//		log.Errorf("[Upload][Image]:数据发送失败 %s", err)
	//		return
	//	}
	//}()
	//<-isClose
	//if err := sendBytes.RecvMsg(&fileInfo); err != nil {
	//	log.Errorf("[Upload][Image]:数据发送失败 %s", err)
	//	middleware.ServerError(ctx, middleware.Response{Error: err})
	//	return
	//}
	log.Info("[Upload][File]:End")
	middleware.Success(ctx, middleware.Response{

	})
}
