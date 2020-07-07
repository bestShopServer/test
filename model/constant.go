package model

const (
	DEFAULT_MSG = iota
	ROOM_ENTER_MSG
	ROLE_CHOOSE_MSG
	GAME_START_MSG
	ROOM_EXIT_MSG
	CLEW_OPEN_MSG
	CLEW_PUB_MSG
	STAGE_NEXT_MSG
	STEP_INTO_VOTE_MSG
	GAME_VOTE_MSG
	ROOM_DELETE_MSG
	GAME_END_MSG
)

const (
	ERR_OK                     = 0
	ERR_DEFAULT                = 4000
	ERR_WRONG_FORMAT           = 4001
	ERR_ROOM_NOT_EXIST         = 4002
	ERR_ROOM_IS_FULL           = 4003
	ERR_ENTERED_ROOM           = 4004
	ERR_BELONG                 = 4005
	ERR_ROLE_SELECT            = 4006
	ERR_ROLE_SELECTED          = 4007
	ERR_SOMEONE_NOT_SELECT     = 4008
	ERR_ROLE_NOT_SELECT        = 4009
	ERR_CLEW_NOT_YOU_OPENED    = 4010
	ERR_CLEW_NOT_FOUND         = 4011
	ERR_CLEW_NOT_OPENED        = 4012
	ERR_CLEW_HAS_OPENED        = 4013
	ERR_CLEW_HAS_PUBBED        = 4014
	ERR_CLEW_PASSWORD_WRONG    = 4015
	ERR_CLEW_AP_NOT_ENOUGH     = 4016
	ERR_PARENT_CLEW_NOT_OPEN   = 4017
	ERR_PARENT_CLEW_NOT_PUB    = 4018
	ERR_GAME_HAS_OVER          = 4019
	ERR_STAGE_NEXT_HAS_CLICKED = 4020
	ERR_HAS_VOTED              = 4021
	ERR_NOT_ALL_VOTED          = 4022
	ERR_SERVER_ABNORMAL        = 4023
	ERR_OPENID_FAILED          = 4024
	ERR_TASK_JSON              = 4025
	ERR_GET_REPORT             = 4026
	ERR_GET_ROLE_INFO          = 4027
	ERR_GET_SCRIPTS            = 4028
	ERR_GET_GAME_INFO          = 4029
	ERR_NOT_IN_ROOM            = 4030
	ERR_NOT_GET_GAME_INFO      = 4031
	ERR_ROOM_LINK              = 4032
	ERR_CANNOT_OPEN_OWN_CLEW   = 4033
	ERR_CHECK_SCRIPT_COST      = 4033
	ERR_NOT_ROOM_OWNER         = 4034
	ERR_NOT_KICK_ROOM_OWNER    = 4035
	ERR_USER_BASE_INFO_SYNC    = 4036
	ERR_ROOM_NOT_OVER          = 4037
	ERR_ROOM_ALREADY_OVER      = 4038
	ERR_ROOM_DELETE            = 4039
	ERR_SCRIPT_ORDER_SYNC      = 4040
)

const (
	CLEW_NORMAL   = 1
	CLEW_URL      = 2
	CLEW_PASSWORD = 3
)

const (
	STATUS_UNOPEN = 0
	STATUS_OPEN   = 1
	STATUS_PUB    = 2
)

var ErrMap = map[int]string{
	ERR_DEFAULT:                "网络错误",
	ERR_WRONG_FORMAT:           "格式有误",
	ERR_ROOM_NOT_EXIST:         "房间不存在",
	ERR_ROOM_IS_FULL:           "房间已满",
	ERR_ENTERED_ROOM:           "您已进入房间",
	ERR_BELONG:                 "You Are Not In This Room",
	ERR_ROLE_SELECT:            "You Has Selected A Role",
	ERR_ROLE_SELECTED:          "Role Has Been Selected",
	ERR_SOMEONE_NOT_SELECT:     "Someone Has Not Select Role",
	ERR_ROLE_NOT_SELECT:        "You Has Not Select A Role",
	ERR_CLEW_NOT_YOU_OPENED:    "Clew Is Not Your's",
	ERR_CLEW_NOT_FOUND:         "Clew Is Not Found",
	ERR_CLEW_NOT_OPENED:        "Clew Has Not Been Opened",
	ERR_CLEW_HAS_OPENED:        "Clew Has Been Opened",
	ERR_CLEW_HAS_PUBBED:        "Clew Has Been Pubbed",
	ERR_CLEW_PASSWORD_WRONG:    "Clew Password Is Not Correct",
	ERR_CLEW_AP_NOT_ENOUGH:     "Ap Is Not Enough",
	ERR_PARENT_CLEW_NOT_OPEN:   "Parent Clew Is Not Open",
	ERR_PARENT_CLEW_NOT_PUB:    "Parent Clew Is Not Pub",
	ERR_GAME_HAS_OVER:          "Game Has Over",
	ERR_STAGE_NEXT_HAS_CLICKED: "You Has Clicked Next Stage",
	ERR_HAS_VOTED:              "You Has Voted",
	ERR_NOT_ALL_VOTED:          "Not All User Has Voted",
	ERR_SERVER_ABNORMAL:        "Server Abnormal",
	ERR_OPENID_FAILED:          "Get OpenId Failed",
	ERR_TASK_JSON:              "Task Json Failed",
	ERR_GET_REPORT:             "Get Report Url By Script Id Failed",
	ERR_GET_ROLE_INFO:          "Get Player Info By Script Id Failed",
	ERR_GET_SCRIPTS:            "Get Scripts Failed",
	ERR_GET_GAME_INFO:          "Get Game Info Failed",
	ERR_NOT_IN_ROOM:            "You Are Not In Room",
	ERR_NOT_GET_GAME_INFO:      "You Have Not Get Game Info Yet",
	ERR_ROOM_LINK:              "Room Link Has Expired",
	ERR_CANNOT_OPEN_OWN_CLEW:   "Can't Open Your Own Clews",
	ERR_NOT_ROOM_OWNER:         "非房主，无权操作",
	ERR_NOT_KICK_ROOM_OWNER:    "不能踢出房主",
	ERR_USER_BASE_INFO_SYNC:    "同步用户信息",
	ERR_ROOM_NOT_OVER:          "上一房间游戏未结束，不允许加入新房间",
	ERR_ROOM_ALREADY_OVER:      "房间游戏已结束，无法再加入房间",
	ERR_ROOM_DELETE:            "删除房间失败",
	ERR_SCRIPT_ORDER_SYNC:      "同步用户剧本订单失败",
}
