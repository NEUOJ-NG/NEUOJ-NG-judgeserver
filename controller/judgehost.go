package controller

import (
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/form"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/mq"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Add a new judgehost to the list of judgehosts
// and also restart unfinished judgings.
func PostJudgehosts(ctx *gin.Context) {
	var judgehostVersion string
	if v, ok := ctx.Request.Header["User-Agent"]; ok {
		judgehostVersion = v[0]
	} else {
		judgehostVersion = "Unknown"
	}
	hostname := ctx.PostForm("hostname")
	log.Debugf("registering judgehost %s version %s",
		hostname, judgehostVersion)
	// TODO: save hostname to the list of judgehosts
	// TODO: restart unfinished judgings
	ctx.JSON(http.StatusOK, nil)
}

// Get judgehost configuration configured in config.toml
func GetJudgehostConfig(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		// return all config
		ctx.Data(
			http.StatusOK,
			"application/json; charset=utf-8",
			[]byte(config.GetConfig().Judgehost.Configuration),
		)
	} else {
		// return specific config
		c, err := config.GetJudgehostConfiguration(name, false)
		if err != nil {
			log.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, nil)
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				name: c,
			})
		}
	}
}

// API for judgehost to report internal errors.
func PostJudgehostsInternalError(ctx *gin.Context) {
	var internalError form.InternalError
	if err := ctx.ShouldBind(&internalError); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// TODO: save internal error and return auto-increment ID
		log.Debugf("%v", internalError)
		log.Warnf("receive internal error from judgehost, description: %v",
			internalError.Description)
		// now just return fake ID 0
		ctx.String(http.StatusOK, "0")
	}
}

// Retrieve judging tasks from message queue
// and give them to judgehost.
func PostJudgings(ctx *gin.Context) {
	judgehost := ctx.PostForm("judgehost")
	log.Debugf("judgehost %s fetching task", judgehost)

	// TODO: update judgehost timestamp in redis

	select {
	case task := <-mq.ConsumerMessages:
		log.Debugf("judgehost %s received a task: %s", judgehost, task.Body)
		_ = task.Ack(false)

		// TODO: perform testcases & executables check and prefetch

		ctx.Data(
			http.StatusOK,
			"application/json; charset=utf-8",
			[]byte(task.Body),
		)
	default:
		log.Debug("no task in channel")
		ctx.Data(
			http.StatusOK,
			"application/json; charset=utf-8",
			[]byte("{}"),
		)
	}
}
