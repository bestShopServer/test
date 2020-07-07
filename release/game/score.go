package game

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func RecordScore(param model.GameScoreReq) (code int) {
	util.Info("Record  Script Score ...")

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = param.ScriptId
	taskRequest["UnionId"] = param.UnionId
	taskRequest["RoomId"] = strconv.Itoa(param.RoomId)
	taskRequest["Score"] = param.GameScore
	taskRequest["ScriptScore"] = param.ScriptScore

	util.Info("UnionId[%+v]", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.wx.ScriptComment", taskRequest))
	if err != nil {
		util.Error("db.ScriptQuestionGet ERROR[%v]", err.Error())
		return model.ERR_TASK_JSON
	}
	util.Info("获取数据成功[%+v]", dbResult)
	dbcode, _ := global.UnwrapArrayPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		return model.ERR_OK
	default:
		return model.ERR_DEFAULT
	}
}

//游戏评分
func GameScore(c *gin.Context) {
	util.Info("GameScore ...")

	var req model.GameScoreReq

	err := c.Bind(&req)
	util.Info("请求参数[%+v]", req)
	//util.Logger(util.INFO_LEVEL, "CalculationScore", "CalculationScore ...")

	// Optional Fields List
	optionalFields := []string{}
	// Check Param
	if !release.CheckParams(req, "CalculationScore", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	////获取用户所在房间剧本
	//// If no user cache
	//room_id, ok := global.UserCache[req.UnionId]
	//if !ok {
	//	util.Error("ERR_GAME_HAS_OVER")
	//	release.ResponseErr(model.ERR_GAME_HAS_OVER, c)
	//	return
	//}
	//
	//// If Game Has Over
	//_, ok = global.RoomCache[room_id]
	//if !ok {
	//	util.Error("ERR_GAME_HAS_OVER")
	//	release.ResponseErr(model.ERR_GAME_HAS_OVER, c)
	//	return
	//}
	//
	//// Find KV in Room Bucket
	//b := boltdb.View([]byte(room_id), "RoomBucket")
	//if b == nil {
	//	release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	//
	//// Decode Room Info
	//rif := model.RoomInfo{}
	//de := json.Unmarshal(b, &rif)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameReconnect", "Decoding Room Info Err:"+de.Error())
	//	util.Error("GameReconnect Decoding Room Info ERROR[%v]", de.Error())
	//}
	//util.Debug("房间信息:%+v", rif)

	code := RecordScore(req)
	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "RecordScore", "Record Score Err")
		release.ResponseErr(model.ERR_GET_SCRIPTS, c)
		return
	}
	util.Info("处理结束...")

	Resp := model.ErrResp{}
	release.ResponseOk(c, &Resp)
}
