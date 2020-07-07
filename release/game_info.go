package release

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"time"

	//"github.com/json-iterator/go"
	"net/http"
	"strconv"
	"strings"
)

// Func: game.info handler
func GameInfo(c *gin.Context) {
	util.Info("GameInfo ...")
	btime := time.Now().UnixNano()
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get RoomId and OpenId From Param
	var req model.GameInfoReq
	keywords := []model.KeyWord{}
	player := model.PlayerInfo{}
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "GameInfo", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	//uid := req.OpenId
	uid := req.UnionId
	rid := req.RoomId

	// Find Room In Room Cache
	//_, ok := global.RoomCache[rid]
	//if !ok {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	conn := gocache.RedisConnPool.Get()
	defer conn.Close()
	//ok, err := gocache.CheckRoomExists(rid)
	ok, err := gocache.ConnCheckRoomExists(conn, rid)
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
	//b := boltdb.View([]byte(rid), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	//查看用户是否在白名单
	isWhite := false
	whiteUsers, err := gocache.GetWhitelist()
	if err != nil {
		util.Error("获取白名单失败, ERROR:%v", err.Error())
	}
	util.Info("白名单数据:%+v", whiteUsers)
	for _, tmp := range whiteUsers {
		if tmp == uid {
			isWhite = true
		}
	}
	util.Info("是否在白名单:%v", isWhite)

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//de := json.Unmarshal(b, &roomInfo)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameInfo", "Decoding Room Info Err:"+de.Error())
	//	util.Error("GameInfo Decoding Room Info ERROR[%v]", de.Error())
	//}
	//读取房间信息
	//err = gocache.GetRoomInfo(rid, &roomInfo)
	err = gocache.ConnGetRoomInfo(conn, rid, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, roomInfo.PlayerSlice)

	// Find Role You Selected
	roleId := 0
	for _, v := range roomInfo.PlayerSlice {
		if v.UnionId == uid {
			roleId = v.Role.Id
			player = v
			break
		}
	}

	// Get Detail Info By RoleId and ScriptId
	if roleId == 0 {
		util.Error("ERR_ROLE_NOT_SELECT")
		ResponseErr(model.ERR_ROLE_NOT_SELECT, c)
		return
	}
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, roomInfo.PlayerSlice)

	gameInfoResp := model.GameInfoResp{}
	params := model.GameInfo{}

	about, task, clews, explore, exploresShip, aps, code := GetAboutAndTaskAndClew(conn, roleId, roomInfo.ScriptId)
	if code != model.ERR_OK {
		//util.Logger(util.ERROR_LEVEL, "GameInfo", "Get About & Task & Clew By ScriptId Err")
		util.Error("GameInfo Get About & Task & Clew By ScriptId Err")
		ResponseErr(model.ERR_GET_GAME_INFO, c)
		return
	}
	//var aps string
	//for _, tmp :=  range roomInfo.PlayerSlice{
	//	if roleId == tmp.Role.Id {
	//		aps = tmp.Role.Pofround
	//	}
	//}

	//story, ap, code := GetApAndStoryFromDB(conn, roomInfo.ScriptId)
	story, ap, code := GetApAndStoryFromDBV2(conn, roomInfo.ScriptId, roleId)
	if code != model.ERR_OK {
		//util.Logger(util.ERROR_LEVEL, "GameInfo", "Get AP & Story By ScriptId Err")
		util.Error("GameInfo Get AP & Story By ScriptId Err")
		ResponseErr(model.ERR_GET_GAME_INFO, c)
		return
	}
	util.Info("ap:%+v", ap)
	/* 兼容每轮AP点存剧本表的情况，剧本表无数据则取人物表AP数据 */
	if len(ap) == 0 {
		for i, tmp := range strings.Split(aps, ",") {
			if i == 0 {
				ap = append(ap, 0) //前置一轮AP点
			}
			if isWhite { //白名单用户
				ap = append(ap, 9999)
			} else {
				num, _ := strconv.Atoi(tmp)
				ap = append(ap, num)
			}
		}
	} else {
		if isWhite { //白名单用户
			for i, _ := range ap {
				ap[i] = 9998
			}
		}
	}

	/* 关键词增加人物特定关键词、公共关键词 */
	for _, tmp := range explore {
		if tmp.SpId == 0 || tmp.SpId == roleId {
			keywords = append(keywords, tmp)
		}
	}
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, roomInfo.PlayerSlice)

	params.Ap = ap
	params.About = about
	params.Task = task
	params.Clew = clews
	params.Story = story
	//params.Explore = explore
	params.Explore = keywords
	params.ExploreRelation = exploresShip
	gameInfoResp.Params.GameInfo = params

	//删除数据轮次按script表配置
	//round := len(gameInfoResp.Params.Clew)

	// Find RoomId in ApBucket
	//ab := boltdb.View([]byte(rid), "ApBucket")
	//if ab == nil {
	//apMap := make(map[string][]int)
	//	apMap[uid] = ap
	//	encodingApMap, _ := json.Marshal(apMap)
	//	boltdb.CreateOrUpdate([]byte(rid), encodingApMap, "ApBucket")
	//} else {
	//	//var apMap map[string][]int
	//	//json.Unmarshal(ab, &apMap)
	//	_, ok := apMap[uid]
	//	if !ok {
	//		apMap[uid] = ap
	//		encodingApMap, _ := json.Marshal(apMap)
	//		boltdb.CreateOrUpdate([]byte(rid), encodingApMap, "ApBucket")
	//	}
	//}

	//apMap, bl, err := gocache.GetAPInfo(rid)

	apMap, bl, err := gocache.ConnGetAPInfo(conn, rid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_GET_GAME_INFO, c)
		return
	}
	if bl {
		_, ok := apMap[uid]
		if !ok {
			apMap[uid] = ap
			util.Debug("apMap:%+v", apMap)
			//err = gocache.SetAPInfo(rid, apMap)
			err = gocache.ConnSetAPInfo(conn, rid, apMap)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_GET_GAME_INFO, c)
				return
			}
		}
	} else {
		//apMap2 := make(map[string][]int)
		//apMap2[uid] = ap
		//util.Debug("apMap:%+v", apMap2)
		//err = gocache.ConnSetAPInfo(conn, rid, apMap2)
		//if err != nil {
		//	util.Error("ERROR:%v", err.Error())
		//	ResponseErr(model.ERR_GET_GAME_INFO, c)
		//	return
		//}
		//同步用户AP点 add by skc at 2020-05-21
		apMap2 := map[string][]int{}
		for _, player := range roomInfo.PlayerSlice {
			util.Debug("%v %v", player.UnionId, player.Role.Pofround)
			ap2 := []int{}
			for i, tmp := range strings.Split(player.Role.Pofround, ",") {
				if i == 0 {
					ap2 = append(ap2, 0) //前置一轮AP点
				}
				if isWhite { //白名单用户
					ap2 = append(ap2, 9999)
				} else {
					num, _ := strconv.Atoi(tmp)
					ap2 = append(ap2, num)
				}
			}
			apMap2[player.UnionId] = ap2
		}
		util.Debug("AP: %+v", apMap2)
		err = gocache.ConnSetAPInfo(conn, rid, apMap2)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_GET_GAME_INFO, c)
			return
		}
	}
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, roomInfo.PlayerSlice)

	// Find RoomId in Game Bucket
	//gb := boltdb.View([]byte(rid), "GameBucket")
	//if gb == nil {
	//	clewInfoEncoding, _ := json.Marshal(gameInfoResp.Params.Clew)
	//	boltdb.CreateOrUpdate([]byte(rid), clewInfoEncoding, "GameBucket")
	//}
	util.Debug("GameBucket...")
	var tmps []model.Clew
	//bl, err = gocache.GetRoomGame(rid, &tmps)
	bl, err = gocache.ConnGetRoomGame(conn, rid, &tmps)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
		return
	}
	if !bl {
		//err = gocache.SetRoomGame(rid, gameInfoResp.Params.Clew)
		err = gocache.ConnSetRoomGame(conn, rid, gameInfoResp.Params.Clew)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
			return
		}
	}
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, roomInfo.PlayerSlice)

	// Add Clews To Cache If it not exist
	//if _, ok := global.ClewCache[rid]; !ok {
	//	global.AddClewToCache(BuildClewMap(clews), rid)
	//}

	//保存线索数据
	//err = gocache.SetRoomClew(rid, BuildClewMap(clews))
	err = gocache.ConnSetRoomClew(conn, rid, BuildClewMap(clews))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
		return
	}

	// Find RoomId in Stage Bucket
	//sb := boltdb.View([]byte(rid), "StageBucket")
	//if sb == nil {
	//	onClickMap := make(map[int][]string)
	//	for i := 1; i <= round; i++ {
	//		onClickMap[i] = []string{}
	//	}
	//	stageInfo := model.StageInfo{
	//		Round:   1,
	//		OnClick: onClickMap,
	//	}
	//	stageInfoEncoding, _ := json.Marshal(stageInfo)
	//	boltdb.CreateOrUpdate([]byte(rid), stageInfoEncoding, "StageBucket")
	//}
	util.Debug("StageBucket...")
	var tmpStage model.StageInfo
	//bl, err = gocache.GetRoomStage(rid, &tmpStage)
	bl, err = gocache.ConnGetRoomStage(conn, rid, &tmpStage)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
		return
	}
	if !bl { //不存在重新赋值
		onClickMap := make(map[int][]string)
		//for i := 1; i <= round; i++ {
		util.Info("房间剧本轮次:%v", roomInfo.Round)
		//for i := 1; i <= roomInfo.Round; i++ {
		for i := 0; i <= roomInfo.Round; i++ { //前置一轮给
			onClickMap[i] = []string{}
		}
		stageInfo := model.StageInfo{
			//Round: 1,
			Round:   0, //默认从第0轮开始
			OnClick: onClickMap,
		}
		//err = gocache.SetRoomStage(rid, stageInfo)
		err = gocache.ConnSetRoomStage(conn, rid, stageInfo)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
			return
		}
	}
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, roomInfo.PlayerSlice)

	//增加返回剧本人物信息
	gameInfoResp.Params.Player = player

	//util.Logger(util.INFO_LEVEL, "GameInfo", "story:"+gameInfoResp.Params.Story)
	//util.Logger(util.INFO_LEVEL, "GameInfo", "about:"+gameInfoResp.Params.About)
	//util.Info("GameInfo story:", gameInfoResp.Params.Story)
	//util.Info("GameInfo about:", gameInfoResp.Params.About)
	//util.Info("响应信息:%+v", gameInfoResp)
	gameInfoResp.Params.GameStage = tmpStage.Round
	//util.Info("耗时:%v 响应信息:%+v", time.Now().UnixNano()-btime, gameInfoResp)
	util.Info("耗时:%v 响应信息AP:%+v", time.Now().UnixNano()-btime, gameInfoResp.Params.GameInfo.Ap)

	// Return Story Task About and Clew
	c.JSON(http.StatusOK, gameInfoResp)
}

