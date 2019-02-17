package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// update compile results
func PutJudgings(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Debug(id)
	judgehost := ctx.PostForm("judgehost")
	compileSuccess := ctx.PostForm("compile_success")
	outputCompile := ctx.PostForm("output_compile")
	log.Debug(judgehost)
	log.Debug(compileSuccess)
	log.Debug(outputCompile)
	// TODO: PUT compile result to backend
}
