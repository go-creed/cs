package ecode

/*
错误码设计
第一位表示错误级别, 1 为系统错误, 2 为其它错误
第二三位表示服务模块代码 比如210是数据库错误，211具体到mysql，221还是reids，241具体到规则话术模块
第四五位表示具体错误代码
*/
var (
	Success = &Errno{Code: 0, Message: "success"}

	//--------
	// 系统错误 x00
	// 400 常见的请求端错误
	ErrRequestServer = &Errno{Code: 40001, Message: "请求端错误"}
	// 500 常用的服务端错误
	ErrInternalServer = &Errno{Code: 50001, Message: "内部服务器错误"}

	// 800 常用的Grpc错误
	ErrGrpcServer = &Errno{Code: 80001, Message: "Grpc服务异常"}
	ErrGrpcOp     = &Errno{Code: 80002, Message: "Grpc服务操作失败"}
)