//获取剧本任务、线索
func GetAboutAndTaskAndClew(conn redis.Conn, roleId, scriptId int) (string, model.MainTask, []model.Clew, []model.KeyWord, []model.KeyWordRelation, string, int) {
	util.Info("GetAboutAndTaskAndClew ... ")
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	var about string
	task := model.MainTask{}
	var clews []model.Clew
	var explores []model.KeyWord
	var exploresShip []model.KeyWordRelation

	var aps string
	dbcode := global.ERR_DB_OK

	//优先缓存数据
	//dbparams, ok, err := gocache.GetAboutAndTaskAndClewAndExplores(scriptId)
	dbparams, ok, err := gocache.ConnGetAboutAndTaskAndClewAndExplores(conn, scriptId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return about, task, clews, explores, exploresShip, aps, model.ERR_TASK_JSON
	}
	util.Debug("ok:%v", ok)

	if !ok {
		taskRequest := make(map[string]interface{})
		taskRequest["ScriptId"] = scriptId

		dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptInfo", taskRequest))
		//dbResult, err := global.Task.TaskJsonNew(global.NewDBRequest("db.ScriptInfo", taskRequest))
		//dbResult, err := zeromq.TaskJsonComm(sock, global.NewDBRequest("db.ScriptInfo", taskRequest))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return about, task, clews, explores, exploresShip, aps, model.ERR_TASK_JSON
		}

		dbcode, dbparams = global.UnwrapObjectPackage(dbResult)
	}
	switch dbcode {
	case global.ERR_DB_OK:
		if !ok {
			//err = gocache.SetAboutAndTaskAndClewAndExplores(scriptId, dbparams)
			err = gocache.ConnSetAboutAndTaskAndClewAndExplores(conn, scriptId, dbparams)
			if err != nil {
				util.Error("ERROR:%v", err)
				return about, task, clews, explores, exploresShip, aps, model.ERR_TASK_JSON
			}
		}

		peoples := dbparams["People"].([]interface{})
		for _, p := range peoples {
			pm := p.(map[string]interface{})
			if (pm["Id"] != nil) && (int(pm["Id"].(float64)) == roleId) {
				if pm["Story"] != nil {
					about = pm["Story"].(string)
				}

				if pm["Pofround"] != nil {
					aps = pm["Pofround"].(string)
				}

				if pm["Infomation"] != nil {
					task.Main = pm["Infomation"].(string)
				}

				if pm["Task"] != nil {
					t := pm["Task"].([]interface{})
					var sideTask []model.SideTask
					for _, v := range t {
						raw_task, ok := v.(map[string]interface{})
						if !ok {
							continue
						}
						sTask := model.SideTask{}

						sTask.Id = int(raw_task["Id"].(float64))
						sTask.Caption = raw_task["Caption"].(string)

						var options []model.Option
						raw_options, ok := raw_task["Options"].([]interface{})
						if ok {
							for _, o := range raw_options {
								raw_option, ok := o.(map[string]interface{})
								if ok {
									option := model.Option{}
									option.Id = int(raw_option["Id"].(float64))
									option.Option = raw_option["Content"].(string)

									options = append(options, option)
								}
							}
							sTask.Options = options
						}

						sideTask = append(sideTask, sTask)
					}
					task.Side = sideTask
				}
				break
			}
		}

		raw_clew, ok := dbparams["Clew"].([]interface{})
		if ok {
			for i, v := range raw_clew {
				clew := model.Clew{}
				clew.Round = i + 1

				keys := make(map[string]model.Key)
				for k, sv := range v.(map[string]interface{}) {
					svm, ok := sv.(map[string]interface{})
					if ok {
						key := model.Key{}

						album, ok := svm["Album"].(string)
						if ok {
							key.Album = album
						}

						clews, ok := svm["Clew"].([]interface{})
						if ok && clews != nil {
							key.TotalNum = len(clews)

							var key_list []model.MainKeyDetail
							for _, c := range clews {
								//util.Debug("线索:%+v", c)
								vm, ok := c.(map[string]interface{})
								if ok {
									mkd := model.MainKeyDetail{}
									kd := model.KeyDetail{}

									//var kc []model.KeyContent
									////err = json.Unmarshal([]byte(vm["Content"].(string)), &kc)
									//cont := strings.ReplaceAll(vm["Content"].(string), "\n", "\\n")
									//err = json.Unmarshal([]byte(cont), &kc)
									//if err != nil {
									//	util.Error("ERROR:%v \n 线索:%+v", err.Error(), vm["Content"].(string))
									//}

									kd.Id = int(vm["Id"].(float64))
									kd.Ap = int(vm["AP"].(float64))
									kd.Type = int(vm["Type"].(float64)) + 1
									//kd.Content = kc
									kd.Content = vm["Content"].(string)
									question, ok := vm["Question"].(string)
									if ok {
										kd.Question = question
									}
									password, ok := vm["Password"].(string)
									if ok {
										kd.Password = password
									}
									if kd.Ap == 0 {
										kd.Status = 2
									}
									//if kd.Id == 11744 {
									//	util.Debug("11744:%+v \n %+v", kc, vm["Content"].(string))
									//}
									title, ok := vm["Title"].(string)
									if ok {
										kd.Title = title
									}
									opids, ok := vm["Opids"].(string)
									if ok {
										kd.Opids = opids
									}
									unopids, ok := vm["UnOpids"].(string)
									if ok {
										kd.UnOpids = unopids
									}

									mkd.KeyDetail = kd

									sub, ok := vm["Sub"].([]interface{})
									if ok && sub != nil {
										var subs []model.SubKeyDetail

										for _, s := range sub {
											sm, ok := s.(map[string]interface{})
											if ok {
												skd := model.SubKeyDetail{}
												kd := model.KeyDetail{}

												//var kc []model.KeyContent
												////util.Info("12503:%+v", sm["Content"])
												//
												//err = json.Unmarshal([]byte(strings.ReplaceAll(sm["Content"].(string), "\n", "\\n")), &kc)
												//if err != nil {
												//	util.Error("ERROR:%v \n ID:%v 线索:%+v", err.Error(), sm["Id"], sm["Content"])
												//}
												kd.Id = int(sm["Id"].(float64))
												kd.Ap = int(sm["AP"].(float64))
												kd.Type = int(sm["Type"].(float64)) + 1
												//kd.Content = kc
												kd.Content = sm["Content"].(string)
												question, ok := sm["Question"].(string)
												if ok {
													kd.Question = question
												}
												password, ok := sm["Password"].(string)
												if ok {
													kd.Password = password
												}
												if kd.Ap == 0 {
													kd.Status = 2
												}
												title, ok := sm["Title"].(string)
												if ok {
													kd.Title = title
												}
												if opids, ok := sm["Opids"].(string); ok {
													kd.Opids = opids
												}
												if unopids, ok := sm["UnOpids"].(string); ok {
													kd.UnOpids = unopids
												}

												sub3, ok := sm["Sub"].([]interface{})
												if ok && sub3 != nil {
													util.Debug("sub3:%v", sub3)
													var subs3 []model.SubKeyDetail

													for _, tmp := range sub3 {
														sm3, ok := tmp.(map[string]interface{})
														if ok {
															skd3 := model.SubKeyDetail{}
															kd3 := model.KeyDetail{}

															//var kc3 []model.KeyContent
															//err = json.Unmarshal([]byte(strings.ReplaceAll(sm3["Content"].(string), "\n", "\\n")), &kc3)
															//if err != nil {
															//	util.Error("ERROR:%v \n ID:%v 线索:%+v", err.Error(), sm3["Id"], sm3["Content"])
															//}

															kd3.Id = int(sm3["Id"].(float64))
															kd3.Ap = int(sm3["AP"].(float64))
															kd3.Type = int(sm3["Type"].(float64)) + 1
															//kd3.Content = kc3
															kd3.Content = sm3["Content"].(string)
															question, ok := sm3["Question"].(string)
															if ok {
																kd3.Question = question
															}
															password, ok := sm3["Password"].(string)
															if ok {
																kd3.Password = password
															}
															if kd.Ap == 0 {
																kd3.Status = 2
															}
															title, ok := sm3["Title"].(string)
															if ok {
																kd3.Title = title
															}
															if opids, ok := sm3["Opids"].(string); ok {
																kd3.Opids = opids
															}
															if unopids, ok := sm3["UnOpids"].(string); ok {
																kd3.UnOpids = unopids
															}
															skd3.KeyDetail = kd3

															subs3 = append(subs3, skd3)
														}
													}

													util.Debug("sub s3:%v", subs3)
													skd.Sub = subs3
												}

												skd.KeyDetail = kd

												subs = append(subs, skd)
											}
										}
										mkd.Sub = subs
									}

									key_list = append(key_list, mkd)
								}
							}

							key.KeyList = key_list
						}

						keys[k] = key
					}
				}
				clew.Key = keys
				clews = append(clews, clew)
			}
		}

		keywords, ok := dbparams["Explore"].([]interface{})
		if ok {
			for _, val := range keywords {
				var keyword model.KeyWord
				//util.Debug("keyword:%v", val)
				pm := val.(map[string]interface{})
				if pm["skid"] != nil {
					keyword.Skid = int(pm["skid"].(float64))
				}
				if pm["round"] != nil {
					keyword.Round = int(pm["round"].(float64))
				}
				if pm["keyword"] != nil {
					keyword.Keyword = pm["keyword"].(string)
				}
				if pm["content"] != nil {
					keyword.Content = pm["content"].(string)
				}
				if pm["album"] != nil {
					//keyword.Content += pm["album"].(string)
					keyword.Album = pm["album"].(string)
				}
				if pm["spid"] != nil {
					keyword.SpId = int(pm["spid"].(float64))
				}
				if pm["end_id"] != nil {
					keyword.EndId = int(pm["end_id"].(float64))
				}
				if pm["resume"] != nil {
					keyword.Resume = pm["resume"].(string)
				}
				explores = append(explores, keyword)
			}
		}

		keywordsRelation, ok := dbparams["ExploreRelation"].([]interface{})
		if ok {
			for _, val := range keywordsRelation {
				//a.sid", "a.pid", "a.round", "a.skids", "a.end_id", "e.resume", "e.remark"
				var keywordShip model.KeyWordRelation
				//util.Debug("keyword:%v", val)
				pm := val.(map[string]interface{})
				if pm["sid"] != nil {
					keywordShip.Sid = int(pm["sid"].(float64))
				}
				if pm["pid"] != nil {
					keywordShip.Pid = int(pm["pid"].(float64))
				}
				if pm["round"] != nil {
					keywordShip.Round = int(pm["round"].(float64))
				}
				if pm["skids"] != nil {
					keywordShip.Skids = pm["skids"].(string)
				}
				if pm["end_id"] != nil {
					keywordShip.EndId = int(pm["end_id"].(float64))
				}
				if pm["resume"] != nil {
					keywordShip.Resume = pm["resume"].(string)
				}
				if pm["remark"] != nil {
					keywordShip.Remark = pm["remark"].(string)
				}
				exploresShip = append(exploresShip, keywordShip)
			}
		}

		return about, task, clews, explores, exploresShip, aps, model.ERR_OK
	default:
		return about, task, clews, explores, exploresShip, aps, model.ERR_DEFAULT
	}
}

