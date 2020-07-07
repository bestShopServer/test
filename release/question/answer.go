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

//查询问题的答案和分值
func QueryQuestionScore(jparm model.ScriptQuestionGetReq) (questions []model.Questions, code int) {
	util.Info("Calculation Score...")
	var pids []int

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = jparm.ScriptId
	taskRequest["UnionId"] = jparm.UnionId
	pids = append(pids, jparm.PId)
	if jparm.PId != 0 {
		pids = append(pids, 0) //0标识全部
	}
	taskRequest["PIds"] = pids
	util.Info("参数[%+v]", taskRequest)

	util.Info("db.QueryScriptQuestion ...")
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptQuestionGet", taskRequest))
	if err != nil {
		util.Error("db.ScriptQuestionGet ERROR[%v]", err.Error())
		return questions, model.ERR_TASK_JSON
	}
	util.Info("获取数据成功[%+v]", dbResult)
	jtmp, err := json.Marshal(dbResult)
	if err != nil {
		util.Error("db.ScriptQuestionGet ERROR[%v]", err.Error())
		return questions, model.ERR_TASK_JSON
	}
	util.Info("获取数据成功JSON[%+v]", string(jtmp))

	dbcode, dbparams := global.UnwrapArrayPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		if len(dbparams) == 0 {
			util.Info("数据数为零!")
			return questions, global.ERR_DB_NOTFOUND_DATA
		}
		//dbparams.([]model.Questions)
		for _, p := range dbparams {
			var tmp model.Questions
			//pm := p.(map[interface{}]interface{})
			util.Debug("循环数据:%+v", p)
			jstmp, err := json.Marshal(p)
			if err != nil {
				util.Error("数据转换出错!")
				return questions, global.ERR_DB_JSON
			}
			//util.Debug("循环数据JSON:%v", string(jstmp))

			err = json.Unmarshal(jstmp, &tmp)
			if err != nil {
				util.Error("数据转换出错!")
				return questions, global.ERR_DB_JSON
			}

			questions = append(questions, tmp)
		}
		return questions, model.ERR_OK
	default:
		return questions, model.ERR_DEFAULT
	}
}

//查询问题的答案和分值
func QueryQuestionEnding(req model.QuestionEndingReq) (res model.Ending, code int) {
	util.Info("Check User Script Cost...")

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = req.ScriptId
	taskRequest["UnionId"] = req.UnionId
	taskRequest["Opt"] = req.Opt
	util.Info("UnionId[%+v]", taskRequest)

	util.Info("db.QueryScriptQuestion ...")
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.QuestionEnding", taskRequest))
	if err != nil {
		util.Error("db.ScriptQuestionGet ERROR[%v]", err.Error())
		return res, model.ERR_TASK_JSON
	}
	util.Info("获取数据成功[%+v]", dbResult)
	jtmp, err := json.Marshal(dbResult)
	if err != nil {
		util.Error("db.QueryQuestionEnding ERROR[%v]", err.Error())
		return res, model.ERR_TASK_JSON
	}
	util.Info("获取数据成功JSON[%+v]", string(jtmp))

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		util.Debug("%+v", dbparams)
		if dbparams["resume"] != nil {
			res.Resume = dbparams["resume"].(string)
		}

		return res, model.ERR_OK
	default:
		return res, model.ERR_DEFAULT
	}
}

//计算得分
func CalculationScore(c *gin.Context) {
	util.Info("CalculationScore ...")

	var req model.AnswerReq
	var score, maxScore int
	var endOpt string
	var roomId string

	err := c.Bind(&req)
	util.Info("请求参数[%+v]", req)
	//util.Logger(util.INFO_LEVEL, "CalculationScore", "CalculationScore ...")

	// Optional Fields List
	optionalFields := []string{"PId"}
	// Check Param
	if !release.CheckParams(req, "CalculationScore", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	var que model.ScriptQuestionGetReq
	que.UnionId = req.UnionId
	que.ScriptId = req.ScriptId
	if que.PId == 0 {
		util.Info("准备获取用户选的角色信息")
		//校验用户房间
		roomId, err = gocache.GetUserRoom(req.UnionId)
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
				que.PId = tmp.Role.Id
				break
			}
		}
		util.Info("角色ID:%v", que.PId)
	}
	res, code := QueryQuestionScore(que)
	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "CalculationScore", "Get CalculationScore Err")
		release.ResponseErr(model.ERR_GET_SCRIPTS, c)
		return
	}

	util.Info("开始计算分值...")
	for _, qtmp := range res { //循环问题
		for i, atmp := range req.Answer { //循环回答的问题
			if qtmp.QueId == atmp.QueId { //同一个问题
				for _, qans := range qtmp.Ans { //继续循环问题的选项
					for _, ans_id := range atmp.Answer { //继续循环问题回答的选项
						if qans.AnsId == ans_id {
							score += qans.Score
							endOpt = qans.EndOpt
						}
					}
				}
			}
			req.Answer = append(req.Answer[:i], req.Answer[i:]...)
		}
		maxScore += qtmp.MaxScore
	}
	util.Debug("所得分数:%v 结局选项:%v 理论最高得分:%v", score, endOpt, maxScore)

	var qend model.QuestionEndingReq
	qend.ScriptId = req.ScriptId
	qend.UnionId = req.UnionId
	qend.Opt = endOpt
	end, code := QueryQuestionEnding(qend)
	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "ScriptGet", "Get Scripts Err")
		release.ResponseErr(model.ERR_GET_SCRIPTS, c)
		return
	}
	util.Info("处理结束...")

	err = gocache.SetRoomUserQuestionScore(roomId, req.UnionId, score, maxScore)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_GET_SCRIPTS, c)
		return
	}

	Resp := model.AnswerResp{}
	Resp.Params.Score = score
	Resp.Params.Resume = end.Resume
	//util.Info("响应信息:%+v", Resp)

	release.ResponseOk(c, &Resp)
}
