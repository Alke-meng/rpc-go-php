package tool

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// 识别手机号码
func CheckMobile(mobile string) bool {
	result, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, mobile)
	if result {
		return true
	}
	return false
}

// 获取正在运行的函数名
func GetFunName(l int) string {
	pc, _, _, _ := runtime.Caller(l)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, ".")
	return split[len(split)-1]
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", BToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", BToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", BToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
