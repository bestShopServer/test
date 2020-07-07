package release

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Func: Gin Err Response
func ResponseErr(errNo int, c *gin.Context) {
	//util.Logger(util.ERROR_LEVEL, "Response", strconv.Itoa(errNo))
	util.Error("Response: ", strconv.Itoa(errNo))
	// Generate Json String Response
	c.JSON(http.StatusOK, model.ErrResp{
		Err: errNo,
		Msg: model.ErrMap[errNo],
	})
}

// Func: Gin Err Response
func ResponseNotice(c *gin.Context, errNo int, notice string) {
	//util.Logger(util.ERROR_LEVEL, "Response", strconv.Itoa(errNo))
	util.Error("Code:%v Notice:%v", errNo, notice)
	// Generate Json String Response
	c.JSON(http.StatusOK, model.ErrResp{
		Err: errNo,
		Msg: notice,
	})
}

// Check All Params From The Request
func CheckParams(v interface{}, cat string, err error, optionalFields []string) bool {
	check := true
	if err != nil {
		//util.Logger(util.ERROR_LEVEL, cat, "Parsing Param Err:"+err.Error())
		util.Error("%s Parsing Param ERROR[%v]", cat, err.Error())
		check = false
	} else if !CheckField(v, optionalFields) {
		check = false
	}
	return check
}

// Check All Field In A Struct
func CheckField(v interface{}, optionalFields []string) bool {
	//ref := reflect.ValueOf(v)
	ref := reflect.TypeOf(v)
	check := true
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)

		var optional bool
		for j := 0; j < len(optionalFields); j++ {
			if field.Name == optionalFields[j] {
				optional = true
				break
			}
		}

		if optional {
			continue
		}

		switch field.Type.String() {
		case "int":
			if reflect.ValueOf(v).Field(i).Int() == 0 {
				check = false
			}
		case "string":
			if reflect.ValueOf(v).Field(i).String() == "" {
				check = false
			}
		case "struct":
			check = CheckField(reflect.ValueOf(v).Field(i).Interface(), []string{})
		}

	}
	return check
}

//// Find Clew Id in Cache
//func FindClewIdInCacheBak(index []int, roomId string, key string) (int, int) {
//	if _, ok := global.ClewCache[roomId]; ok {
//		clews := global.ClewCache[roomId]
//		if len(index) == 1 {
//			if _, ok := clews[key]; ok {
//				if len(clews[key]) > index[0] {
//					return clews[key][index[0]].Id, clews[key][index[0]].Round
//				}
//				return 0, 0
//			}
//			return 0, 0
//		}
//
//		if len(index) == 2 {
//			if _, ok := clews[key]; ok {
//				if len(clews[key]) > index[0] {
//					main := clews[key][index[0]]
//					if len(main.Sub) > index[1] {
//						return main.Sub[index[1]].Id, main.Sub[index[1]].Round
//					}
//					return 0, 0
//				}
//				return 0, 0
//			}
//			return 0, 0
//		}
//	}
//	return 0, 0
//}

//Find Clew Id in Cache
func FindClewIdInCache(conn redis.Conn, index []int, roomId string, key string) (int, int) {
	util.Debug("FindClewIdInCache index:%+v roomId:%v key:%v", index, roomId, key)
	var clews map[string][]global.MainClew
	//ok, err := gocache.GetRoomClew(roomId, &clews)
	ok, err := gocache.ConnGetRoomClew(conn, roomId, &clews)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return 0, 0
	}
	if !ok {
		util.Error("获取线索失败:%v", ok)
		return 0, 0
	}
	util.Debug("len:%v", len(index))
	//util.Debug("%+v", clews)

	if _, ok := clews[key]; ok {
		//util.Debug("len:%v %+v", len(clews), clews[key])
		if len(index) == 1 {
			if len(clews[key]) > index[0] {
				util.Debug("id:%v round:%v", clews[key][index[0]].Id, clews[key][index[0]].Round)
				return clews[key][index[0]].Id, clews[key][index[0]].Round
			}
			return 0, 0
		} else if len(index) == 2 {
			if len(clews[key]) > index[0] {
				main := clews[key][index[0]]
				if len(main.Sub) > index[1] {
					util.Debug("id:%v round:%v", main.Sub[index[1]].Id, main.Sub[index[1]].Round)
					return main.Sub[index[1]].Id, main.Sub[index[1]].Round
				}
				return 0, 0
			}
			return 0, 0
			/* add by skc at 2020-04-17 begin */
		} else if len(index) == 3 {
			if len(clews[key]) > index[0] {
				//util.Debug("三级线索:%+v", clews[key][index[0]])
				main := clews[key][index[0]]
				if len(main.Sub) > index[1] {
					util.Debug("三级线索:%+v len:%v", main, len(main.Sub[index[1]].Sub))
					if len(main.Sub[index[1]].Sub) > index[2] {
						util.Debug("id:%v round:%v", main.Sub[index[1]].Sub[index[2]].Id, main.Sub[index[1]].Sub[index[2]].Round)
						return main.Sub[index[1]].Sub[index[2]].Id, main.Sub[index[1]].Sub[index[2]].Round
					}
				}
				return 0, 0
			}
			return 0, 0
		}
		/* add by skc at 2020-04-17 end */

	}
	util.Info("没找到线索:%v", key)
	return 0, 0
}

