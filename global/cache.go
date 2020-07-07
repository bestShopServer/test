package global

import (
	//"DetectiveMasterServer/util"
	"sync"
	"time"
)

//var UserCache map[string]string
//var RoomCache map[string][]string
//var ClewCache map[string]map[string][]MainClew

var RoomCacheDeleteTask map[string]*RoomCacheTimer
var UserCacheMutex *sync.Mutex
var RoomCacheMutex *sync.Mutex
var ClewCacheMutex *sync.Mutex

var RoomDeleteTaskMutex *sync.Mutex

//var MongoRoomCache *mgo.Collection
//var MongoUserCache *mgo.Collection

type ClewRound struct {
	Id    int         `json:"id"`
	Round int         `json:"round"`
	Sub   []ClewRound `json:"sub"`
}

type MainClew struct {
	ClewRound
	Sub []ClewRound `json:"sub"`
}

type SubClew struct {
	ClewRound
}

type RoomCacheTimer struct {
	RoomId string
	Timer  *time.Timer
}

func (r *RoomCacheTimer) StartTimer() {
	r.Timer = time.NewTimer(12 * time.Hour)
	go func() {
		select {
		case <-r.Timer.C:
			r.StartTimer()
			//DeleteRoomCache(r.RoomId)
			//DeleteClewCache(r.RoomId)
			break
		}
	}()
}

func (r *RoomCacheTimer) StopTimer() {
	r.Timer.Stop()
	r.Timer = nil
}

func init() {
	//UserCache = make(map[string]string)
	//RoomCache = make(map[string][]string)
	//ClewCache = make(map[string]map[string][]MainClew)
	RoomCacheDeleteTask = make(map[string]*RoomCacheTimer)
	UserCacheMutex = new(sync.Mutex)
	RoomCacheMutex = new(sync.Mutex)
	ClewCacheMutex = new(sync.Mutex)
	RoomDeleteTaskMutex = new(sync.Mutex)
	RoomInfoMutex = new(sync.Mutex)
	ClewOpenMutex = new(sync.Mutex)

	//add mongo
	//MongoRoomCache = MongoDB.C("room")
	//MongoUserCache = MongoDB.C("user")

}

//
//func SetUserCache(unionId string, roomId string) {
//	UserCacheMutex.Lock()
//	UserCache[unionId] = roomId
//	UserCacheMutex.Unlock()
//}
//
//func AddUserToRoomCache(unionId string, roomId string) {
//	RoomCacheMutex.Lock()
//	var unionIdSli []string
//	if _, ok := RoomCache[roomId]; ok {
//		unionIdSli = RoomCache[roomId]
//	}
//	exist := false
//	for _, v := range unionIdSli {
//		if v == unionId {
//			exist = true
//			break
//		}
//	}
//	if !exist {
//		unionIdSli = append(unionIdSli, unionId)
//	}
//	RoomCache[roomId] = unionIdSli
//	RoomCacheMutex.Unlock()
//
//}

func AddRoomDeleteTask(roomId string) {
	var t RoomCacheTimer
	t.RoomId = roomId
	t.StartTimer()
	RoomDeleteTaskMutex.Lock()
	RoomCacheDeleteTask[roomId] = &t
	RoomDeleteTaskMutex.Unlock()

}

func AddClewToCache(m map[string][]MainClew, roomId string) {
	ClewCacheMutex.Lock()
	//ClewCache[roomId] = m
	//clewStr, _ :=json.Marshal(&m)
	ClewCacheMutex.Unlock()
}

//func DeleteUserCache(unionId string, roomId string) {
//	UserCacheMutex.Lock()
//	if UserCache[unionId] != "" && UserCache[unionId] == roomId {
//		delete(UserCache, unionId)
//	}
//	UserCacheMutex.Unlock()
//	util.Debug("删除用户[%v]", unionId)
//}

//func DeleteRoomCache(roomId string) {
//	RoomCacheMutex.Lock()
//	if _, ok := RoomCache[roomId]; ok {
//		for _, v := range RoomCache[roomId] {
//			DeleteUserCache(v, roomId)
//		}
//		delete(RoomCache, roomId)
//	}
//	RoomCacheMutex.Unlock()
//
//}

func DeleteRoomDeleteTask(roomId string) {
	RoomDeleteTaskMutex.Lock()
	if t, ok := RoomCacheDeleteTask[roomId]; ok {
		t.StopTimer()
		delete(RoomCacheDeleteTask, roomId)
	}
	RoomDeleteTaskMutex.Unlock()

}

//
//func DeleteClewCache(roomId string) {
//	ClewCacheMutex.Lock()
//	if _, ok := ClewCache[roomId]; ok {
//		delete(ClewCache, roomId)
//	}
//	ClewCacheMutex.Unlock()
//}
