package redis

import (
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/model"
	"github.com/gin-gonic/gin/json"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

const INT_MAX = int(^uint(0) >> 1)

// init a judging info in redis
// the judging info is then used through the whole judging process
// IMPORTANT: we assume that the rank of test cases are consecutive
func InitJudging(judging model.Task) error {
	judgingID := KEY_PREFIX_JUDGING + strconv.FormatInt(int64(judging.JudgingID), 10)

	// clean up stale judging info
	// stale judging info may exist only because of rejudging
	cnt, err := Client.Del(judgingID).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		log.Warnf("removed stale info for judging %s in redis", judgingID)
	}

	// save judging info into a hash table
	startRank := INT_MAX
	for _, t := range judging.TestCases {
		tJson, err := json.Marshal(t)
		if err != nil {
			return err
		}

		err = Client.HSet(
			judgingID,
			KEY_PREFIX_TESTCASE+strconv.FormatInt(int64(t.Rank), 10),
			tJson,
		).Err()
		if err != nil {
			return err
		}

		startRank = int(math.Min(float64(startRank), float64(t.Rank)))
	}
	err = Client.HMSet(judgingID, map[string]interface{}{
		KEY_TESTCASE_NOW:    startRank,
		KEY_TESTCASES_TOTAL: len(judging.TestCases),
	}).Err()
	if err != nil {
		return err
	}

	return nil
}
