package release

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"fmt"
	"strconv"
)

//
//// Func: room.new handler
//func RoomNew(c *gin.Context) {
//	fmt.Println("RoomNew ...")
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	var user model.UserInfo
//	util.Logger(util.INFO_LEVEL, "RoomNew", "RoomNew ...")
//
//	// Get ScriptId & OpenId From Param
//	var req model.RoomNewReq
//	err := c.Bind(&req)
//
//	// Optional Fields List
//	var optionalFields []string
//
//	// Check Param
//	if !CheckParams(req, "RoomNew", err, optionalFields) {
//		util.Error("ERR_WRONG_FORMAT")
//		ResponseErr(model.ERR_WRONG_FORMAT, c)
//		return
//	}
//
//	//uid := req.OpenId
//	uid := req.UnionId
//	sid := req.ScriptId
//	user.UnionId = req.UnionId
//
//	// Create KV in Room Bucket: roomId -> roomInfo
//	// Generate RoomId
//	// Generate RoomInfo
//	rid := util.GetRandomString(32)
//
//	var unionIdSlice []string
//	unionIdSlice = append(unionIdSlice, uid)
//
//	playerSlice, code := GetPlayerInfoByScriptId(sid)
//	if code != model.ERR_OK {
//		//util.Logger(util.ERROR_LEVEL, "RoomNew", "Get Role Info By Script Id Err")
//		util.Error("RoomNew Get Role Info By Script Id Err")
//		ResponseErr(model.ERR_GET_ROLE_INFO, c)
//		return
//	}
//
//	rif := model.RoomInfo{
//		ScriptId: sid,
//		//OpenIdSlice: openIdSlice,
//		UnionIdSlice: unionIdSlice,
//		PlayerSlice:  playerSlice,
//	}
//
//	roomInfoEncoding, err := json.Marshal(rif)
//	if err != nil {
//		//util.Logger(util.ERROR_LEVEL, "RoomNew", "Room Info Encoding Err:"+err.Error())
//		util.Error("RoomNew Room Info Encoding ERROR[%v]", err.Error())
//	}
//
//	boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
//
//	// Update UserCache & RoomCache
//	global.SetUserCache(uid, rid)
//	global.AddUserToRoomCache(uid, rid)
//
//	// Send RoomEnter Message
//	//go websocket.SendRoomEnterMessage(unionIdSlice, uid)
//	go websocket.SendRoomEnterMessage(unionIdSlice, user)
//
//	// Return roomId & playerSlice
//	resp := model.RoomNewResp{}
//	resp.Params.RoomId = rid
//	resp.Params.RoomInfo = rif
//	c.JSON(http.StatusOK, resp)
//}

//type Screen struct {
//	Round   int
//	Content string
//}

func MatchScreenSliceModify(s *[]model.ScreenInfo, index int, value model.ScreenInfo) {
	rear := append([]model.ScreenInfo{}, (*s)[index+1:]...)
	*s = append(append((*s)[:index], value), rear...)
}

