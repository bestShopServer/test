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
	"time"
)

func MatchDataSliceModify(s *[]model.UserInfo, index int, value model.UserInfo) {
	rear := append([]model.UserInfo{}, (*s)[index+1:]...)
	*s = append(append((*s)[:index], value), rear...)
}

//玩家重连游戏
func GameReconnect(c *gin.Context) {
	util.Info("GameReconnect ...")
	btime := time.Now().UnixNano()

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	user := model.UserInfo{}
	keywords := []model.KeyWord{}

	// Get RoomId and OpenId From Param
	var req model.GameReconnectReq
	err := c.Bind(&req)

	// Optional Fields List
	//var optionalFields []string
	optionalFields := []string{}

	// Check Param
	if !CheckParams(req, "GameReconnect", err, optionalFields) {
		util.Error("报文格式有误:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	//uid := req.OpenId
	uid := req.UnionId
	rid := req.RoomId
	user.UnionId = req.UnionId
	user.Name = req.UserName
	util.Info("用户信息:%+v", user)

	// If no user cache
	//cache_room_id, ok := global.UserCache[uid]
	//if !ok {
	//	util.Error("ERR_GAME_HAS_OVER")
	//	ResponseErr(model.ERR_GAME_HAS_OVER, c)
	//	return
	//}
	util.Debug("耗时:%v", time.Now().UnixNano()-btime)

	conn := gocache.RedisConnPool.Get()
	defer conn.Close()
	//cache_room_id, err := gocache.GetUserRoom(uid)
	cache_room_id, err := gocache.ConnGetUserRoom(conn, uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_GAME_HAS_OVER, c)
		return
	}
	util.Debug("cache_room_id:%v", cache_room_id)

	// If user cache not eq room id
	if cache_room_id != rid {
		util.Error("ERR_ROOM_LINK")
		ResponseErr(model.ERR_ROOM_LINK, c)
		return
	}

	// If Game Has Over
	//_, ok = global.RoomCache[rid]
	//if !ok {
	//	util.Error("ERR_GAME_HAS_OVER")
	//	ResponseErr(model.ERR_GAME_HAS_OVER, c)
	//	return
	//}
	util.Debug("耗时:%v", time.Now().UnixNano()-btime)

	//ok, err := gocache.CheckRoomExists(rid)
	ok, err := gocache.ConnCheckRoomExists(conn, rid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Debug("CheckRoomExists:%v", ok)

	if !ok {
		util.Error("ERR_ROOM_NOT_EXIST")
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find KV in Room Bucket
	//b := boltdb.View([]byte(rid), "RoomBucket")
	//if b == nil {
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Room Info
	rif := model.RoomInfo{}
	//de := json.Unmarshal(b, &rif)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameReconnect", "Decoding Room Info Err:"+de.Error())
	//	util.Error("GameReconnect Decoding Room Info ERROR[%v]", de.Error())
	//}
	util.Debug("耗时:%v, 轮次:%v", time.Now().UnixNano()-btime, rif.Round)

	//读取房间信息
	//err = gocache.GetRoomInfo(rid, &rif)
	err = gocache.ConnGetRoomInfo(conn, rid, &rif)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Debug("房间信息剧本:%+v 轮次:%v", rif.ScriptId, rif.Round)

	//房间游戏已结束，则无法再加入房间
	if rif.Status == 2 { //标识房间已结束
		util.Error("房间:%v 状态:%v", rid, rif.Status)
		ResponseErr(model.ERR_ROOM_LINK, c)
		return
	}

	// If openIdSlice Contains openId
	exist := false
	for _, v := range rif.UnionIdSlice {
		util.Debug("房间用户:%v", v)
		if v == uid {
			exist = true
			break
		}
	}
	util.Debug("耗时:%v 轮次:%v", time.Now().UnixNano()-btime, rif.Round)

	// Add openId to openIdSlice
	if !exist {
		util.Info("处理用户不存在情况...")
		// If openIdSlice length eq playerSlice
		if len(rif.UnionIdSlice) == len(rif.PlayerSlice) {
			util.Error("ERR_ROOM_IS_FULL")
			ResponseErr(model.ERR_ROOM_IS_FULL, c)
			return
		}
		rif.UnionIdSlice = append(rif.UnionIdSlice, uid)
		//rif.UserInfoSlice = append(rif.UserInfoSlice, use)
		//util.Debug("新增用户:%+v", use)
	}
	util.Debug("耗时:%v 轮次:%v", time.Now().UnixNano()-btime, rif.Round)

	userExist := false
	//同步用户名称 user name 取上送的值
	for i, u := range rif.UserInfoSlice {
		util.Debug("房间的用户信息:%v", u.UnionId)
		if u.UnionId == uid {
			userExist = true
			user.CoverVoteNum = u.CoverVoteNum
			user.VoteUser = u.VoteUser
			user.IsPay = u.IsPay
			//add by skc at 2020-04-27
			//查看用户是否在白名单
			isWhite := false
			whiteUsers, err := gocache.GetWhitelist()
			if err != nil {
				util.Error("获取白名单失败, ERROR:%v", err.Error())
			}
			util.Info("白名单数据:%+v", whiteUsers)
			for _, tmp := range whiteUsers {
				if tmp == req.UnionId {
					isWhite = true
				}
			}
			util.Info("是否在白名单:%v", isWhite)
			if isWhite {
				user.Member = 1
				util.Debug("会员级别:%v", user.Member)
			} else {
				//user.Member = u.Member
				res, err := GetUserMemberBase(req.UnionId)
				if err != nil {
					util.Error("ERROR:%v", err.Error())
					ResponseErr(model.ERR_DEFAULT, c)
					return
				}
				//util.Debug("会员信息:%+v", res)
				user.Member = res.Member
			}
			user.Surrender = u.Surrender
			util.Debug("同步用户信息[%+v]", user)
			MatchDataSliceModify(&rif.UserInfoSlice, i, user)
			break
		}
	}
	util.Debug("用户信息是否在房间内:%v", userExist)
	util.Debug("耗时:%v, 轮次:%v", time.Now().UnixNano()-btime, rif.Round)

	//用户信息不存在时，添加用户混淆
	//if !userExist {
	util.Info("处理用户信息不存在情况...")
	util.Info("剧本价格:%f", rif.Price)
	//if rif.Price > 0.001 {
	if rif.Price > 0.001 && user.IsPay != 2 { //付款标识 1 未付款 2 已付款
		var cost model.ScriptUserCostReq
		cost.ScriptId = rif.ScriptId
		cost.UnionIds = append(cost.UnionIds, req.UnionId)

		//查询付款
		us, code := CheckUserScriptCost(cost)
		util.Debug("查询用户剧本付费情况,code:%v", code)
		if code != model.ERR_OK && code != global.ERR_DB_NOTFOUND_DATA {
			util.Error("Get User ScriptId Err [%v]", code)
			ResponseErr(model.ERR_CHECK_SCRIPT_COST, c)
			return
		}
		if code == global.ERR_DB_NOTFOUND_DATA { //未付款
			user.IsPay = 1
		} else {
			user.IsPay = us[0].CostFlag
		}
	}
	util.Debug("耗时:%v", time.Now().UnixNano()-btime)
	/* add by skc at 2020-04-24 begin */
	if userExist { //用户存在
		//同步用户付款状态
		for i, u := range rif.UserInfoSlice {
			if u.UnionId == uid {
				util.Debug("修改房间的用户信息:%v", u.UnionId)
				MatchDataSliceModify(&rif.UserInfoSlice, i, user)
				break
			}
		}
	} else {
		rif.UserInfoSlice = append(rif.UserInfoSlice, user)
		util.Debug("新增用户信息:%+v", user)
	}
	/* add by skc at 2020-04-24 end */
	//}
	util.Debug("耗时:%v 轮次:%v", time.Now().UnixNano()-btime, rif.Round)

	// Find RoleId
	var roleId int
	for _, p := range rif.PlayerSlice {
		//if p.OpenId == uid {
		if p.UnionId == uid {
			roleId = p.Role.Id
			break
		}
	}
	util.Info("处理存储数据...")

	//// Update BoltDB
	//roomInfoEncoding, err := json.Marshal(rif)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameReconnect", "Room Info Encoding Err:"+err.Error())
	//	util.Error("GameReconnect Room Info Encoding ERROR[%v]", err.Error())
	//
	//}
	//boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
	util.Debug("耗时:%v, 轮次:%v, 多幕:%+v", time.Now().UnixNano()-btime, rif.Round, rif.PlayerSlice)

	//err = gocache.SetRoomInfo(rid, rif)
	err = gocache.ConnSetRoomInfo(conn, rid, rif)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find Other Bucket
	//gb := boltdb.View([]byte(rid), "GameBucket")
	//sb := boltdb.View([]byte(rid), "StageBucket")
	//vb := boltdb.View([]byte(rid), "VoteBucket")
	//ab := boltdb.View([]byte(rid), "ApBucket")

	reconnectInfoResp := model.ReconnectResp{}
	params := model.ReconnectInfo{}

	params.RoomInfo = rif
	sendGameVote := false
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, rif.PlayerSlice)

	// If Game Has Start
	if roleId != 0 {
		util.Info("房间已经开始游戏...")
		//sock, err := zeromq.InitZeroMQOneClient()
		//if err != nil {
		//	util.Error("ERROR:%v", err.Error())
		//	ResponseErr(model.ERR_GET_GAME_INFO, c)
		//	return
		//}
		//defer sock.Close()

		// Fetch Game Info
		gif := model.GameInfo{}

		util.Debug("耗时:%v", time.Now().UnixNano()-btime)
		about, task, clews, explore, exploresShip, aps, code := GetAboutAndTaskAndClew(conn, roleId, rif.ScriptId)
		if code != model.ERR_OK {
			//util.Logger(util.ERROR_LEVEL, "GameReconnect", "Get About & Task & Clew By ScriptId Err")
			util.Error("GameReconnect: Get About & Task & Clew By ScriptId Err:%v", code)
			ResponseErr(model.ERR_GET_GAME_INFO, c)
			return
		}
		util.Debug("耗时:%v", time.Now().UnixNano()-btime)
		//story, ap, code := GetApAndStoryFromDB(conn, rif.ScriptId)
		story, ap, code := GetApAndStoryFromDBV2(conn, rif.ScriptId, roleId)
		if code != model.ERR_OK {
			//util.Logger(util.ERROR_LEVEL, "GameReconnect", "Get AP & Story By ScriptId Err")
			util.Error("GameReconnect Error: Get AP & Story By ScriptId Err")
			ResponseErr(model.ERR_GET_GAME_INFO, c)
			return
		}
		util.Info("ap:%+v", ap)
		/* 兼容每轮AP点存剧本表的情况，剧本表无数据则取人物表AP数据 */
		if len(ap) == 0 {
			for _, tmp := range strings.Split(aps, ",") {
				num, _ := strconv.Atoi(tmp)
				ap = append(ap, num)
			}
		}
		/* 关键词增加人物特定关键词、公共关键词 */
		for _, tmp := range explore {
			if tmp.SpId == 0 || tmp.SpId == roleId {
				keywords = append(keywords, tmp)
			}
		}

		gif.Clew = clews
		gif.About = about
		gif.Ap = ap
		gif.Story = story
		gif.Task = task
		//gif.Explore = explore
		gif.Explore = keywords
		gif.ExploreRelation = exploresShip

		//if gb != nil {
		//	json.Unmarshal(gb, &clews)
		//	// ClewTransform(clews, uid)
		//	gif.Clew = clews
		//}
		util.Debug("耗时:%v", time.Now().UnixNano()-btime)
		util.Debug("GetRoomGame...")
		//bl, err := gocache.GetRoomGame(rid, &clews)
		bl, err := gocache.ConnGetRoomGame(conn, rid, &clews)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
			return
		}
		if bl {
			gif.Clew = clews
		}

		var stage int
		//if sb != nil {
		//	sif := model.StageInfo{}
		//	json.Unmarshal(sb, &sif)
		//	util.Debug("stage info:%+v", sif)
		//	stage = sif.Round
		//	params.GameStage = stage
		//	params.GameClew, _ = GetGameClew(rid, stage)
		//}
		util.Debug("耗时:%v", time.Now().UnixNano()-btime)
		sif := model.StageInfo{}
		util.Debug("GetRoomStage...")
		//bl, err = gocache.GetRoomStage(rid, &sif)
		bl, err = gocache.ConnGetRoomStage(conn, rid, &sif)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
			return
		}
		if bl {
			util.Debug("stage info:%+v", sif)
			stage = sif.Round
			params.GameStage = stage
			params.GameClew, _ = GetGameClew(rid, stage)

			//返回前端，当前阶段用户是否已点击下一阶段
			params.StageClick = false
			for _, unid := range sif.OnClick[stage] {
				if unid == uid {
					params.StageClick = true
					break
				}
			}
		}

		//if ab != nil {
		//	var apMap map[string][]int
		//
		//	err = json.Unmarshal(ab, &apMap)
		//	if err != nil {
		//		//util.Logger(util.ERROR_LEVEL, "GameReconnect", err.Error())
		//		util.Error("GameReconnect ERROR[%v]", err.Error())
		//	}
		//
		//	ap, ok := apMap[uid]
		//	if !ok {
		//		ap = gif.Ap
		//	}
		//
		//	if stage > 1 {
		//		var leftAp int
		//		s := util.Min(stage-1, len(ap)-1)
		//		for i := 0; i <= s; i++ {
		//			leftAp = leftAp + ap[i]
		//			ap[i] = 0
		//		}
		//		ap[s] = leftAp + ap[s]
		//	}
		//
		//	gif.Ap = ap
		//	apMap[uid] = ap
		//
		//	params.LeftAp = ap[util.Min(stage-1, len(ap)-1)]
		//
		//	// Update Ap Bucket
		//	encodingAp, err := json.Marshal(apMap)
		//	if err != nil {
		//		//util.Logger(util.ERROR_LEVEL, "GameReconnect", "Encoding Ap Map Err:"+err.Error())
		//		util.Error("GameReconnect Encoding Ap Map ERROR[%v]", err.Error())
		//	}
		//	boltdb.CreateOrUpdate([]byte(rid), encodingAp, "ApBucket")
		//}
		//var apMap map[string][]int
		util.Debug("耗时:%v", time.Now().UnixNano()-btime)
		util.Debug("GetAPInfo...")
		//apMap, _, err := gocache.GetAPInfo(rid)
		apMap, _, err := gocache.ConnGetAPInfo(conn, rid)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_CLEW_AP_NOT_ENOUGH, c)
			return
		}
		if len(apMap) > 0 {
			util.Debug("轮次:%v AP点:%+v", stage, apMap)
			ap, ok := apMap[uid]
			if !ok {
				util.Debug("%+v", gif.Ap)
				ap = gif.Ap
			}
			params.LeftAp = ap[util.Min(stage, len(ap)-1)]

			//if stage > 1 {
			//	var leftAp int
			//	s := util.Min(stage-1, len(ap)-1)
			//	for i := 0; i <= s; i++ {
			//		leftAp = leftAp + ap[i]
			//		ap[i] = 0
			//	}
			//	ap[s] = leftAp + ap[s]
			//}
			//gif.Ap = ap
			//apMap[uid] = ap
			//
			//params.LeftAp = ap[util.Min(stage-1, len(ap)-1)]

			//util.Debug("SetAPInfo ...")
			//err = gocache.ConnSetAPInfo(conn, rid, apMap)
			//if err != nil {
			//	util.Error("ERROR:%v", err.Error())
			//	ResponseErr(model.ERR_CLEW_AP_NOT_ENOUGH, c)
			//	return
			//}
		}
		util.Debug("耗时:%v", time.Now().UnixNano()-btime)

		params.GameInfo = gif

		//if vb != nil {
		//	// Decode Vote Map
		var votes map[string]bool
		//	err = json.Unmarshal(vb, &votes)
		//	if err != nil {
		//		//util.Logger(util.ERROR_LEVEL, "GameVote", "Decoding Votes Err:"+err.Error())
		//		util.Error("GameVote Decoding Votes ERROR[%v]", err.Error())
		//	}
		//
		//	_, ok := votes[uid]
		//	if !ok {
		//		votes[uid] = false
		//		sendGameVote = true
		//	}
		//
		//	notVoteNum := global.CalcNotVoteNum(votes)
		//	params.NotVoteNum = notVoteNum
		//
		//	// Update Vote Bucket
		//	//encodingVotes, err := json.Marshal(votes)
		//	//if err != nil {
		//	//	//util.Logger(util.ERROR_LEVEL, "GameReconnect", "Encoding Votes Err:"+err.Error())
		//	//	util.Error("GameReconnect Encoding Votes ERROR[%v]", err.Error())
		//	//}
		//	//boltdb.CreateOrUpdate([]byte(rid), encodingVotes, "VoteBucket")
		//
		//}
		util.Debug("GetVoteInfo...")
		//votes, bl, err = gocache.GetVoteInfo(rid)
		votes, bl, err = gocache.ConnGetVoteInfo(conn, rid)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_NOT_ALL_VOTED, c)
			return
		}
		if len(votes) > 0 && bl {
			_, ok := votes[uid]
			if !ok {
				votes[uid] = false
				sendGameVote = true
			}

			notVoteNum := global.CalcNotVoteNum(votes)
			params.NotVoteNum = notVoteNum
			util.Debug("SetVoteInfo...")
			//err = gocache.SetVoteInfo(rid, votes)
			err = gocache.ConnSetVoteInfo(conn, rid, votes)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_GET_GAME_INFO, c)
				return
			}
		}
		util.Debug("耗时:%v", time.Now().UnixNano()-btime)
		//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, rif.PlayerSlice)

	}
	util.Debug("耗时:%v", time.Now().UnixNano()-btime)

	// Send Room Enter Msg
	//go websocket.SendRoomEnterMessage(rif.UnionIdSlice, uid)
	//测试一下并发，先注释掉
	go websocket.SendRoomEnterMessage(rif.UnionIdSlice, user)
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, rif.PlayerSlice)

	util.Debug("耗时:%v", time.Now().UnixNano()-btime)
	//// Send Game Vote Msg
	if sendGameVote {
		go websocket.SendGameVoteMsg(rif.UnionIdSlice, params.NotVoteNum.(int))
	}
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, rif.PlayerSlice)

	reconnectInfoResp.Params = params
	etime := time.Now().UnixNano()
	//util.Debug("耗时:%v 多幕:%+v", time.Now().UnixNano()-btime, rif.PlayerSlice)

	util.Debug("重连剧本:%+v,耗时:%v", reconnectInfoResp.Params.RoomInfo.ScriptId, etime-btime)
	c.JSON(http.StatusOK, reconnectInfoResp)
}

func ClewTransform(clews []model.Clew, unionId string) {
	for i, v := range clews {
		for j, keys := range v.Key {
			for k, key := range keys.KeyList {
				if key.Status == 1 {
					if key.UnionId != unionId {
						clews[i].Key[j].KeyList[k].Status = 2
					}
				} else if key.Status == 2 {
					clews[i].Key[j].KeyList[k].Status = 3
				}
				for l, sub := range key.Sub {
					if sub.Status == 1 {
						if sub.UnionId != unionId {
							clews[i].Key[j].KeyList[k].Sub[l].Status = 2
						}
					} else if sub.Status == 2 {
						clews[i].Key[j].KeyList[k].Sub[l].Status = 3
					}
				}
			}
		}
	}
}
