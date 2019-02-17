package controller

import (
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	myRedis "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/redis"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

// base function for judgehost to get files.
// block on downloadingMap[fullID] until
// files are prepared
func GetFiles(ctx *gin.Context, key string, id string, path string, jsonBase64Encode bool) {
	fullID := util.GetFileFullID(key, id)

	if util.IsPreparing(fullID) {
		rst := util.WaitForFile(fullID)
		util.UpdatePrepareResult(fullID, rst)
		if !rst {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	// check redis for file existence
	_, err := myRedis.Client.HGet(key, id).Result()
	if err == redis.Nil {
		log.Errorf("file %s not found in redis", fullID)
		ctx.Status(http.StatusInternalServerError)
		return
	} else if err != nil {
		log.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	} else {
		// file exist, ready to serve
		if jsonBase64Encode {
			if content, err := util.GetFileBase64(path); err != nil {
				log.Error(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			} else {
				if s, err := json.Marshal(content); err != nil {
					log.Error(err.Error())
					ctx.Status(http.StatusInternalServerError)
					return
				} else {
					ctx.Data(http.StatusOK, "application/json", []byte(s))
				}
			}
		} else {
			ctx.File(path)
			return
		}
	}
}

func GetSubmissionFiles(ctx *gin.Context) {
	id := ctx.Query("submission_id")
	GetFiles(
		ctx,
		myRedis.KEY_SUBMISSIONS,
		id,
		filepath.Join(config.GetSubmissionStoragePath(), id),
		false,
	)
}

func GetExecutable(ctx *gin.Context) {
	id := ctx.Query("execid")
	GetFiles(
		ctx,
		myRedis.KEY_EXECUTABLES,
		id,
		filepath.Join(config.GetExecutableStoragePath(), id+util.POSTFIX_EXECUTABLES),
		true,
	)
}
