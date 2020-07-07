package model

type RoomEnterMsg struct {
	RespType int `json:"resp_type"`
	//Data     string `json:"data"`
	Data UserInfo `json:"data"`
}

type RoleChooseMsg struct {
	RespType int            `json:"resp_type"`
	Data     RoleChooseData `json:"data"`
}

type RoleChooseData struct {
	OpenId string `json:"open_id"`
	RoleId int    `json:"role_id"`
}

type GameStartMsg struct {
	RespType int    `json:"resp_type"`
	Data     string `json:"data"`
}

type RoomExitMsg struct {
	RespType int    `json:"resp_type"`
	Data     string `json:"data"`
}

type ClewOpenMsg struct {
	RespType int         `json:"resp_type"`
	Data     ClewMsgData `json:"data"`
}

type ClewPubMsg struct {
	RespType int         `json:"resp_type"`
	Data     ClewMsgData `json:"data"`
}

type ClewMsgData struct {
	//OpenId string `json:"open_id"`
	UnionId string `json:"union_id"`
	Key     string `json:"key"`
	Index   []int  `json:"index"`
}

type StageNextMsg struct {
	RespType int `json:"resp_type"`
	Data     int `json:"data"`
}

type StepIntoVoteMsg struct {
	RespType int `json:"resp_type"`
	Data     int `json:"data"`
}

type GameVoteMsg struct {
	RespType int `json:"resp_type"`
	Data     int `json:"data"`
}

type RoomDeleteMsg struct {
	RespType int    `json:"resp_type"`
	Data     string `json:"data"`
}

func InitRoomEnterMsg() RoomEnterMsg {
	return RoomEnterMsg{
		RespType: ROOM_ENTER_MSG,
	}
}

func InitRoleChooseMsg() RoleChooseMsg {
	return RoleChooseMsg{
		RespType: ROLE_CHOOSE_MSG,
	}
}

func InitGameStartMsg() GameStartMsg {
	return GameStartMsg{
		RespType: GAME_START_MSG,
	}
}

func InitRoomExitMsg() RoomExitMsg {
	return RoomExitMsg{
		RespType: ROOM_EXIT_MSG,
	}
}

func InitClewOpenMsg() ClewOpenMsg {
	return ClewOpenMsg{
		RespType: CLEW_OPEN_MSG,
	}
}

func InitClewPubMsg() ClewPubMsg {
	return ClewPubMsg{
		RespType: CLEW_PUB_MSG,
	}
}

func InitStageNextMsg() StageNextMsg {
	return StageNextMsg{
		RespType: STAGE_NEXT_MSG,
	}
}

func InitStepIntoVoteMsg() StepIntoVoteMsg {
	return StepIntoVoteMsg{
		RespType: STEP_INTO_VOTE_MSG,
	}
}

func InitGameVoteMsg() GameVoteMsg {
	return GameVoteMsg{
		RespType: GAME_VOTE_MSG,
		Data:     -1,
	}
}

//房间解散
func InitRoomDeleteMsg() RoomDeleteMsg {
	return RoomDeleteMsg{
		RespType: ROOM_DELETE_MSG,
	}
}

//房间游戏结束
func InitGameEndMsg() RoomDeleteMsg {
	return RoomDeleteMsg{
		RespType: GAME_END_MSG,
	}
}
