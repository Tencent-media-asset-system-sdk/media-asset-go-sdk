package utils

import (
	"math/rand"
	"time"
)

// TimeTick 计时器，单位 ms
type TimeTick struct {
	TimeList  []int64 // 打点时间点
	BuildTime int64   // 开始时间
}

// BuildTimeTick 构建一个计时器
func BuildTimeTick() TimeTick {
	return TimeTick{BuildTime: time.Now().UnixNano() / 1000000}
}

// TimeTick.Tick 计时器打点
func (tick *TimeTick) Tick() int64 {
	now := time.Now().UnixNano() / 1000000
	defer func() { tick.TimeList = append(tick.TimeList, now) }()
	if len(tick.TimeList) == 0 {
		return now - tick.BuildTime
	} else {
		return now - tick.TimeList[len(tick.TimeList)-1]
	}
}

// TimeTick.GetTotalTime 获取计计时器总时间
func (tick *TimeTick) GetTotalTime() int64 {
	now := time.Now().UnixNano() / 1000000
	return now - tick.BuildTime
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(result)
}

// GenerateObjectID 随机文件ID
func GenerateObjectID(flowID string, suffix string) string {
	return "/ai_media/" + flowID + "/" + GenerateRandomString(32) + suffix
}
