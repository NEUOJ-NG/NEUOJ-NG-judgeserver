package controller

import (
	"encoding/json"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	myRedis "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/redis"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strconv"
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

func GetTestCases(ctx *gin.Context) {
	id := ctx.Query("judgingid")

	// get current test case ID from redis
	nowTID, err := myRedis.Client.HGet(
		myRedis.KEY_PREFIX_JUDGING+id,
		myRedis.KEY_TESTCASE_NOW,
	).Result()
	if err != nil {
		log.Errorf("failed to get current test case ID for judging %s: %s", id, err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	nowTIDInt, _ := strconv.Atoi(nowTID)
	log.Debugf("current test case ID for judging %s is %s", id, nowTID)

	// get test case rank list from redis
	rankListJson, err := myRedis.Client.HGet(
		myRedis.KEY_PREFIX_JUDGING+id,
		myRedis.KEY_TESTCASE_RANK_LIST,
	).Result()
	if err != nil {
		log.Errorf("failed to get test case rank list for judging %s: %s", id, err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	var ranks []int
	err = json.Unmarshal([]byte(rankListJson), &ranks)
	if err != nil {
		log.Errorf("failed to get test case rank list for judging %s: %s", id, err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	log.Debugf("rank list for judging %s is %v", id, ranks)

	// get next test case info from redis
	if nowTIDInt >= len(ranks) {
		// no more test case, just return empty array
		log.Debugf("no more test case for judging %s", id)
		ctx.Data(http.StatusOK, "application/json", []byte("[]"))
		return
	} else {
		tInfo, err := myRedis.Client.HGet(
			myRedis.KEY_PREFIX_JUDGING+id,
			myRedis.KEY_PREFIX_TESTCASE+strconv.FormatInt(int64(ranks[nowTIDInt]), 10),
		).Result()
		if err != nil {
			log.Errorf("failed to get next test case for judging %s: %s", id, err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		} else {
			// increase current test case ID
			err := myRedis.Client.HIncrBy(
				myRedis.KEY_PREFIX_JUDGING+id,
				myRedis.KEY_TESTCASE_NOW,
				1,
			).Err()
			if err != nil {
				log.Errorf("failed to increase current test case ID for judging %s: %s", id, err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}

			ctx.Data(http.StatusOK, "application/json", []byte(tInfo))
			return
		}
	}
}

func GetTestCaseFiles(ctx *gin.Context) {
	id := ctx.Query("testcaseid")
	_, isInput := ctx.Request.URL.Query()["input"]
	_, isOutput := ctx.Request.URL.Query()["output"]

	if isInput {
		GetFiles(
			ctx,
			myRedis.KEY_TESTCASES,
			id+myRedis.KEY_POSTFIX_INPUT,
			filepath.Join(config.GetTestCaseStoragePath(), id+util.POSTFIX_INPUT),
			true,
		)
		return
	} else if isOutput {
		GetFiles(
			ctx,
			myRedis.KEY_TESTCASES,
			id+myRedis.KEY_POSTFIX_OUTPUT,
			filepath.Join(config.GetTestCaseStoragePath(), id+util.POSTFIX_OUTPUT),
			true,
		)
		return
	} else {
		log.Errorf("unknown test case file type, query is %v", ctx.Request.URL.Query())
		ctx.Status(http.StatusBadRequest)
		return
	}
}