func GetApAndStoryFromDB(conn redis.Conn, scriptId int) (string, []int, int) {
	util.Info("GetApAndStoryFromDB ...")
	var story string
	var final string
	var ap []int

	//优先缓存数据
	//ap, story, ok, err := gocache.GetApAndStory(scriptId)
	ap, story, final, ok, err := gocache.ConnGetApAndStory(conn, scriptId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return story, ap, model.ERR_TASK_JSON
	}
	util.Debug("ok:%v", ok)

	if ok {
		//util.Debug("OK:ap:%v story:%v", ap, story)
		return story, ap, model.ERR_OK
	}

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = scriptId
	taskRequest["UserId"] = 0

	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptSimpleGet", taskRequest))
	//dbResult, err := global.Task.TaskJsonNew(global.NewDBRequest("db.ScriptSimpleGet", taskRequest))
	//dbResult, err := zeromq.TaskJsonNew(global.NewDBRequest("db.ScriptInfo", taskRequest))
	//dbResult, err := zeromq.TaskJsonComm(sock, global.NewDBRequest("db.ScriptSimpleGet", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return story, ap, model.ERR_TASK_JSON
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		//if dbparams["Content"] != nil {
		//	story = dbparams["Content"].(string)
		//}
		if dbparams["About"] != nil {
			story = dbparams["About"].(string)
		}
		if dbparams["PointOfRound"] != nil {
			ap = append(ap, 0)
			apSli := strings.Split(dbparams["PointOfRound"].(string), ",")
			for _, s := range apSli {
				i, err := strconv.Atoi(s)
				if err != nil {
					i = 0
				}
				if i == 0 && len(apSli) == 1 {
					util.Info("本轮ap点为0，调过不处理!")
					continue
				}
				ap = append(ap, i)
			}
		}
		util.Info("AP:%+v", ap)
		if dbparams["Final"] != nil {
			final = dbparams["Final"].(string)
		}
		//err = gocache.SetApAndStory(scriptId, ap, story)
		err = gocache.ConnSetApAndStory(conn, scriptId, ap, story, final)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return story, ap, model.ERR_DEFAULT
		}
		return story, ap, model.ERR_OK
	default:
		return story, ap, model.ERR_DEFAULT
	}
}

