package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"
)

type MicroContext struct {
	UserId int64 `json:"user_id"`
	*gin.Context
}

type HandlerFunc func(*MicroContext)

func C(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mc := &MicroContext{}
		mc.Context = ctx
		if get, exists := ctx.Get("userId"); !exists {
			log.Warn("userId is not found!")
		} else {
			mc.UserId = get.(int64)
		}
		handlerFunc(mc)
	}
}
