package release

import (
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: Gin Ok Response
func ResponseOk(c *gin.Context, it interface{}) {
	//utils.ErrorLog("Response Error [%v][%v]", errNo, models.RspMsg[errNo] )
	// Generate Json String Response
	jit, err := json.Marshal(it)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		panic(err)
	}
	util.Info("Response[%v]", string(jit))
	//util.Logger(util.INFO_LEVEL, "response[%v]", it)
	c.JSON(http.StatusOK, it)
}

// Func: Gin Ok Response
func ResponseSuccess(c *gin.Context) {
	//utils.ErrorLog("Response Error [%v][%v]", errNo, models.RspMsg[errNo] )
	// Generate Json String Response
	c.JSON(http.StatusOK, &model.ErrResp{model.ERR_OK, model.ErrMap[model.ERR_OK]})
}

// Func: Gin Ok Response
func ResponseErrMsg(c *gin.Context, msg string) {
	//util.Logger(util.INFO_LEVEL, "Response Error:", msg)
	util.Info("MSG:%v", msg)
	// Generate Json String Response
	c.JSON(http.StatusOK, &model.ErrResp{model.ERR_DEFAULT, msg})
}