func BuildClewMap(clews []model.Clew) map[string][]global.MainClew {
	ret := make(map[string][]global.MainClew)

	for _, c := range clews {
		for k, v := range c.Key {
			for _, j := range v.KeyList {
				var mains global.MainClew
				mains.Id = j.Id
				mains.Round = c.Round

				subs := []global.ClewRound{}
				for _, s := range j.Sub {
					var sub global.ClewRound
					sub.Round = c.Round
					sub.Id = s.Id

					// add by skc at 2020-04-17 begin
					for _, m := range s.Sub {
						var sub3 global.ClewRound
						sub3.Round = c.Round
						sub3.Id = m.Id
						sub.Sub = append(sub.Sub, sub3)
					}
					// add by skc at 2020-04-17 end

					subs = append(subs, sub)
				}
				mains.Sub = subs

				ret[k] = append(ret[k], mains)
			}
		}
	}

	return ret
}

//响应任务的 story, ap, code
func GetApAndStoryFromDBV2(conn redis.Conn, scriptId, roleId int) (string, []int, int) {
	util.Info("GetApAndStoryFromDB ...")
	var story string
	var ap []int

	info := model.ScriptPeopleInfo{}

	//优先缓存数据
	ok, err := gocache.ConnGetScriptPeopleInfo(conn, scriptId, roleId, &info)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return story, ap, model.ERR_TASK_JSON
	}
	util.Debug("ok:%v", ok)

	if ok {
		//util.Debug("OK:ap:%v story:%v", ap, story)
		story = info.About
		for i, tmp := range strings.Split(info.Ap, ",") {
			if i == 0 {
				ap = append(ap, 0) //前置一轮AP点
			}
			num, _ := strconv.Atoi(tmp)
			ap = append(ap, num)
		}
		return story, ap, model.ERR_OK
	}

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = scriptId
	taskRequest["RoleId"] = roleId
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptPeopleBaseInfo", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return story, ap, model.ERR_TASK_JSON
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)

	if dbcode == global.ERR_DB_OK {
		err = gocache.ConnSetScriptPeopleInfo(conn, scriptId, roleId, dbparams)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return story, ap, model.ERR_DEFAULT
		}
		if dbparams["about"] != nil {
			story = dbparams["about"].(string)
		}

		if dbparams["ap"] != nil {
			info.Ap = dbparams["ap"].(string)
			for i, tmp := range strings.Split(info.Ap, ",") {
				if i == 0 {
					ap = append(ap, 0) //前置一轮AP点
				}
				num, _ := strconv.Atoi(tmp)
				ap = append(ap, num)
			}
		}
		util.Info("story:%v ap:%v", story, ap)
	}

	return story, ap, model.ERR_OK
}
