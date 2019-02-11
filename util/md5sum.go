package util

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"os"
)

const FILECHUNK = 8192

// get file md5sum with chunked method
func GetFileMD5Sum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	info, _ := file.Stat()
	fileSize := info.Size()
	blockNumber := uint64(math.Ceil(float64(fileSize) / float64(FILECHUNK)))

	hash := md5.New()
	for i := uint64(0); i < blockNumber; i++ {
		blockSize := int(math.Min(FILECHUNK, float64(fileSize-int64(i*FILECHUNK))))
		buf := make([]byte, blockSize)
		_, err := file.Read(buf)
		if err != nil {
			return "", err
		}
		hash.Write(buf)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
