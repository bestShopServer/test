package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StageNext(c *gin.Context) {
	util.Info("StageNext ...")
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get openId & roomId From Param
	var req model.StageNextReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "StageNext", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	roomId := req.RoomId
	//userId := req.OpenId
	userId := req.UnionId

	// Find Room In Room Cache
	//_, ok := global.RoomCache[roomId]
	//if !ok {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	ok, err := gocache.CheckRoomExists(roomId)
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

	// Find RoomId in Room Bucket
	//rb := boltdb.View([]byte(roomId), "RoomBucket")
	//if rb == nil {
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
		if tmp == userId {
			isWhite = true
		}
	}
	util.Info("是否在白名单:%v", isWhite)

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//err = json.Unmarshal(rb, &roomInfo)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "StageNext", "Decoding Room Info Err:"+err.Error())
	//	util.Error("StageNext Decoding Room Info ERROR[%v]", err.Error())
	//}
	//读取房间信息
	err = gocache.GetRoomInfo(roomId, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// If You Are Not In this Room
	exist := false
	for _, v := range roomInfo.UnionIdSlice {
		if v == userId {
			exist = true
			break
		}
	}
	util.Debug("exist:%v", exist)
	if !exist {
		util.Error("ERR_BELONG")
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	// Find RoomId in Ap Bucket
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
	//	//util.Logger(util.ERROR_LEVEL, "StageNext", "Decoding Ap Info Err:"+err.Error())
	//	util.Error("StageNext Decoding Ap Info ERROR[%v]", err.Error())
	//}
	//获取
	apMap, _, err = gocache.GetAPInfo(roomId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find RoomId in Stage Bucket
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
	//	//util.Logger(util.ERROR_LEVEL, "StageNext", "Decoding Stage Info Err:"+err.Error())
	//	util.Error("StageNext Decoding Stage Info ERROR[%v]", err.Error())
	//}
	_, err = gocache.GetRoomStage(roomId, &stageInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Debug("stage:%v", stageInfo.Round)

	// Find This Round OpenId Click State
	if stageInfo.OnClick[stageInfo.Round] == nil {
		stageInfo.OnClick[stageInfo.Round] = append(stageInfo.OnClick[stageInfo.Round], userId)
	} else {
		for _, v := range stageInfo.OnClick[stageInfo.Round] {
			if v == userId {
				util.Error("ERR_STAGE_NEXT_HAS_CLICKED:%v", userId)
				ResponseErr(model.ERR_STAGE_NEXT_HAS_CLICKED, c)
				return
			}
		}
		stageInfo.OnClick[stageInfo.Round] = append(stageInfo.OnClick[stageInfo.Round], userId)

		if isWhite { //白名单处理
			util.Debug("白名单:%+v", roomInfo.UnionIdSlice)
			for _, v := range roomInfo.UnionIdSlice {
				isIn := false //是否在点击数组中
				for _, tmp := range stageInfo.OnClick[stageInfo.Round] {
					//util.Debug("v:%v tmp:%v", v, tmp)
					if v == tmp {
						isIn = true
						break
					}
				}
				util.Debug("%v", isIn)
				if !isIn {
					util.Debug("自动填充点击:%v", v)
					stageInfo.OnClick[stageInfo.Round] = append(stageInfo.OnClick[stageInfo.Round], v)
				}
			}
		}
	}
	util.Debug("len click:%v", len(stageInfo.OnClick))

	// If Round > Map Key Len
	if stageInfo.Round > len(stageInfo.OnClick) {
		util.Error("ERR_GAME_HAS_OVER")
		ResponseErr(model.ERR_GAME_HAS_OVER, c)
		return
	}

	// Update Ap Info	关闭下一阶段限制问题，点击下一阶段可以搜证
	//_, ok = apMap[userId]
	//if ok {
	//	var leftAp int
	//	for i := 0; i < stageInfo.Round; i++ {
	//		leftAp += apMap[userId][i]
	//		apMap[userId][i] = 0 //上一轮的AP点清零
	//	}
	//
	//	if stageInfo.Round < len(apMap[userId]) {
	//		apMap[userId][stageInfo.Round] = apMap[userId][stageInfo.Round] + leftAp
	//	}
	//}

	util.Info("onclick:%+v  unionids:%+v", stageInfo.OnClick[stageInfo.Round], roomInfo.UnionIdSlice)
	// If All People Click Next Stage Then Send Stage Next Msg
	//if len(stageInfo.OnClick[stageInfo.Round]) >= len(roomInfo.UnionIdSlice) {
	if len(stageInfo.OnClick[stageInfo.Round]) >= len(roomInfo.UserInfoSlice) { //房间用户数
		stageInfo.Round++
		//if stageInfo.Round > len(stageInfo.OnClick) {
		//if stageInfo.Round >= len(stageInfo.OnClick) { //处理前置0轮情况
		if stageInfo.Round > roomInfo.Round {
			votes := make(map[string]bool)
			for _, v := range roomInfo.UnionIdSlice {
				votes[v] = false
			}
			//encodingVotes, _ := json.Marshal(votes)
			//boltdb.CreateOrUpdate([]byte(roomId), encodingVotes, "VoteBucket")
			err = gocache.SetVoteInfo(roomId, votes)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_HAS_VOTED, c)
				return
			}
			util.Debug("SendStepIntoGameVoteMsg...")
			go websocket.SendStepIntoGameVoteMsg(roomInfo.UnionIdSlice)
		} else {
			util.Debug("SendStageNextMsg...")
			//AP点挪到下一阶段
			for _, tmp := range roomInfo.UserInfoSlice {
				_, ok = apMap[tmp.UnionId]
				if ok {
					var leftAp int
					for i := 0; i < stageInfo.Round; i++ {
						leftAp += apMap[tmp.UnionId][i]
						apMap[tmp.UnionId][i] = 0 //上一轮的AP点清零
					}

					if stageInfo.Round < len(apMap[tmp.UnionId]) {
						apMap[tmp.UnionId][stageInfo.Round] = apMap[tmp.UnionId][stageInfo.Round] + leftAp
					}
				}
			}
			util.Debug("SetAPInfo...")
			err = gocache.SetAPInfo(roomId, apMap)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_DEFAULT, c)
				return
			}

			go websocket.SendStageNextMsg(roomInfo.UnionIdSlice)
		}
	}

	// Update Bolt DB
	//encodingApInfo, _ := json.Marshal(apMap)
	//boltdb.CreateOrUpdate([]byte(roomId), encodingApInfo, "ApBucket")
	//encodingStageInfo, _ := json.Marshal(stageInfo)
	//boltdb.CreateOrUpdate([]byte(roomId), encodingStageInfo, "StageBucket")
	util.Debug("%+v", stageInfo)

	err = gocache.SetRoomStage(roomId, stageInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_DEFAULT, c)
		return
	}

	util.Info("success")
	resp := model.StageNextResp{}
	resp.Params.GameStage = stageInfo.Round
	//resp.Params = "success"
	c.JSON(http.StatusOK, &resp)
}