func GetPlayerInfoByScriptId(scriptId int) ([]model.PlayerInfo, int) {
	util.Info("GetPlayerInfoByScriptId:%v ...", scriptId)
	var playerInfoList []model.PlayerInfo
	var round, num int
	//var contents []string

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = scriptId

	fmt.Println("GetPlayerInfoByScriptId taskRequest:", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptInfo", taskRequest))
	if err != nil {
		return playerInfoList, model.ERR_TASK_JSON
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)

	fmt.Println("UnwrapObjectPackage:", dbcode, len(dbparams))
	switch dbcode {
	case global.ERR_DB_OK:
		if dbparams["Round"] != nil {
			round = int(dbparams["Round"].(float64))
			util.Debug("剧本总轮数:%v", round)
		}
		peoples := dbparams["People"].([]interface{})
		for _, p := range peoples {
			//util.Debug("people:%+v", p)
			var player model.PlayerInfo
			pm := p.(map[string]interface{})
			if pm["Id"] != nil {
				player.Role.Id = int(pm["Id"].(float64))
			}
			if pm["Name"] != nil {
				player.Role.Name = pm["Name"].(string)
			}
			if pm["Album"] != nil {
				player.Role.Album = pm["Album"].(string)
			}
			if pm["Tall"] != nil {
				player.Role.Tall = int(pm["Tall"].(float64))
			}
			if pm["Age"] != nil {
				player.Role.Age, _ = strconv.Atoi(pm["Age"].(string))
			}
			if pm["Sex"] != nil {
				player.Role.Sex = int(pm["Sex"].(float64))
			}
			if pm["Vote"] != nil {
				player.Role.Vote = pm["Vote"].(bool)
			}
			//能否被选择
			if pm["Choice"] != nil {
				player.Role.Choice = pm["Choice"].(bool)
			}

			if pm["About"] != nil {
				player.Role.About = pm["About"].(string)
			}
			//每轮行动点1,2
			if pm["Pofround"] != nil {
				player.Role.Pofround = pm["Pofround"].(string)
			}
			//是否是凶手
			if pm["Murderer"] != nil {
				player.Role.Murderer = pm["Murderer"].(bool)
			}
			player.VoteResult = false //默认投票结果是错误的

			//人物解析
			if pm["Final"] != nil {
				player.Role.Final = pm["Final"].(string)
			}

			//多幕
			if pm["Screens"] != nil {
				util.Debug("多幕:%+v", pm["Screens"])
				//contents = make([]string, round)
				//screenContents := model.ScreenInfo{}
				screenContents := make([]model.ScreenInfo, round+1)
				strs := pm["Screens"].([]interface{})
				for _, scr := range strs {
					//util.Debug("screen:%+v", scr)
					ps := scr.(map[string]interface{})
					if ps["Round"] != nil {
						//if ps["Content"] != nil {
						//	conent := ps["Content"].(string)
						//	if len(conent) > 0 {
						//		idx := int(ps["Round"].(float64))
						//		util.Debug("idx:%+v", idx)
						//		//contents[idx-1] = conent
						//		num += 1
						//	}
						//}
						screenBase := model.ScreenInfo{}
						screenBase.ScreenNum = int(ps["Round"].(float64))
						if ps["screens"] != nil {
							//util.Debug("screens:%+v", ps["screens"])
							contents := ps["screens"].([]interface{})
							//util.Debug("contents:%+v", contents)
							for _, x := range contents {
								//util.Debug("11111")
								tmp := x.(map[string]interface{})
								//util.Debug("11111")
								sq := model.ScreenQuestion{}
								if tmp["question"] != nil {
									//util.Debug("yyyyy:%+v", tmp["question"])
									//util.Debug("question:%+v", tmp["question"])
									sq.Question = tmp["question"].(string)
									//sq.Ssqid = int(tmp["ssqid"].(int64))
									if tmp["ssqid"] != nil {
										sq.Ssqid = int(tmp["ssqid"].(float64))
									}
									if tmp["flag"] != nil {
										sq.Flag = int(tmp["flag"].(float64))
									}

									//util.Debug("answers:%+v", tmp["answers"])
									if tmp["answers"] != nil {
										ass := tmp["answers"].([]interface{})
										for _, y := range ass {
											//util.Debug("yyyyy:%+v", y)
											tmp2 := y.(map[string]interface{})
											as := model.ScreenAnswer{}
											if tmp2["answer"] != nil {
												as.Answer = tmp2["answer"].(string)
											}
											if tmp2["screen"] != nil {
												as.Screen = tmp2["screen"].(string)
											}
											//as.Ssaid = int(tmp2["ssaid"].(int64))
											if tmp2["ssaid"] != nil {
												as.Ssaid = int(tmp2["ssaid"].(float64))
											}

											sq.Answers = append(sq.Answers, as)
										}
									}
								}
								num += 1
								util.Debug("num:%v", num)
								screenBase.Content = append(screenBase.Content, sq)
							}
							//util.Debug("screens:%v", content)
						}
						util.Debug("idx:%+v round:%v", screenBase.ScreenNum, round)
						MatchScreenSliceModify(&screenContents, round, screenBase)
					}
				}
				util.Debug("%+v", screenContents)
				if num > 0 {
					//player.ScreenInfo.ScreenNum = round
					//player.ScreenInfo.Content = contents
					player.Screens = append(player.Screens, screenContents...)
				}
			}

			playerInfoList = append(playerInfoList, player)
		}
		return playerInfoList, model.ERR_OK
	default:
		return playerInfoList, model.ERR_DEFAULT
	}
}
