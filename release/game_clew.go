package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: game.clew handler
func GameClew(c *gin.Context) {
	util.Info("GameClew ...")
	var req model.GameClewReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "GameClew", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	//uid := req.OpenId
	uid := req.UnionId
	rid := req.RoomId

	// Find Room in RoomBucket
	//b := boltdb.View([]byte(rid), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//de := json.Unmarshal(b, &roomInfo)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameClew", "Decoding Room Info Err:"+de.Error())
	//	util.Error("GameClew Decoding Room Info ERROR[%v]", de.Error())
	//}
	//读取房间信息
	err = gocache.GetRoomInfo(rid, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	exit := false
	for _, v := range roomInfo.UnionIdSlice {
		if v == uid {
			exit = true
			break
		}
	}

	if !exit {
		util.Error("ERR_NOT_IN_ROOM")
		ResponseErr(model.ERR_NOT_IN_ROOM, c)
		return
	}

	round := 0

	// Find RoomId in Stage Bucket
	var stageInfo model.StageInfo
	//sb := boltdb.View([]byte(rid), "StageBucket")
	//if sb != nil {
	//	if err := json.Unmarshal(sb, &stageInfo); err != nil {
	//		//util.Logger(util.ERROR_LEVEL, "GameClew", "Decoding Stage Info Err:"+err.Error())
	//		util.Error("GameClew Decoding Stage Info ERROR[%v]", err.Error())
	//	} else {
	//		round = stageInfo.Round
	//	}
	//} else {
	//	util.Error("ERR_NOT_GET_GAME_INFO")
	//	ResponseErr(model.ERR_NOT_GET_GAME_INFO, c)
	//	return
	//}
	_, err = gocache.GetRoomStage(rid, &stageInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_NOT_GET_GAME_INFO, c)
		return
	}
	util.Info("Round:%v", stageInfo.Round)
	round = stageInfo.Round

	gameClewResp := model.GameClewResp{}
	params := model.GameClew{}

	params, ok := GetGameClew(rid, round)
	if ok != model.ERR_OK {
		ResponseErr(ok, c)
		return
	}

	bytes, _ := json.Marshal(params.Roles)

	//util.Logger(util.INFO_LEVEL, "GameClew", string(bytes))
	util.Info("GameClew", string(bytes))

	gameClewResp.Params = params
	c.JSON(http.StatusOK, gameClewResp)
}