// Func: Get Game Clew By RoomId And round
func GetGameClew(rid string, round int) (model.GameClew, int) {
	params := model.GameClew{}

	// Find RoomId in Game Bucket
	var clews []model.Clew
	//gb := boltdb.View([]byte(rid), "GameBucket")
	//if gb != nil {
	//	if err := json.Unmarshal(gb, &clews); err != nil {
	//		util.Logger(util.ERROR_LEVEL, "GameClew", "Decoding Clew Err:"+err.Error())
	//	}
	//} else {
	//	return params, model.ERR_NOT_GET_GAME_INFO
	//}
	_, err := gocache.GetRoomGame(rid, &clews)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return params, model.ERR_NOT_GET_GAME_INFO
	}

	roleMap := make(map[string]model.Role)
	newKey := make(map[string]model.Key)

	for i := 0; i < len(clews); i++ {
		if clews[i].Round <= round {
			for k, v := range clews[i].Key {
				key, ok := newKey[k]
				if !ok {
					key = model.Key{}
				}
				key.Album = v.Album
				key.TotalNum += v.TotalNum
				key.FetchNum += v.FetchNum
				key.KeyList = append(key.KeyList, v.KeyList...)
				newKey[k] = key

				role, ok := roleMap[k]
				if !ok {
					role = model.Role{}
				}

				var count int
				for _, clew := range v.KeyList {
					if clew.Status != 0 {
						count++
					}
				}

				role.TotalNum += len(v.KeyList)
				role.OpenNum += count
				role.Album = key.Album
				role.Name = k

				roleMap[k] = role
			}
		}
	}

	var roles []model.Role
	for _, v := range roleMap {
		roles = append(roles, v)
	}

	sort.Slice(roles, func(i, j int) bool {
		return strings.Compare(roles[i].Name, roles[j].Name) <= 0
	})

	params.Roles = roles
	params.Key = newKey

	return params, model.ERR_OK
}

//sandbox回传的数据转结构体
func InterfaceToJsonStruct(p interface{}, v interface{}) (err error) {
	jstmp, err := json.Marshal(p)
	if err != nil {
		util.Error("数据转换出错!")
		return err
	}
	util.Debug("循环数据JSON:%v", string(jstmp))

	err = json.Unmarshal(jstmp, &v)
	if err != nil {
		util.Error("数据转换出错!")
		return err
	}
	return err
}

//校验能否进入新房间
func CheckJoinNewRoom(uid string) (ok bool, err error) {
	roomCode, err := gocache.GetUserRoom(uid)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		return ok, err
	}
	util.Info("用户:%v在房间:%v", uid, roomCode)

	if len(roomCode) == 0 {
		util.Info("用户不在任何房间,可以继续!")
		ok = true
		return ok, err
	}

	//判断之前的房间是否已结束
	rif2 := model.RoomInfo{}
	ok, err = gocache.CheckRoomExists(roomCode)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return ok, err
	}
	if ok {
		err = gocache.GetRoomInfo(roomCode, &rif2)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			return ok, err
		}
		if rif2.Status != 2 { //标识房间已结束
			util.Error("上一房间%v游戏未结束，不允许加入新房间! 状态:%v", roomCode, rif2.Status)
			err = errors.New(fmt.Sprintf("上一房间%v游戏未结束，不允许加入新房间!", roomCode))
			return ok, err
		}
	}
	return ok, err
}

//获取用户是否为会员
func GetUserMemberBase(union_id string) (info model.MemberBaseInfo, err error) {
	conn := gocache.RedisConnPool.Get()
	defer conn.Close()
	_, err = conn.Do("select", config.GetConfig().RedisDataDb)
	if err != nil {
		_, err = conn.Do("select", config.GetConfig().RedisDb)
		util.Error("ERROR[%v]", err.Error())
		return info, err
	}
	exists, err := redis.Bool(conn.Do("exists", union_id))
	if err != nil {
		_, err = conn.Do("select", config.GetConfig().RedisDb)
		util.Error("ERROR:%v", err.Error())
		return info, err
	}
	if exists {
		ty, err := redis.String(conn.Do("type", union_id))
		util.Info("type:%+v", ty)
		//reply, err := redis.Bytes(conn.Do("hgetall", union_id))
		////reply, err := redis.StringMap(conn.Do("hgetall", union_id))
		//if err != nil {
		//	util.Error("ERROR[%v]", err.Error())
		//	_, err = conn.Do("select", config.GetConfig().RedisDb)
		//	return info, err
		//}

		reply, err := redis.Values(conn.Do("hgetall", union_id))
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			_, err = conn.Do("select", config.GetConfig().RedisDb)
			return info, err
		}

		err = redis.ScanStruct(reply, &info)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			_, err = conn.Do("select", config.GetConfig().RedisDb)
			return info, err
		}

		nowTime := time.Now().Format("2006-01-02 15:04:05")
		util.Debug("会员等级:%v 到期时间:%v", info.Member, info.InvTime)
		if info.InvTime < nowTime {
			info.Member = 0
		}
	}
	_, err = conn.Do("select", config.GetConfig().RedisDb)

	util.Debug("会员信息:%+v", info)

	return info, err
}
