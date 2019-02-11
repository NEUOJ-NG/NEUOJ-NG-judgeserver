package util

import (
	myRedis "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/redis"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"sync"
)

const (
	DOWNLOAD_MAX_RETRY  = 5
	POSTFIX_EXECUTABLES = ".zip"
	POSTFIX_INPUT       = ".in"
	POSTFIX_OUTPUT      = ".out"
)

var (
	// map used for passing async download results
	downloadingMap     = make(map[string]chan bool)
	downloadingMapLock = sync.RWMutex{}
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	// get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// download file async and save md5sum to redis
// you can leave targetMD5 blank to skip md5sum check
func PrepareFileAsync(key string, id string, url string, dest string, targetMD5Sum string) error {
	// check if downloading is in process
	fullID := key + "_" + id
	downloadingMapLock.RLock()
	if _, ok := downloadingMap[fullID]; ok {
		// avoid re-download
		log.Debugf("file %s is already downloading", id)
		downloadingMapLock.RUnlock()
		return nil
	}
	downloadingMapLock.RUnlock()

	// check redis for file existence
	md5, err := myRedis.Client.HGet(key, id).Result()
	missingFile := err == redis.Nil
	invalidFile := err == nil && targetMD5Sum != "" && md5 != targetMD5Sum
	if missingFile || invalidFile {
		if missingFile {
			log.Infof("missing file %s, download async", fullID)
		} else if invalidFile {
			log.Infof("invalid file %s (targetMD5 = %s while "+
				"current md5 = %s), download async", fullID, targetMD5Sum, md5)
		}
		// start async download
		go func() {
			// check if destination path already exists
			// if true, delete it first
			if _, err := os.Stat(dest); err == nil {
				log.Debugf("file %s exists, removing", dest)
				err = os.Remove(dest)
				if err != nil {
					log.Errorf("failed to remove file %s: %s", dest, err.Error())
					panic(err)
				}
			}

			log.Debugf("start downloading file %s from %s to %s", fullID, url, dest)
			downloadingMapLock.Lock()
			downloadingMap[fullID] = make(chan bool)
			downloadingMapLock.Unlock()

			rst := true
			retries := 1
			for retries <= DOWNLOAD_MAX_RETRY {
				err := DownloadFile(dest, url)
				if err != nil {
					rst = false
					log.Errorf("failed to download file %s from %s after %x tries: %s",
						id, url, retries, err.Error(),
					)
				} else {
					rst = true
					log.Debugf("file %s successfully downloaded", id)
					break
				}
				retries++
			}

			// update md5sum to redis
			sum, err := GetFileMD5Sum(dest)
			if err != nil {
				rst = false
				log.Errorf("failed to get md5sum for file %s: %s", dest, err.Error())
			} else {
				// check md5sum if provided
				if targetMD5Sum != "" && targetMD5Sum != sum {
					rst = false
					log.Errorf("file %s corrupted while downloading (target md5sum = %s "+
						"while downloaded md5sum = %s", fullID, targetMD5Sum, sum)
				} else {
					// update md5sum to redis
					err := myRedis.Client.HSet(key, id, sum).Err()
					if err != nil {
						rst = false
						log.Errorf("failed to update md5sum for file %s to redis: %s",
							fullID, err.Error(),
						)
					}
				}
			}

			// TODO: should we use Lock here?
			downloadingMapLock.RLock()
			// send rst to channel with no-blocking way
			select {
			case downloadingMap[fullID] <- rst:
				downloadingMapLock.RUnlock()
				log.Debugf("notify blocking goroutines with result %x", rst)
			default:
				downloadingMapLock.RUnlock()
				log.Debugf("no blocking goroutines waiting for file %s", fullID)
				downloadingMapLock.Lock()
				delete(downloadingMap, fullID)
				downloadingMapLock.Unlock()
			}
		}()
	} else if err != nil {
		return err
	} else {
		log.Debugf("file %s ready", fullID)
	}

	return nil
}
