package model

type RoomRecordBase struct {
	RoomId   string
	Owner    string
	UnionId  string
	ScriptId int
	RoleId   int
	Status   int
	Score    int
}

type MemberBaseInfo struct {
	Uid     int    `json:"uid"  redis:"uid"`
	Member  int    `json:"member"  redis:"member"`
	EffTime string `json:"eff_time"  redis:"eff_time"`
	InvTime string `json:"inv_time"  redis:"inv_time"`
}
