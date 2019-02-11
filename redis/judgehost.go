package redis

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func UpdateJudgehostHeartbeat(hostname string) {
	err := Client.HMSet(KEY_JUDGEHOST, map[string]interface{}{
		hostname: time.Now().Unix(),
	}).Err()
	if err != nil {
		log.Errorf(
			"failed to update heartbeat for judgehost %s: %s",
			hostname,
			err.Error(),
		)
	} else {
		log.Debugf("successfully update heartbeat for judgehost %s", hostname)
	}
}
