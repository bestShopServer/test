package release

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func ClewOpen(c *gin.Context) {
	util.Info("ClewOpen ...")
	defer global.ClewOpenMutex.Unlock()
	global.ClewOpenMutex.Lock()

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get RoomId and OpenId From Param
	var req model.ClewOpenReq
	var roleId int
	err := c.Bind(&req)

	// Optional Fields List
	optionalFields := []string{"Password"}

	// Check Param
	if !CheckParams(req, "ClewOpen", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	if len(req.Index) > 3 || len(req.Index) == 0 {
		util.Error("ERR_WRONG_FORMAT")
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	//userId := req.OpenId
	userId := req.UnionId
	roomId := req.RoomId
	password := req.Password
	//clewId, round := FindClewIdInCache2(req.Index, roomId, req.Key)
	//if round == 0 {
	//	util.Error("ERR_WRONG_FORMAT")
	//}

	conn := gocache.RedisConnPool.Get()
	defer conn.Close()
	util.Debug("FindClewIdInCache ...")
	clewId, round := FindClewIdInCache(conn, req.Index, roomId, req.Key)
	if round == 0 {
		util.Error("ERR_WRONG_FORMAT")
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("线索ID:%v", clewId)

	// Find Room In Room Cache
	//_, ok := global.RoomCache[roomId]
	//if !ok {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	//ok, err := gocache.CheckRoomExists(roomId)
	ok, err := gocache.ConnCheckRoomExists(conn, roomId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	if !ok {
		util.Error("ERR_ROOM_NOT_EXIST")
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find Room in RoomBucket
	//b := boltdb.View([]byte(roomId), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//err = json.Unmarshal(b, &roomInfo)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "ClewOpen", "Decoding Room Info Err:"+err.Error())
	//	util.Error("ClewOpen Decoding Room Info ERROR[%v]", err.Error())
	//}
	//读取房间信息
	//err = gocache.GetRoomInfo(roomId, &roomInfo)
	err = gocache.ConnGetRoomInfo(conn, roomId, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// If You are not in this room
	exist := false
	for _, v := range roomInfo.UnionIdSlice {
		if v == userId {
			exist = true
		}
	}
	if !exist {
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	// Find Role You Choose
	var role string
	for _, v := range roomInfo.PlayerSlice {
		if v.UnionId == userId {
			role = v.Role.Name
			roleId = v.Role.Id
		}
	}
	util.Info("role id:%v", roleId)

	if req.Key == role {
		util.Error("ERR_CANNOT_OPEN_OWN_CLEW")
		ResponseErr(model.ERR_CANNOT_OPEN_OWN_CLEW, c)
		return
	}

	// Find RoomId In GameBucket
	//gb := boltdb.View([]byte(roomId), "GameBucket")
	//if gb == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode ClewInfo
	var clew []model.Clew
	//err = json.Unmarshal(gb, &clew)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "ClewOpen", "Decoding Clew Info Err:"+err.Error())
	//	util.Error("ClewOpen Decoding Clew Info ERROR[%v]", err.Error())
	//}
	//_, err = gocache.GetRoomGame(roomId, &clew)
	_, err = gocache.ConnGetRoomGame(conn, roomId, &clew)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find RoomId in StageBucket
	//sb := boltdb.View([]byte(roomId), "StageBucket")
	//if sb == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Stage Info
	stageInfo := model.StageInfo{}
	//err = json.Unmarshal(sb, &stageInfo)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "ClewOpen", "Decoding Stage Info Err:"+err.Error())
	//	util.Error("ClewOpen Decoding Stage Info ERROR[%v]", err.Error())
	//}
	//_, err = gocache.GetRoomStage(roomId, &stageInfo)
	_, err = gocache.ConnGetRoomStage(conn, roomId, &stageInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	if round > stageInfo.Round {
		util.Error("ERR_WRONG_FORMAT")
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	// Get User Ap In ApBucket
	//ab := boltdb.View([]byte(roomId), "ApBucket")
	//if ab == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Ap Info
	var apMap map[string][]int
	//err = json.Unmarshal(ab, &apMap)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "ClewOpen", "Decoding Ap Info Err:"+err.Error())
	//	util.Error("ClewOpen Decoding Ap Info ERROR[%v]", err.Error())
	//}
	//获取
	//apMap, _, err = gocache.GetAPInfo(roomId)
	apMap, _, err = gocache.ConnGetAPInfo(conn, roomId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Info("处理[%v]轮线索...", stageInfo.Round)

	currentStage := stageInfo.Round
	if currentStage < 1 {
		util.Error("线索轮次有误:%v", currentStage)
		ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
		return
	}
	if currentStage > roomInfo.Round {
		util.Error("搜证已结束:%v %v", currentStage, roomInfo.Round)
		ResponseErr(model.ERR_CLEW_NOT_OPENED, c)
		return
	}
	//userAp := apMap[userId][currentStage-1]
	userAp := apMap[userId][currentStage] //增加前置0轮情况
	util.Info("userAp:%v", userAp)

	// Find And Open Clew
	found := false
	var costAp int
	for i := 0; i < currentStage; i++ {
		for k, v := range clew[i].Key {
			for j := 0; j < len(v.KeyList); j++ {
				if v.KeyList[j].Id == clewId {
					found = true
					err := CanClewOpen(v.KeyList[j].KeyDetail, userId, password, userAp, roleId)
					if err == 0 {
						clew[i].Key[k].KeyList[j].UnionId = userId
						clew[i].Key[k].KeyList[j].Status = 1
						costAp = v.KeyList[j].KeyDetail.Ap
						goto quit
					} else {
						util.Error("这个报错了...")
						ResponseErr(err, c)
						return
					}
				} else if v.KeyList[j].Sub != nil {
					for l := 0; l < len(v.KeyList[j].Sub); l++ {
						//util.Debug("寻找线索:%+v", v.KeyList[j].Sub[l])
						if v.KeyList[j].Sub[l].Id == clewId {
							found = true
							err := CanClewOpen(v.KeyList[j].Sub[l].KeyDetail, userId, password, userAp, roleId)
							if err == 0 {
								if v.KeyList[j].Status != model.STATUS_PUB && v.KeyList[j].UnionId != userId {
									util.Error("ERR_PARENT_CLEW_NOT_PUB")
									ResponseErr(model.ERR_PARENT_CLEW_NOT_PUB, c)
									return
								}
								clew[i].Key[k].KeyList[j].Sub[l].UnionId = userId
								clew[i].Key[k].KeyList[j].Sub[l].Status = 1
								costAp = v.KeyList[j].Sub[l].KeyDetail.Ap //AP点获取线索对应值
								util.Info("goto quit ...")
								goto quit
							} else {
								util.Error("这个报错了...")
								ResponseErr(err, c)
								return
							}
						}
					}
					//支持三级线索搜索
					if len(req.Index) == 3 {
						for l := 0; l < len(v.KeyList[j].Sub); l++ {
							if v.KeyList[j].Sub[l].Sub != nil {
								util.Debug("寻找线索:%+v", v.KeyList[j].Sub[l])
								for m := 0; m < len(v.KeyList[j].Sub[l].Sub); m++ {
									if v.KeyList[j].Sub[l].Sub[m].Id == clewId {
										found = true
										err := CanClewOpen(v.KeyList[j].Sub[l].Sub[m].KeyDetail, userId, password, userAp, roleId)
										if err == 0 {
											if v.KeyList[j].Status != model.STATUS_PUB && v.KeyList[j].UnionId != userId {
												util.Error("ERR_PARENT_CLEW_NOT_PUB")
												ResponseErr(model.ERR_PARENT_CLEW_NOT_PUB, c)
												return
											}
											clew[i].Key[k].KeyList[j].Sub[l].Sub[m].UnionId = userId
											clew[i].Key[k].KeyList[j].Sub[l].Sub[m].Status = 1
											costAp = v.KeyList[j].Sub[l].Sub[m].KeyDetail.Ap
											util.Info("goto quit ...")
											goto quit
										} else {
											util.Error("这个报错了...")
											ResponseErr(err, c)
											return
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

quit:
	// Update Game Bucket And Ap Bucket
	if found {
		//clewEncoding, _ := json.Marshal(clew)
		//boltdb.CreateOrUpdate([]byte(roomId), clewEncoding, "GameBucket")
		//err = gocache.SetRoomGame(roomId, clew)
		err = gocache.ConnSetRoomGame(conn, roomId, clew)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
			return
		}

		//apMap[userId][currentStage-1] = apMap[userId][currentStage-1] - costAp
		apMap[userId][currentStage] = apMap[userId][currentStage] - costAp //增加前置0轮情况
		//apEncoding, _ := json.Marshal(apMap)
		//boltdb.CreateOrUpdate([]byte(roomId), apEncoding, "ApBucket")
		//err = gocache.SetAPInfo(roomId, apMap)
		err = gocache.ConnSetAPInfo(conn, roomId, apMap)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
			return
		}

		go websocket.SendClewOpenMessage(roomInfo.UnionIdSlice, userId, req.Key, req.Index)

		clewOpenResp := model.ClewOpenResp{}
		clewOpenResp.Params = "success"
		c.JSON(http.StatusOK, &clewOpenResp)
		return
	} else {
		util.Error("ERR_CLEW_NOT_FOUND")
		ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
		return
	}
}

// Func: Judge Clew Can Be Opened
func CanClewOpen(k model.KeyDetail, unionId string, password string, userAp, roleId int) int {
	if k.UnionId != "" && k.UnionId != unionId {
		util.Error("不能打开非自己的线索")
		return model.ERR_CLEW_NOT_YOU_OPENED
	}
	if k.UnionId == unionId {
		util.Error("不能打开自己的线索")
		return model.ERR_CLEW_HAS_OPENED
	}
	if k.Type == model.CLEW_PASSWORD && k.Password != password {
		util.Error("密码错误")
		return model.ERR_CLEW_PASSWORD_WRONG
	}
	if userAp < k.Ap {
		util.Error("AP点不够了")
		return model.ERR_CLEW_AP_NOT_ENOUGH
	}
	//增加角色权限判断roleId
	if len(strings.TrimSpace(k.Opids)) > 0 {
		for _, tmp := range strings.Split(k.Opids, ",") { //可以搜索的角色
			util.Debug("opids:%v tmp:%v role:%v", k.Opids, tmp, roleId)
			id, _ := strconv.Atoi(tmp)
			if id != roleId {
				util.Error("限制:%v 才能打开", id)
				return model.ERR_CLEW_NOT_YOU_OPENED
			}
		}
	}

	if len(strings.TrimSpace(k.UnOpids)) > 0 {
		for _, tmp := range strings.Split(k.UnOpids, ",") { //不可以搜索的角色
			util.Debug("unopids:%v tmp:%v role:%v", k.UnOpids, tmp, roleId)
			id, _ := strconv.Atoi(tmp)
			if id == roleId {
				util.Error("限制:%v 不能打开", id)
				return model.ERR_CLEW_NOT_YOU_OPENED
			}
		}
	}

	return 0
}
