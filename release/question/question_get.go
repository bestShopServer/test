package question

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

//func QueryScriptQuestion(jparm model.ScriptQuestionGetReq) (questions []model.Questions, code int) {
func QueryScriptQuestion(jparm model.ScriptQuestionGetReq) (questions model.QuestionAndAnswer, code int) {
	util.Info("Check User Script Cost...")
	var pids []int

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = jparm.ScriptId
	taskRequest["UnionId"] = jparm.UnionId
	pids = append(pids, jparm.PId)
	if jparm.PId != 0 {
		pids = append(pids, 0) //0标识全部
	}
	taskRequest["PIds"] = pids
	util.Info("请求参数[%+v]", taskRequest)

	util.Info("db.ScriptQuestionRelation ...")
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptQuestionRelation", taskRequest))
	if err != nil {
		util.Error("db.ScriptQuestionGet ERROR[%v]", err.Error())
		return questions, model.ERR_TASK_JSON
	}
	util.Info("获取数据成功[%+v]", len(dbResult))
	//jtmp, err := json.Marshal(dbResult)
	//if err != nil {
	//	util.Error("db.ScriptQuestionGet ERROR[%v]", err.Error())
	//	return questions, model.ERR_TASK_JSON
	//}
	//util.Info("获取数据成功JSON[%+v]", string(jtmp))

	//dbcode, dbparams := global.UnwrapArrayPackage(dbResult)
	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		//if len(dbparams) == 0 {
		//	util.Info("数据数为零!")
		//	return questions, global.ERR_DB_NOTFOUND_DATA
		//}
		//dbparams.([]model.Questions)
		//for _, p := range dbparams {
		//	var tmp model.Questions
		//	//pm := p.(map[interface{}]interface{})
		//	//util.Debug("循环数据:%+v", p)
		//	jstmp, err := json.Marshal(p)
		//	if err != nil {
		//		util.Error("数据转换出错!")
		//		return questions, global.ERR_DB_JSON
		//	}
		//	//util.Debug("循环数据JSON:%v", string(jstmp))
		//
		//	err = json.Unmarshal(jstmp, &tmp)
		//	if err != nil {
		//		util.Error("数据转换出错!")
		//		return questions, global.ERR_DB_JSON
		//	}
		//
		//	questions = append(questions, tmp)
		//}

		jstmp, err := json.Marshal(dbparams)
		if err != nil {
			util.Error("数据转换出错!")
			return questions, global.ERR_DB_JSON
		}
		//util.Debug("循环数据JSON:%v", string(jstmp))

		err = json.Unmarshal(jstmp, &questions)
		if err != nil {
			util.Error("数据转换出错!")
			return questions, global.ERR_DB_JSON
		}
		return questions, model.ERR_OK
	default:
		return questions, model.ERR_DEFAULT
	}
}

//获取问题
func ScriptQuestionGet(c *gin.Context) {
	util.Info("ScriptQuestionGet ...")

	var req model.ScriptQuestionGetReq

	err := c.Bind(&req)
	util.Info("请求参数[%+v]", req)
	util.Logger(util.INFO_LEVEL, "ScriptQuestionGet", "ScriptQuestionGet ...")

	// Optional Fields List
	optionalFields := []string{"PId"}
	// Check Param
	if !release.CheckParams(req, "ScriptQuestionGet", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	if req.PId == 0 {
		util.Info("准备获取用户选的角色信息")
		//校验用户房间
		roomId, err := gocache.GetUserRoom(req.UnionId)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			release.ResponseErr(model.ERR_GAME_HAS_OVER, c)
			return
		}

		//获取房间信息
		roomInfo := model.RoomInfo{}
		err = gocache.GetRoomInfo(roomId, &roomInfo)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
			return
		}

		for _, tmp := range roomInfo.PlayerSlice {
			if tmp.UnionId == req.UnionId {
				req.PId = tmp.Role.Id
				break
			}
		}
		util.Info("角色ID:%v", req.PId)
	}

	res, code := QueryScriptQuestion(req)
	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "ScriptQuestionGet", "Get Script Question Err")
		release.ResponseErr(model.ERR_GET_SCRIPTS, c)
		return
	}
	util.Debug("处理结束")

	Resp := model.ScriptQuestionGetResp{}
	Resp.Params = res
	//util.Info("响应信息:%+v", Resp)

	release.ResponseOk(c, &Resp)
}
