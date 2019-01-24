package router

import (
	c "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
	// test
	r.GET("/ping", c.Ping)
}
