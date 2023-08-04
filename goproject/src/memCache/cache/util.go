package cache

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
)

func parseSize(size string) (uint64, string) {
	//默认情况下大小为100MB
	var (
		maxSize    uint64
		maxSizeStr string
	)
	re, _ := regexp.Compile("[0-9]+")
	unit := re.ReplaceAll([]byte(size), []byte(""))
	num := re.FindAll([]byte(size), len(size))
	if len(num) != 1 || string(unit) == "" {
		maxSize = 100 * MB
		maxSizeStr = "100MB"
	} else {
		number, _ := strconv.ParseInt(string(num[0]), 10, 64)
		switch strings.ToUpper(string(unit)) {
		case "B":
			maxSize = uint64(number)
		case "KB":
			maxSize = uint64(number) * KB
		case "MB":
			maxSize = uint64(number) * MB
		case "GB":
			maxSize = uint64(number) * GB
		case "TB":
			maxSize = uint64(number) * TB
		case "PB":
			maxSize = uint64(number) * PB
		}
		maxSizeStr = string(num[0]) + string(unit)
	}
	return maxSize, maxSizeStr
}

func ParseOcSize(val interface{}) uint64 {
	js, _ := json.Marshal(val) //json包提供的序列化函数会自动递归，将每一个指针都解析出来，因此采用该方式估算长度较为准确
	return uint64(len(js))
}
