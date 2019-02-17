package router

import (
	c "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
	r.POST("/judgehosts", c.PostJudgehosts)
	r.POST("/judgehosts/internal-error", c.PostJudgehostsInternalError)
	r.GET("/config", c.GetJudgehostConfig)
	r.POST("/judgings", c.PostJudgings)
	r.GET("/submission_files", c.GetSubmissionFiles)
	r.GET("/executable", c.GetExecutable)
	r.PUT("/judgings/:id", c.PutJudgings)
}
