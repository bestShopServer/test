package release

//
//// Func: user.login handler
//func UserLogin(c *gin.Context) {
//	util.Info("UserLogin ...")
//
//	// Get Code Param
//	var req model.UserLoginReq
//	err := c.Bind(&req)
//
//	// Optional Fields List
//	var optionalFields []string
//
//	// Check Param
//	if !CheckParams(req, "UserLogin", err, optionalFields) {
//		ResponseErr(model.ERR_WRONG_FORMAT, c)
//		return
//	}
//
//	// Get OpenId From Tencent
//	config := conf.GetConfig()
//	openId, err := util.GetOpenId(config.AppId, config.AppSecret, req.Code, req.EncryptData, req.IV)
//	if err != nil {
//		ResponseErr(model.ERR_OPENID_FAILED, c)
//		return
//	}
//
//	// Save OpenId in Database if it Not Exist
//	ob := boltdb.View([]byte(openId), "UserBucket")
//	if ob == nil {
//		ct := time.Now().Format("2006-01-02 15:04:05")
//		boltdb.CreateOrUpdate([]byte(openId), []byte(ct), "UserBucket")
//	}
//
//	// Check If User In A Room Not Over Yet
//	roomId := global.UserCache[openId]
//
//	// Notify WebSocket User Login
//	go websocket.SendUserLoginMessage(openId)
//
//	// Return OpenId
//	c.JSON(http.StatusOK, gin.H{
//		"open_id": openId,
//		"room_id": roomId,
//	})
//}
