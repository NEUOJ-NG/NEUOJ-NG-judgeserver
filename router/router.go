package router

import (
	c "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
	r.POST("/judgehosts", c.PostJudgehosts)
	r.POST("/internal_error", c.PostInternalError)
	r.GET("/config", c.GetJudgehostConfig)
	r.POST("/judgings", c.PostJudgings)
	r.GET("/submission_files", c.GetSubmissionFiles)
	r.GET("/executable", c.GetExecutable)
	r.PUT("/judgings/:id", c.PutJudgings)
	r.GET("/testcases", c.GetTestCases)
	r.GET("/testcase_files", c.GetTestCaseFiles)
	r.POST("/judging_runs", c.PostJudgingRuns)
}
