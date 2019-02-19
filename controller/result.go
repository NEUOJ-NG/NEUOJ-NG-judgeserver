package controller

import (
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/form"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func PostJudgingRuns(ctx *gin.Context) {
	var judgingRunResult form.JudgingRunResult
	if err := ctx.ShouldBind(&judgingRunResult); err != nil {
		log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// TODO: send judging run result to backend
		log.Debugf(
			"judging run result for judging %s test case %s is %s",
			judgingRunResult.JudgingID,
			judgingRunResult.TestCaseID,
			judgingRunResult.RunResult,
		)
		log.Debugf("runtime: %s", judgingRunResult.Runtime)
		log.Debugf("output_system: %s", judgingRunResult.OutputSystem)
	}
}
