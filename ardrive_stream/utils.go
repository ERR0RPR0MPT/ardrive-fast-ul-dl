package ardrive_stream

import (
	"math/rand"
	"time"
)

func randomStringFromSlice(slice []string) string {
	// 设置随机种子
	rand.NewSource(time.Now().UnixNano())
	// 生成一个随机索引
	index := rand.Intn(len(slice))
	return slice[index]
}
