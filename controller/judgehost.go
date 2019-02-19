package controller

import (
	"encoding/json"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/form"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/model"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/mq"
	myRedis "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/redis"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strconv"
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
	log.Debugf(
		"registering judgehost %s version %s",
		hostname,
		judgehostVersion,
	)
	myRedis.UpdateJudgehostHeartbeat(hostname)
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
func PostInternalError(ctx *gin.Context) {
	var internalError form.InternalError
	if err := ctx.ShouldBind(&internalError); err != nil {
		log.Error(err.Error())
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

	myRedis.UpdateJudgehostHeartbeat(judgehost)

	select {
	case task := <-mq.ConsumerMessages:
		log.Debugf("judgehost %s received a task: %s", judgehost, task.Body)
		err := task.Ack(false)
		if err != nil {
			log.Errorf("failed to send ACK for task %s: %s", task.Body, err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// perform submissions & testcases & executables check and prefetch
		var taskObj model.Task
		err = json.Unmarshal(task.Body, &taskObj)
		if err != nil {
			log.Errorf("failed to unmarshal task json: %s", err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// prepare submissions
		sid := strconv.Itoa(taskObj.SubmitID)
		err = util.PrepareFileAsync(
			myRedis.KEY_SUBMISSIONS,
			sid,
			config.GetConfig().URL.Submissions+sid,
			filepath.Join(config.GetSubmissionStoragePath(), sid),
			"",
		)
		if err != nil {
			log.Error(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// prepare executables
		// prepare run
		err = util.PrepareFileAsync(
			myRedis.KEY_EXECUTABLES,
			taskObj.Run,
			config.GetConfig().URL.Executables+taskObj.Run,
			filepath.Join(config.GetExecutableStoragePath(), taskObj.Run+util.POSTFIX_EXECUTABLES),
			taskObj.RunMD5Sum,
		)
		if err != nil {
			log.Error(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// prepare compare
		err = util.PrepareFileAsync(
			myRedis.KEY_EXECUTABLES,
			taskObj.Compare,
			config.GetConfig().URL.Executables+taskObj.Compare,
			filepath.Join(config.GetExecutableStoragePath(), taskObj.Compare+util.POSTFIX_EXECUTABLES),
			taskObj.CompareMD5Sum,
		)
		if err != nil {
			log.Error(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// prepare compile_script
		err = util.PrepareFileAsync(
			myRedis.KEY_EXECUTABLES,
			taskObj.CompileScript,
			config.GetConfig().URL.Executables+taskObj.CompileScript,
			filepath.Join(config.GetExecutableStoragePath(), taskObj.CompileScript+util.POSTFIX_EXECUTABLES),
			taskObj.CompileScriptMD5Sum,
		)
		if err != nil {
			log.Error(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// prepare testcases
		for i, t := range taskObj.TestCases {
			// propagate problem ID to testcases
			taskObj.TestCases[i].ProbID = taskObj.ProbID

			tid := strconv.Itoa(t.TestCaseID)

			// prepare input
			err = util.PrepareFileAsync(
				myRedis.KEY_TESTCASES,
				tid+myRedis.KEY_POSTFIX_INPUT,
				config.GetConfig().URL.TestCases+tid+"?type=input",
				filepath.Join(config.GetTestCaseStoragePath(), tid+util.POSTFIX_INPUT),
				t.MD5SumInput,
			)
			if err != nil {
				log.Error(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}

			// prepare output
			err = util.PrepareFileAsync(
				myRedis.KEY_TESTCASES,
				tid+myRedis.KEY_POSTFIX_OUTPUT,
				config.GetConfig().URL.TestCases+tid+"?type=output",
				filepath.Join(config.GetTestCaseStoragePath(), tid+util.POSTFIX_OUTPUT),
				t.MD5SumOutput,
			)
			if err != nil {
				log.Error(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}

		// save judging info to redis
		err = myRedis.InitJudging(taskObj)
		if err != nil {
			log.Errorf("failed to save judging info to redis: %s", err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// TODO: perform judging info GC

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
