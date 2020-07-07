package global

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/util"
	"sync"
)

var RoomInfoMutex *sync.Mutex
var ClewOpenMutex *sync.Mutex

// Func: Calc Not Voted People Num
func CalcNotVoteNum(m map[string]bool) int {
	var count int
	for _, v := range m {
		if v == false {
			count++
		}
	}
	return count
}

// Func: Create RoomId that Not In Cache
func CreateRoomId(len int) string {
	for {
		rid := util.GetRandomInt(len)
		//if RoomCache[rid] == nil {
		//	return rid
		//}
		if !gocache.Exists(rid) {
			return rid
		}
	}
}
