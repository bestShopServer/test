package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ClewPub(c *gin.Context) {
	util.Info("ClewPub ...")

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get openId & roomId & clewId From Param
	var req model.ClewPubReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "ClewPub", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	conn := gocache.RedisConnPool.Get()
	defer conn.Close()
	//userId := req.OpenId
	userId := req.UnionId
	roomId := req.RoomId
	//clewId, round := FindClewIdInCache(req.Index, roomId, req.Key)
	clewId, round := FindClewIdInCache(conn, req.Index, roomId, req.Key)
	util.Info("线索ID:%v", clewId)
	// Find Room In Room Cache
	//_, ok := global.RoomCache[roomId]
	//if !ok || round == 0 {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	//ok, err := gocache.CheckRoomExists(roomId)
	ok, err := gocache.ConnCheckRoomExists(conn, roomId)
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

	// Find ClewInfo By RoomId In GameBucket
	//gb := boltdb.View([]byte(roomId), "GameBucket")
	//if gb == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode ClewInfo
	var clew []model.Clew
	//err = json.Unmarshal(gb, &clew)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "ClewPub", "Decode Clew Err:"+err.Error())
	//	util.Error("ClewPub Decode Clew ERROR[%v]", err.Error())
	//	return
	//}
	//_, err = gocache.GetRoomGame(roomId, &clew)
	_, err = gocache.ConnGetRoomGame(conn, roomId, &clew)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find KV in Room Bucket
	//b := boltdb.View([]byte(roomId), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Room Info
	rif := model.RoomInfo{}
	//de := json.Unmarshal(b, &rif)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "RoomEnter", "Decoding Room Info Err:"+de.Error())
	//	util.Error("RoomEnter Decoding Room Info ERROR[%v]", de.Error())
	//}
	//读取房间信息
	//err = gocache.GetRoomInfo(roomId, &rif)
	err = gocache.ConnGetRoomInfo(conn, roomId, &rif)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find RoomId in StageBucket
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
	//	//util.Logger(util.ERROR_LEVEL, "ClewOpen", "Decoding Stage Info Err:"+err.Error())
	//	util.Error("ClewOpen Decoding Stage Info ERROR[%v]", err.Error())
	//}
	//_, err = gocache.GetRoomStage(roomId, &stageInfo)
	_, err = gocache.ConnGetRoomStage(conn, roomId, &stageInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	if round > stageInfo.Round {
		util.Error("ERR_WRONG_FORMAT")
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	for iParent := 0; iParent < len(clew); iParent++ {
		for pKey, pValue := range clew[iParent].Key {
			for jParent := 0; jParent < len(pValue.KeyList); jParent++ {
				// Compare Parent Clew
				if clewId == pValue.KeyList[jParent].Id {
					// Parent Clew Not Open
					if pValue.KeyList[jParent].Status == model.STATUS_UNOPEN {
						util.Error("ERR_CLEW_NOT_OPENED")
						ResponseErr(model.ERR_CLEW_NOT_OPENED, c)
						return
					}

					// Clew Not You Opened
					if userId != pValue.KeyList[jParent].UnionId {
						util.Error("ERR_CLEW_NOT_YOU_OPENED")
						ResponseErr(model.ERR_CLEW_NOT_YOU_OPENED, c)
						return
					}

					// Parent Clew has Pubbed
					if pValue.KeyList[jParent].Status == model.STATUS_PUB {
						util.Error("ERR_CLEW_HAS_PUBBED")
						ResponseErr(model.ERR_CLEW_HAS_PUBBED, c)
						return
					}

					// Update GameBucket
					clew[iParent].Key[pKey].KeyList[jParent].Status = model.STATUS_PUB
					//clewEncoding, _ := json.Marshal(clew)
					//boltdb.CreateOrUpdate([]byte(roomId), clewEncoding, "GameBucket")
					//err = gocache.SetRoomGame(roomId, clew)
					err = gocache.ConnSetRoomGame(conn, roomId, clew)
					if err != nil {
						util.Error("ERROR:%v", err.Error())
						ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
						return
					}

					// Send Clew Open Msg
					go websocket.SendClewPubMessage(rif.UnionIdSlice, userId, req.Key, req.Index)

					// Return Room Info
					util.Info("success")
					resp := model.ClewPubResp{}
					resp.Params = "success"
					c.JSON(http.StatusOK, resp)
					return
				}

				// Compare Child Clew If Exist
				childClew := pValue.KeyList[jParent].Sub
				if childClew != nil {
					for iChild := 0; iChild < len(childClew); iChild++ {
						if clewId == childClew[iChild].Id {
							// Parent Clew Not Open
							if pValue.KeyList[jParent].Status == model.STATUS_UNOPEN {
								util.Error("ERR_PARENT_CLEW_NOT_OPEN")
								ResponseErr(model.ERR_PARENT_CLEW_NOT_OPEN, c)
								return
							}

							// Parent Clew Not Pub
							if pValue.KeyList[jParent].Status != model.STATUS_PUB {
								util.Error("ERR_PARENT_CLEW_NOT_PUB")
								ResponseErr(model.ERR_PARENT_CLEW_NOT_PUB, c)
								return
							}

							// Child Clew Not Open
							if childClew[iChild].Status == model.STATUS_UNOPEN {
								util.Error("ERR_CLEW_NOT_OPENED")
								ResponseErr(model.ERR_CLEW_NOT_OPENED, c)
								return
							}

							// Clew Not You Opened
							if userId != childClew[iChild].UnionId {
								util.Error("ERR_CLEW_NOT_YOU_OPENED")
								ResponseErr(model.ERR_CLEW_NOT_YOU_OPENED, c)
								return
							}

							// Clew Has Pubbed
							if childClew[iChild].Status == model.STATUS_PUB {
								util.Error("ERR_CLEW_HAS_PUBBED")
								ResponseErr(model.ERR_CLEW_HAS_PUBBED, c)
								return
							}

							// Update GameBucket
							clew[iParent].Key[pKey].KeyList[jParent].Sub[iChild].Status = model.STATUS_PUB
							//clewEncoding, _ := json.Marshal(clew)
							//boltdb.CreateOrUpdate([]byte(roomId), clewEncoding, "GameBucket")
							//err = gocache.SetRoomGame(roomId, clew)
							err = gocache.ConnSetRoomGame(conn, roomId, clew)
							if err != nil {
								util.Error("ERROR:%v", err.Error())
								ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
								return
							}

							// Send Clew Open Msg
							go websocket.SendClewPubMessage(rif.UnionIdSlice, userId, req.Key, req.Index)

							// Return Room Info
							util.Info("success")
							resp := model.ClewPubResp{}
							resp.Params = "success"
							c.JSON(http.StatusOK, resp)
							return
						}

						//兼容三级线索
						if childClew[iChild].Sub != nil {
							util.Debug("三级线索:%v", childClew[iChild])
							for m := 0; m < len(childClew[iChild].Sub); m++ {
								if clewId == childClew[iChild].Sub[m].Id {
									// Parent Clew Not Open
									if pValue.KeyList[jParent].Status == model.STATUS_UNOPEN {
										util.Error("ERR_PARENT_CLEW_NOT_OPEN")
										ResponseErr(model.ERR_PARENT_CLEW_NOT_OPEN, c)
										return
									}

									// Parent Clew Not Pub
									if pValue.KeyList[jParent].Status != model.STATUS_PUB {
										util.Error("ERR_PARENT_CLEW_NOT_PUB")
										ResponseErr(model.ERR_PARENT_CLEW_NOT_PUB, c)
										return
									}

									// Child Clew Not Open
									if childClew[iChild].Status == model.STATUS_UNOPEN {
										util.Error("ERR_CLEW_NOT_OPENED")
										ResponseErr(model.ERR_CLEW_NOT_OPENED, c)
										return
									}

									// Clew Has Pubbed
									if childClew[iChild].Status != model.STATUS_PUB {
										util.Error("ERR_CLEW_HAS_PUBBED")
										ResponseErr(model.ERR_CLEW_HAS_PUBBED, c)
										return
									}

									// Grandson Clew Not Open
									if childClew[iChild].Sub[m].Status == model.STATUS_UNOPEN {
										util.Error("ERR_CLEW_NOT_OPENED")
										ResponseErr(model.ERR_CLEW_NOT_OPENED, c)
										return
									}

									// Grandson Not You Opened
									if userId != childClew[iChild].Sub[m].UnionId {
										util.Error("ERR_CLEW_NOT_YOU_OPENED")
										ResponseErr(model.ERR_CLEW_NOT_YOU_OPENED, c)
										return
									}

									// Grandson Has Pubbed
									if childClew[iChild].Sub[m].Status == model.STATUS_PUB {
										util.Error("ERR_CLEW_HAS_PUBBED")
										ResponseErr(model.ERR_CLEW_HAS_PUBBED, c)
										return
									}

									// Update GameBucket
									clew[iParent].Key[pKey].KeyList[jParent].Sub[iChild].Sub[m].Status = model.STATUS_PUB
									//clewEncoding, _ := json.Marshal(clew)
									//boltdb.CreateOrUpdate([]byte(roomId), clewEncoding, "GameBucket")
									//err = gocache.SetRoomGame(roomId, clew)
									err = gocache.ConnSetRoomGame(conn, roomId, clew)
									if err != nil {
										util.Error("ERROR:%v", err.Error())
										ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
										return
									}

									// Send Clew Open Msg
									go websocket.SendClewPubMessage(rif.UnionIdSlice, userId, req.Key, req.Index)

									// Return Room Info
									util.Info("success")
									resp := model.ClewPubResp{}
									resp.Params = "success"
									c.JSON(http.StatusOK, resp)
									return
								}
							}
						}
					}
				}
			}
		}
	}
	util.Info("FAIL")
	// Clew Not Found
	ResponseErr(model.ERR_CLEW_NOT_FOUND, c)
	return
}
