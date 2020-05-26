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
	// 系统错误 100
	// 110 常用的系统错
	ErrInternalServer = &Errno{Code: 50001, Message: "内部服务器错误"}
)
